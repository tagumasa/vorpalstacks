// This is a generated file - do not edit.
//
// Generated from ssm.proto.

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

@$core.Deprecated('Use accessRequestStatusDescriptor instead')
const AccessRequestStatus$json = {
  '1': 'AccessRequestStatus',
  '2': [
    {'1': 'ACCESS_REQUEST_STATUS_PENDING', '2': 0},
    {'1': 'ACCESS_REQUEST_STATUS_REVOKED', '2': 1},
    {'1': 'ACCESS_REQUEST_STATUS_REJECTED', '2': 2},
    {'1': 'ACCESS_REQUEST_STATUS_APPROVED', '2': 3},
    {'1': 'ACCESS_REQUEST_STATUS_EXPIRED', '2': 4},
  ],
};

/// Descriptor for `AccessRequestStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accessRequestStatusDescriptor = $convert.base64Decode(
    'ChNBY2Nlc3NSZXF1ZXN0U3RhdHVzEiEKHUFDQ0VTU19SRVFVRVNUX1NUQVRVU19QRU5ESU5HEA'
    'ASIQodQUNDRVNTX1JFUVVFU1RfU1RBVFVTX1JFVk9LRUQQARIiCh5BQ0NFU1NfUkVRVUVTVF9T'
    'VEFUVVNfUkVKRUNURUQQAhIiCh5BQ0NFU1NfUkVRVUVTVF9TVEFUVVNfQVBQUk9WRUQQAxIhCh'
    '1BQ0NFU1NfUkVRVUVTVF9TVEFUVVNfRVhQSVJFRBAE');

@$core.Deprecated('Use accessTypeDescriptor instead')
const AccessType$json = {
  '1': 'AccessType',
  '2': [
    {'1': 'ACCESS_TYPE_JUSTINTIME', '2': 0},
    {'1': 'ACCESS_TYPE_STANDARD', '2': 1},
  ],
};

/// Descriptor for `AccessType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List accessTypeDescriptor = $convert.base64Decode(
    'CgpBY2Nlc3NUeXBlEhoKFkFDQ0VTU19UWVBFX0pVU1RJTlRJTUUQABIYChRBQ0NFU1NfVFlQRV'
    '9TVEFOREFSRBAB');

@$core.Deprecated('Use associationComplianceSeverityDescriptor instead')
const AssociationComplianceSeverity$json = {
  '1': 'AssociationComplianceSeverity',
  '2': [
    {'1': 'ASSOCIATION_COMPLIANCE_SEVERITY_MEDIUM', '2': 0},
    {'1': 'ASSOCIATION_COMPLIANCE_SEVERITY_UNSPECIFIED', '2': 1},
    {'1': 'ASSOCIATION_COMPLIANCE_SEVERITY_CRITICAL', '2': 2},
    {'1': 'ASSOCIATION_COMPLIANCE_SEVERITY_LOW', '2': 3},
    {'1': 'ASSOCIATION_COMPLIANCE_SEVERITY_HIGH', '2': 4},
  ],
};

/// Descriptor for `AssociationComplianceSeverity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationComplianceSeverityDescriptor = $convert.base64Decode(
    'Ch1Bc3NvY2lhdGlvbkNvbXBsaWFuY2VTZXZlcml0eRIqCiZBU1NPQ0lBVElPTl9DT01QTElBTk'
    'NFX1NFVkVSSVRZX01FRElVTRAAEi8KK0FTU09DSUFUSU9OX0NPTVBMSUFOQ0VfU0VWRVJJVFlf'
    'VU5TUEVDSUZJRUQQARIsCihBU1NPQ0lBVElPTl9DT01QTElBTkNFX1NFVkVSSVRZX0NSSVRJQ0'
    'FMEAISJwojQVNTT0NJQVRJT05fQ09NUExJQU5DRV9TRVZFUklUWV9MT1cQAxIoCiRBU1NPQ0lB'
    'VElPTl9DT01QTElBTkNFX1NFVkVSSVRZX0hJR0gQBA==');

@$core.Deprecated('Use associationExecutionFilterKeyDescriptor instead')
const AssociationExecutionFilterKey$json = {
  '1': 'AssociationExecutionFilterKey',
  '2': [
    {'1': 'ASSOCIATION_EXECUTION_FILTER_KEY_CREATEDTIME', '2': 0},
    {'1': 'ASSOCIATION_EXECUTION_FILTER_KEY_STATUS', '2': 1},
    {'1': 'ASSOCIATION_EXECUTION_FILTER_KEY_EXECUTIONID', '2': 2},
  ],
};

/// Descriptor for `AssociationExecutionFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationExecutionFilterKeyDescriptor = $convert.base64Decode(
    'Ch1Bc3NvY2lhdGlvbkV4ZWN1dGlvbkZpbHRlcktleRIwCixBU1NPQ0lBVElPTl9FWEVDVVRJT0'
    '5fRklMVEVSX0tFWV9DUkVBVEVEVElNRRAAEisKJ0FTU09DSUFUSU9OX0VYRUNVVElPTl9GSUxU'
    'RVJfS0VZX1NUQVRVUxABEjAKLEFTU09DSUFUSU9OX0VYRUNVVElPTl9GSUxURVJfS0VZX0VYRU'
    'NVVElPTklEEAI=');

@$core.Deprecated('Use associationExecutionTargetsFilterKeyDescriptor instead')
const AssociationExecutionTargetsFilterKey$json = {
  '1': 'AssociationExecutionTargetsFilterKey',
  '2': [
    {'1': 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_STATUS', '2': 0},
    {'1': 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCETYPE', '2': 1},
    {'1': 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCEID', '2': 2},
  ],
};

/// Descriptor for `AssociationExecutionTargetsFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationExecutionTargetsFilterKeyDescriptor =
    $convert.base64Decode(
        'CiRBc3NvY2lhdGlvbkV4ZWN1dGlvblRhcmdldHNGaWx0ZXJLZXkSMwovQVNTT0NJQVRJT05fRV'
        'hFQ1VUSU9OX1RBUkdFVFNfRklMVEVSX0tFWV9TVEFUVVMQABI5CjVBU1NPQ0lBVElPTl9FWEVD'
        'VVRJT05fVEFSR0VUU19GSUxURVJfS0VZX1JFU09VUkNFVFlQRRABEjcKM0FTU09DSUFUSU9OX0'
        'VYRUNVVElPTl9UQVJHRVRTX0ZJTFRFUl9LRVlfUkVTT1VSQ0VJRBAC');

@$core.Deprecated('Use associationFilterKeyDescriptor instead')
const AssociationFilterKey$json = {
  '1': 'AssociationFilterKey',
  '2': [
    {'1': 'ASSOCIATION_FILTER_KEY_STATUS', '2': 0},
    {'1': 'ASSOCIATION_FILTER_KEY_LASTEXECUTEDBEFORE', '2': 1},
    {'1': 'ASSOCIATION_FILTER_KEY_INSTANCEID', '2': 2},
    {'1': 'ASSOCIATION_FILTER_KEY_LASTEXECUTEDAFTER', '2': 3},
    {'1': 'ASSOCIATION_FILTER_KEY_RESOURCEGROUPNAME', '2': 4},
    {'1': 'ASSOCIATION_FILTER_KEY_NAME', '2': 5},
    {'1': 'ASSOCIATION_FILTER_KEY_ASSOCIATIONNAME', '2': 6},
    {'1': 'ASSOCIATION_FILTER_KEY_ASSOCIATIONID', '2': 7},
  ],
};

/// Descriptor for `AssociationFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationFilterKeyDescriptor = $convert.base64Decode(
    'ChRBc3NvY2lhdGlvbkZpbHRlcktleRIhCh1BU1NPQ0lBVElPTl9GSUxURVJfS0VZX1NUQVRVUx'
    'AAEi0KKUFTU09DSUFUSU9OX0ZJTFRFUl9LRVlfTEFTVEVYRUNVVEVEQkVGT1JFEAESJQohQVNT'
    'T0NJQVRJT05fRklMVEVSX0tFWV9JTlNUQU5DRUlEEAISLAooQVNTT0NJQVRJT05fRklMVEVSX0'
    'tFWV9MQVNURVhFQ1VURURBRlRFUhADEiwKKEFTU09DSUFUSU9OX0ZJTFRFUl9LRVlfUkVTT1VS'
    'Q0VHUk9VUE5BTUUQBBIfChtBU1NPQ0lBVElPTl9GSUxURVJfS0VZX05BTUUQBRIqCiZBU1NPQ0'
    'lBVElPTl9GSUxURVJfS0VZX0FTU09DSUFUSU9OTkFNRRAGEigKJEFTU09DSUFUSU9OX0ZJTFRF'
    'Ul9LRVlfQVNTT0NJQVRJT05JRBAH');

@$core.Deprecated('Use associationFilterOperatorTypeDescriptor instead')
const AssociationFilterOperatorType$json = {
  '1': 'AssociationFilterOperatorType',
  '2': [
    {'1': 'ASSOCIATION_FILTER_OPERATOR_TYPE_EQUAL', '2': 0},
    {'1': 'ASSOCIATION_FILTER_OPERATOR_TYPE_LESSTHAN', '2': 1},
    {'1': 'ASSOCIATION_FILTER_OPERATOR_TYPE_GREATERTHAN', '2': 2},
  ],
};

/// Descriptor for `AssociationFilterOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationFilterOperatorTypeDescriptor = $convert.base64Decode(
    'Ch1Bc3NvY2lhdGlvbkZpbHRlck9wZXJhdG9yVHlwZRIqCiZBU1NPQ0lBVElPTl9GSUxURVJfT1'
    'BFUkFUT1JfVFlQRV9FUVVBTBAAEi0KKUFTU09DSUFUSU9OX0ZJTFRFUl9PUEVSQVRPUl9UWVBF'
    'X0xFU1NUSEFOEAESMAosQVNTT0NJQVRJT05fRklMVEVSX09QRVJBVE9SX1RZUEVfR1JFQVRFUl'
    'RIQU4QAg==');

@$core.Deprecated('Use associationStatusNameDescriptor instead')
const AssociationStatusName$json = {
  '1': 'AssociationStatusName',
  '2': [
    {'1': 'ASSOCIATION_STATUS_NAME_FAILED', '2': 0},
    {'1': 'ASSOCIATION_STATUS_NAME_SUCCESS', '2': 1},
    {'1': 'ASSOCIATION_STATUS_NAME_PENDING', '2': 2},
  ],
};

/// Descriptor for `AssociationStatusName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationStatusNameDescriptor = $convert.base64Decode(
    'ChVBc3NvY2lhdGlvblN0YXR1c05hbWUSIgoeQVNTT0NJQVRJT05fU1RBVFVTX05BTUVfRkFJTE'
    'VEEAASIwofQVNTT0NJQVRJT05fU1RBVFVTX05BTUVfU1VDQ0VTUxABEiMKH0FTU09DSUFUSU9O'
    'X1NUQVRVU19OQU1FX1BFTkRJTkcQAg==');

@$core.Deprecated('Use associationSyncComplianceDescriptor instead')
const AssociationSyncCompliance$json = {
  '1': 'AssociationSyncCompliance',
  '2': [
    {'1': 'ASSOCIATION_SYNC_COMPLIANCE_MANUAL', '2': 0},
    {'1': 'ASSOCIATION_SYNC_COMPLIANCE_AUTO', '2': 1},
  ],
};

/// Descriptor for `AssociationSyncCompliance`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List associationSyncComplianceDescriptor =
    $convert.base64Decode(
        'ChlBc3NvY2lhdGlvblN5bmNDb21wbGlhbmNlEiYKIkFTU09DSUFUSU9OX1NZTkNfQ09NUExJQU'
        '5DRV9NQU5VQUwQABIkCiBBU1NPQ0lBVElPTl9TWU5DX0NPTVBMSUFOQ0VfQVVUTxAB');

@$core.Deprecated('Use attachmentHashTypeDescriptor instead')
const AttachmentHashType$json = {
  '1': 'AttachmentHashType',
  '2': [
    {'1': 'ATTACHMENT_HASH_TYPE_SHA256', '2': 0},
  ],
};

/// Descriptor for `AttachmentHashType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attachmentHashTypeDescriptor = $convert.base64Decode(
    'ChJBdHRhY2htZW50SGFzaFR5cGUSHwobQVRUQUNITUVOVF9IQVNIX1RZUEVfU0hBMjU2EAA=');

@$core.Deprecated('Use attachmentsSourceKeyDescriptor instead')
const AttachmentsSourceKey$json = {
  '1': 'AttachmentsSourceKey',
  '2': [
    {'1': 'ATTACHMENTS_SOURCE_KEY_ATTACHMENTREFERENCE', '2': 0},
    {'1': 'ATTACHMENTS_SOURCE_KEY_SOURCEURL', '2': 1},
    {'1': 'ATTACHMENTS_SOURCE_KEY_S3FILEURL', '2': 2},
  ],
};

/// Descriptor for `AttachmentsSourceKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attachmentsSourceKeyDescriptor = $convert.base64Decode(
    'ChRBdHRhY2htZW50c1NvdXJjZUtleRIuCipBVFRBQ0hNRU5UU19TT1VSQ0VfS0VZX0FUVEFDSE'
    '1FTlRSRUZFUkVOQ0UQABIkCiBBVFRBQ0hNRU5UU19TT1VSQ0VfS0VZX1NPVVJDRVVSTBABEiQK'
    'IEFUVEFDSE1FTlRTX1NPVVJDRV9LRVlfUzNGSUxFVVJMEAI=');

@$core.Deprecated('Use automationExecutionFilterKeyDescriptor instead')
const AutomationExecutionFilterKey$json = {
  '1': 'AutomationExecutionFilterKey',
  '2': [
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_AFTER', '2': 0},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_TAG_KEY', '2': 1},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_PARENT_EXECUTION_ID', '2': 2},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_TARGET_RESOURCE_GROUP', '2': 3},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_ID', '2': 4},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_OPS_ITEM_ID', '2': 5},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_CURRENT_ACTION', '2': 6},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_BEFORE', '2': 7},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_STATUS', '2': 8},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_DOCUMENT_NAME_PREFIX', '2': 9},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_TYPE', '2': 10},
    {'1': 'AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_SUBTYPE', '2': 11},
  ],
};

/// Descriptor for `AutomationExecutionFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List automationExecutionFilterKeyDescriptor = $convert.base64Decode(
    'ChxBdXRvbWF0aW9uRXhlY3V0aW9uRmlsdGVyS2V5EjQKMEFVVE9NQVRJT05fRVhFQ1VUSU9OX0'
    'ZJTFRFUl9LRVlfU1RBUlRfVElNRV9BRlRFUhAAEisKJ0FVVE9NQVRJT05fRVhFQ1VUSU9OX0ZJ'
    'TFRFUl9LRVlfVEFHX0tFWRABEjcKM0FVVE9NQVRJT05fRVhFQ1VUSU9OX0ZJTFRFUl9LRVlfUE'
    'FSRU5UX0VYRUNVVElPTl9JRBACEjkKNUFVVE9NQVRJT05fRVhFQ1VUSU9OX0ZJTFRFUl9LRVlf'
    'VEFSR0VUX1JFU09VUkNFX0dST1VQEAMSMAosQVVUT01BVElPTl9FWEVDVVRJT05fRklMVEVSX0'
    'tFWV9FWEVDVVRJT05fSUQQBBIvCitBVVRPTUFUSU9OX0VYRUNVVElPTl9GSUxURVJfS0VZX09Q'
    'U19JVEVNX0lEEAUSMgouQVVUT01BVElPTl9FWEVDVVRJT05fRklMVEVSX0tFWV9DVVJSRU5UX0'
    'FDVElPThAGEjUKMUFVVE9NQVRJT05fRVhFQ1VUSU9OX0ZJTFRFUl9LRVlfU1RBUlRfVElNRV9C'
    'RUZPUkUQBxI0CjBBVVRPTUFUSU9OX0VYRUNVVElPTl9GSUxURVJfS0VZX0VYRUNVVElPTl9TVE'
    'FUVVMQCBI4CjRBVVRPTUFUSU9OX0VYRUNVVElPTl9GSUxURVJfS0VZX0RPQ1VNRU5UX05BTUVf'
    'UFJFRklYEAkSMwovQVVUT01BVElPTl9FWEVDVVRJT05fRklMVEVSX0tFWV9BVVRPTUFUSU9OX1'
    'RZUEUQChI2CjJBVVRPTUFUSU9OX0VYRUNVVElPTl9GSUxURVJfS0VZX0FVVE9NQVRJT05fU1VC'
    'VFlQRRAL');

@$core.Deprecated('Use automationExecutionStatusDescriptor instead')
const AutomationExecutionStatus$json = {
  '1': 'AutomationExecutionStatus',
  '2': [
    {
      '1': 'AUTOMATION_EXECUTION_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE',
      '2': 0
    },
    {'1': 'AUTOMATION_EXECUTION_STATUS_PENDING', '2': 1},
    {'1': 'AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_SUCCESS', '2': 2},
    {'1': 'AUTOMATION_EXECUTION_STATUS_TIMEDOUT', '2': 3},
    {'1': 'AUTOMATION_EXECUTION_STATUS_PENDING_APPROVAL', '2': 4},
    {
      '1': 'AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED',
      '2': 5
    },
    {'1': 'AUTOMATION_EXECUTION_STATUS_SUCCESS', '2': 6},
    {'1': 'AUTOMATION_EXECUTION_STATUS_INPROGRESS', '2': 7},
    {'1': 'AUTOMATION_EXECUTION_STATUS_RUNBOOK_INPROGRESS', '2': 8},
    {'1': 'AUTOMATION_EXECUTION_STATUS_EXITED', '2': 9},
    {'1': 'AUTOMATION_EXECUTION_STATUS_CANCELLED', '2': 10},
    {'1': 'AUTOMATION_EXECUTION_STATUS_REJECTED', '2': 11},
    {'1': 'AUTOMATION_EXECUTION_STATUS_APPROVED', '2': 12},
    {'1': 'AUTOMATION_EXECUTION_STATUS_SCHEDULED', '2': 13},
    {
      '1': 'AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED',
      '2': 14
    },
    {'1': 'AUTOMATION_EXECUTION_STATUS_WAITING', '2': 15},
    {'1': 'AUTOMATION_EXECUTION_STATUS_CANCELLING', '2': 16},
    {'1': 'AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_FAILURE', '2': 17},
    {'1': 'AUTOMATION_EXECUTION_STATUS_FAILED', '2': 18},
  ],
};

/// Descriptor for `AutomationExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List automationExecutionStatusDescriptor = $convert.base64Decode(
    'ChlBdXRvbWF0aW9uRXhlY3V0aW9uU3RhdHVzEkAKPEFVVE9NQVRJT05fRVhFQ1VUSU9OX1NUQV'
    'RVU19QRU5ESU5HX0NIQU5HRV9DQUxFTkRBUl9PVkVSUklERRAAEicKI0FVVE9NQVRJT05fRVhF'
    'Q1VUSU9OX1NUQVRVU19QRU5ESU5HEAESNgoyQVVUT01BVElPTl9FWEVDVVRJT05fU1RBVFVTX0'
    'NPTVBMRVRFRF9XSVRIX1NVQ0NFU1MQAhIoCiRBVVRPTUFUSU9OX0VYRUNVVElPTl9TVEFUVVNf'
    'VElNRURPVVQQAxIwCixBVVRPTUFUSU9OX0VYRUNVVElPTl9TVEFUVVNfUEVORElOR19BUFBST1'
    'ZBTBAEEkEKPUFVVE9NQVRJT05fRVhFQ1VUSU9OX1NUQVRVU19DSEFOR0VfQ0FMRU5EQVJfT1ZF'
    'UlJJREVfQVBQUk9WRUQQBRInCiNBVVRPTUFUSU9OX0VYRUNVVElPTl9TVEFUVVNfU1VDQ0VTUx'
    'AGEioKJkFVVE9NQVRJT05fRVhFQ1VUSU9OX1NUQVRVU19JTlBST0dSRVNTEAcSMgouQVVUT01B'
    'VElPTl9FWEVDVVRJT05fU1RBVFVTX1JVTkJPT0tfSU5QUk9HUkVTUxAIEiYKIkFVVE9NQVRJT0'
    '5fRVhFQ1VUSU9OX1NUQVRVU19FWElURUQQCRIpCiVBVVRPTUFUSU9OX0VYRUNVVElPTl9TVEFU'
    'VVNfQ0FOQ0VMTEVEEAoSKAokQVVUT01BVElPTl9FWEVDVVRJT05fU1RBVFVTX1JFSkVDVEVEEA'
    'sSKAokQVVUT01BVElPTl9FWEVDVVRJT05fU1RBVFVTX0FQUFJPVkVEEAwSKQolQVVUT01BVElP'
    'Tl9FWEVDVVRJT05fU1RBVFVTX1NDSEVEVUxFRBANEkEKPUFVVE9NQVRJT05fRVhFQ1VUSU9OX1'
    'NUQVRVU19DSEFOR0VfQ0FMRU5EQVJfT1ZFUlJJREVfUkVKRUNURUQQDhInCiNBVVRPTUFUSU9O'
    'X0VYRUNVVElPTl9TVEFUVVNfV0FJVElORxAPEioKJkFVVE9NQVRJT05fRVhFQ1VUSU9OX1NUQV'
    'RVU19DQU5DRUxMSU5HEBASNgoyQVVUT01BVElPTl9FWEVDVVRJT05fU1RBVFVTX0NPTVBMRVRF'
    'RF9XSVRIX0ZBSUxVUkUQERImCiJBVVRPTUFUSU9OX0VYRUNVVElPTl9TVEFUVVNfRkFJTEVEEB'
    'I=');

@$core.Deprecated('Use automationSubtypeDescriptor instead')
const AutomationSubtype$json = {
  '1': 'AutomationSubtype',
  '2': [
    {'1': 'AUTOMATION_SUBTYPE_ACCESSREQUEST', '2': 0},
    {'1': 'AUTOMATION_SUBTYPE_CHANGEREQUEST', '2': 1},
  ],
};

/// Descriptor for `AutomationSubtype`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List automationSubtypeDescriptor = $convert.base64Decode(
    'ChFBdXRvbWF0aW9uU3VidHlwZRIkCiBBVVRPTUFUSU9OX1NVQlRZUEVfQUNDRVNTUkVRVUVTVB'
    'AAEiQKIEFVVE9NQVRJT05fU1VCVFlQRV9DSEFOR0VSRVFVRVNUEAE=');

@$core.Deprecated('Use automationTypeDescriptor instead')
const AutomationType$json = {
  '1': 'AutomationType',
  '2': [
    {'1': 'AUTOMATION_TYPE_LOCAL', '2': 0},
    {'1': 'AUTOMATION_TYPE_CROSSACCOUNT', '2': 1},
  ],
};

/// Descriptor for `AutomationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List automationTypeDescriptor = $convert.base64Decode(
    'Cg5BdXRvbWF0aW9uVHlwZRIZChVBVVRPTUFUSU9OX1RZUEVfTE9DQUwQABIgChxBVVRPTUFUSU'
    '9OX1RZUEVfQ1JPU1NBQ0NPVU5UEAE=');

@$core.Deprecated('Use calendarStateDescriptor instead')
const CalendarState$json = {
  '1': 'CalendarState',
  '2': [
    {'1': 'CALENDAR_STATE_CLOSED', '2': 0},
    {'1': 'CALENDAR_STATE_OPEN', '2': 1},
  ],
};

/// Descriptor for `CalendarState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List calendarStateDescriptor = $convert.base64Decode(
    'Cg1DYWxlbmRhclN0YXRlEhkKFUNBTEVOREFSX1NUQVRFX0NMT1NFRBAAEhcKE0NBTEVOREFSX1'
    'NUQVRFX09QRU4QAQ==');

@$core.Deprecated('Use commandFilterKeyDescriptor instead')
const CommandFilterKey$json = {
  '1': 'CommandFilterKey',
  '2': [
    {'1': 'COMMAND_FILTER_KEY_EXECUTION_STAGE', '2': 0},
    {'1': 'COMMAND_FILTER_KEY_DOCUMENT_NAME', '2': 1},
    {'1': 'COMMAND_FILTER_KEY_STATUS', '2': 2},
    {'1': 'COMMAND_FILTER_KEY_INVOKED_BEFORE', '2': 3},
    {'1': 'COMMAND_FILTER_KEY_INVOKED_AFTER', '2': 4},
  ],
};

/// Descriptor for `CommandFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List commandFilterKeyDescriptor = $convert.base64Decode(
    'ChBDb21tYW5kRmlsdGVyS2V5EiYKIkNPTU1BTkRfRklMVEVSX0tFWV9FWEVDVVRJT05fU1RBR0'
    'UQABIkCiBDT01NQU5EX0ZJTFRFUl9LRVlfRE9DVU1FTlRfTkFNRRABEh0KGUNPTU1BTkRfRklM'
    'VEVSX0tFWV9TVEFUVVMQAhIlCiFDT01NQU5EX0ZJTFRFUl9LRVlfSU5WT0tFRF9CRUZPUkUQAx'
    'IkCiBDT01NQU5EX0ZJTFRFUl9LRVlfSU5WT0tFRF9BRlRFUhAE');

@$core.Deprecated('Use commandInvocationStatusDescriptor instead')
const CommandInvocationStatus$json = {
  '1': 'CommandInvocationStatus',
  '2': [
    {'1': 'COMMAND_INVOCATION_STATUS_PENDING', '2': 0},
    {'1': 'COMMAND_INVOCATION_STATUS_TIMED_OUT', '2': 1},
    {'1': 'COMMAND_INVOCATION_STATUS_SUCCESS', '2': 2},
    {'1': 'COMMAND_INVOCATION_STATUS_IN_PROGRESS', '2': 3},
    {'1': 'COMMAND_INVOCATION_STATUS_DELAYED', '2': 4},
    {'1': 'COMMAND_INVOCATION_STATUS_CANCELLED', '2': 5},
    {'1': 'COMMAND_INVOCATION_STATUS_CANCELLING', '2': 6},
    {'1': 'COMMAND_INVOCATION_STATUS_FAILED', '2': 7},
  ],
};

/// Descriptor for `CommandInvocationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List commandInvocationStatusDescriptor = $convert.base64Decode(
    'ChdDb21tYW5kSW52b2NhdGlvblN0YXR1cxIlCiFDT01NQU5EX0lOVk9DQVRJT05fU1RBVFVTX1'
    'BFTkRJTkcQABInCiNDT01NQU5EX0lOVk9DQVRJT05fU1RBVFVTX1RJTUVEX09VVBABEiUKIUNP'
    'TU1BTkRfSU5WT0NBVElPTl9TVEFUVVNfU1VDQ0VTUxACEikKJUNPTU1BTkRfSU5WT0NBVElPTl'
    '9TVEFUVVNfSU5fUFJPR1JFU1MQAxIlCiFDT01NQU5EX0lOVk9DQVRJT05fU1RBVFVTX0RFTEFZ'
    'RUQQBBInCiNDT01NQU5EX0lOVk9DQVRJT05fU1RBVFVTX0NBTkNFTExFRBAFEigKJENPTU1BTk'
    'RfSU5WT0NBVElPTl9TVEFUVVNfQ0FOQ0VMTElORxAGEiQKIENPTU1BTkRfSU5WT0NBVElPTl9T'
    'VEFUVVNfRkFJTEVEEAc=');

@$core.Deprecated('Use commandPluginStatusDescriptor instead')
const CommandPluginStatus$json = {
  '1': 'CommandPluginStatus',
  '2': [
    {'1': 'COMMAND_PLUGIN_STATUS_PENDING', '2': 0},
    {'1': 'COMMAND_PLUGIN_STATUS_TIMED_OUT', '2': 1},
    {'1': 'COMMAND_PLUGIN_STATUS_SUCCESS', '2': 2},
    {'1': 'COMMAND_PLUGIN_STATUS_IN_PROGRESS', '2': 3},
    {'1': 'COMMAND_PLUGIN_STATUS_CANCELLED', '2': 4},
    {'1': 'COMMAND_PLUGIN_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `CommandPluginStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List commandPluginStatusDescriptor = $convert.base64Decode(
    'ChNDb21tYW5kUGx1Z2luU3RhdHVzEiEKHUNPTU1BTkRfUExVR0lOX1NUQVRVU19QRU5ESU5HEA'
    'ASIwofQ09NTUFORF9QTFVHSU5fU1RBVFVTX1RJTUVEX09VVBABEiEKHUNPTU1BTkRfUExVR0lO'
    'X1NUQVRVU19TVUNDRVNTEAISJQohQ09NTUFORF9QTFVHSU5fU1RBVFVTX0lOX1BST0dSRVNTEA'
    'MSIwofQ09NTUFORF9QTFVHSU5fU1RBVFVTX0NBTkNFTExFRBAEEiAKHENPTU1BTkRfUExVR0lO'
    'X1NUQVRVU19GQUlMRUQQBQ==');

@$core.Deprecated('Use commandStatusDescriptor instead')
const CommandStatus$json = {
  '1': 'CommandStatus',
  '2': [
    {'1': 'COMMAND_STATUS_PENDING', '2': 0},
    {'1': 'COMMAND_STATUS_TIMED_OUT', '2': 1},
    {'1': 'COMMAND_STATUS_SUCCESS', '2': 2},
    {'1': 'COMMAND_STATUS_IN_PROGRESS', '2': 3},
    {'1': 'COMMAND_STATUS_CANCELLED', '2': 4},
    {'1': 'COMMAND_STATUS_CANCELLING', '2': 5},
    {'1': 'COMMAND_STATUS_FAILED', '2': 6},
  ],
};

/// Descriptor for `CommandStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List commandStatusDescriptor = $convert.base64Decode(
    'Cg1Db21tYW5kU3RhdHVzEhoKFkNPTU1BTkRfU1RBVFVTX1BFTkRJTkcQABIcChhDT01NQU5EX1'
    'NUQVRVU19USU1FRF9PVVQQARIaChZDT01NQU5EX1NUQVRVU19TVUNDRVNTEAISHgoaQ09NTUFO'
    'RF9TVEFUVVNfSU5fUFJPR1JFU1MQAxIcChhDT01NQU5EX1NUQVRVU19DQU5DRUxMRUQQBBIdCh'
    'lDT01NQU5EX1NUQVRVU19DQU5DRUxMSU5HEAUSGQoVQ09NTUFORF9TVEFUVVNfRkFJTEVEEAY=');

@$core.Deprecated('Use complianceQueryOperatorTypeDescriptor instead')
const ComplianceQueryOperatorType$json = {
  '1': 'ComplianceQueryOperatorType',
  '2': [
    {'1': 'COMPLIANCE_QUERY_OPERATOR_TYPE_NOTEQUAL', '2': 0},
    {'1': 'COMPLIANCE_QUERY_OPERATOR_TYPE_EQUAL', '2': 1},
    {'1': 'COMPLIANCE_QUERY_OPERATOR_TYPE_LESSTHAN', '2': 2},
    {'1': 'COMPLIANCE_QUERY_OPERATOR_TYPE_GREATERTHAN', '2': 3},
    {'1': 'COMPLIANCE_QUERY_OPERATOR_TYPE_BEGINWITH', '2': 4},
  ],
};

/// Descriptor for `ComplianceQueryOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List complianceQueryOperatorTypeDescriptor = $convert.base64Decode(
    'ChtDb21wbGlhbmNlUXVlcnlPcGVyYXRvclR5cGUSKwonQ09NUExJQU5DRV9RVUVSWV9PUEVSQV'
    'RPUl9UWVBFX05PVEVRVUFMEAASKAokQ09NUExJQU5DRV9RVUVSWV9PUEVSQVRPUl9UWVBFX0VR'
    'VUFMEAESKwonQ09NUExJQU5DRV9RVUVSWV9PUEVSQVRPUl9UWVBFX0xFU1NUSEFOEAISLgoqQ0'
    '9NUExJQU5DRV9RVUVSWV9PUEVSQVRPUl9UWVBFX0dSRUFURVJUSEFOEAMSLAooQ09NUExJQU5D'
    'RV9RVUVSWV9PUEVSQVRPUl9UWVBFX0JFR0lOV0lUSBAE');

@$core.Deprecated('Use complianceSeverityDescriptor instead')
const ComplianceSeverity$json = {
  '1': 'ComplianceSeverity',
  '2': [
    {'1': 'COMPLIANCE_SEVERITY_INFORMATIONAL', '2': 0},
    {'1': 'COMPLIANCE_SEVERITY_MEDIUM', '2': 1},
    {'1': 'COMPLIANCE_SEVERITY_UNSPECIFIED', '2': 2},
    {'1': 'COMPLIANCE_SEVERITY_CRITICAL', '2': 3},
    {'1': 'COMPLIANCE_SEVERITY_LOW', '2': 4},
    {'1': 'COMPLIANCE_SEVERITY_HIGH', '2': 5},
  ],
};

/// Descriptor for `ComplianceSeverity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List complianceSeverityDescriptor = $convert.base64Decode(
    'ChJDb21wbGlhbmNlU2V2ZXJpdHkSJQohQ09NUExJQU5DRV9TRVZFUklUWV9JTkZPUk1BVElPTk'
    'FMEAASHgoaQ09NUExJQU5DRV9TRVZFUklUWV9NRURJVU0QARIjCh9DT01QTElBTkNFX1NFVkVS'
    'SVRZX1VOU1BFQ0lGSUVEEAISIAocQ09NUExJQU5DRV9TRVZFUklUWV9DUklUSUNBTBADEhsKF0'
    'NPTVBMSUFOQ0VfU0VWRVJJVFlfTE9XEAQSHAoYQ09NUExJQU5DRV9TRVZFUklUWV9ISUdIEAU=');

@$core.Deprecated('Use complianceStatusDescriptor instead')
const ComplianceStatus$json = {
  '1': 'ComplianceStatus',
  '2': [
    {'1': 'COMPLIANCE_STATUS_COMPLIANT', '2': 0},
    {'1': 'COMPLIANCE_STATUS_NONCOMPLIANT', '2': 1},
  ],
};

/// Descriptor for `ComplianceStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List complianceStatusDescriptor = $convert.base64Decode(
    'ChBDb21wbGlhbmNlU3RhdHVzEh8KG0NPTVBMSUFOQ0VfU1RBVFVTX0NPTVBMSUFOVBAAEiIKHk'
    'NPTVBMSUFOQ0VfU1RBVFVTX05PTkNPTVBMSUFOVBAB');

@$core.Deprecated('Use complianceUploadTypeDescriptor instead')
const ComplianceUploadType$json = {
  '1': 'ComplianceUploadType',
  '2': [
    {'1': 'COMPLIANCE_UPLOAD_TYPE_PARTIAL', '2': 0},
    {'1': 'COMPLIANCE_UPLOAD_TYPE_COMPLETE', '2': 1},
  ],
};

/// Descriptor for `ComplianceUploadType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List complianceUploadTypeDescriptor = $convert.base64Decode(
    'ChRDb21wbGlhbmNlVXBsb2FkVHlwZRIiCh5DT01QTElBTkNFX1VQTE9BRF9UWVBFX1BBUlRJQU'
    'wQABIjCh9DT01QTElBTkNFX1VQTE9BRF9UWVBFX0NPTVBMRVRFEAE=');

@$core.Deprecated('Use connectionStatusDescriptor instead')
const ConnectionStatus$json = {
  '1': 'ConnectionStatus',
  '2': [
    {'1': 'CONNECTION_STATUS_NOT_CONNECTED', '2': 0},
    {'1': 'CONNECTION_STATUS_CONNECTED', '2': 1},
  ],
};

/// Descriptor for `ConnectionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionStatusDescriptor = $convert.base64Decode(
    'ChBDb25uZWN0aW9uU3RhdHVzEiMKH0NPTk5FQ1RJT05fU1RBVFVTX05PVF9DT05ORUNURUQQAB'
    'IfChtDT05ORUNUSU9OX1NUQVRVU19DT05ORUNURUQQAQ==');

@$core.Deprecated('Use describeActivationsFilterKeysDescriptor instead')
const DescribeActivationsFilterKeys$json = {
  '1': 'DescribeActivationsFilterKeys',
  '2': [
    {'1': 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_ACTIVATION_IDS', '2': 0},
    {'1': 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_IAM_ROLE', '2': 1},
    {'1': 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_DEFAULT_INSTANCE_NAME', '2': 2},
  ],
};

/// Descriptor for `DescribeActivationsFilterKeys`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List describeActivationsFilterKeysDescriptor = $convert.base64Decode(
    'Ch1EZXNjcmliZUFjdGl2YXRpb25zRmlsdGVyS2V5cxIzCi9ERVNDUklCRV9BQ1RJVkFUSU9OU1'
    '9GSUxURVJfS0VZU19BQ1RJVkFUSU9OX0lEUxAAEi0KKURFU0NSSUJFX0FDVElWQVRJT05TX0ZJ'
    'TFRFUl9LRVlTX0lBTV9ST0xFEAESOgo2REVTQ1JJQkVfQUNUSVZBVElPTlNfRklMVEVSX0tFWV'
    'NfREVGQVVMVF9JTlNUQU5DRV9OQU1FEAI=');

@$core.Deprecated('Use documentFilterKeyDescriptor instead')
const DocumentFilterKey$json = {
  '1': 'DocumentFilterKey',
  '2': [
    {'1': 'DOCUMENT_FILTER_KEY_OWNER', '2': 0},
    {'1': 'DOCUMENT_FILTER_KEY_PLATFORMTYPES', '2': 1},
    {'1': 'DOCUMENT_FILTER_KEY_NAME', '2': 2},
    {'1': 'DOCUMENT_FILTER_KEY_DOCUMENTTYPE', '2': 3},
  ],
};

/// Descriptor for `DocumentFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentFilterKeyDescriptor = $convert.base64Decode(
    'ChFEb2N1bWVudEZpbHRlcktleRIdChlET0NVTUVOVF9GSUxURVJfS0VZX09XTkVSEAASJQohRE'
    '9DVU1FTlRfRklMVEVSX0tFWV9QTEFURk9STVRZUEVTEAESHAoYRE9DVU1FTlRfRklMVEVSX0tF'
    'WV9OQU1FEAISJAogRE9DVU1FTlRfRklMVEVSX0tFWV9ET0NVTUVOVFRZUEUQAw==');

@$core.Deprecated('Use documentFormatDescriptor instead')
const DocumentFormat$json = {
  '1': 'DocumentFormat',
  '2': [
    {'1': 'DOCUMENT_FORMAT_JSON', '2': 0},
    {'1': 'DOCUMENT_FORMAT_YAML', '2': 1},
    {'1': 'DOCUMENT_FORMAT_TEXT', '2': 2},
  ],
};

/// Descriptor for `DocumentFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentFormatDescriptor = $convert.base64Decode(
    'Cg5Eb2N1bWVudEZvcm1hdBIYChRET0NVTUVOVF9GT1JNQVRfSlNPThAAEhgKFERPQ1VNRU5UX0'
    'ZPUk1BVF9ZQU1MEAESGAoURE9DVU1FTlRfRk9STUFUX1RFWFQQAg==');

@$core.Deprecated('Use documentHashTypeDescriptor instead')
const DocumentHashType$json = {
  '1': 'DocumentHashType',
  '2': [
    {'1': 'DOCUMENT_HASH_TYPE_SHA256', '2': 0},
    {'1': 'DOCUMENT_HASH_TYPE_SHA1', '2': 1},
  ],
};

/// Descriptor for `DocumentHashType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentHashTypeDescriptor = $convert.base64Decode(
    'ChBEb2N1bWVudEhhc2hUeXBlEh0KGURPQ1VNRU5UX0hBU0hfVFlQRV9TSEEyNTYQABIbChdET0'
    'NVTUVOVF9IQVNIX1RZUEVfU0hBMRAB');

@$core.Deprecated('Use documentMetadataEnumDescriptor instead')
const DocumentMetadataEnum$json = {
  '1': 'DocumentMetadataEnum',
  '2': [
    {'1': 'DOCUMENT_METADATA_ENUM_DOCUMENTREVIEWS', '2': 0},
  ],
};

/// Descriptor for `DocumentMetadataEnum`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentMetadataEnumDescriptor = $convert.base64Decode(
    'ChREb2N1bWVudE1ldGFkYXRhRW51bRIqCiZET0NVTUVOVF9NRVRBREFUQV9FTlVNX0RPQ1VNRU'
    '5UUkVWSUVXUxAA');

@$core.Deprecated('Use documentParameterTypeDescriptor instead')
const DocumentParameterType$json = {
  '1': 'DocumentParameterType',
  '2': [
    {'1': 'DOCUMENT_PARAMETER_TYPE_STRINGLIST', '2': 0},
    {'1': 'DOCUMENT_PARAMETER_TYPE_STRING', '2': 1},
  ],
};

/// Descriptor for `DocumentParameterType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentParameterTypeDescriptor = $convert.base64Decode(
    'ChVEb2N1bWVudFBhcmFtZXRlclR5cGUSJgoiRE9DVU1FTlRfUEFSQU1FVEVSX1RZUEVfU1RSSU'
    '5HTElTVBAAEiIKHkRPQ1VNRU5UX1BBUkFNRVRFUl9UWVBFX1NUUklORxAB');

@$core.Deprecated('Use documentPermissionTypeDescriptor instead')
const DocumentPermissionType$json = {
  '1': 'DocumentPermissionType',
  '2': [
    {'1': 'DOCUMENT_PERMISSION_TYPE_SHARE', '2': 0},
  ],
};

/// Descriptor for `DocumentPermissionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentPermissionTypeDescriptor =
    $convert.base64Decode(
        'ChZEb2N1bWVudFBlcm1pc3Npb25UeXBlEiIKHkRPQ1VNRU5UX1BFUk1JU1NJT05fVFlQRV9TSE'
        'FSRRAA');

@$core.Deprecated('Use documentReviewActionDescriptor instead')
const DocumentReviewAction$json = {
  '1': 'DocumentReviewAction',
  '2': [
    {'1': 'DOCUMENT_REVIEW_ACTION_SENDFORREVIEW', '2': 0},
    {'1': 'DOCUMENT_REVIEW_ACTION_REJECT', '2': 1},
    {'1': 'DOCUMENT_REVIEW_ACTION_APPROVE', '2': 2},
    {'1': 'DOCUMENT_REVIEW_ACTION_UPDATEREVIEW', '2': 3},
  ],
};

/// Descriptor for `DocumentReviewAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentReviewActionDescriptor = $convert.base64Decode(
    'ChREb2N1bWVudFJldmlld0FjdGlvbhIoCiRET0NVTUVOVF9SRVZJRVdfQUNUSU9OX1NFTkRGT1'
    'JSRVZJRVcQABIhCh1ET0NVTUVOVF9SRVZJRVdfQUNUSU9OX1JFSkVDVBABEiIKHkRPQ1VNRU5U'
    'X1JFVklFV19BQ1RJT05fQVBQUk9WRRACEicKI0RPQ1VNRU5UX1JFVklFV19BQ1RJT05fVVBEQV'
    'RFUkVWSUVXEAM=');

@$core.Deprecated('Use documentReviewCommentTypeDescriptor instead')
const DocumentReviewCommentType$json = {
  '1': 'DocumentReviewCommentType',
  '2': [
    {'1': 'DOCUMENT_REVIEW_COMMENT_TYPE_COMMENT', '2': 0},
  ],
};

/// Descriptor for `DocumentReviewCommentType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentReviewCommentTypeDescriptor =
    $convert.base64Decode(
        'ChlEb2N1bWVudFJldmlld0NvbW1lbnRUeXBlEigKJERPQ1VNRU5UX1JFVklFV19DT01NRU5UX1'
        'RZUEVfQ09NTUVOVBAA');

@$core.Deprecated('Use documentStatusDescriptor instead')
const DocumentStatus$json = {
  '1': 'DocumentStatus',
  '2': [
    {'1': 'DOCUMENT_STATUS_ACTIVE', '2': 0},
    {'1': 'DOCUMENT_STATUS_UPDATING', '2': 1},
    {'1': 'DOCUMENT_STATUS_FAILED', '2': 2},
    {'1': 'DOCUMENT_STATUS_CREATING', '2': 3},
    {'1': 'DOCUMENT_STATUS_DELETING', '2': 4},
  ],
};

/// Descriptor for `DocumentStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentStatusDescriptor = $convert.base64Decode(
    'Cg5Eb2N1bWVudFN0YXR1cxIaChZET0NVTUVOVF9TVEFUVVNfQUNUSVZFEAASHAoYRE9DVU1FTl'
    'RfU1RBVFVTX1VQREFUSU5HEAESGgoWRE9DVU1FTlRfU1RBVFVTX0ZBSUxFRBACEhwKGERPQ1VN'
    'RU5UX1NUQVRVU19DUkVBVElORxADEhwKGERPQ1VNRU5UX1NUQVRVU19ERUxFVElORxAE');

@$core.Deprecated('Use documentTypeDescriptor instead')
const DocumentType$json = {
  '1': 'DocumentType',
  '2': [
    {'1': 'DOCUMENT_TYPE_PROBLEMANALYSISTEMPLATE', '2': 0},
    {'1': 'DOCUMENT_TYPE_AUTOAPPROVALPOLICY', '2': 1},
    {'1': 'DOCUMENT_TYPE_APPLICATIONCONFIGURATION', '2': 2},
    {'1': 'DOCUMENT_TYPE_MANUALAPPROVALPOLICY', '2': 3},
    {'1': 'DOCUMENT_TYPE_DEPLOYMENTSTRATEGY', '2': 4},
    {'1': 'DOCUMENT_TYPE_PROBLEMANALYSIS', '2': 5},
    {'1': 'DOCUMENT_TYPE_AUTOMATION', '2': 6},
    {'1': 'DOCUMENT_TYPE_CHANGECALENDAR', '2': 7},
    {'1': 'DOCUMENT_TYPE_APPLICATIONCONFIGURATIONSCHEMA', '2': 8},
    {'1': 'DOCUMENT_TYPE_QUICKSETUP', '2': 9},
    {'1': 'DOCUMENT_TYPE_CONFORMANCEPACKTEMPLATE', '2': 10},
    {'1': 'DOCUMENT_TYPE_POLICY', '2': 11},
    {'1': 'DOCUMENT_TYPE_PACKAGE', '2': 12},
    {'1': 'DOCUMENT_TYPE_SESSION', '2': 13},
    {'1': 'DOCUMENT_TYPE_COMMAND', '2': 14},
    {'1': 'DOCUMENT_TYPE_CLOUDFORMATION', '2': 15},
    {'1': 'DOCUMENT_TYPE_CHANGETEMPLATE', '2': 16},
  ],
};

/// Descriptor for `DocumentType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List documentTypeDescriptor = $convert.base64Decode(
    'CgxEb2N1bWVudFR5cGUSKQolRE9DVU1FTlRfVFlQRV9QUk9CTEVNQU5BTFlTSVNURU1QTEFURR'
    'AAEiQKIERPQ1VNRU5UX1RZUEVfQVVUT0FQUFJPVkFMUE9MSUNZEAESKgomRE9DVU1FTlRfVFlQ'
    'RV9BUFBMSUNBVElPTkNPTkZJR1VSQVRJT04QAhImCiJET0NVTUVOVF9UWVBFX01BTlVBTEFQUF'
    'JPVkFMUE9MSUNZEAMSJAogRE9DVU1FTlRfVFlQRV9ERVBMT1lNRU5UU1RSQVRFR1kQBBIhCh1E'
    'T0NVTUVOVF9UWVBFX1BST0JMRU1BTkFMWVNJUxAFEhwKGERPQ1VNRU5UX1RZUEVfQVVUT01BVE'
    'lPThAGEiAKHERPQ1VNRU5UX1RZUEVfQ0hBTkdFQ0FMRU5EQVIQBxIwCixET0NVTUVOVF9UWVBF'
    'X0FQUExJQ0FUSU9OQ09ORklHVVJBVElPTlNDSEVNQRAIEhwKGERPQ1VNRU5UX1RZUEVfUVVJQ0'
    'tTRVRVUBAJEikKJURPQ1VNRU5UX1RZUEVfQ09ORk9STUFOQ0VQQUNLVEVNUExBVEUQChIYChRE'
    'T0NVTUVOVF9UWVBFX1BPTElDWRALEhkKFURPQ1VNRU5UX1RZUEVfUEFDS0FHRRAMEhkKFURPQ1'
    'VNRU5UX1RZUEVfU0VTU0lPThANEhkKFURPQ1VNRU5UX1RZUEVfQ09NTUFORBAOEiAKHERPQ1VN'
    'RU5UX1RZUEVfQ0xPVURGT1JNQVRJT04QDxIgChxET0NVTUVOVF9UWVBFX0NIQU5HRVRFTVBMQV'
    'RFEBA=');

@$core.Deprecated('Use executionModeDescriptor instead')
const ExecutionMode$json = {
  '1': 'ExecutionMode',
  '2': [
    {'1': 'EXECUTION_MODE_AUTO', '2': 0},
    {'1': 'EXECUTION_MODE_INTERACTIVE', '2': 1},
  ],
};

/// Descriptor for `ExecutionMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionModeDescriptor = $convert.base64Decode(
    'Cg1FeGVjdXRpb25Nb2RlEhcKE0VYRUNVVElPTl9NT0RFX0FVVE8QABIeChpFWEVDVVRJT05fTU'
    '9ERV9JTlRFUkFDVElWRRAB');

@$core.Deprecated('Use executionPreviewStatusDescriptor instead')
const ExecutionPreviewStatus$json = {
  '1': 'ExecutionPreviewStatus',
  '2': [
    {'1': 'EXECUTION_PREVIEW_STATUS_PENDING', '2': 0},
    {'1': 'EXECUTION_PREVIEW_STATUS_SUCCESS', '2': 1},
    {'1': 'EXECUTION_PREVIEW_STATUS_IN_PROGRESS', '2': 2},
    {'1': 'EXECUTION_PREVIEW_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `ExecutionPreviewStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List executionPreviewStatusDescriptor = $convert.base64Decode(
    'ChZFeGVjdXRpb25QcmV2aWV3U3RhdHVzEiQKIEVYRUNVVElPTl9QUkVWSUVXX1NUQVRVU19QRU'
    '5ESU5HEAASJAogRVhFQ1VUSU9OX1BSRVZJRVdfU1RBVFVTX1NVQ0NFU1MQARIoCiRFWEVDVVRJ'
    'T05fUFJFVklFV19TVEFUVVNfSU5fUFJPR1JFU1MQAhIjCh9FWEVDVVRJT05fUFJFVklFV19TVE'
    'FUVVNfRkFJTEVEEAM=');

@$core.Deprecated('Use externalAlarmStateDescriptor instead')
const ExternalAlarmState$json = {
  '1': 'ExternalAlarmState',
  '2': [
    {'1': 'EXTERNAL_ALARM_STATE_ALARM', '2': 0},
    {'1': 'EXTERNAL_ALARM_STATE_UNKNOWN', '2': 1},
  ],
};

/// Descriptor for `ExternalAlarmState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List externalAlarmStateDescriptor = $convert.base64Decode(
    'ChJFeHRlcm5hbEFsYXJtU3RhdGUSHgoaRVhURVJOQUxfQUxBUk1fU1RBVEVfQUxBUk0QABIgCh'
    'xFWFRFUk5BTF9BTEFSTV9TVEFURV9VTktOT1dOEAE=');

@$core.Deprecated('Use faultDescriptor instead')
const Fault$json = {
  '1': 'Fault',
  '2': [
    {'1': 'FAULT_CLIENT', '2': 0},
    {'1': 'FAULT_UNKNOWN', '2': 1},
    {'1': 'FAULT_SERVER', '2': 2},
  ],
};

/// Descriptor for `Fault`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List faultDescriptor = $convert.base64Decode(
    'CgVGYXVsdBIQCgxGQVVMVF9DTElFTlQQABIRCg1GQVVMVF9VTktOT1dOEAESEAoMRkFVTFRfU0'
    'VSVkVSEAI=');

@$core.Deprecated('Use impactTypeDescriptor instead')
const ImpactType$json = {
  '1': 'ImpactType',
  '2': [
    {'1': 'IMPACT_TYPE_UNDETERMINED', '2': 0},
    {'1': 'IMPACT_TYPE_NON_MUTATING', '2': 1},
    {'1': 'IMPACT_TYPE_MUTATING', '2': 2},
  ],
};

/// Descriptor for `ImpactType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List impactTypeDescriptor = $convert.base64Decode(
    'CgpJbXBhY3RUeXBlEhwKGElNUEFDVF9UWVBFX1VOREVURVJNSU5FRBAAEhwKGElNUEFDVF9UWV'
    'BFX05PTl9NVVRBVElORxABEhgKFElNUEFDVF9UWVBFX01VVEFUSU5HEAI=');

@$core.Deprecated('Use instanceInformationFilterKeyDescriptor instead')
const InstanceInformationFilterKey$json = {
  '1': 'InstanceInformationFilterKey',
  '2': [
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_ASSOCIATION_STATUS', '2': 0},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_INSTANCE_IDS', '2': 1},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_ACTIVATION_IDS', '2': 2},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_IAM_ROLE', '2': 3},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_AGENT_VERSION', '2': 4},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_PING_STATUS', '2': 5},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_PLATFORM_TYPES', '2': 6},
    {'1': 'INSTANCE_INFORMATION_FILTER_KEY_RESOURCE_TYPE', '2': 7},
  ],
};

/// Descriptor for `InstanceInformationFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List instanceInformationFilterKeyDescriptor = $convert.base64Decode(
    'ChxJbnN0YW5jZUluZm9ybWF0aW9uRmlsdGVyS2V5EjYKMklOU1RBTkNFX0lORk9STUFUSU9OX0'
    'ZJTFRFUl9LRVlfQVNTT0NJQVRJT05fU1RBVFVTEAASMAosSU5TVEFOQ0VfSU5GT1JNQVRJT05f'
    'RklMVEVSX0tFWV9JTlNUQU5DRV9JRFMQARIyCi5JTlNUQU5DRV9JTkZPUk1BVElPTl9GSUxURV'
    'JfS0VZX0FDVElWQVRJT05fSURTEAISLAooSU5TVEFOQ0VfSU5GT1JNQVRJT05fRklMVEVSX0tF'
    'WV9JQU1fUk9MRRADEjEKLUlOU1RBTkNFX0lORk9STUFUSU9OX0ZJTFRFUl9LRVlfQUdFTlRfVk'
    'VSU0lPThAEEi8KK0lOU1RBTkNFX0lORk9STUFUSU9OX0ZJTFRFUl9LRVlfUElOR19TVEFUVVMQ'
    'BRIyCi5JTlNUQU5DRV9JTkZPUk1BVElPTl9GSUxURVJfS0VZX1BMQVRGT1JNX1RZUEVTEAYSMQ'
    'otSU5TVEFOQ0VfSU5GT1JNQVRJT05fRklMVEVSX0tFWV9SRVNPVVJDRV9UWVBFEAc=');

@$core.Deprecated('Use instancePatchStateOperatorTypeDescriptor instead')
const InstancePatchStateOperatorType$json = {
  '1': 'InstancePatchStateOperatorType',
  '2': [
    {'1': 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_LESS_THAN', '2': 0},
    {'1': 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_GREATER_THAN', '2': 1},
    {'1': 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_NOT_EQUAL', '2': 2},
    {'1': 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_EQUAL', '2': 3},
  ],
};

/// Descriptor for `InstancePatchStateOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List instancePatchStateOperatorTypeDescriptor = $convert.base64Decode(
    'Ch5JbnN0YW5jZVBhdGNoU3RhdGVPcGVyYXRvclR5cGUSMAosSU5TVEFOQ0VfUEFUQ0hfU1RBVE'
    'VfT1BFUkFUT1JfVFlQRV9MRVNTX1RIQU4QABIzCi9JTlNUQU5DRV9QQVRDSF9TVEFURV9PUEVS'
    'QVRPUl9UWVBFX0dSRUFURVJfVEhBThABEjAKLElOU1RBTkNFX1BBVENIX1NUQVRFX09QRVJBVE'
    '9SX1RZUEVfTk9UX0VRVUFMEAISLAooSU5TVEFOQ0VfUEFUQ0hfU1RBVEVfT1BFUkFUT1JfVFlQ'
    'RV9FUVVBTBAD');

@$core.Deprecated('Use instancePropertyFilterKeyDescriptor instead')
const InstancePropertyFilterKey$json = {
  '1': 'InstancePropertyFilterKey',
  '2': [
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_ASSOCIATION_STATUS', '2': 0},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_INSTANCE_IDS', '2': 1},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_DOCUMENT_NAME', '2': 2},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_ACTIVATION_IDS', '2': 3},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_IAM_ROLE', '2': 4},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_AGENT_VERSION', '2': 5},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_PING_STATUS', '2': 6},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_PLATFORM_TYPES', '2': 7},
    {'1': 'INSTANCE_PROPERTY_FILTER_KEY_RESOURCE_TYPE', '2': 8},
  ],
};

/// Descriptor for `InstancePropertyFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List instancePropertyFilterKeyDescriptor = $convert.base64Decode(
    'ChlJbnN0YW5jZVByb3BlcnR5RmlsdGVyS2V5EjMKL0lOU1RBTkNFX1BST1BFUlRZX0ZJTFRFUl'
    '9LRVlfQVNTT0NJQVRJT05fU1RBVFVTEAASLQopSU5TVEFOQ0VfUFJPUEVSVFlfRklMVEVSX0tF'
    'WV9JTlNUQU5DRV9JRFMQARIuCipJTlNUQU5DRV9QUk9QRVJUWV9GSUxURVJfS0VZX0RPQ1VNRU'
    '5UX05BTUUQAhIvCitJTlNUQU5DRV9QUk9QRVJUWV9GSUxURVJfS0VZX0FDVElWQVRJT05fSURT'
    'EAMSKQolSU5TVEFOQ0VfUFJPUEVSVFlfRklMVEVSX0tFWV9JQU1fUk9MRRAEEi4KKklOU1RBTk'
    'NFX1BST1BFUlRZX0ZJTFRFUl9LRVlfQUdFTlRfVkVSU0lPThAFEiwKKElOU1RBTkNFX1BST1BF'
    'UlRZX0ZJTFRFUl9LRVlfUElOR19TVEFUVVMQBhIvCitJTlNUQU5DRV9QUk9QRVJUWV9GSUxURV'
    'JfS0VZX1BMQVRGT1JNX1RZUEVTEAcSLgoqSU5TVEFOQ0VfUFJPUEVSVFlfRklMVEVSX0tFWV9S'
    'RVNPVVJDRV9UWVBFEAg=');

@$core.Deprecated('Use instancePropertyFilterOperatorDescriptor instead')
const InstancePropertyFilterOperator$json = {
  '1': 'InstancePropertyFilterOperator',
  '2': [
    {'1': 'INSTANCE_PROPERTY_FILTER_OPERATOR_LESS_THAN', '2': 0},
    {'1': 'INSTANCE_PROPERTY_FILTER_OPERATOR_GREATER_THAN', '2': 1},
    {'1': 'INSTANCE_PROPERTY_FILTER_OPERATOR_BEGIN_WITH', '2': 2},
    {'1': 'INSTANCE_PROPERTY_FILTER_OPERATOR_NOT_EQUAL', '2': 3},
    {'1': 'INSTANCE_PROPERTY_FILTER_OPERATOR_EQUAL', '2': 4},
  ],
};

/// Descriptor for `InstancePropertyFilterOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List instancePropertyFilterOperatorDescriptor = $convert.base64Decode(
    'Ch5JbnN0YW5jZVByb3BlcnR5RmlsdGVyT3BlcmF0b3ISLworSU5TVEFOQ0VfUFJPUEVSVFlfRk'
    'lMVEVSX09QRVJBVE9SX0xFU1NfVEhBThAAEjIKLklOU1RBTkNFX1BST1BFUlRZX0ZJTFRFUl9P'
    'UEVSQVRPUl9HUkVBVEVSX1RIQU4QARIwCixJTlNUQU5DRV9QUk9QRVJUWV9GSUxURVJfT1BFUk'
    'FUT1JfQkVHSU5fV0lUSBACEi8KK0lOU1RBTkNFX1BST1BFUlRZX0ZJTFRFUl9PUEVSQVRPUl9O'
    'T1RfRVFVQUwQAxIrCidJTlNUQU5DRV9QUk9QRVJUWV9GSUxURVJfT1BFUkFUT1JfRVFVQUwQBA'
    '==');

@$core.Deprecated('Use inventoryAttributeDataTypeDescriptor instead')
const InventoryAttributeDataType$json = {
  '1': 'InventoryAttributeDataType',
  '2': [
    {'1': 'INVENTORY_ATTRIBUTE_DATA_TYPE_STRING', '2': 0},
    {'1': 'INVENTORY_ATTRIBUTE_DATA_TYPE_NUMBER', '2': 1},
  ],
};

/// Descriptor for `InventoryAttributeDataType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryAttributeDataTypeDescriptor =
    $convert.base64Decode(
        'ChpJbnZlbnRvcnlBdHRyaWJ1dGVEYXRhVHlwZRIoCiRJTlZFTlRPUllfQVRUUklCVVRFX0RBVE'
        'FfVFlQRV9TVFJJTkcQABIoCiRJTlZFTlRPUllfQVRUUklCVVRFX0RBVEFfVFlQRV9OVU1CRVIQ'
        'AQ==');

@$core.Deprecated('Use inventoryDeletionStatusDescriptor instead')
const InventoryDeletionStatus$json = {
  '1': 'InventoryDeletionStatus',
  '2': [
    {'1': 'INVENTORY_DELETION_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'INVENTORY_DELETION_STATUS_COMPLETE', '2': 1},
  ],
};

/// Descriptor for `InventoryDeletionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryDeletionStatusDescriptor = $convert.base64Decode(
    'ChdJbnZlbnRvcnlEZWxldGlvblN0YXR1cxIpCiVJTlZFTlRPUllfREVMRVRJT05fU1RBVFVTX0'
    'lOX1BST0dSRVNTEAASJgoiSU5WRU5UT1JZX0RFTEVUSU9OX1NUQVRVU19DT01QTEVURRAB');

@$core.Deprecated('Use inventoryQueryOperatorTypeDescriptor instead')
const InventoryQueryOperatorType$json = {
  '1': 'InventoryQueryOperatorType',
  '2': [
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_LESS_THAN', '2': 0},
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_EXISTS', '2': 1},
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_GREATER_THAN', '2': 2},
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_BEGIN_WITH', '2': 3},
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_NOT_EQUAL', '2': 4},
    {'1': 'INVENTORY_QUERY_OPERATOR_TYPE_EQUAL', '2': 5},
  ],
};

/// Descriptor for `InventoryQueryOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventoryQueryOperatorTypeDescriptor = $convert.base64Decode(
    'ChpJbnZlbnRvcnlRdWVyeU9wZXJhdG9yVHlwZRIrCidJTlZFTlRPUllfUVVFUllfT1BFUkFUT1'
    'JfVFlQRV9MRVNTX1RIQU4QABIoCiRJTlZFTlRPUllfUVVFUllfT1BFUkFUT1JfVFlQRV9FWElT'
    'VFMQARIuCipJTlZFTlRPUllfUVVFUllfT1BFUkFUT1JfVFlQRV9HUkVBVEVSX1RIQU4QAhIsCi'
    'hJTlZFTlRPUllfUVVFUllfT1BFUkFUT1JfVFlQRV9CRUdJTl9XSVRIEAMSKwonSU5WRU5UT1JZ'
    'X1FVRVJZX09QRVJBVE9SX1RZUEVfTk9UX0VRVUFMEAQSJwojSU5WRU5UT1JZX1FVRVJZX09QRV'
    'JBVE9SX1RZUEVfRVFVQUwQBQ==');

@$core.Deprecated('Use inventorySchemaDeleteOptionDescriptor instead')
const InventorySchemaDeleteOption$json = {
  '1': 'InventorySchemaDeleteOption',
  '2': [
    {'1': 'INVENTORY_SCHEMA_DELETE_OPTION_DISABLE_SCHEMA', '2': 0},
    {'1': 'INVENTORY_SCHEMA_DELETE_OPTION_DELETE_SCHEMA', '2': 1},
  ],
};

/// Descriptor for `InventorySchemaDeleteOption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List inventorySchemaDeleteOptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZlbnRvcnlTY2hlbWFEZWxldGVPcHRpb24SMQotSU5WRU5UT1JZX1NDSEVNQV9ERUxFVE'
        'VfT1BUSU9OX0RJU0FCTEVfU0NIRU1BEAASMAosSU5WRU5UT1JZX1NDSEVNQV9ERUxFVEVfT1BU'
        'SU9OX0RFTEVURV9TQ0hFTUEQAQ==');

@$core.Deprecated('Use lastResourceDataSyncStatusDescriptor instead')
const LastResourceDataSyncStatus$json = {
  '1': 'LastResourceDataSyncStatus',
  '2': [
    {'1': 'LAST_RESOURCE_DATA_SYNC_STATUS_SUCCESSFUL', '2': 0},
    {'1': 'LAST_RESOURCE_DATA_SYNC_STATUS_INPROGRESS', '2': 1},
    {'1': 'LAST_RESOURCE_DATA_SYNC_STATUS_FAILED', '2': 2},
  ],
};

/// Descriptor for `LastResourceDataSyncStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List lastResourceDataSyncStatusDescriptor = $convert.base64Decode(
    'ChpMYXN0UmVzb3VyY2VEYXRhU3luY1N0YXR1cxItCilMQVNUX1JFU09VUkNFX0RBVEFfU1lOQ1'
    '9TVEFUVVNfU1VDQ0VTU0ZVTBAAEi0KKUxBU1RfUkVTT1VSQ0VfREFUQV9TWU5DX1NUQVRVU19J'
    'TlBST0dSRVNTEAESKQolTEFTVF9SRVNPVVJDRV9EQVRBX1NZTkNfU1RBVFVTX0ZBSUxFRBAC');

@$core.Deprecated('Use maintenanceWindowExecutionStatusDescriptor instead')
const MaintenanceWindowExecutionStatus$json = {
  '1': 'MaintenanceWindowExecutionStatus',
  '2': [
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_TIMEDOUT', '2': 0},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_SKIPPEDOVERLAPPING', '2': 1},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_FAILED', '2': 2},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLED', '2': 3},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLING', '2': 4},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_SUCCESS', '2': 5},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_INPROGRESS', '2': 6},
    {'1': 'MAINTENANCE_WINDOW_EXECUTION_STATUS_PENDING', '2': 7},
  ],
};

/// Descriptor for `MaintenanceWindowExecutionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List maintenanceWindowExecutionStatusDescriptor = $convert.base64Decode(
    'CiBNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblN0YXR1cxIwCixNQUlOVEVOQU5DRV9XSU5ET1'
    'dfRVhFQ1VUSU9OX1NUQVRVU19USU1FRE9VVBAAEjoKNk1BSU5URU5BTkNFX1dJTkRPV19FWEVD'
    'VVRJT05fU1RBVFVTX1NLSVBQRURPVkVSTEFQUElORxABEi4KKk1BSU5URU5BTkNFX1dJTkRPV1'
    '9FWEVDVVRJT05fU1RBVFVTX0ZBSUxFRBACEjEKLU1BSU5URU5BTkNFX1dJTkRPV19FWEVDVVRJ'
    'T05fU1RBVFVTX0NBTkNFTExFRBADEjIKLk1BSU5URU5BTkNFX1dJTkRPV19FWEVDVVRJT05fU1'
    'RBVFVTX0NBTkNFTExJTkcQBBIvCitNQUlOVEVOQU5DRV9XSU5ET1dfRVhFQ1VUSU9OX1NUQVRV'
    'U19TVUNDRVNTEAUSMgouTUFJTlRFTkFOQ0VfV0lORE9XX0VYRUNVVElPTl9TVEFUVVNfSU5QUk'
    '9HUkVTUxAGEi8KK01BSU5URU5BTkNFX1dJTkRPV19FWEVDVVRJT05fU1RBVFVTX1BFTkRJTkcQ'
    'Bw==');

@$core.Deprecated('Use maintenanceWindowResourceTypeDescriptor instead')
const MaintenanceWindowResourceType$json = {
  '1': 'MaintenanceWindowResourceType',
  '2': [
    {'1': 'MAINTENANCE_WINDOW_RESOURCE_TYPE_RESOURCEGROUP', '2': 0},
    {'1': 'MAINTENANCE_WINDOW_RESOURCE_TYPE_INSTANCE', '2': 1},
  ],
};

/// Descriptor for `MaintenanceWindowResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List maintenanceWindowResourceTypeDescriptor =
    $convert.base64Decode(
        'Ch1NYWludGVuYW5jZVdpbmRvd1Jlc291cmNlVHlwZRIyCi5NQUlOVEVOQU5DRV9XSU5ET1dfUk'
        'VTT1VSQ0VfVFlQRV9SRVNPVVJDRUdST1VQEAASLQopTUFJTlRFTkFOQ0VfV0lORE9XX1JFU09V'
        'UkNFX1RZUEVfSU5TVEFOQ0UQAQ==');

@$core.Deprecated('Use maintenanceWindowTaskCutoffBehaviorDescriptor instead')
const MaintenanceWindowTaskCutoffBehavior$json = {
  '1': 'MaintenanceWindowTaskCutoffBehavior',
  '2': [
    {'1': 'MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CANCELTASK', '2': 0},
    {'1': 'MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CONTINUETASK', '2': 1},
  ],
};

/// Descriptor for `MaintenanceWindowTaskCutoffBehavior`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List maintenanceWindowTaskCutoffBehaviorDescriptor =
    $convert.base64Decode(
        'CiNNYWludGVuYW5jZVdpbmRvd1Rhc2tDdXRvZmZCZWhhdmlvchI2CjJNQUlOVEVOQU5DRV9XSU'
        '5ET1dfVEFTS19DVVRPRkZfQkVIQVZJT1JfQ0FOQ0VMVEFTSxAAEjgKNE1BSU5URU5BTkNFX1dJ'
        'TkRPV19UQVNLX0NVVE9GRl9CRUhBVklPUl9DT05USU5VRVRBU0sQAQ==');

@$core.Deprecated('Use maintenanceWindowTaskTypeDescriptor instead')
const MaintenanceWindowTaskType$json = {
  '1': 'MaintenanceWindowTaskType',
  '2': [
    {'1': 'MAINTENANCE_WINDOW_TASK_TYPE_STEPFUNCTIONS', '2': 0},
    {'1': 'MAINTENANCE_WINDOW_TASK_TYPE_RUNCOMMAND', '2': 1},
    {'1': 'MAINTENANCE_WINDOW_TASK_TYPE_AUTOMATION', '2': 2},
    {'1': 'MAINTENANCE_WINDOW_TASK_TYPE_LAMBDA', '2': 3},
  ],
};

/// Descriptor for `MaintenanceWindowTaskType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List maintenanceWindowTaskTypeDescriptor = $convert.base64Decode(
    'ChlNYWludGVuYW5jZVdpbmRvd1Rhc2tUeXBlEi4KKk1BSU5URU5BTkNFX1dJTkRPV19UQVNLX1'
    'RZUEVfU1RFUEZVTkNUSU9OUxAAEisKJ01BSU5URU5BTkNFX1dJTkRPV19UQVNLX1RZUEVfUlVO'
    'Q09NTUFORBABEisKJ01BSU5URU5BTkNFX1dJTkRPV19UQVNLX1RZUEVfQVVUT01BVElPThACEi'
    'cKI01BSU5URU5BTkNFX1dJTkRPV19UQVNLX1RZUEVfTEFNQkRBEAM=');

@$core.Deprecated('Use managedStatusDescriptor instead')
const ManagedStatus$json = {
  '1': 'ManagedStatus',
  '2': [
    {'1': 'MANAGED_STATUS_UNMANAGED', '2': 0},
    {'1': 'MANAGED_STATUS_MANAGED', '2': 1},
    {'1': 'MANAGED_STATUS_ALL', '2': 2},
  ],
};

/// Descriptor for `ManagedStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List managedStatusDescriptor = $convert.base64Decode(
    'Cg1NYW5hZ2VkU3RhdHVzEhwKGE1BTkFHRURfU1RBVFVTX1VOTUFOQUdFRBAAEhoKFk1BTkFHRU'
    'RfU1RBVFVTX01BTkFHRUQQARIWChJNQU5BR0VEX1NUQVRVU19BTEwQAg==');

@$core.Deprecated('Use nodeAggregatorTypeDescriptor instead')
const NodeAggregatorType$json = {
  '1': 'NodeAggregatorType',
  '2': [
    {'1': 'NODE_AGGREGATOR_TYPE_COUNT', '2': 0},
  ],
};

/// Descriptor for `NodeAggregatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List nodeAggregatorTypeDescriptor = $convert.base64Decode(
    'ChJOb2RlQWdncmVnYXRvclR5cGUSHgoaTk9ERV9BR0dSRUdBVE9SX1RZUEVfQ09VTlQQAA==');

@$core.Deprecated('Use nodeAttributeNameDescriptor instead')
const NodeAttributeName$json = {
  '1': 'NodeAttributeName',
  '2': [
    {'1': 'NODE_ATTRIBUTE_NAME_PLATFORM_VERSION', '2': 0},
    {'1': 'NODE_ATTRIBUTE_NAME_REGION', '2': 1},
    {'1': 'NODE_ATTRIBUTE_NAME_PLATFORM_TYPE', '2': 2},
    {'1': 'NODE_ATTRIBUTE_NAME_AGENT_VERSION', '2': 3},
    {'1': 'NODE_ATTRIBUTE_NAME_PLATFORM_NAME', '2': 4},
    {'1': 'NODE_ATTRIBUTE_NAME_RESOURCE_TYPE', '2': 5},
  ],
};

/// Descriptor for `NodeAttributeName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List nodeAttributeNameDescriptor = $convert.base64Decode(
    'ChFOb2RlQXR0cmlidXRlTmFtZRIoCiROT0RFX0FUVFJJQlVURV9OQU1FX1BMQVRGT1JNX1ZFUl'
    'NJT04QABIeChpOT0RFX0FUVFJJQlVURV9OQU1FX1JFR0lPThABEiUKIU5PREVfQVRUUklCVVRF'
    'X05BTUVfUExBVEZPUk1fVFlQRRACEiUKIU5PREVfQVRUUklCVVRFX05BTUVfQUdFTlRfVkVSU0'
    'lPThADEiUKIU5PREVfQVRUUklCVVRFX05BTUVfUExBVEZPUk1fTkFNRRAEEiUKIU5PREVfQVRU'
    'UklCVVRFX05BTUVfUkVTT1VSQ0VfVFlQRRAF');

@$core.Deprecated('Use nodeFilterKeyDescriptor instead')
const NodeFilterKey$json = {
  '1': 'NodeFilterKey',
  '2': [
    {'1': 'NODE_FILTER_KEY_PLATFORM_VERSION', '2': 0},
    {'1': 'NODE_FILTER_KEY_REGION', '2': 1},
    {'1': 'NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_ID', '2': 2},
    {'1': 'NODE_FILTER_KEY_IP_ADDRESS', '2': 3},
    {'1': 'NODE_FILTER_KEY_PLATFORM_TYPE', '2': 4},
    {'1': 'NODE_FILTER_KEY_ACCOUNT_ID', '2': 5},
    {'1': 'NODE_FILTER_KEY_MANAGED_STATUS', '2': 6},
    {'1': 'NODE_FILTER_KEY_INSTANCE_STATUS', '2': 7},
    {'1': 'NODE_FILTER_KEY_AGENT_VERSION', '2': 8},
    {'1': 'NODE_FILTER_KEY_COMPUTER_NAME', '2': 9},
    {'1': 'NODE_FILTER_KEY_AGENT_TYPE', '2': 10},
    {'1': 'NODE_FILTER_KEY_PLATFORM_NAME', '2': 11},
    {'1': 'NODE_FILTER_KEY_RESOURCE_TYPE', '2': 12},
    {'1': 'NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_PATH', '2': 13},
    {'1': 'NODE_FILTER_KEY_INSTANCE_ID', '2': 14},
  ],
};

/// Descriptor for `NodeFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List nodeFilterKeyDescriptor = $convert.base64Decode(
    'Cg1Ob2RlRmlsdGVyS2V5EiQKIE5PREVfRklMVEVSX0tFWV9QTEFURk9STV9WRVJTSU9OEAASGg'
    'oWTk9ERV9GSUxURVJfS0VZX1JFR0lPThABEioKJk5PREVfRklMVEVSX0tFWV9PUkdBTklaQVRJ'
    'T05BTF9VTklUX0lEEAISHgoaTk9ERV9GSUxURVJfS0VZX0lQX0FERFJFU1MQAxIhCh1OT0RFX0'
    'ZJTFRFUl9LRVlfUExBVEZPUk1fVFlQRRAEEh4KGk5PREVfRklMVEVSX0tFWV9BQ0NPVU5UX0lE'
    'EAUSIgoeTk9ERV9GSUxURVJfS0VZX01BTkFHRURfU1RBVFVTEAYSIwofTk9ERV9GSUxURVJfS0'
    'VZX0lOU1RBTkNFX1NUQVRVUxAHEiEKHU5PREVfRklMVEVSX0tFWV9BR0VOVF9WRVJTSU9OEAgS'
    'IQodTk9ERV9GSUxURVJfS0VZX0NPTVBVVEVSX05BTUUQCRIeChpOT0RFX0ZJTFRFUl9LRVlfQU'
    'dFTlRfVFlQRRAKEiEKHU5PREVfRklMVEVSX0tFWV9QTEFURk9STV9OQU1FEAsSIQodTk9ERV9G'
    'SUxURVJfS0VZX1JFU09VUkNFX1RZUEUQDBIsCihOT0RFX0ZJTFRFUl9LRVlfT1JHQU5JWkFUSU'
    '9OQUxfVU5JVF9QQVRIEA0SHwobTk9ERV9GSUxURVJfS0VZX0lOU1RBTkNFX0lEEA4=');

@$core.Deprecated('Use nodeFilterOperatorTypeDescriptor instead')
const NodeFilterOperatorType$json = {
  '1': 'NodeFilterOperatorType',
  '2': [
    {'1': 'NODE_FILTER_OPERATOR_TYPE_BEGIN_WITH', '2': 0},
    {'1': 'NODE_FILTER_OPERATOR_TYPE_NOT_EQUAL', '2': 1},
    {'1': 'NODE_FILTER_OPERATOR_TYPE_EQUAL', '2': 2},
  ],
};

/// Descriptor for `NodeFilterOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List nodeFilterOperatorTypeDescriptor = $convert.base64Decode(
    'ChZOb2RlRmlsdGVyT3BlcmF0b3JUeXBlEigKJE5PREVfRklMVEVSX09QRVJBVE9SX1RZUEVfQk'
    'VHSU5fV0lUSBAAEicKI05PREVfRklMVEVSX09QRVJBVE9SX1RZUEVfTk9UX0VRVUFMEAESIwof'
    'Tk9ERV9GSUxURVJfT1BFUkFUT1JfVFlQRV9FUVVBTBAC');

@$core.Deprecated('Use nodeTypeNameDescriptor instead')
const NodeTypeName$json = {
  '1': 'NodeTypeName',
  '2': [
    {'1': 'NODE_TYPE_NAME_INSTANCE', '2': 0},
  ],
};

/// Descriptor for `NodeTypeName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List nodeTypeNameDescriptor = $convert.base64Decode(
    'CgxOb2RlVHlwZU5hbWUSGwoXTk9ERV9UWVBFX05BTUVfSU5TVEFOQ0UQAA==');

@$core.Deprecated('Use notificationEventDescriptor instead')
const NotificationEvent$json = {
  '1': 'NotificationEvent',
  '2': [
    {'1': 'NOTIFICATION_EVENT_TIMED_OUT', '2': 0},
    {'1': 'NOTIFICATION_EVENT_SUCCESS', '2': 1},
    {'1': 'NOTIFICATION_EVENT_IN_PROGRESS', '2': 2},
    {'1': 'NOTIFICATION_EVENT_CANCELLED', '2': 3},
    {'1': 'NOTIFICATION_EVENT_ALL', '2': 4},
    {'1': 'NOTIFICATION_EVENT_FAILED', '2': 5},
  ],
};

/// Descriptor for `NotificationEvent`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notificationEventDescriptor = $convert.base64Decode(
    'ChFOb3RpZmljYXRpb25FdmVudBIgChxOT1RJRklDQVRJT05fRVZFTlRfVElNRURfT1VUEAASHg'
    'oaTk9USUZJQ0FUSU9OX0VWRU5UX1NVQ0NFU1MQARIiCh5OT1RJRklDQVRJT05fRVZFTlRfSU5f'
    'UFJPR1JFU1MQAhIgChxOT1RJRklDQVRJT05fRVZFTlRfQ0FOQ0VMTEVEEAMSGgoWTk9USUZJQ0'
    'FUSU9OX0VWRU5UX0FMTBAEEh0KGU5PVElGSUNBVElPTl9FVkVOVF9GQUlMRUQQBQ==');

@$core.Deprecated('Use notificationTypeDescriptor instead')
const NotificationType$json = {
  '1': 'NotificationType',
  '2': [
    {'1': 'NOTIFICATION_TYPE_INVOCATION', '2': 0},
    {'1': 'NOTIFICATION_TYPE_COMMAND', '2': 1},
  ],
};

/// Descriptor for `NotificationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notificationTypeDescriptor = $convert.base64Decode(
    'ChBOb3RpZmljYXRpb25UeXBlEiAKHE5PVElGSUNBVElPTl9UWVBFX0lOVk9DQVRJT04QABIdCh'
    'lOT1RJRklDQVRJT05fVFlQRV9DT01NQU5EEAE=');

@$core.Deprecated('Use operatingSystemDescriptor instead')
const OperatingSystem$json = {
  '1': 'OperatingSystem',
  '2': [
    {'1': 'OPERATING_SYSTEM_MACOS', '2': 0},
    {'1': 'OPERATING_SYSTEM_REDHATENTERPRISELINUX', '2': 1},
    {'1': 'OPERATING_SYSTEM_ALMALINUX', '2': 2},
    {'1': 'OPERATING_SYSTEM_DEBIAN', '2': 3},
    {'1': 'OPERATING_SYSTEM_WINDOWS', '2': 4},
    {'1': 'OPERATING_SYSTEM_AMAZONLINUX', '2': 5},
    {'1': 'OPERATING_SYSTEM_AMAZONLINUX2022', '2': 6},
    {'1': 'OPERATING_SYSTEM_CENTOS', '2': 7},
    {'1': 'OPERATING_SYSTEM_AMAZONLINUX2', '2': 8},
    {'1': 'OPERATING_SYSTEM_ORACLELINUX', '2': 9},
    {'1': 'OPERATING_SYSTEM_SUSE', '2': 10},
    {'1': 'OPERATING_SYSTEM_RASPBIAN', '2': 11},
    {'1': 'OPERATING_SYSTEM_UBUNTU', '2': 12},
    {'1': 'OPERATING_SYSTEM_ROCKY_LINUX', '2': 13},
    {'1': 'OPERATING_SYSTEM_AMAZONLINUX2023', '2': 14},
  ],
};

/// Descriptor for `OperatingSystem`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List operatingSystemDescriptor = $convert.base64Decode(
    'Cg9PcGVyYXRpbmdTeXN0ZW0SGgoWT1BFUkFUSU5HX1NZU1RFTV9NQUNPUxAAEioKJk9QRVJBVE'
    'lOR19TWVNURU1fUkVESEFURU5URVJQUklTRUxJTlVYEAESHgoaT1BFUkFUSU5HX1NZU1RFTV9B'
    'TE1BTElOVVgQAhIbChdPUEVSQVRJTkdfU1lTVEVNX0RFQklBThADEhwKGE9QRVJBVElOR19TWV'
    'NURU1fV0lORE9XUxAEEiAKHE9QRVJBVElOR19TWVNURU1fQU1BWk9OTElOVVgQBRIkCiBPUEVS'
    'QVRJTkdfU1lTVEVNX0FNQVpPTkxJTlVYMjAyMhAGEhsKF09QRVJBVElOR19TWVNURU1fQ0VOVE'
    '9TEAcSIQodT1BFUkFUSU5HX1NZU1RFTV9BTUFaT05MSU5VWDIQCBIgChxPUEVSQVRJTkdfU1lT'
    'VEVNX09SQUNMRUxJTlVYEAkSGQoVT1BFUkFUSU5HX1NZU1RFTV9TVVNFEAoSHQoZT1BFUkFUSU'
    '5HX1NZU1RFTV9SQVNQQklBThALEhsKF09QRVJBVElOR19TWVNURU1fVUJVTlRVEAwSIAocT1BF'
    'UkFUSU5HX1NZU1RFTV9ST0NLWV9MSU5VWBANEiQKIE9QRVJBVElOR19TWVNURU1fQU1BWk9OTE'
    'lOVVgyMDIzEA4=');

@$core.Deprecated('Use opsFilterOperatorTypeDescriptor instead')
const OpsFilterOperatorType$json = {
  '1': 'OpsFilterOperatorType',
  '2': [
    {'1': 'OPS_FILTER_OPERATOR_TYPE_LESS_THAN', '2': 0},
    {'1': 'OPS_FILTER_OPERATOR_TYPE_EXISTS', '2': 1},
    {'1': 'OPS_FILTER_OPERATOR_TYPE_GREATER_THAN', '2': 2},
    {'1': 'OPS_FILTER_OPERATOR_TYPE_BEGIN_WITH', '2': 3},
    {'1': 'OPS_FILTER_OPERATOR_TYPE_NOT_EQUAL', '2': 4},
    {'1': 'OPS_FILTER_OPERATOR_TYPE_EQUAL', '2': 5},
  ],
};

/// Descriptor for `OpsFilterOperatorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsFilterOperatorTypeDescriptor = $convert.base64Decode(
    'ChVPcHNGaWx0ZXJPcGVyYXRvclR5cGUSJgoiT1BTX0ZJTFRFUl9PUEVSQVRPUl9UWVBFX0xFU1'
    'NfVEhBThAAEiMKH09QU19GSUxURVJfT1BFUkFUT1JfVFlQRV9FWElTVFMQARIpCiVPUFNfRklM'
    'VEVSX09QRVJBVE9SX1RZUEVfR1JFQVRFUl9USEFOEAISJwojT1BTX0ZJTFRFUl9PUEVSQVRPUl'
    '9UWVBFX0JFR0lOX1dJVEgQAxImCiJPUFNfRklMVEVSX09QRVJBVE9SX1RZUEVfTk9UX0VRVUFM'
    'EAQSIgoeT1BTX0ZJTFRFUl9PUEVSQVRPUl9UWVBFX0VRVUFMEAU=');

@$core.Deprecated('Use opsItemDataTypeDescriptor instead')
const OpsItemDataType$json = {
  '1': 'OpsItemDataType',
  '2': [
    {'1': 'OPS_ITEM_DATA_TYPE_SEARCHABLE_STRING', '2': 0},
    {'1': 'OPS_ITEM_DATA_TYPE_STRING', '2': 1},
  ],
};

/// Descriptor for `OpsItemDataType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemDataTypeDescriptor = $convert.base64Decode(
    'Cg9PcHNJdGVtRGF0YVR5cGUSKAokT1BTX0lURU1fREFUQV9UWVBFX1NFQVJDSEFCTEVfU1RSSU'
    '5HEAASHQoZT1BTX0lURU1fREFUQV9UWVBFX1NUUklORxAB');

@$core.Deprecated('Use opsItemEventFilterKeyDescriptor instead')
const OpsItemEventFilterKey$json = {
  '1': 'OpsItemEventFilterKey',
  '2': [
    {'1': 'OPS_ITEM_EVENT_FILTER_KEY_OPSITEM_ID', '2': 0},
  ],
};

/// Descriptor for `OpsItemEventFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemEventFilterKeyDescriptor = $convert.base64Decode(
    'ChVPcHNJdGVtRXZlbnRGaWx0ZXJLZXkSKAokT1BTX0lURU1fRVZFTlRfRklMVEVSX0tFWV9PUF'
    'NJVEVNX0lEEAA=');

@$core.Deprecated('Use opsItemEventFilterOperatorDescriptor instead')
const OpsItemEventFilterOperator$json = {
  '1': 'OpsItemEventFilterOperator',
  '2': [
    {'1': 'OPS_ITEM_EVENT_FILTER_OPERATOR_EQUAL', '2': 0},
  ],
};

/// Descriptor for `OpsItemEventFilterOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemEventFilterOperatorDescriptor =
    $convert.base64Decode(
        'ChpPcHNJdGVtRXZlbnRGaWx0ZXJPcGVyYXRvchIoCiRPUFNfSVRFTV9FVkVOVF9GSUxURVJfT1'
        'BFUkFUT1JfRVFVQUwQAA==');

@$core.Deprecated('Use opsItemFilterKeyDescriptor instead')
const OpsItemFilterKey$json = {
  '1': 'OpsItemFilterKey',
  '2': [
    {'1': 'OPS_ITEM_FILTER_KEY_ACTUAL_START_TIME', '2': 0},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_OPS_ITEM_ID', '2': 1},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TARGETS_RESOURCE_GROUP', '2': 2},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TEMPLATE', '2': 3},
    {'1': 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_VALUE', '2': 4},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_NAME', '2': 5},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_ARN', '2': 6},
    {'1': 'OPS_ITEM_FILTER_KEY_PRIORITY', '2': 7},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_ARN', '2': 8},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_REGION', '2': 9},
    {'1': 'OPS_ITEM_FILTER_KEY_CREATED_TIME', '2': 10},
    {'1': 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA', '2': 11},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCOUNT_ID', '2': 12},
    {'1': 'OPS_ITEM_FILTER_KEY_RESOURCE_ID', '2': 13},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_ACCOUNT_ID', '2': 14},
    {'1': 'OPS_ITEM_FILTER_KEY_STATUS', '2': 15},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ID', '2': 16},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ID', '2': 17},
    {'1': 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_KEY', '2': 18},
    {'1': 'OPS_ITEM_FILTER_KEY_SOURCE', '2': 19},
    {'1': 'OPS_ITEM_FILTER_KEY_AUTOMATION_ID', '2': 20},
    {'1': 'OPS_ITEM_FILTER_KEY_CREATED_BY', '2': 21},
    {'1': 'OPS_ITEM_FILTER_KEY_PLANNED_START_TIME', '2': 22},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ARN', '2': 23},
    {'1': 'OPS_ITEM_FILTER_KEY_TITLE', '2': 24},
    {'1': 'OPS_ITEM_FILTER_KEY_PLANNED_END_TIME', '2': 25},
    {'1': 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_NAME', '2': 26},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_IS_REPLICA', '2': 27},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ARN', '2': 28},
    {'1': 'OPS_ITEM_FILTER_KEY_INSIGHT_TYPE', '2': 29},
    {'1': 'OPS_ITEM_FILTER_KEY_LAST_MODIFIED_TIME', '2': 30},
    {'1': 'OPS_ITEM_FILTER_KEY_SEVERITY', '2': 31},
    {'1': 'OPS_ITEM_FILTER_KEY_ACTUAL_END_TIME', '2': 32},
    {'1': 'OPS_ITEM_FILTER_KEY_CATEGORY', '2': 33},
    {'1': 'OPS_ITEM_FILTER_KEY_OPSITEM_TYPE', '2': 34},
    {'1': 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_TARGET_RESOURCE_ID', '2': 35},
    {'1': 'OPS_ITEM_FILTER_KEY_OPSITEM_ID', '2': 36},
  ],
};

/// Descriptor for `OpsItemFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemFilterKeyDescriptor = $convert.base64Decode(
    'ChBPcHNJdGVtRmlsdGVyS2V5EikKJU9QU19JVEVNX0ZJTFRFUl9LRVlfQUNUVUFMX1NUQVJUX1'
    'RJTUUQABI5CjVPUFNfSVRFTV9GSUxURVJfS0VZX0FDQ0VTU19SRVFVRVNUX1NPVVJDRV9PUFNf'
    'SVRFTV9JRBABEj0KOU9QU19JVEVNX0ZJTFRFUl9LRVlfQ0hBTkdFX1JFUVVFU1RfVEFSR0VUU1'
    '9SRVNPVVJDRV9HUk9VUBACEi8KK09QU19JVEVNX0ZJTFRFUl9LRVlfQ0hBTkdFX1JFUVVFU1Rf'
    'VEVNUExBVEUQAxIuCipPUFNfSVRFTV9GSUxURVJfS0VZX09QRVJBVElPTkFMX0RBVEFfVkFMVU'
    'UQBBI1CjFPUFNfSVRFTV9GSUxURVJfS0VZX0NIQU5HRV9SRVFVRVNUX1JFUVVFU1RFUl9OQU1F'
    'EAUSMwovT1BTX0lURU1fRklMVEVSX0tFWV9DSEFOR0VfUkVRVUVTVF9BUFBST1ZFUl9BUk4QBh'
    'IgChxPUFNfSVRFTV9GSUxURVJfS0VZX1BSSU9SSVRZEAcSNAowT1BTX0lURU1fRklMVEVSX0tF'
    'WV9DSEFOR0VfUkVRVUVTVF9SRVFVRVNURVJfQVJOEAgSNAowT1BTX0lURU1fRklMVEVSX0tFWV'
    '9BQ0NFU1NfUkVRVUVTVF9TT1VSQ0VfUkVHSU9OEAkSJAogT1BTX0lURU1fRklMVEVSX0tFWV9D'
    'UkVBVEVEX1RJTUUQChIoCiRPUFNfSVRFTV9GSUxURVJfS0VZX09QRVJBVElPTkFMX0RBVEEQCx'
    'IiCh5PUFNfSVRFTV9GSUxURVJfS0VZX0FDQ09VTlRfSUQQDBIjCh9PUFNfSVRFTV9GSUxURVJf'
    'S0VZX1JFU09VUkNFX0lEEA0SOAo0T1BTX0lURU1fRklMVEVSX0tFWV9BQ0NFU1NfUkVRVUVTVF'
    '9TT1VSQ0VfQUNDT1VOVF9JRBAOEh4KGk9QU19JVEVNX0ZJTFRFUl9LRVlfU1RBVFVTEA8SMwov'
    'T1BTX0lURU1fRklMVEVSX0tFWV9BQ0NFU1NfUkVRVUVTVF9SRVFVRVNURVJfSUQQEBIyCi5PUF'
    'NfSVRFTV9GSUxURVJfS0VZX0FDQ0VTU19SRVFVRVNUX0FQUFJPVkVSX0lEEBESLAooT1BTX0lU'
    'RU1fRklMVEVSX0tFWV9PUEVSQVRJT05BTF9EQVRBX0tFWRASEh4KGk9QU19JVEVNX0ZJTFRFUl'
    '9LRVlfU09VUkNFEBMSJQohT1BTX0lURU1fRklMVEVSX0tFWV9BVVRPTUFUSU9OX0lEEBQSIgoe'
    'T1BTX0lURU1fRklMVEVSX0tFWV9DUkVBVEVEX0JZEBUSKgomT1BTX0lURU1fRklMVEVSX0tFWV'
    '9QTEFOTkVEX1NUQVJUX1RJTUUQFhIzCi9PUFNfSVRFTV9GSUxURVJfS0VZX0FDQ0VTU19SRVFV'
    'RVNUX0FQUFJPVkVSX0FSThAXEh0KGU9QU19JVEVNX0ZJTFRFUl9LRVlfVElUTEUQGBIoCiRPUF'
    'NfSVRFTV9GSUxURVJfS0VZX1BMQU5ORURfRU5EX1RJTUUQGRI0CjBPUFNfSVRFTV9GSUxURVJf'
    'S0VZX0NIQU5HRV9SRVFVRVNUX0FQUFJPVkVSX05BTUUQGhIxCi1PUFNfSVRFTV9GSUxURVJfS0'
    'VZX0FDQ0VTU19SRVFVRVNUX0lTX1JFUExJQ0EQGxI0CjBPUFNfSVRFTV9GSUxURVJfS0VZX0FD'
    'Q0VTU19SRVFVRVNUX1JFUVVFU1RFUl9BUk4QHBIkCiBPUFNfSVRFTV9GSUxURVJfS0VZX0lOU0'
    'lHSFRfVFlQRRAdEioKJk9QU19JVEVNX0ZJTFRFUl9LRVlfTEFTVF9NT0RJRklFRF9USU1FEB4S'
    'IAocT1BTX0lURU1fRklMVEVSX0tFWV9TRVZFUklUWRAfEicKI09QU19JVEVNX0ZJTFRFUl9LRV'
    'lfQUNUVUFMX0VORF9USU1FECASIAocT1BTX0lURU1fRklMVEVSX0tFWV9DQVRFR09SWRAhEiQK'
    'IE9QU19JVEVNX0ZJTFRFUl9LRVlfT1BTSVRFTV9UWVBFECISOQo1T1BTX0lURU1fRklMVEVSX0'
    'tFWV9BQ0NFU1NfUkVRVUVTVF9UQVJHRVRfUkVTT1VSQ0VfSUQQIxIiCh5PUFNfSVRFTV9GSUxU'
    'RVJfS0VZX09QU0lURU1fSUQQJA==');

@$core.Deprecated('Use opsItemFilterOperatorDescriptor instead')
const OpsItemFilterOperator$json = {
  '1': 'OpsItemFilterOperator',
  '2': [
    {'1': 'OPS_ITEM_FILTER_OPERATOR_LESS_THAN', '2': 0},
    {'1': 'OPS_ITEM_FILTER_OPERATOR_GREATER_THAN', '2': 1},
    {'1': 'OPS_ITEM_FILTER_OPERATOR_CONTAINS', '2': 2},
    {'1': 'OPS_ITEM_FILTER_OPERATOR_EQUAL', '2': 3},
  ],
};

/// Descriptor for `OpsItemFilterOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemFilterOperatorDescriptor = $convert.base64Decode(
    'ChVPcHNJdGVtRmlsdGVyT3BlcmF0b3ISJgoiT1BTX0lURU1fRklMVEVSX09QRVJBVE9SX0xFU1'
    'NfVEhBThAAEikKJU9QU19JVEVNX0ZJTFRFUl9PUEVSQVRPUl9HUkVBVEVSX1RIQU4QARIlCiFP'
    'UFNfSVRFTV9GSUxURVJfT1BFUkFUT1JfQ09OVEFJTlMQAhIiCh5PUFNfSVRFTV9GSUxURVJfT1'
    'BFUkFUT1JfRVFVQUwQAw==');

@$core.Deprecated('Use opsItemRelatedItemsFilterKeyDescriptor instead')
const OpsItemRelatedItemsFilterKey$json = {
  '1': 'OpsItemRelatedItemsFilterKey',
  '2': [
    {'1': 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_URI', '2': 0},
    {'1': 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_ASSOCIATION_ID', '2': 1},
    {'1': 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_TYPE', '2': 2},
  ],
};

/// Descriptor for `OpsItemRelatedItemsFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemRelatedItemsFilterKeyDescriptor = $convert.base64Decode(
    'ChxPcHNJdGVtUmVsYXRlZEl0ZW1zRmlsdGVyS2V5EjIKLk9QU19JVEVNX1JFTEFURURfSVRFTV'
    'NfRklMVEVSX0tFWV9SRVNPVVJDRV9VUkkQABI0CjBPUFNfSVRFTV9SRUxBVEVEX0lURU1TX0ZJ'
    'TFRFUl9LRVlfQVNTT0NJQVRJT05fSUQQARIzCi9PUFNfSVRFTV9SRUxBVEVEX0lURU1TX0ZJTF'
    'RFUl9LRVlfUkVTT1VSQ0VfVFlQRRAC');

@$core.Deprecated('Use opsItemRelatedItemsFilterOperatorDescriptor instead')
const OpsItemRelatedItemsFilterOperator$json = {
  '1': 'OpsItemRelatedItemsFilterOperator',
  '2': [
    {'1': 'OPS_ITEM_RELATED_ITEMS_FILTER_OPERATOR_EQUAL', '2': 0},
  ],
};

/// Descriptor for `OpsItemRelatedItemsFilterOperator`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemRelatedItemsFilterOperatorDescriptor =
    $convert.base64Decode(
        'CiFPcHNJdGVtUmVsYXRlZEl0ZW1zRmlsdGVyT3BlcmF0b3ISMAosT1BTX0lURU1fUkVMQVRFRF'
        '9JVEVNU19GSUxURVJfT1BFUkFUT1JfRVFVQUwQAA==');

@$core.Deprecated('Use opsItemStatusDescriptor instead')
const OpsItemStatus$json = {
  '1': 'OpsItemStatus',
  '2': [
    {'1': 'OPS_ITEM_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE', '2': 0},
    {'1': 'OPS_ITEM_STATUS_PENDING', '2': 1},
    {'1': 'OPS_ITEM_STATUS_COMPLETED_WITH_SUCCESS', '2': 2},
    {'1': 'OPS_ITEM_STATUS_RUNBOOK_IN_PROGRESS', '2': 3},
    {'1': 'OPS_ITEM_STATUS_REVOKED', '2': 4},
    {'1': 'OPS_ITEM_STATUS_TIMED_OUT', '2': 5},
    {'1': 'OPS_ITEM_STATUS_CLOSED', '2': 6},
    {'1': 'OPS_ITEM_STATUS_PENDING_APPROVAL', '2': 7},
    {'1': 'OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED', '2': 8},
    {'1': 'OPS_ITEM_STATUS_OPEN', '2': 9},
    {'1': 'OPS_ITEM_STATUS_IN_PROGRESS', '2': 10},
    {'1': 'OPS_ITEM_STATUS_RESOLVED', '2': 11},
    {'1': 'OPS_ITEM_STATUS_CANCELLED', '2': 12},
    {'1': 'OPS_ITEM_STATUS_REJECTED', '2': 13},
    {'1': 'OPS_ITEM_STATUS_APPROVED', '2': 14},
    {'1': 'OPS_ITEM_STATUS_SCHEDULED', '2': 15},
    {'1': 'OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED', '2': 16},
    {'1': 'OPS_ITEM_STATUS_CANCELLING', '2': 17},
    {'1': 'OPS_ITEM_STATUS_COMPLETED_WITH_FAILURE', '2': 18},
    {'1': 'OPS_ITEM_STATUS_FAILED', '2': 19},
  ],
};

/// Descriptor for `OpsItemStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List opsItemStatusDescriptor = $convert.base64Decode(
    'Cg1PcHNJdGVtU3RhdHVzEjQKME9QU19JVEVNX1NUQVRVU19QRU5ESU5HX0NIQU5HRV9DQUxFTk'
    'RBUl9PVkVSUklERRAAEhsKF09QU19JVEVNX1NUQVRVU19QRU5ESU5HEAESKgomT1BTX0lURU1f'
    'U1RBVFVTX0NPTVBMRVRFRF9XSVRIX1NVQ0NFU1MQAhInCiNPUFNfSVRFTV9TVEFUVVNfUlVOQk'
    '9PS19JTl9QUk9HUkVTUxADEhsKF09QU19JVEVNX1NUQVRVU19SRVZPS0VEEAQSHQoZT1BTX0lU'
    'RU1fU1RBVFVTX1RJTUVEX09VVBAFEhoKFk9QU19JVEVNX1NUQVRVU19DTE9TRUQQBhIkCiBPUF'
    'NfSVRFTV9TVEFUVVNfUEVORElOR19BUFBST1ZBTBAHEjUKMU9QU19JVEVNX1NUQVRVU19DSEFO'
    'R0VfQ0FMRU5EQVJfT1ZFUlJJREVfQVBQUk9WRUQQCBIYChRPUFNfSVRFTV9TVEFUVVNfT1BFTh'
    'AJEh8KG09QU19JVEVNX1NUQVRVU19JTl9QUk9HUkVTUxAKEhwKGE9QU19JVEVNX1NUQVRVU19S'
    'RVNPTFZFRBALEh0KGU9QU19JVEVNX1NUQVRVU19DQU5DRUxMRUQQDBIcChhPUFNfSVRFTV9TVE'
    'FUVVNfUkVKRUNURUQQDRIcChhPUFNfSVRFTV9TVEFUVVNfQVBQUk9WRUQQDhIdChlPUFNfSVRF'
    'TV9TVEFUVVNfU0NIRURVTEVEEA8SNQoxT1BTX0lURU1fU1RBVFVTX0NIQU5HRV9DQUxFTkRBUl'
    '9PVkVSUklERV9SRUpFQ1RFRBAQEh4KGk9QU19JVEVNX1NUQVRVU19DQU5DRUxMSU5HEBESKgom'
    'T1BTX0lURU1fU1RBVFVTX0NPTVBMRVRFRF9XSVRIX0ZBSUxVUkUQEhIaChZPUFNfSVRFTV9TVE'
    'FUVVNfRkFJTEVEEBM=');

@$core.Deprecated('Use parameterTierDescriptor instead')
const ParameterTier$json = {
  '1': 'ParameterTier',
  '2': [
    {'1': 'PARAMETER_TIER_STANDARD', '2': 0},
    {'1': 'PARAMETER_TIER_ADVANCED', '2': 1},
    {'1': 'PARAMETER_TIER_INTELLIGENT_TIERING', '2': 2},
  ],
};

/// Descriptor for `ParameterTier`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List parameterTierDescriptor = $convert.base64Decode(
    'Cg1QYXJhbWV0ZXJUaWVyEhsKF1BBUkFNRVRFUl9USUVSX1NUQU5EQVJEEAASGwoXUEFSQU1FVE'
    'VSX1RJRVJfQURWQU5DRUQQARImCiJQQVJBTUVURVJfVElFUl9JTlRFTExJR0VOVF9USUVSSU5H'
    'EAI=');

@$core.Deprecated('Use parameterTypeDescriptor instead')
const ParameterType$json = {
  '1': 'ParameterType',
  '2': [
    {'1': 'PARAMETER_TYPE_SECURE_STRING', '2': 0},
    {'1': 'PARAMETER_TYPE_STRING_LIST', '2': 1},
    {'1': 'PARAMETER_TYPE_STRING', '2': 2},
  ],
};

/// Descriptor for `ParameterType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List parameterTypeDescriptor = $convert.base64Decode(
    'Cg1QYXJhbWV0ZXJUeXBlEiAKHFBBUkFNRVRFUl9UWVBFX1NFQ1VSRV9TVFJJTkcQABIeChpQQV'
    'JBTUVURVJfVFlQRV9TVFJJTkdfTElTVBABEhkKFVBBUkFNRVRFUl9UWVBFX1NUUklORxAC');

@$core.Deprecated('Use parametersFilterKeyDescriptor instead')
const ParametersFilterKey$json = {
  '1': 'ParametersFilterKey',
  '2': [
    {'1': 'PARAMETERS_FILTER_KEY_KEY_ID', '2': 0},
    {'1': 'PARAMETERS_FILTER_KEY_NAME', '2': 1},
    {'1': 'PARAMETERS_FILTER_KEY_TYPE', '2': 2},
  ],
};

/// Descriptor for `ParametersFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List parametersFilterKeyDescriptor = $convert.base64Decode(
    'ChNQYXJhbWV0ZXJzRmlsdGVyS2V5EiAKHFBBUkFNRVRFUlNfRklMVEVSX0tFWV9LRVlfSUQQAB'
    'IeChpQQVJBTUVURVJTX0ZJTFRFUl9LRVlfTkFNRRABEh4KGlBBUkFNRVRFUlNfRklMVEVSX0tF'
    'WV9UWVBFEAI=');

@$core.Deprecated('Use patchActionDescriptor instead')
const PatchAction$json = {
  '1': 'PatchAction',
  '2': [
    {'1': 'PATCH_ACTION_ALLOWASDEPENDENCY', '2': 0},
    {'1': 'PATCH_ACTION_BLOCK', '2': 1},
  ],
};

/// Descriptor for `PatchAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchActionDescriptor = $convert.base64Decode(
    'CgtQYXRjaEFjdGlvbhIiCh5QQVRDSF9BQ1RJT05fQUxMT1dBU0RFUEVOREVOQ1kQABIWChJQQV'
    'RDSF9BQ1RJT05fQkxPQ0sQAQ==');

@$core.Deprecated('Use patchComplianceDataStateDescriptor instead')
const PatchComplianceDataState$json = {
  '1': 'PatchComplianceDataState',
  '2': [
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_AVAILABLESECURITYUPDATE', '2': 0},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_INSTALLED', '2': 1},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_FAILED', '2': 2},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_MISSING', '2': 3},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDPENDINGREBOOT', '2': 4},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDREJECTED', '2': 5},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_NOTAPPLICABLE', '2': 6},
    {'1': 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDOTHER', '2': 7},
  ],
};

/// Descriptor for `PatchComplianceDataState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchComplianceDataStateDescriptor = $convert.base64Decode(
    'ChhQYXRjaENvbXBsaWFuY2VEYXRhU3RhdGUSNwozUEFUQ0hfQ09NUExJQU5DRV9EQVRBX1NUQV'
    'RFX0FWQUlMQUJMRVNFQ1VSSVRZVVBEQVRFEAASKQolUEFUQ0hfQ09NUExJQU5DRV9EQVRBX1NU'
    'QVRFX0lOU1RBTExFRBABEiYKIlBBVENIX0NPTVBMSUFOQ0VfREFUQV9TVEFURV9GQUlMRUQQAh'
    'InCiNQQVRDSF9DT01QTElBTkNFX0RBVEFfU1RBVEVfTUlTU0lORxADEjYKMlBBVENIX0NPTVBM'
    'SUFOQ0VfREFUQV9TVEFURV9JTlNUQUxMRURQRU5ESU5HUkVCT09UEAQSMQotUEFUQ0hfQ09NUE'
    'xJQU5DRV9EQVRBX1NUQVRFX0lOU1RBTExFRFJFSkVDVEVEEAUSLQopUEFUQ0hfQ09NUExJQU5D'
    'RV9EQVRBX1NUQVRFX05PVEFQUExJQ0FCTEUQBhIuCipQQVRDSF9DT01QTElBTkNFX0RBVEFfU1'
    'RBVEVfSU5TVEFMTEVET1RIRVIQBw==');

@$core.Deprecated('Use patchComplianceLevelDescriptor instead')
const PatchComplianceLevel$json = {
  '1': 'PatchComplianceLevel',
  '2': [
    {'1': 'PATCH_COMPLIANCE_LEVEL_INFORMATIONAL', '2': 0},
    {'1': 'PATCH_COMPLIANCE_LEVEL_MEDIUM', '2': 1},
    {'1': 'PATCH_COMPLIANCE_LEVEL_UNSPECIFIED', '2': 2},
    {'1': 'PATCH_COMPLIANCE_LEVEL_CRITICAL', '2': 3},
    {'1': 'PATCH_COMPLIANCE_LEVEL_LOW', '2': 4},
    {'1': 'PATCH_COMPLIANCE_LEVEL_HIGH', '2': 5},
  ],
};

/// Descriptor for `PatchComplianceLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchComplianceLevelDescriptor = $convert.base64Decode(
    'ChRQYXRjaENvbXBsaWFuY2VMZXZlbBIoCiRQQVRDSF9DT01QTElBTkNFX0xFVkVMX0lORk9STU'
    'FUSU9OQUwQABIhCh1QQVRDSF9DT01QTElBTkNFX0xFVkVMX01FRElVTRABEiYKIlBBVENIX0NP'
    'TVBMSUFOQ0VfTEVWRUxfVU5TUEVDSUZJRUQQAhIjCh9QQVRDSF9DT01QTElBTkNFX0xFVkVMX0'
    'NSSVRJQ0FMEAMSHgoaUEFUQ0hfQ09NUExJQU5DRV9MRVZFTF9MT1cQBBIfChtQQVRDSF9DT01Q'
    'TElBTkNFX0xFVkVMX0hJR0gQBQ==');

@$core.Deprecated('Use patchComplianceStatusDescriptor instead')
const PatchComplianceStatus$json = {
  '1': 'PatchComplianceStatus',
  '2': [
    {'1': 'PATCH_COMPLIANCE_STATUS_COMPLIANT', '2': 0},
    {'1': 'PATCH_COMPLIANCE_STATUS_NONCOMPLIANT', '2': 1},
  ],
};

/// Descriptor for `PatchComplianceStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchComplianceStatusDescriptor = $convert.base64Decode(
    'ChVQYXRjaENvbXBsaWFuY2VTdGF0dXMSJQohUEFUQ0hfQ09NUExJQU5DRV9TVEFUVVNfQ09NUE'
    'xJQU5UEAASKAokUEFUQ0hfQ09NUExJQU5DRV9TVEFUVVNfTk9OQ09NUExJQU5UEAE=');

@$core.Deprecated('Use patchDeploymentStatusDescriptor instead')
const PatchDeploymentStatus$json = {
  '1': 'PatchDeploymentStatus',
  '2': [
    {'1': 'PATCH_DEPLOYMENT_STATUS_APPROVED', '2': 0},
    {'1': 'PATCH_DEPLOYMENT_STATUS_PENDINGAPPROVAL', '2': 1},
    {'1': 'PATCH_DEPLOYMENT_STATUS_EXPLICITREJECTED', '2': 2},
    {'1': 'PATCH_DEPLOYMENT_STATUS_EXPLICITAPPROVED', '2': 3},
  ],
};

/// Descriptor for `PatchDeploymentStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchDeploymentStatusDescriptor = $convert.base64Decode(
    'ChVQYXRjaERlcGxveW1lbnRTdGF0dXMSJAogUEFUQ0hfREVQTE9ZTUVOVF9TVEFUVVNfQVBQUk'
    '9WRUQQABIrCidQQVRDSF9ERVBMT1lNRU5UX1NUQVRVU19QRU5ESU5HQVBQUk9WQUwQARIsCihQ'
    'QVRDSF9ERVBMT1lNRU5UX1NUQVRVU19FWFBMSUNJVFJFSkVDVEVEEAISLAooUEFUQ0hfREVQTE'
    '9ZTUVOVF9TVEFUVVNfRVhQTElDSVRBUFBST1ZFRBAD');

@$core.Deprecated('Use patchFilterKeyDescriptor instead')
const PatchFilterKey$json = {
  '1': 'PatchFilterKey',
  '2': [
    {'1': 'PATCH_FILTER_KEY_PRIORITY', '2': 0},
    {'1': 'PATCH_FILTER_KEY_PRODUCT', '2': 1},
    {'1': 'PATCH_FILTER_KEY_MSRCSEVERITY', '2': 2},
    {'1': 'PATCH_FILTER_KEY_ADVISORYID', '2': 3},
    {'1': 'PATCH_FILTER_KEY_PATCHID', '2': 4},
    {'1': 'PATCH_FILTER_KEY_PATCHSET', '2': 5},
    {'1': 'PATCH_FILTER_KEY_CVEID', '2': 6},
    {'1': 'PATCH_FILTER_KEY_RELEASE', '2': 7},
    {'1': 'PATCH_FILTER_KEY_REPOSITORY', '2': 8},
    {'1': 'PATCH_FILTER_KEY_ARCH', '2': 9},
    {'1': 'PATCH_FILTER_KEY_CLASSIFICATION', '2': 10},
    {'1': 'PATCH_FILTER_KEY_VERSION', '2': 11},
    {'1': 'PATCH_FILTER_KEY_SECTION', '2': 12},
    {'1': 'PATCH_FILTER_KEY_EPOCH', '2': 13},
    {'1': 'PATCH_FILTER_KEY_PRODUCTFAMILY', '2': 14},
    {'1': 'PATCH_FILTER_KEY_SECURITY', '2': 15},
    {'1': 'PATCH_FILTER_KEY_NAME', '2': 16},
    {'1': 'PATCH_FILTER_KEY_BUGZILLAID', '2': 17},
    {'1': 'PATCH_FILTER_KEY_SEVERITY', '2': 18},
  ],
};

/// Descriptor for `PatchFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchFilterKeyDescriptor = $convert.base64Decode(
    'Cg5QYXRjaEZpbHRlcktleRIdChlQQVRDSF9GSUxURVJfS0VZX1BSSU9SSVRZEAASHAoYUEFUQ0'
    'hfRklMVEVSX0tFWV9QUk9EVUNUEAESIQodUEFUQ0hfRklMVEVSX0tFWV9NU1JDU0VWRVJJVFkQ'
    'AhIfChtQQVRDSF9GSUxURVJfS0VZX0FEVklTT1JZSUQQAxIcChhQQVRDSF9GSUxURVJfS0VZX1'
    'BBVENISUQQBBIdChlQQVRDSF9GSUxURVJfS0VZX1BBVENIU0VUEAUSGgoWUEFUQ0hfRklMVEVS'
    'X0tFWV9DVkVJRBAGEhwKGFBBVENIX0ZJTFRFUl9LRVlfUkVMRUFTRRAHEh8KG1BBVENIX0ZJTF'
    'RFUl9LRVlfUkVQT1NJVE9SWRAIEhkKFVBBVENIX0ZJTFRFUl9LRVlfQVJDSBAJEiMKH1BBVENI'
    'X0ZJTFRFUl9LRVlfQ0xBU1NJRklDQVRJT04QChIcChhQQVRDSF9GSUxURVJfS0VZX1ZFUlNJT0'
    '4QCxIcChhQQVRDSF9GSUxURVJfS0VZX1NFQ1RJT04QDBIaChZQQVRDSF9GSUxURVJfS0VZX0VQ'
    'T0NIEA0SIgoeUEFUQ0hfRklMVEVSX0tFWV9QUk9EVUNURkFNSUxZEA4SHQoZUEFUQ0hfRklMVE'
    'VSX0tFWV9TRUNVUklUWRAPEhkKFVBBVENIX0ZJTFRFUl9LRVlfTkFNRRAQEh8KG1BBVENIX0ZJ'
    'TFRFUl9LRVlfQlVHWklMTEFJRBAREh0KGVBBVENIX0ZJTFRFUl9LRVlfU0VWRVJJVFkQEg==');

@$core.Deprecated('Use patchOperationTypeDescriptor instead')
const PatchOperationType$json = {
  '1': 'PatchOperationType',
  '2': [
    {'1': 'PATCH_OPERATION_TYPE_SCAN', '2': 0},
    {'1': 'PATCH_OPERATION_TYPE_INSTALL', '2': 1},
  ],
};

/// Descriptor for `PatchOperationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchOperationTypeDescriptor = $convert.base64Decode(
    'ChJQYXRjaE9wZXJhdGlvblR5cGUSHQoZUEFUQ0hfT1BFUkFUSU9OX1RZUEVfU0NBThAAEiAKHF'
    'BBVENIX09QRVJBVElPTl9UWVBFX0lOU1RBTEwQAQ==');

@$core.Deprecated('Use patchPropertyDescriptor instead')
const PatchProperty$json = {
  '1': 'PatchProperty',
  '2': [
    {'1': 'PATCH_PROPERTY_PRODUCT', '2': 0},
    {'1': 'PATCH_PROPERTY_PATCHPRODUCTFAMILY', '2': 1},
    {'1': 'PATCH_PROPERTY_PATCHCLASSIFICATION', '2': 2},
    {'1': 'PATCH_PROPERTY_PATCHMSRCSEVERITY', '2': 3},
    {'1': 'PATCH_PROPERTY_PATCHPRIORITY', '2': 4},
    {'1': 'PATCH_PROPERTY_PATCHSEVERITY', '2': 5},
  ],
};

/// Descriptor for `PatchProperty`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchPropertyDescriptor = $convert.base64Decode(
    'Cg1QYXRjaFByb3BlcnR5EhoKFlBBVENIX1BST1BFUlRZX1BST0RVQ1QQABIlCiFQQVRDSF9QUk'
    '9QRVJUWV9QQVRDSFBST0RVQ1RGQU1JTFkQARImCiJQQVRDSF9QUk9QRVJUWV9QQVRDSENMQVNT'
    'SUZJQ0FUSU9OEAISJAogUEFUQ0hfUFJPUEVSVFlfUEFUQ0hNU1JDU0VWRVJJVFkQAxIgChxQQV'
    'RDSF9QUk9QRVJUWV9QQVRDSFBSSU9SSVRZEAQSIAocUEFUQ0hfUFJPUEVSVFlfUEFUQ0hTRVZF'
    'UklUWRAF');

@$core.Deprecated('Use patchSetDescriptor instead')
const PatchSet$json = {
  '1': 'PatchSet',
  '2': [
    {'1': 'PATCH_SET_OS', '2': 0},
    {'1': 'PATCH_SET_APPLICATION', '2': 1},
  ],
};

/// Descriptor for `PatchSet`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List patchSetDescriptor = $convert.base64Decode(
    'CghQYXRjaFNldBIQCgxQQVRDSF9TRVRfT1MQABIZChVQQVRDSF9TRVRfQVBQTElDQVRJT04QAQ'
    '==');

@$core.Deprecated('Use pingStatusDescriptor instead')
const PingStatus$json = {
  '1': 'PingStatus',
  '2': [
    {'1': 'PING_STATUS_ONLINE', '2': 0},
    {'1': 'PING_STATUS_CONNECTION_LOST', '2': 1},
    {'1': 'PING_STATUS_INACTIVE', '2': 2},
  ],
};

/// Descriptor for `PingStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List pingStatusDescriptor = $convert.base64Decode(
    'CgpQaW5nU3RhdHVzEhYKElBJTkdfU1RBVFVTX09OTElORRAAEh8KG1BJTkdfU1RBVFVTX0NPTk'
    '5FQ1RJT05fTE9TVBABEhgKFFBJTkdfU1RBVFVTX0lOQUNUSVZFEAI=');

@$core.Deprecated('Use platformTypeDescriptor instead')
const PlatformType$json = {
  '1': 'PlatformType',
  '2': [
    {'1': 'PLATFORM_TYPE_LINUX', '2': 0},
    {'1': 'PLATFORM_TYPE_MACOS', '2': 1},
    {'1': 'PLATFORM_TYPE_WINDOWS', '2': 2},
  ],
};

/// Descriptor for `PlatformType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List platformTypeDescriptor = $convert.base64Decode(
    'CgxQbGF0Zm9ybVR5cGUSFwoTUExBVEZPUk1fVFlQRV9MSU5VWBAAEhcKE1BMQVRGT1JNX1RZUE'
    'VfTUFDT1MQARIZChVQTEFURk9STV9UWVBFX1dJTkRPV1MQAg==');

@$core.Deprecated('Use rebootOptionDescriptor instead')
const RebootOption$json = {
  '1': 'RebootOption',
  '2': [
    {'1': 'REBOOT_OPTION_REBOOT_IF_NEEDED', '2': 0},
    {'1': 'REBOOT_OPTION_NO_REBOOT', '2': 1},
  ],
};

/// Descriptor for `RebootOption`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List rebootOptionDescriptor = $convert.base64Decode(
    'CgxSZWJvb3RPcHRpb24SIgoeUkVCT09UX09QVElPTl9SRUJPT1RfSUZfTkVFREVEEAASGwoXUk'
    'VCT09UX09QVElPTl9OT19SRUJPT1QQAQ==');

@$core.Deprecated('Use resourceDataSyncS3FormatDescriptor instead')
const ResourceDataSyncS3Format$json = {
  '1': 'ResourceDataSyncS3Format',
  '2': [
    {'1': 'RESOURCE_DATA_SYNC_S3_FORMAT_JSON_SERDE', '2': 0},
  ],
};

/// Descriptor for `ResourceDataSyncS3Format`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceDataSyncS3FormatDescriptor =
    $convert.base64Decode(
        'ChhSZXNvdXJjZURhdGFTeW5jUzNGb3JtYXQSKwonUkVTT1VSQ0VfREFUQV9TWU5DX1MzX0ZPUk'
        '1BVF9KU09OX1NFUkRFEAA=');

@$core.Deprecated('Use resourceTypeDescriptor instead')
const ResourceType$json = {
  '1': 'ResourceType',
  '2': [
    {'1': 'RESOURCE_TYPE_EC2_INSTANCE', '2': 0},
    {'1': 'RESOURCE_TYPE_MANAGED_INSTANCE', '2': 1},
  ],
};

/// Descriptor for `ResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceTypeDescriptor = $convert.base64Decode(
    'CgxSZXNvdXJjZVR5cGUSHgoaUkVTT1VSQ0VfVFlQRV9FQzJfSU5TVEFOQ0UQABIiCh5SRVNPVV'
    'JDRV9UWVBFX01BTkFHRURfSU5TVEFOQ0UQAQ==');

@$core.Deprecated('Use resourceTypeForTaggingDescriptor instead')
const ResourceTypeForTagging$json = {
  '1': 'ResourceTypeForTagging',
  '2': [
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_MAINTENANCE_WINDOW', '2': 0},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_OPS_ITEM', '2': 1},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_DOCUMENT', '2': 2},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_AUTOMATION', '2': 3},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_PATCH_BASELINE', '2': 4},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_OPSMETADATA', '2': 5},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_PARAMETER', '2': 6},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_MANAGED_INSTANCE', '2': 7},
    {'1': 'RESOURCE_TYPE_FOR_TAGGING_ASSOCIATION', '2': 8},
  ],
};

/// Descriptor for `ResourceTypeForTagging`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceTypeForTaggingDescriptor = $convert.base64Decode(
    'ChZSZXNvdXJjZVR5cGVGb3JUYWdnaW5nEjAKLFJFU09VUkNFX1RZUEVfRk9SX1RBR0dJTkdfTU'
    'FJTlRFTkFOQ0VfV0lORE9XEAASJgoiUkVTT1VSQ0VfVFlQRV9GT1JfVEFHR0lOR19PUFNfSVRF'
    'TRABEiYKIlJFU09VUkNFX1RZUEVfRk9SX1RBR0dJTkdfRE9DVU1FTlQQAhIoCiRSRVNPVVJDRV'
    '9UWVBFX0ZPUl9UQUdHSU5HX0FVVE9NQVRJT04QAxIsCihSRVNPVVJDRV9UWVBFX0ZPUl9UQUdH'
    'SU5HX1BBVENIX0JBU0VMSU5FEAQSKQolUkVTT1VSQ0VfVFlQRV9GT1JfVEFHR0lOR19PUFNNRV'
    'RBREFUQRAFEicKI1JFU09VUkNFX1RZUEVfRk9SX1RBR0dJTkdfUEFSQU1FVEVSEAYSLgoqUkVT'
    'T1VSQ0VfVFlQRV9GT1JfVEFHR0lOR19NQU5BR0VEX0lOU1RBTkNFEAcSKQolUkVTT1VSQ0VfVF'
    'lQRV9GT1JfVEFHR0lOR19BU1NPQ0lBVElPThAI');

@$core.Deprecated('Use reviewStatusDescriptor instead')
const ReviewStatus$json = {
  '1': 'ReviewStatus',
  '2': [
    {'1': 'REVIEW_STATUS_PENDING', '2': 0},
    {'1': 'REVIEW_STATUS_REJECTED', '2': 1},
    {'1': 'REVIEW_STATUS_NOT_REVIEWED', '2': 2},
    {'1': 'REVIEW_STATUS_APPROVED', '2': 3},
  ],
};

/// Descriptor for `ReviewStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reviewStatusDescriptor = $convert.base64Decode(
    'CgxSZXZpZXdTdGF0dXMSGQoVUkVWSUVXX1NUQVRVU19QRU5ESU5HEAASGgoWUkVWSUVXX1NUQV'
    'RVU19SRUpFQ1RFRBABEh4KGlJFVklFV19TVEFUVVNfTk9UX1JFVklFV0VEEAISGgoWUkVWSUVX'
    'X1NUQVRVU19BUFBST1ZFRBAD');

@$core.Deprecated('Use sessionFilterKeyDescriptor instead')
const SessionFilterKey$json = {
  '1': 'SessionFilterKey',
  '2': [
    {'1': 'SESSION_FILTER_KEY_SESSION_ID', '2': 0},
    {'1': 'SESSION_FILTER_KEY_STATUS', '2': 1},
    {'1': 'SESSION_FILTER_KEY_ACCESS_TYPE', '2': 2},
    {'1': 'SESSION_FILTER_KEY_INVOKED_BEFORE', '2': 3},
    {'1': 'SESSION_FILTER_KEY_OWNER', '2': 4},
    {'1': 'SESSION_FILTER_KEY_TARGET_ID', '2': 5},
    {'1': 'SESSION_FILTER_KEY_INVOKED_AFTER', '2': 6},
  ],
};

/// Descriptor for `SessionFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sessionFilterKeyDescriptor = $convert.base64Decode(
    'ChBTZXNzaW9uRmlsdGVyS2V5EiEKHVNFU1NJT05fRklMVEVSX0tFWV9TRVNTSU9OX0lEEAASHQ'
    'oZU0VTU0lPTl9GSUxURVJfS0VZX1NUQVRVUxABEiIKHlNFU1NJT05fRklMVEVSX0tFWV9BQ0NF'
    'U1NfVFlQRRACEiUKIVNFU1NJT05fRklMVEVSX0tFWV9JTlZPS0VEX0JFRk9SRRADEhwKGFNFU1'
    'NJT05fRklMVEVSX0tFWV9PV05FUhAEEiAKHFNFU1NJT05fRklMVEVSX0tFWV9UQVJHRVRfSUQQ'
    'BRIkCiBTRVNTSU9OX0ZJTFRFUl9LRVlfSU5WT0tFRF9BRlRFUhAG');

@$core.Deprecated('Use sessionStateDescriptor instead')
const SessionState$json = {
  '1': 'SessionState',
  '2': [
    {'1': 'SESSION_STATE_HISTORY', '2': 0},
    {'1': 'SESSION_STATE_ACTIVE', '2': 1},
  ],
};

/// Descriptor for `SessionState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sessionStateDescriptor = $convert.base64Decode(
    'CgxTZXNzaW9uU3RhdGUSGQoVU0VTU0lPTl9TVEFURV9ISVNUT1JZEAASGAoUU0VTU0lPTl9TVE'
    'FURV9BQ1RJVkUQAQ==');

@$core.Deprecated('Use sessionStatusDescriptor instead')
const SessionStatus$json = {
  '1': 'SessionStatus',
  '2': [
    {'1': 'SESSION_STATUS_TERMINATING', '2': 0},
    {'1': 'SESSION_STATUS_CONNECTING', '2': 1},
    {'1': 'SESSION_STATUS_TERMINATED', '2': 2},
    {'1': 'SESSION_STATUS_DISCONNECTED', '2': 3},
    {'1': 'SESSION_STATUS_CONNECTED', '2': 4},
    {'1': 'SESSION_STATUS_FAILED', '2': 5},
  ],
};

/// Descriptor for `SessionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sessionStatusDescriptor = $convert.base64Decode(
    'Cg1TZXNzaW9uU3RhdHVzEh4KGlNFU1NJT05fU1RBVFVTX1RFUk1JTkFUSU5HEAASHQoZU0VTU0'
    'lPTl9TVEFUVVNfQ09OTkVDVElORxABEh0KGVNFU1NJT05fU1RBVFVTX1RFUk1JTkFURUQQAhIf'
    'ChtTRVNTSU9OX1NUQVRVU19ESVNDT05ORUNURUQQAxIcChhTRVNTSU9OX1NUQVRVU19DT05ORU'
    'NURUQQBBIZChVTRVNTSU9OX1NUQVRVU19GQUlMRUQQBQ==');

@$core.Deprecated('Use signalTypeDescriptor instead')
const SignalType$json = {
  '1': 'SignalType',
  '2': [
    {'1': 'SIGNAL_TYPE_REVOKE', '2': 0},
    {'1': 'SIGNAL_TYPE_STOP_STEP', '2': 1},
    {'1': 'SIGNAL_TYPE_START_STEP', '2': 2},
    {'1': 'SIGNAL_TYPE_REJECT', '2': 3},
    {'1': 'SIGNAL_TYPE_APPROVE', '2': 4},
    {'1': 'SIGNAL_TYPE_RESUME', '2': 5},
  ],
};

/// Descriptor for `SignalType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List signalTypeDescriptor = $convert.base64Decode(
    'CgpTaWduYWxUeXBlEhYKElNJR05BTF9UWVBFX1JFVk9LRRAAEhkKFVNJR05BTF9UWVBFX1NUT1'
    'BfU1RFUBABEhoKFlNJR05BTF9UWVBFX1NUQVJUX1NURVAQAhIWChJTSUdOQUxfVFlQRV9SRUpF'
    'Q1QQAxIXChNTSUdOQUxfVFlQRV9BUFBST1ZFEAQSFgoSU0lHTkFMX1RZUEVfUkVTVU1FEAU=');

@$core.Deprecated('Use sourceTypeDescriptor instead')
const SourceType$json = {
  '1': 'SourceType',
  '2': [
    {'1': 'SOURCE_TYPE_AWS_SSM_MANAGEDINSTANCE', '2': 0},
    {'1': 'SOURCE_TYPE_AWS_EC2_INSTANCE', '2': 1},
    {'1': 'SOURCE_TYPE_AWS_IOT_THING', '2': 2},
  ],
};

/// Descriptor for `SourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sourceTypeDescriptor = $convert.base64Decode(
    'CgpTb3VyY2VUeXBlEicKI1NPVVJDRV9UWVBFX0FXU19TU01fTUFOQUdFRElOU1RBTkNFEAASIA'
    'ocU09VUkNFX1RZUEVfQVdTX0VDMl9JTlNUQU5DRRABEh0KGVNPVVJDRV9UWVBFX0FXU19JT1Rf'
    'VEhJTkcQAg==');

@$core.Deprecated('Use stepExecutionFilterKeyDescriptor instead')
const StepExecutionFilterKey$json = {
  '1': 'StepExecutionFilterKey',
  '2': [
    {'1': 'STEP_EXECUTION_FILTER_KEY_START_TIME_AFTER', '2': 0},
    {'1': 'STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_STATUS', '2': 1},
    {'1': 'STEP_EXECUTION_FILTER_KEY_STEP_NAME', '2': 2},
    {'1': 'STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_ID', '2': 3},
    {'1': 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATION', '2': 4},
    {'1': 'STEP_EXECUTION_FILTER_KEY_START_TIME_BEFORE', '2': 5},
    {'1': 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATOR_VALUE', '2': 6},
    {'1': 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_EXECUTION_ID', '2': 7},
    {'1': 'STEP_EXECUTION_FILTER_KEY_ACTION', '2': 8},
  ],
};

/// Descriptor for `StepExecutionFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stepExecutionFilterKeyDescriptor = $convert.base64Decode(
    'ChZTdGVwRXhlY3V0aW9uRmlsdGVyS2V5Ei4KKlNURVBfRVhFQ1VUSU9OX0ZJTFRFUl9LRVlfU1'
    'RBUlRfVElNRV9BRlRFUhAAEjMKL1NURVBfRVhFQ1VUSU9OX0ZJTFRFUl9LRVlfU1RFUF9FWEVD'
    'VVRJT05fU1RBVFVTEAESJwojU1RFUF9FWEVDVVRJT05fRklMVEVSX0tFWV9TVEVQX05BTUUQAh'
    'IvCitTVEVQX0VYRUNVVElPTl9GSUxURVJfS0VZX1NURVBfRVhFQ1VUSU9OX0lEEAMSMwovU1RF'
    'UF9FWEVDVVRJT05fRklMVEVSX0tFWV9QQVJFTlRfU1RFUF9JVEVSQVRJT04QBBIvCitTVEVQX0'
    'VYRUNVVElPTl9GSUxURVJfS0VZX1NUQVJUX1RJTUVfQkVGT1JFEAUSOAo0U1RFUF9FWEVDVVRJ'
    'T05fRklMVEVSX0tFWV9QQVJFTlRfU1RFUF9JVEVSQVRPUl9WQUxVRRAGEjYKMlNURVBfRVhFQ1'
    'VUSU9OX0ZJTFRFUl9LRVlfUEFSRU5UX1NURVBfRVhFQ1VUSU9OX0lEEAcSJAogU1RFUF9FWEVD'
    'VVRJT05fRklMVEVSX0tFWV9BQ1RJT04QCA==');

@$core.Deprecated('Use stopTypeDescriptor instead')
const StopType$json = {
  '1': 'StopType',
  '2': [
    {'1': 'STOP_TYPE_COMPLETE', '2': 0},
    {'1': 'STOP_TYPE_CANCEL', '2': 1},
  ],
};

/// Descriptor for `StopType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List stopTypeDescriptor = $convert.base64Decode(
    'CghTdG9wVHlwZRIWChJTVE9QX1RZUEVfQ09NUExFVEUQABIUChBTVE9QX1RZUEVfQ0FOQ0VMEA'
    'E=');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use accountSharingInfoDescriptor instead')
const AccountSharingInfo$json = {
  '1': 'AccountSharingInfo',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'shareddocumentversion',
      '3': 240155144,
      '4': 1,
      '5': 9,
      '10': 'shareddocumentversion'
    },
  ],
};

/// Descriptor for `AccountSharingInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountSharingInfoDescriptor = $convert.base64Decode(
    'ChJBY2NvdW50U2hhcmluZ0luZm8SHwoJYWNjb3VudGlkGNLBuR8gASgJUglhY2NvdW50aWQSNw'
    'oVc2hhcmVkZG9jdW1lbnR2ZXJzaW9uGIj0wXIgASgJUhVzaGFyZWRkb2N1bWVudHZlcnNpb24=');

@$core.Deprecated('Use activationDescriptor instead')
const Activation$json = {
  '1': 'Activation',
  '2': [
    {'1': 'activationid', '3': 146858601, '4': 1, '5': 9, '10': 'activationid'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {
      '1': 'defaultinstancename',
      '3': 27070947,
      '4': 1,
      '5': 9,
      '10': 'defaultinstancename'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'expirationdate',
      '3': 17577725,
      '4': 1,
      '5': 9,
      '10': 'expirationdate'
    },
    {'1': 'expired', '3': 480900051, '4': 1, '5': 8, '10': 'expired'},
    {'1': 'iamrole', '3': 242424351, '4': 1, '5': 9, '10': 'iamrole'},
    {
      '1': 'registrationlimit',
      '3': 312393964,
      '4': 1,
      '5': 5,
      '10': 'registrationlimit'
    },
    {
      '1': 'registrationscount',
      '3': 54916841,
      '4': 1,
      '5': 5,
      '10': 'registrationscount'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `Activation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List activationDescriptor = $convert.base64Decode(
    'CgpBY3RpdmF0aW9uEiUKDGFjdGl2YXRpb25pZBjpxINGIAEoCVIMYWN0aXZhdGlvbmlkEiQKC2'
    'NyZWF0ZWRkYXRlGLCw58YBIAEoCVILY3JlYXRlZGRhdGUSMwoTZGVmYXVsdGluc3RhbmNlbmFt'
    'ZRjjo/QMIAEoCVITZGVmYXVsdGluc3RhbmNlbmFtZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCV'
    'ILZGVzY3JpcHRpb24SKQoOZXhwaXJhdGlvbmRhdGUY/e2wCCABKAlSDmV4cGlyYXRpb25kYXRl'
    'EhwKB2V4cGlyZWQY0+en5QEgASgIUgdleHBpcmVkEhsKB2lhbXJvbGUYn7TMcyABKAlSB2lhbX'
    'JvbGUSMAoRcmVnaXN0cmF0aW9ubGltaXQY7IH7lAEgASgFUhFyZWdpc3RyYXRpb25saW1pdBIx'
    'ChJyZWdpc3RyYXRpb25zY291bnQY6e2XGiABKAVSEnJlZ2lzdHJhdGlvbnNjb3VudBIgCgR0YW'
    'dzGMHB9rUBIAMoCzIILnNzbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use addTagsToResourceRequestDescriptor instead')
const AddTagsToResourceRequest$json = {
  '1': 'AddTagsToResourceRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceTypeForTagging',
      '10': 'resourcetype'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `AddTagsToResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsToResourceRequestDescriptor = $convert.base64Decode(
    'ChhBZGRUYWdzVG9SZXNvdXJjZVJlcXVlc3QSIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc2'
    '91cmNlaWQSQwoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoDjIbLnNzbS5SZXNvdXJjZVR5cGVGb3JU'
    'YWdnaW5nUgxyZXNvdXJjZXR5cGUSIAoEdGFncxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdz');

@$core.Deprecated('Use addTagsToResourceResultDescriptor instead')
const AddTagsToResourceResult$json = {
  '1': 'AddTagsToResourceResult',
};

/// Descriptor for `AddTagsToResourceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsToResourceResultDescriptor =
    $convert.base64Decode('ChdBZGRUYWdzVG9SZXNvdXJjZVJlc3VsdA==');

@$core.Deprecated('Use alarmDescriptor instead')
const Alarm$json = {
  '1': 'Alarm',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `Alarm`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmDescriptor =
    $convert.base64Decode('CgVBbGFybRIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use alarmConfigurationDescriptor instead')
const AlarmConfiguration$json = {
  '1': 'AlarmConfiguration',
  '2': [
    {
      '1': 'alarms',
      '3': 70013740,
      '4': 3,
      '5': 11,
      '6': '.ssm.Alarm',
      '10': 'alarms'
    },
    {
      '1': 'ignorepollalarmfailure',
      '3': 331523304,
      '4': 1,
      '5': 8,
      '10': 'ignorepollalarmfailure'
    },
  ],
};

/// Descriptor for `AlarmConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmConfigurationDescriptor = $convert.base64Decode(
    'ChJBbGFybUNvbmZpZ3VyYXRpb24SJQoGYWxhcm1zGKymsSEgAygLMgouc3NtLkFsYXJtUgZhbG'
    'FybXMSOgoWaWdub3JlcG9sbGFsYXJtZmFpbHVyZRjoyYqeASABKAhSFmlnbm9yZXBvbGxhbGFy'
    'bWZhaWx1cmU=');

@$core.Deprecated('Use alarmStateInformationDescriptor instead')
const AlarmStateInformation$json = {
  '1': 'AlarmStateInformation',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.ssm.ExternalAlarmState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `AlarmStateInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alarmStateInformationDescriptor = $convert.base64Decode(
    'ChVBbGFybVN0YXRlSW5mb3JtYXRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIxCgVzdGF0ZR'
    'iXybLvASABKA4yFy5zc20uRXh0ZXJuYWxBbGFybVN0YXRlUgVzdGF0ZQ==');

@$core.Deprecated('Use alreadyExistsExceptionDescriptor instead')
const AlreadyExistsException$json = {
  '1': 'AlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChZBbHJlYWR5RXhpc3RzRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use associateOpsItemRelatedItemRequestDescriptor instead')
const AssociateOpsItemRelatedItemRequest$json = {
  '1': 'AssociateOpsItemRelatedItemRequest',
  '2': [
    {
      '1': 'associationtype',
      '3': 165803189,
      '4': 1,
      '5': 9,
      '10': 'associationtype'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'resourceuri', '3': 442608014, '4': 1, '5': 9, '10': 'resourceuri'},
  ],
};

/// Descriptor for `AssociateOpsItemRelatedItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateOpsItemRelatedItemRequestDescriptor =
    $convert.base64Decode(
        'CiJBc3NvY2lhdGVPcHNJdGVtUmVsYXRlZEl0ZW1SZXF1ZXN0EisKD2Fzc29jaWF0aW9udHlwZR'
        'i16YdPIAEoCVIPYXNzb2NpYXRpb250eXBlEh8KCW9wc2l0ZW1pZBjS0pUMIAEoCVIJb3BzaXRl'
        'bWlkEiYKDHJlc291cmNldHlwZRjevtiPASABKAlSDHJlc291cmNldHlwZRIkCgtyZXNvdXJjZX'
        'VyaRiO04bTASABKAlSC3Jlc291cmNldXJp');

@$core.Deprecated('Use associateOpsItemRelatedItemResponseDescriptor instead')
const AssociateOpsItemRelatedItemResponse$json = {
  '1': 'AssociateOpsItemRelatedItemResponse',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
  ],
};

/// Descriptor for `AssociateOpsItemRelatedItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associateOpsItemRelatedItemResponseDescriptor =
    $convert.base64Decode(
        'CiNBc3NvY2lhdGVPcHNJdGVtUmVsYXRlZEl0ZW1SZXNwb25zZRInCg1hc3NvY2lhdGlvbmlkGJ'
        'L8lUIgASgJUg1hc3NvY2lhdGlvbmlk');

@$core.Deprecated('Use associatedInstancesDescriptor instead')
const AssociatedInstances$json = {
  '1': 'AssociatedInstances',
};

/// Descriptor for `AssociatedInstances`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associatedInstancesDescriptor =
    $convert.base64Decode('ChNBc3NvY2lhdGVkSW5zdGFuY2Vz');

@$core.Deprecated('Use associationDescriptor instead')
const Association$json = {
  '1': 'Association',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'lastexecutiondate',
      '3': 492296302,
      '4': 1,
      '5': 9,
      '10': 'lastexecutiondate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'overview',
      '3': 172576699,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationOverview',
      '10': 'overview'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.Association.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [Association_TargetmapsEntry$json],
};

@$core.Deprecated('Use associationDescriptor instead')
const Association_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Association`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationDescriptor = $convert.base64Decode(
    'CgtBc3NvY2lhdGlvbhInCg1hc3NvY2lhdGlvbmlkGJL8lUIgASgJUg1hc3NvY2lhdGlvbmlkEi'
    'wKD2Fzc29jaWF0aW9ubmFtZRiYkMWVASABKAlSD2Fzc29jaWF0aW9ubmFtZRIyChJhc3NvY2lh'
    'dGlvbnZlcnNpb24YkYrJ1QEgASgJUhJhc3NvY2lhdGlvbnZlcnNpb24SKwoPZG9jdW1lbnR2ZX'
    'JzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SHgoIZHVyYXRpb24YrpKdpgEgASgFUghk'
    'dXJhdGlvbhIhCgppbnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEjAKEWxhc3RleGVjdX'
    'Rpb25kYXRlGO6w3+oBIAEoCVIRbGFzdGV4ZWN1dGlvbmRhdGUSFQoEbmFtZRiH5oF/IAEoCVIE'
    'bmFtZRI3CghvdmVydmlldxi7n6VSIAEoCzIYLnNzbS5Bc3NvY2lhdGlvbk92ZXJ2aWV3Ughvdm'
    'VydmlldxIyChJzY2hlZHVsZWV4cHJlc3Npb24Y/5Hb1AEgASgJUhJzY2hlZHVsZWV4cHJlc3Np'
    'b24SKQoOc2NoZWR1bGVvZmZzZXQY2JHqSiABKAVSDnNjaGVkdWxlb2Zmc2V0EkMKCnRhcmdldG'
    '1hcHMYuLzVIyADKAsyIC5zc20uQXNzb2NpYXRpb24uVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRt'
    'YXBzEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzGj0KD1RhcmdldG'
    '1hcHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use associationAlreadyExistsDescriptor instead')
const AssociationAlreadyExists$json = {
  '1': 'AssociationAlreadyExists',
};

/// Descriptor for `AssociationAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationAlreadyExistsDescriptor =
    $convert.base64Decode('ChhBc3NvY2lhdGlvbkFscmVhZHlFeGlzdHM=');

@$core.Deprecated('Use associationDescriptionDescriptor instead')
const AssociationDescription$json = {
  '1': 'AssociationDescription',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'applyonlyatcroninterval',
      '3': 285867158,
      '4': 1,
      '5': 8,
      '10': 'applyonlyatcroninterval'
    },
    {
      '1': 'associationdispatchassumerole',
      '3': 281627425,
      '4': 1,
      '5': 9,
      '10': 'associationdispatchassumerole'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'automationtargetparametername',
      '3': 348584826,
      '4': 1,
      '5': 9,
      '10': 'automationtargetparametername'
    },
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
    {
      '1': 'complianceseverity',
      '3': 278891158,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationComplianceSeverity',
      '10': 'complianceseverity'
    },
    {'1': 'date', '3': 458388346, '4': 1, '5': 9, '10': 'date'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'lastexecutiondate',
      '3': 492296302,
      '4': 1,
      '5': 9,
      '10': 'lastexecutiondate'
    },
    {
      '1': 'lastsuccessfulexecutiondate',
      '3': 160123084,
      '4': 1,
      '5': 9,
      '10': 'lastsuccessfulexecutiondate'
    },
    {
      '1': 'lastupdateassociationdate',
      '3': 223483812,
      '4': 1,
      '5': 9,
      '10': 'lastupdateassociationdate'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'overview',
      '3': 172576699,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationOverview',
      '10': 'overview'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationDescription.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationStatus',
      '10': 'status'
    },
    {
      '1': 'synccompliance',
      '3': 500469318,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationSyncCompliance',
      '10': 'synccompliance'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationDescription.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
  ],
  '3': [
    AssociationDescription_ParametersEntry$json,
    AssociationDescription_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use associationDescriptionDescriptor instead')
const AssociationDescription_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use associationDescriptionDescriptor instead')
const AssociationDescription_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AssociationDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationDescriptionDescriptor = $convert.base64Decode(
    'ChZBc3NvY2lhdGlvbkRlc2NyaXB0aW9uEkoKEmFsYXJtY29uZmlndXJhdGlvbhiJmbkhIAEoCz'
    'IXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlvbhI8ChdhcHBseW9u'
    'bHlhdGNyb25pbnRlcnZhbBiW+aeIASABKAhSF2FwcGx5b25seWF0Y3JvbmludGVydmFsEkgKHW'
    'Fzc29jaWF0aW9uZGlzcGF0Y2hhc3N1bWVyb2xlGKGWpYYBIAEoCVIdYXNzb2NpYXRpb25kaXNw'
    'YXRjaGFzc3VtZXJvbGUSJwoNYXNzb2NpYXRpb25pZBiS/JVCIAEoCVINYXNzb2NpYXRpb25pZB'
    'IsCg9hc3NvY2lhdGlvbm5hbWUYmJDFlQEgASgJUg9hc3NvY2lhdGlvbm5hbWUSMgoSYXNzb2Np'
    'YXRpb252ZXJzaW9uGJGKydUBIAEoCVISYXNzb2NpYXRpb252ZXJzaW9uEkgKHWF1dG9tYXRpb2'
    '50YXJnZXRwYXJhbWV0ZXJuYW1lGPr2m6YBIAEoCVIdYXV0b21hdGlvbnRhcmdldHBhcmFtZXRl'
    'cm5hbWUSJwoNY2FsZW5kYXJuYW1lcxi+85kRIAMoCVINY2FsZW5kYXJuYW1lcxJWChJjb21wbG'
    'lhbmNlc2V2ZXJpdHkYlpX+hAEgASgOMiIuc3NtLkFzc29jaWF0aW9uQ29tcGxpYW5jZVNldmVy'
    'aXR5UhJjb21wbGlhbmNlc2V2ZXJpdHkSFgoEZGF0ZRj65snaASABKAlSBGRhdGUSKwoPZG9jdW'
    '1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SHgoIZHVyYXRpb24YrpKdpgEg'
    'ASgFUghkdXJhdGlvbhIhCgppbnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEjAKEWxhc3'
    'RleGVjdXRpb25kYXRlGO6w3+oBIAEoCVIRbGFzdGV4ZWN1dGlvbmRhdGUSQwobbGFzdHN1Y2Nl'
    'c3NmdWxleGVjdXRpb25kYXRlGMyRrUwgASgJUhtsYXN0c3VjY2Vzc2Z1bGV4ZWN1dGlvbmRhdG'
    'USPwoZbGFzdHVwZGF0ZWFzc29jaWF0aW9uZGF0ZRikr8hqIAEoCVIZbGFzdHVwZGF0ZWFzc29j'
    'aWF0aW9uZGF0ZRIpCg5tYXhjb25jdXJyZW5jeRj9wY4OIAEoCVIObWF4Y29uY3VycmVuY3kSHw'
    'oJbWF4ZXJyb3JzGKvC9T0gASgJUgltYXhlcnJvcnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRJR'
    'Cg5vdXRwdXRsb2NhdGlvbhj067UgIAEoCzImLnNzbS5JbnN0YW5jZUFzc29jaWF0aW9uT3V0cH'
    'V0TG9jYXRpb25SDm91dHB1dGxvY2F0aW9uEjcKCG92ZXJ2aWV3GLufpVIgASgLMhguc3NtLkFz'
    'c29jaWF0aW9uT3ZlcnZpZXdSCG92ZXJ2aWV3Ek8KCnBhcmFtZXRlcnMY+qf+6wEgAygLMisuc3'
    'NtLkFzc29jaWF0aW9uRGVzY3JpcHRpb24uUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEjIK'
    'EnNjaGVkdWxlZXhwcmVzc2lvbhj/kdvUASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhIpCg5zY2'
    'hlZHVsZW9mZnNldBjYkepKIAEoBVIOc2NoZWR1bGVvZmZzZXQSMQoGc3RhdHVzGJDk+wIgASgL'
    'MhYuc3NtLkFzc29jaWF0aW9uU3RhdHVzUgZzdGF0dXMSSgoOc3luY2NvbXBsaWFuY2UYxpzS7g'
    'EgASgOMh4uc3NtLkFzc29jaWF0aW9uU3luY0NvbXBsaWFuY2VSDnN5bmNjb21wbGlhbmNlEkEK'
    'D3RhcmdldGxvY2F0aW9ucxilu/GJASADKAsyEy5zc20uVGFyZ2V0TG9jYXRpb25SD3RhcmdldG'
    'xvY2F0aW9ucxJOCgp0YXJnZXRtYXBzGLi81SMgAygLMisuc3NtLkFzc29jaWF0aW9uRGVzY3Jp'
    'cHRpb24uVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEigKB3RhcmdldHMYgpuCfSADKAsyCy'
    '5zc20uVGFyZ2V0Ugd0YXJnZXRzEkcKD3RyaWdnZXJlZGFsYXJtcxiF7cF9IAMoCzIaLnNzbS5B'
    'bGFybVN0YXRlSW5mb3JtYXRpb25SD3RyaWdnZXJlZGFsYXJtcxo9Cg9QYXJhbWV0ZXJzRW50cn'
    'kSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARo9Cg9UYXJnZXRt'
    'YXBzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use associationDoesNotExistDescriptor instead')
const AssociationDoesNotExist$json = {
  '1': 'AssociationDoesNotExist',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AssociationDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationDoesNotExistDescriptor =
    $convert.base64Decode(
        'ChdBc3NvY2lhdGlvbkRvZXNOb3RFeGlzdBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use associationExecutionDescriptor instead')
const AssociationExecution$json = {
  '1': 'AssociationExecution',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {
      '1': 'detailedstatus',
      '3': 223171376,
      '4': 1,
      '5': 9,
      '10': 'detailedstatus'
    },
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {
      '1': 'lastexecutiondate',
      '3': 492296302,
      '4': 1,
      '5': 9,
      '10': 'lastexecutiondate'
    },
    {
      '1': 'resourcecountbystatus',
      '3': 129812756,
      '4': 1,
      '5': 9,
      '10': 'resourcecountbystatus'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
  ],
};

/// Descriptor for `AssociationExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationExecutionDescriptor = $convert.base64Decode(
    'ChRBc3NvY2lhdGlvbkV4ZWN1dGlvbhJKChJhbGFybWNvbmZpZ3VyYXRpb24YiZm5ISABKAsyFy'
    '5zc20uQWxhcm1Db25maWd1cmF0aW9uUhJhbGFybWNvbmZpZ3VyYXRpb24SJwoNYXNzb2NpYXRp'
    'b25pZBiS/JVCIAEoCVINYXNzb2NpYXRpb25pZBIyChJhc3NvY2lhdGlvbnZlcnNpb24YkYrJ1Q'
    'EgASgJUhJhc3NvY2lhdGlvbnZlcnNpb24SIwoLY3JlYXRlZHRpbWUY8+vzOSABKAlSC2NyZWF0'
    'ZWR0aW1lEikKDmRldGFpbGVkc3RhdHVzGLCmtWogASgJUg5kZXRhaWxlZHN0YXR1cxIjCgtleG'
    'VjdXRpb25pZBixz69GIAEoCVILZXhlY3V0aW9uaWQSMAoRbGFzdGV4ZWN1dGlvbmRhdGUY7rDf'
    '6gEgASgJUhFsYXN0ZXhlY3V0aW9uZGF0ZRI3ChVyZXNvdXJjZWNvdW50YnlzdGF0dXMYlJLzPS'
    'ABKAlSFXJlc291cmNlY291bnRieXN0YXR1cxIZCgZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cxJH'
    'Cg90cmlnZ2VyZWRhbGFybXMYhe3BfSADKAsyGi5zc20uQWxhcm1TdGF0ZUluZm9ybWF0aW9uUg'
    '90cmlnZ2VyZWRhbGFybXM=');

@$core.Deprecated('Use associationExecutionDoesNotExistDescriptor instead')
const AssociationExecutionDoesNotExist$json = {
  '1': 'AssociationExecutionDoesNotExist',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AssociationExecutionDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationExecutionDoesNotExistDescriptor =
    $convert.base64Decode(
        'CiBBc3NvY2lhdGlvbkV4ZWN1dGlvbkRvZXNOb3RFeGlzdBIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use associationExecutionFilterDescriptor instead')
const AssociationExecutionFilter$json = {
  '1': 'AssociationExecutionFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationExecutionFilterKey',
      '10': 'key'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationFilterOperatorType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `AssociationExecutionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationExecutionFilterDescriptor = $convert.base64Decode(
    'ChpBc3NvY2lhdGlvbkV4ZWN1dGlvbkZpbHRlchI3CgNrZXkYjZLraCABKA4yIi5zc20uQXNzb2'
    'NpYXRpb25FeGVjdXRpb25GaWx0ZXJLZXlSA2tleRI6CgR0eXBlGO6g14oBIAEoDjIiLnNzbS5B'
    'c3NvY2lhdGlvbkZpbHRlck9wZXJhdG9yVHlwZVIEdHlwZRIYCgV2YWx1ZRjr8p+KASABKAlSBX'
    'ZhbHVl');

@$core.Deprecated('Use associationExecutionTargetDescriptor instead')
const AssociationExecutionTarget$json = {
  '1': 'AssociationExecutionTarget',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'detailedstatus',
      '3': 223171376,
      '4': 1,
      '5': 9,
      '10': 'detailedstatus'
    },
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {
      '1': 'lastexecutiondate',
      '3': 492296302,
      '4': 1,
      '5': 9,
      '10': 'lastexecutiondate'
    },
    {
      '1': 'outputsource',
      '3': 125784002,
      '4': 1,
      '5': 11,
      '6': '.ssm.OutputSource',
      '10': 'outputsource'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `AssociationExecutionTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationExecutionTargetDescriptor = $convert.base64Decode(
    'ChpBc3NvY2lhdGlvbkV4ZWN1dGlvblRhcmdldBInCg1hc3NvY2lhdGlvbmlkGJL8lUIgASgJUg'
    '1hc3NvY2lhdGlvbmlkEjIKEmFzc29jaWF0aW9udmVyc2lvbhiRisnVASABKAlSEmFzc29jaWF0'
    'aW9udmVyc2lvbhIpCg5kZXRhaWxlZHN0YXR1cxiwprVqIAEoCVIOZGV0YWlsZWRzdGF0dXMSIw'
    'oLZXhlY3V0aW9uaWQYsc+vRiABKAlSC2V4ZWN1dGlvbmlkEjAKEWxhc3RleGVjdXRpb25kYXRl'
    'GO6w3+oBIAEoCVIRbGFzdGV4ZWN1dGlvbmRhdGUSOAoMb3V0cHV0c291cmNlGMKf/TsgASgLMh'
    'Euc3NtLk91dHB1dFNvdXJjZVIMb3V0cHV0c291cmNlEiIKCnJlc291cmNlaWQYkbrx+gEgASgJ'
    'UgpyZXNvdXJjZWlkEiYKDHJlc291cmNldHlwZRjevtiPASABKAlSDHJlc291cmNldHlwZRIZCg'
    'ZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use associationExecutionTargetsFilterDescriptor instead')
const AssociationExecutionTargetsFilter$json = {
  '1': 'AssociationExecutionTargetsFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationExecutionTargetsFilterKey',
      '10': 'key'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `AssociationExecutionTargetsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationExecutionTargetsFilterDescriptor =
    $convert.base64Decode(
        'CiFBc3NvY2lhdGlvbkV4ZWN1dGlvblRhcmdldHNGaWx0ZXISPgoDa2V5GI2S62ggASgOMikuc3'
        'NtLkFzc29jaWF0aW9uRXhlY3V0aW9uVGFyZ2V0c0ZpbHRlcktleVIDa2V5EhgKBXZhbHVlGOvy'
        'n4oBIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use associationFilterDescriptor instead')
const AssociationFilter$json = {
  '1': 'AssociationFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationFilterKey',
      '10': 'key'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `AssociationFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationFilterDescriptor = $convert.base64Decode(
    'ChFBc3NvY2lhdGlvbkZpbHRlchIuCgNrZXkY7ZDXQCABKA4yGS5zc20uQXNzb2NpYXRpb25GaW'
    'x0ZXJLZXlSA2tleRIXCgV2YWx1ZRjLp/sSIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use associationLimitExceededDescriptor instead')
const AssociationLimitExceeded$json = {
  '1': 'AssociationLimitExceeded',
};

/// Descriptor for `AssociationLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationLimitExceededDescriptor =
    $convert.base64Decode('ChhBc3NvY2lhdGlvbkxpbWl0RXhjZWVkZWQ=');

@$core.Deprecated('Use associationOverviewDescriptor instead')
const AssociationOverview$json = {
  '1': 'AssociationOverview',
  '2': [
    {
      '1': 'associationstatusaggregatedcount',
      '3': 159361291,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationOverview.AssociationstatusaggregatedcountEntry',
      '10': 'associationstatusaggregatedcount'
    },
    {
      '1': 'detailedstatus',
      '3': 223171376,
      '4': 1,
      '5': 9,
      '10': 'detailedstatus'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
  '3': [AssociationOverview_AssociationstatusaggregatedcountEntry$json],
};

@$core.Deprecated('Use associationOverviewDescriptor instead')
const AssociationOverview_AssociationstatusaggregatedcountEntry$json = {
  '1': 'AssociationstatusaggregatedcountEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AssociationOverview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationOverviewDescriptor = $convert.base64Decode(
    'ChNBc3NvY2lhdGlvbk92ZXJ2aWV3Eo0BCiBhc3NvY2lhdGlvbnN0YXR1c2FnZ3JlZ2F0ZWRjb3'
    'VudBiL0v5LIAMoCzI+LnNzbS5Bc3NvY2lhdGlvbk92ZXJ2aWV3LkFzc29jaWF0aW9uc3RhdHVz'
    'YWdncmVnYXRlZGNvdW50RW50cnlSIGFzc29jaWF0aW9uc3RhdHVzYWdncmVnYXRlZGNvdW50Ei'
    'kKDmRldGFpbGVkc3RhdHVzGLCmtWogASgJUg5kZXRhaWxlZHN0YXR1cxIZCgZzdGF0dXMYkOT7'
    'AiABKAlSBnN0YXR1cxpTCiVBc3NvY2lhdGlvbnN0YXR1c2FnZ3JlZ2F0ZWRjb3VudEVudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgFUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use associationStatusDescriptor instead')
const AssociationStatus$json = {
  '1': 'AssociationStatus',
  '2': [
    {
      '1': 'additionalinfo',
      '3': 62481793,
      '4': 1,
      '5': 9,
      '10': 'additionalinfo'
    },
    {'1': 'date', '3': 458388346, '4': 1, '5': 9, '10': 'date'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationStatusName',
      '10': 'name'
    },
  ],
};

/// Descriptor for `AssociationStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationStatusDescriptor = $convert.base64Decode(
    'ChFBc3NvY2lhdGlvblN0YXR1cxIpCg5hZGRpdGlvbmFsaW5mbxiBy+UdIAEoCVIOYWRkaXRpb2'
    '5hbGluZm8SFgoEZGF0ZRj65snaASABKAlSBGRhdGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVz'
    'c2FnZRIxCgRuYW1lGIfmgX8gASgOMhouc3NtLkFzc29jaWF0aW9uU3RhdHVzTmFtZVIEbmFtZQ'
    '==');

@$core.Deprecated('Use associationVersionInfoDescriptor instead')
const AssociationVersionInfo$json = {
  '1': 'AssociationVersionInfo',
  '2': [
    {
      '1': 'applyonlyatcroninterval',
      '3': 285867158,
      '4': 1,
      '5': 8,
      '10': 'applyonlyatcroninterval'
    },
    {
      '1': 'associationdispatchassumerole',
      '3': 281627425,
      '4': 1,
      '5': 9,
      '10': 'associationdispatchassumerole'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
    {
      '1': 'complianceseverity',
      '3': 278891158,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationComplianceSeverity',
      '10': 'complianceseverity'
    },
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationVersionInfo.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'synccompliance',
      '3': 500469318,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationSyncCompliance',
      '10': 'synccompliance'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationVersionInfo.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    AssociationVersionInfo_ParametersEntry$json,
    AssociationVersionInfo_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use associationVersionInfoDescriptor instead')
const AssociationVersionInfo_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use associationVersionInfoDescriptor instead')
const AssociationVersionInfo_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AssociationVersionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationVersionInfoDescriptor = $convert.base64Decode(
    'ChZBc3NvY2lhdGlvblZlcnNpb25JbmZvEjwKF2FwcGx5b25seWF0Y3JvbmludGVydmFsGJb5p4'
    'gBIAEoCFIXYXBwbHlvbmx5YXRjcm9uaW50ZXJ2YWwSSAodYXNzb2NpYXRpb25kaXNwYXRjaGFz'
    'c3VtZXJvbGUYoZalhgEgASgJUh1hc3NvY2lhdGlvbmRpc3BhdGNoYXNzdW1lcm9sZRInCg1hc3'
    'NvY2lhdGlvbmlkGJL8lUIgASgJUg1hc3NvY2lhdGlvbmlkEiwKD2Fzc29jaWF0aW9ubmFtZRiY'
    'kMWVASABKAlSD2Fzc29jaWF0aW9ubmFtZRIyChJhc3NvY2lhdGlvbnZlcnNpb24YkYrJ1QEgAS'
    'gJUhJhc3NvY2lhdGlvbnZlcnNpb24SJwoNY2FsZW5kYXJuYW1lcxi+85kRIAMoCVINY2FsZW5k'
    'YXJuYW1lcxJWChJjb21wbGlhbmNlc2V2ZXJpdHkYlpX+hAEgASgOMiIuc3NtLkFzc29jaWF0aW'
    '9uQ29tcGxpYW5jZVNldmVyaXR5UhJjb21wbGlhbmNlc2V2ZXJpdHkSJAoLY3JlYXRlZGRhdGUY'
    'sLDnxgEgASgJUgtjcmVhdGVkZGF0ZRIrCg9kb2N1bWVudHZlcnNpb24Yye+pKCABKAlSD2RvY3'
    'VtZW50dmVyc2lvbhIeCghkdXJhdGlvbhiukp2mASABKAVSCGR1cmF0aW9uEikKDm1heGNvbmN1'
    'cnJlbmN5GP3Bjg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcnMYq8L1PSABKAlSCW'
    '1heGVycm9ycxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElEKDm91dHB1dGxvY2F0aW9uGPTrtSAg'
    'ASgLMiYuc3NtLkluc3RhbmNlQXNzb2NpYXRpb25PdXRwdXRMb2NhdGlvblIOb3V0cHV0bG9jYX'
    'Rpb24STwoKcGFyYW1ldGVycxj6p/7rASADKAsyKy5zc20uQXNzb2NpYXRpb25WZXJzaW9uSW5m'
    'by5QYXJhbWV0ZXJzRW50cnlSCnBhcmFtZXRlcnMSMgoSc2NoZWR1bGVleHByZXNzaW9uGP+R29'
    'QBIAEoCVISc2NoZWR1bGVleHByZXNzaW9uEikKDnNjaGVkdWxlb2Zmc2V0GNiR6kogASgFUg5z'
    'Y2hlZHVsZW9mZnNldBJKCg5zeW5jY29tcGxpYW5jZRjGnNLuASABKA4yHi5zc20uQXNzb2NpYX'
    'Rpb25TeW5jQ29tcGxpYW5jZVIOc3luY2NvbXBsaWFuY2USQQoPdGFyZ2V0bG9jYXRpb25zGKW7'
    '8YkBIAMoCzITLnNzbS5UYXJnZXRMb2NhdGlvblIPdGFyZ2V0bG9jYXRpb25zEk4KCnRhcmdldG'
    '1hcHMYuLzVIyADKAsyKy5zc20uQXNzb2NpYXRpb25WZXJzaW9uSW5mby5UYXJnZXRtYXBzRW50'
    'cnlSCnRhcmdldG1hcHMSKAoHdGFyZ2V0cxiCm4J9IAMoCzILLnNzbS5UYXJnZXRSB3RhcmdldH'
    'MaPQoPUGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2'
    'YWx1ZToCOAEaPQoPVGFyZ2V0bWFwc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGA'
    'IgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use associationVersionLimitExceededDescriptor instead')
const AssociationVersionLimitExceeded$json = {
  '1': 'AssociationVersionLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AssociationVersionLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List associationVersionLimitExceededDescriptor =
    $convert.base64Decode(
        'Ch9Bc3NvY2lhdGlvblZlcnNpb25MaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use attachmentContentDescriptor instead')
const AttachmentContent$json = {
  '1': 'AttachmentContent',
  '2': [
    {'1': 'hash', '3': 250828530, '4': 1, '5': 9, '10': 'hash'},
    {
      '1': 'hashtype',
      '3': 172838330,
      '4': 1,
      '5': 14,
      '6': '.ssm.AttachmentHashType',
      '10': 'hashtype'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'size', '3': 105352829, '4': 1, '5': 3, '10': 'size'},
    {'1': 'url', '3': 354018239, '4': 1, '5': 9, '10': 'url'},
  ],
};

/// Descriptor for `AttachmentContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachmentContentDescriptor = $convert.base64Decode(
    'ChFBdHRhY2htZW50Q29udGVudBIVCgRoYXNoGPKtzXcgASgJUgRoYXNoEjYKCGhhc2h0eXBlGL'
    'qbtVIgASgOMhcuc3NtLkF0dGFjaG1lbnRIYXNoVHlwZVIIaGFzaHR5cGUSFQoEbmFtZRiH5oF/'
    'IAEoCVIEbmFtZRIVCgRzaXplGP2cnjIgASgDUgRzaXplEhQKA3VybBi/x+eoASABKAlSA3VybA'
    '==');

@$core.Deprecated('Use attachmentInformationDescriptor instead')
const AttachmentInformation$json = {
  '1': 'AttachmentInformation',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `AttachmentInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachmentInformationDescriptor =
    $convert.base64Decode(
        'ChVBdHRhY2htZW50SW5mb3JtYXRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use attachmentsSourceDescriptor instead')
const AttachmentsSource$json = {
  '1': 'AttachmentsSource',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.AttachmentsSourceKey',
      '10': 'key'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `AttachmentsSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachmentsSourceDescriptor = $convert.base64Decode(
    'ChFBdHRhY2htZW50c1NvdXJjZRIuCgNrZXkYjZLraCABKA4yGS5zc20uQXR0YWNobWVudHNTb3'
    'VyY2VLZXlSA2tleRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhkKBnZhbHVlcxjcxLRqIAMoCVIG'
    'dmFsdWVz');

@$core.Deprecated(
    'Use automationDefinitionNotApprovedExceptionDescriptor instead')
const AutomationDefinitionNotApprovedException$json = {
  '1': 'AutomationDefinitionNotApprovedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationDefinitionNotApprovedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationDefinitionNotApprovedExceptionDescriptor =
    $convert.base64Decode(
        'CihBdXRvbWF0aW9uRGVmaW5pdGlvbk5vdEFwcHJvdmVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use automationDefinitionNotFoundExceptionDescriptor instead')
const AutomationDefinitionNotFoundException$json = {
  '1': 'AutomationDefinitionNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationDefinitionNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationDefinitionNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiVBdXRvbWF0aW9uRGVmaW5pdGlvbk5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cC'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use automationDefinitionVersionNotFoundExceptionDescriptor instead')
const AutomationDefinitionVersionNotFoundException$json = {
  '1': 'AutomationDefinitionVersionNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationDefinitionVersionNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    automationDefinitionVersionNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CixBdXRvbWF0aW9uRGVmaW5pdGlvblZlcnNpb25Ob3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYW'
        'dlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use automationExecutionDescriptor instead')
const AutomationExecution$json = {
  '1': 'AutomationExecution',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
    {
      '1': 'automationexecutionstatus',
      '3': 531364411,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationExecutionStatus',
      '10': 'automationexecutionstatus'
    },
    {
      '1': 'automationsubtype',
      '3': 153542097,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationSubtype',
      '10': 'automationsubtype'
    },
    {
      '1': 'changerequestname',
      '3': 468779362,
      '4': 1,
      '5': 9,
      '10': 'changerequestname'
    },
    {
      '1': 'currentaction',
      '3': 337681497,
      '4': 1,
      '5': 9,
      '10': 'currentaction'
    },
    {
      '1': 'currentstepname',
      '3': 126991996,
      '4': 1,
      '5': 9,
      '10': 'currentstepname'
    },
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'executedby', '3': 186546754, '4': 1, '5': 9, '10': 'executedby'},
    {
      '1': 'executionendtime',
      '3': 139859196,
      '4': 1,
      '5': 9,
      '10': 'executionendtime'
    },
    {
      '1': 'executionstarttime',
      '3': 429847391,
      '4': 1,
      '5': 9,
      '10': 'executionstarttime'
    },
    {
      '1': 'failuremessage',
      '3': 353556937,
      '4': 1,
      '5': 9,
      '10': 'failuremessage'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.ssm.ExecutionMode',
      '10': 'mode'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {
      '1': 'outputs',
      '3': 455868918,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecution.OutputsEntry',
      '10': 'outputs'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecution.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'parentautomationexecutionid',
      '3': 232537390,
      '4': 1,
      '5': 9,
      '10': 'parentautomationexecutionid'
    },
    {
      '1': 'progresscounters',
      '3': 162419216,
      '4': 1,
      '5': 11,
      '6': '.ssm.ProgressCounters',
      '10': 'progresscounters'
    },
    {
      '1': 'resolvedtargets',
      '3': 361325602,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResolvedTargets',
      '10': 'resolvedtargets'
    },
    {
      '1': 'runbooks',
      '3': 514418725,
      '4': 3,
      '5': 11,
      '6': '.ssm.Runbook',
      '10': 'runbooks'
    },
    {
      '1': 'scheduledtime',
      '3': 334708242,
      '4': 1,
      '5': 9,
      '10': 'scheduledtime'
    },
    {
      '1': 'stepexecutions',
      '3': 164157981,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecution',
      '10': 'stepexecutions'
    },
    {
      '1': 'stepexecutionstruncated',
      '3': 472961151,
      '4': 1,
      '5': 8,
      '10': 'stepexecutionstruncated'
    },
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetlocationsurl',
      '3': 107583422,
      '4': 1,
      '5': 9,
      '10': 'targetlocationsurl'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecution.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targetparametername',
      '3': 351056597,
      '4': 1,
      '5': 9,
      '10': 'targetparametername'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
    {
      '1': 'variables',
      '3': 429322339,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecution.VariablesEntry',
      '10': 'variables'
    },
  ],
  '3': [
    AutomationExecution_OutputsEntry$json,
    AutomationExecution_ParametersEntry$json,
    AutomationExecution_TargetmapsEntry$json,
    AutomationExecution_VariablesEntry$json
  ],
};

@$core.Deprecated('Use automationExecutionDescriptor instead')
const AutomationExecution_OutputsEntry$json = {
  '1': 'OutputsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use automationExecutionDescriptor instead')
const AutomationExecution_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use automationExecutionDescriptor instead')
const AutomationExecution_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use automationExecutionDescriptor instead')
const AutomationExecution_VariablesEntry$json = {
  '1': 'VariablesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AutomationExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionDescriptor = $convert.base64Decode(
    'ChNBdXRvbWF0aW9uRXhlY3V0aW9uEkoKEmFsYXJtY29uZmlndXJhdGlvbhiJmbkhIAEoCzIXLn'
    'NzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlvbhInCg1hc3NvY2lhdGlv'
    'bmlkGJL8lUIgASgJUg1hc3NvY2lhdGlvbmlkEjcKFWF1dG9tYXRpb25leGVjdXRpb25pZBjm7/'
    'cFIAEoCVIVYXV0b21hdGlvbmV4ZWN1dGlvbmlkEmAKGWF1dG9tYXRpb25leGVjdXRpb25zdGF0'
    'dXMYu/Sv/QEgASgOMh4uc3NtLkF1dG9tYXRpb25FeGVjdXRpb25TdGF0dXNSGWF1dG9tYXRpb2'
    '5leGVjdXRpb25zdGF0dXMSRwoRYXV0b21hdGlvbnN1YnR5cGUY0bubSSABKA4yFi5zc20uQXV0'
    'b21hdGlvblN1YnR5cGVSEWF1dG9tYXRpb25zdWJ0eXBlEjAKEWNoYW5nZXJlcXVlc3RuYW1lGO'
    'KCxN8BIAEoCVIRY2hhbmdlcmVxdWVzdG5hbWUSKAoNY3VycmVudGFjdGlvbhjZuIKhASABKAlS'
    'DWN1cnJlbnRhY3Rpb24SKwoPY3VycmVudHN0ZXBuYW1lGPz8xjwgASgJUg9jdXJyZW50c3RlcG'
    '5hbWUSJQoMZG9jdW1lbnRuYW1lGNCjxzkgASgJUgxkb2N1bWVudG5hbWUSKwoPZG9jdW1lbnR2'
    'ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SIQoKZXhlY3V0ZWRieRjC9PlYIAEoCV'
    'IKZXhlY3V0ZWRieRItChBleGVjdXRpb25lbmR0aW1lGPyp2EIgASgJUhBleGVjdXRpb25lbmR0'
    'aW1lEjIKEmV4ZWN1dGlvbnN0YXJ0dGltZRjf5vvMASABKAlSEmV4ZWN1dGlvbnN0YXJ0dGltZR'
    'IqCg5mYWlsdXJlbWVzc2FnZRjJs8uoASABKAlSDmZhaWx1cmVtZXNzYWdlEikKDm1heGNvbmN1'
    'cnJlbmN5GP3Bjg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcnMYq8L1PSABKAlSCW'
    '1heGVycm9ycxIqCgRtb2RlGLPuuZoBIAEoDjISLnNzbS5FeGVjdXRpb25Nb2RlUgRtb2RlEh8K'
    'CW9wc2l0ZW1pZBjS0pUMIAEoCVIJb3BzaXRlbWlkEkMKB291dHB1dHMY9oOw2QEgAygLMiUuc3'
    'NtLkF1dG9tYXRpb25FeGVjdXRpb24uT3V0cHV0c0VudHJ5UgdvdXRwdXRzEkwKCnBhcmFtZXRl'
    'cnMY+qf+6wEgAygLMiguc3NtLkF1dG9tYXRpb25FeGVjdXRpb24uUGFyYW1ldGVyc0VudHJ5Ug'
    'pwYXJhbWV0ZXJzEkMKG3BhcmVudGF1dG9tYXRpb25leGVjdXRpb25pZBiu+vBuIAEoCVIbcGFy'
    'ZW50YXV0b21hdGlvbmV4ZWN1dGlvbmlkEkQKEHByb2dyZXNzY291bnRlcnMYkKS5TSABKAsyFS'
    '5zc20uUHJvZ3Jlc3NDb3VudGVyc1IQcHJvZ3Jlc3Njb3VudGVycxJCCg9yZXNvbHZlZHRhcmdl'
    'dHMYosilrAEgASgLMhQuc3NtLlJlc29sdmVkVGFyZ2V0c1IPcmVzb2x2ZWR0YXJnZXRzEiwKCH'
    'J1bmJvb2tzGKXQpfUBIAMoCzIMLnNzbS5SdW5ib29rUghydW5ib29rcxIoCg1zY2hlZHVsZWR0'
    'aW1lGJL8zJ8BIAEoCVINc2NoZWR1bGVkdGltZRI9Cg5zdGVwZXhlY3V0aW9ucxidtKNOIAMoCz'
    'ISLnNzbS5TdGVwRXhlY3V0aW9uUg5zdGVwZXhlY3V0aW9ucxI8ChdzdGVwZXhlY3V0aW9uc3Ry'
    'dW5jYXRlZBj/oMPhASABKAhSF3N0ZXBleGVjdXRpb25zdHJ1bmNhdGVkEhkKBnRhcmdldBjp4p'
    '9bIAEoCVIGdGFyZ2V0EkEKD3RhcmdldGxvY2F0aW9ucxilu/GJASADKAsyEy5zc20uVGFyZ2V0'
    'TG9jYXRpb25SD3RhcmdldGxvY2F0aW9ucxIxChJ0YXJnZXRsb2NhdGlvbnN1cmwYvq+mMyABKA'
    'lSEnRhcmdldGxvY2F0aW9uc3VybBJLCgp0YXJnZXRtYXBzGLi81SMgAygLMiguc3NtLkF1dG9t'
    'YXRpb25FeGVjdXRpb24uVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEjQKE3RhcmdldHBhcm'
    'FtZXRlcm5hbWUY1eWypwEgASgJUhN0YXJnZXRwYXJhbWV0ZXJuYW1lEigKB3RhcmdldHMYgpuC'
    'fSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEkcKD3RyaWdnZXJlZGFsYXJtcxiF7cF9IAMoCz'
    'IaLnNzbS5BbGFybVN0YXRlSW5mb3JtYXRpb25SD3RyaWdnZXJlZGFsYXJtcxJJCgl2YXJpYWJs'
    'ZXMY4+DbzAEgAygLMicuc3NtLkF1dG9tYXRpb25FeGVjdXRpb24uVmFyaWFibGVzRW50cnlSCX'
    'ZhcmlhYmxlcxo6CgxPdXRwdXRzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiAB'
    'KAlSBXZhbHVlOgI4ARo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdm'
    'FsdWUYAiABKAlSBXZhbHVlOgI4ARo9Cg9UYXJnZXRtYXBzRW50cnkSEAoDa2V5GAEgASgJUgNr'
    'ZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARo8Cg5WYXJpYWJsZXNFbnRyeRIQCgNrZXkYAS'
    'ABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use automationExecutionFilterDescriptor instead')
const AutomationExecutionFilter$json = {
  '1': 'AutomationExecutionFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationExecutionFilterKey',
      '10': 'key'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `AutomationExecutionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionFilterDescriptor = $convert.base64Decode(
    'ChlBdXRvbWF0aW9uRXhlY3V0aW9uRmlsdGVyEjYKA2tleRiNkutoIAEoDjIhLnNzbS5BdXRvbW'
    'F0aW9uRXhlY3V0aW9uRmlsdGVyS2V5UgNrZXkSGQoGdmFsdWVzGNzEtGogAygJUgZ2YWx1ZXM=');

@$core.Deprecated('Use automationExecutionInputsDescriptor instead')
const AutomationExecutionInputs$json = {
  '1': 'AutomationExecutionInputs',
  '2': [
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionInputs.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetlocationsurl',
      '3': 107583422,
      '4': 1,
      '5': 9,
      '10': 'targetlocationsurl'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionInputs.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targetparametername',
      '3': 351056597,
      '4': 1,
      '5': 9,
      '10': 'targetparametername'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    AutomationExecutionInputs_ParametersEntry$json,
    AutomationExecutionInputs_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use automationExecutionInputsDescriptor instead')
const AutomationExecutionInputs_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use automationExecutionInputsDescriptor instead')
const AutomationExecutionInputs_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AutomationExecutionInputs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionInputsDescriptor = $convert.base64Decode(
    'ChlBdXRvbWF0aW9uRXhlY3V0aW9uSW5wdXRzElIKCnBhcmFtZXRlcnMY+qf+6wEgAygLMi4uc3'
    'NtLkF1dG9tYXRpb25FeGVjdXRpb25JbnB1dHMuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJz'
    'EkEKD3RhcmdldGxvY2F0aW9ucxilu/GJASADKAsyEy5zc20uVGFyZ2V0TG9jYXRpb25SD3Rhcm'
    'dldGxvY2F0aW9ucxIxChJ0YXJnZXRsb2NhdGlvbnN1cmwYvq+mMyABKAlSEnRhcmdldGxvY2F0'
    'aW9uc3VybBJRCgp0YXJnZXRtYXBzGLi81SMgAygLMi4uc3NtLkF1dG9tYXRpb25FeGVjdXRpb2'
    '5JbnB1dHMuVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEjQKE3RhcmdldHBhcmFtZXRlcm5h'
    'bWUY1eWypwEgASgJUhN0YXJnZXRwYXJhbWV0ZXJuYW1lEigKB3RhcmdldHMYgpuCfSADKAsyCy'
    '5zc20uVGFyZ2V0Ugd0YXJnZXRzGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGj0KD1RhcmdldG1hcHNFbnRyeRIQCgNrZXkYAS'
    'ABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated(
    'Use automationExecutionLimitExceededExceptionDescriptor instead')
const AutomationExecutionLimitExceededException$json = {
  '1': 'AutomationExecutionLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationExecutionLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    automationExecutionLimitExceededExceptionDescriptor = $convert.base64Decode(
        'CilBdXRvbWF0aW9uRXhlY3V0aW9uTGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use automationExecutionMetadataDescriptor instead')
const AutomationExecutionMetadata$json = {
  '1': 'AutomationExecutionMetadata',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
    {
      '1': 'automationexecutionstatus',
      '3': 531364411,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationExecutionStatus',
      '10': 'automationexecutionstatus'
    },
    {
      '1': 'automationsubtype',
      '3': 153542097,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationSubtype',
      '10': 'automationsubtype'
    },
    {
      '1': 'automationtype',
      '3': 72473291,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationType',
      '10': 'automationtype'
    },
    {
      '1': 'changerequestname',
      '3': 468779362,
      '4': 1,
      '5': 9,
      '10': 'changerequestname'
    },
    {
      '1': 'currentaction',
      '3': 337681497,
      '4': 1,
      '5': 9,
      '10': 'currentaction'
    },
    {
      '1': 'currentstepname',
      '3': 126991996,
      '4': 1,
      '5': 9,
      '10': 'currentstepname'
    },
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'executedby', '3': 186546754, '4': 1, '5': 9, '10': 'executedby'},
    {
      '1': 'executionendtime',
      '3': 139859196,
      '4': 1,
      '5': 9,
      '10': 'executionendtime'
    },
    {
      '1': 'executionstarttime',
      '3': 429847391,
      '4': 1,
      '5': 9,
      '10': 'executionstarttime'
    },
    {
      '1': 'failuremessage',
      '3': 353556937,
      '4': 1,
      '5': 9,
      '10': 'failuremessage'
    },
    {'1': 'logfile', '3': 13056280, '4': 1, '5': 9, '10': 'logfile'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.ssm.ExecutionMode',
      '10': 'mode'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {
      '1': 'outputs',
      '3': 455868918,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionMetadata.OutputsEntry',
      '10': 'outputs'
    },
    {
      '1': 'parentautomationexecutionid',
      '3': 232537390,
      '4': 1,
      '5': 9,
      '10': 'parentautomationexecutionid'
    },
    {
      '1': 'resolvedtargets',
      '3': 361325602,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResolvedTargets',
      '10': 'resolvedtargets'
    },
    {
      '1': 'runbooks',
      '3': 514418725,
      '4': 3,
      '5': 11,
      '6': '.ssm.Runbook',
      '10': 'runbooks'
    },
    {
      '1': 'scheduledtime',
      '3': 334708242,
      '4': 1,
      '5': 9,
      '10': 'scheduledtime'
    },
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
    {
      '1': 'targetlocationsurl',
      '3': 107583422,
      '4': 1,
      '5': 9,
      '10': 'targetlocationsurl'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionMetadata.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targetparametername',
      '3': 351056597,
      '4': 1,
      '5': 9,
      '10': 'targetparametername'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
  ],
  '3': [
    AutomationExecutionMetadata_OutputsEntry$json,
    AutomationExecutionMetadata_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use automationExecutionMetadataDescriptor instead')
const AutomationExecutionMetadata_OutputsEntry$json = {
  '1': 'OutputsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use automationExecutionMetadataDescriptor instead')
const AutomationExecutionMetadata_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AutomationExecutionMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionMetadataDescriptor = $convert.base64Decode(
    'ChtBdXRvbWF0aW9uRXhlY3V0aW9uTWV0YWRhdGESSgoSYWxhcm1jb25maWd1cmF0aW9uGImZuS'
    'EgASgLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uEicKDWFz'
    'c29jaWF0aW9uaWQYkvyVQiABKAlSDWFzc29jaWF0aW9uaWQSNwoVYXV0b21hdGlvbmV4ZWN1dG'
    'lvbmlkGObv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQSYAoZYXV0b21hdGlvbmV4ZWN1'
    'dGlvbnN0YXR1cxi79K/9ASABKA4yHi5zc20uQXV0b21hdGlvbkV4ZWN1dGlvblN0YXR1c1IZYX'
    'V0b21hdGlvbmV4ZWN1dGlvbnN0YXR1cxJHChFhdXRvbWF0aW9uc3VidHlwZRjRu5tJIAEoDjIW'
    'LnNzbS5BdXRvbWF0aW9uU3VidHlwZVIRYXV0b21hdGlvbnN1YnR5cGUSPgoOYXV0b21hdGlvbn'
    'R5cGUYy7XHIiABKA4yEy5zc20uQXV0b21hdGlvblR5cGVSDmF1dG9tYXRpb250eXBlEjAKEWNo'
    'YW5nZXJlcXVlc3RuYW1lGOKCxN8BIAEoCVIRY2hhbmdlcmVxdWVzdG5hbWUSKAoNY3VycmVudG'
    'FjdGlvbhjZuIKhASABKAlSDWN1cnJlbnRhY3Rpb24SKwoPY3VycmVudHN0ZXBuYW1lGPz8xjwg'
    'ASgJUg9jdXJyZW50c3RlcG5hbWUSJQoMZG9jdW1lbnRuYW1lGNCjxzkgASgJUgxkb2N1bWVudG'
    '5hbWUSKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SIQoKZXhl'
    'Y3V0ZWRieRjC9PlYIAEoCVIKZXhlY3V0ZWRieRItChBleGVjdXRpb25lbmR0aW1lGPyp2EIgAS'
    'gJUhBleGVjdXRpb25lbmR0aW1lEjIKEmV4ZWN1dGlvbnN0YXJ0dGltZRjf5vvMASABKAlSEmV4'
    'ZWN1dGlvbnN0YXJ0dGltZRIqCg5mYWlsdXJlbWVzc2FnZRjJs8uoASABKAlSDmZhaWx1cmVtZX'
    'NzYWdlEhsKB2xvZ2ZpbGUYmPKcBiABKAlSB2xvZ2ZpbGUSKQoObWF4Y29uY3VycmVuY3kY/cGO'
    'DiABKAlSDm1heGNvbmN1cnJlbmN5Eh8KCW1heGVycm9ycxirwvU9IAEoCVIJbWF4ZXJyb3JzEi'
    'oKBG1vZGUYs+65mgEgASgOMhIuc3NtLkV4ZWN1dGlvbk1vZGVSBG1vZGUSHwoJb3BzaXRlbWlk'
    'GNLSlQwgASgJUglvcHNpdGVtaWQSSwoHb3V0cHV0cxj2g7DZASADKAsyLS5zc20uQXV0b21hdG'
    'lvbkV4ZWN1dGlvbk1ldGFkYXRhLk91dHB1dHNFbnRyeVIHb3V0cHV0cxJDChtwYXJlbnRhdXRv'
    'bWF0aW9uZXhlY3V0aW9uaWQYrvrwbiABKAlSG3BhcmVudGF1dG9tYXRpb25leGVjdXRpb25pZB'
    'JCCg9yZXNvbHZlZHRhcmdldHMYosilrAEgASgLMhQuc3NtLlJlc29sdmVkVGFyZ2V0c1IPcmVz'
    'b2x2ZWR0YXJnZXRzEiwKCHJ1bmJvb2tzGKXQpfUBIAMoCzIMLnNzbS5SdW5ib29rUghydW5ib2'
    '9rcxIoCg1zY2hlZHVsZWR0aW1lGJL8zJ8BIAEoCVINc2NoZWR1bGVkdGltZRIZCgZ0YXJnZXQY'
    '6eKfWyABKAlSBnRhcmdldBIxChJ0YXJnZXRsb2NhdGlvbnN1cmwYvq+mMyABKAlSEnRhcmdldG'
    'xvY2F0aW9uc3VybBJTCgp0YXJnZXRtYXBzGLi81SMgAygLMjAuc3NtLkF1dG9tYXRpb25FeGVj'
    'dXRpb25NZXRhZGF0YS5UYXJnZXRtYXBzRW50cnlSCnRhcmdldG1hcHMSNAoTdGFyZ2V0cGFyYW'
    '1ldGVybmFtZRjV5bKnASABKAlSE3RhcmdldHBhcmFtZXRlcm5hbWUSKAoHdGFyZ2V0cxiCm4J9'
    'IAMoCzILLnNzbS5UYXJnZXRSB3RhcmdldHMSRwoPdHJpZ2dlcmVkYWxhcm1zGIXtwX0gAygLMh'
    'ouc3NtLkFsYXJtU3RhdGVJbmZvcm1hdGlvblIPdHJpZ2dlcmVkYWxhcm1zGjoKDE91dHB1dHNF'
    'bnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGj0KD1Rhcm'
    'dldG1hcHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use automationExecutionNotFoundExceptionDescriptor instead')
const AutomationExecutionNotFoundException$json = {
  '1': 'AutomationExecutionNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationExecutionNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiRBdXRvbWF0aW9uRXhlY3V0aW9uTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use automationExecutionPreviewDescriptor instead')
const AutomationExecutionPreview$json = {
  '1': 'AutomationExecutionPreview',
  '2': [
    {'1': 'regions', '3': 36200107, '4': 3, '5': 9, '10': 'regions'},
    {
      '1': 'steppreviews',
      '3': 182130687,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionPreview.SteppreviewsEntry',
      '10': 'steppreviews'
    },
    {
      '1': 'targetpreviews',
      '3': 53018126,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetPreview',
      '10': 'targetpreviews'
    },
    {
      '1': 'totalaccounts',
      '3': 108841974,
      '4': 1,
      '5': 5,
      '10': 'totalaccounts'
    },
  ],
  '3': [AutomationExecutionPreview_SteppreviewsEntry$json],
};

@$core.Deprecated('Use automationExecutionPreviewDescriptor instead')
const AutomationExecutionPreview_SteppreviewsEntry$json = {
  '1': 'SteppreviewsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AutomationExecutionPreview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationExecutionPreviewDescriptor = $convert.base64Decode(
    'ChpBdXRvbWF0aW9uRXhlY3V0aW9uUHJldmlldxIbCgdyZWdpb25zGKu9oREgAygJUgdyZWdpb2'
    '5zElgKDHN0ZXBwcmV2aWV3cxj/r+xWIAMoCzIxLnNzbS5BdXRvbWF0aW9uRXhlY3V0aW9uUHJl'
    'dmlldy5TdGVwcHJldmlld3NFbnRyeVIMc3RlcHByZXZpZXdzEj0KDnRhcmdldHByZXZpZXdzGI'
    '78oxkgAygLMhIuc3NtLlRhcmdldFByZXZpZXdSDnRhcmdldHByZXZpZXdzEicKDXRvdGFsYWNj'
    'b3VudHMY9pfzMyABKAVSDXRvdGFsYWNjb3VudHMaPwoRU3RlcHByZXZpZXdzRW50cnkSEAoDa2'
    'V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use automationStepNotFoundExceptionDescriptor instead')
const AutomationStepNotFoundException$json = {
  '1': 'AutomationStepNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AutomationStepNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List automationStepNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch9BdXRvbWF0aW9uU3RlcE5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use baselineOverrideDescriptor instead')
const BaselineOverride$json = {
  '1': 'BaselineOverride',
  '2': [
    {
      '1': 'approvalrules',
      '3': 71466346,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchRuleGroup',
      '10': 'approvalrules'
    },
    {
      '1': 'approvedpatches',
      '3': 199384709,
      '4': 3,
      '5': 9,
      '10': 'approvedpatches'
    },
    {
      '1': 'approvedpatchescompliancelevel',
      '3': 63924432,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'approvedpatchescompliancelevel'
    },
    {
      '1': 'approvedpatchesenablenonsecurity',
      '3': 295555901,
      '4': 1,
      '5': 8,
      '10': 'approvedpatchesenablenonsecurity'
    },
    {
      '1': 'availablesecurityupdatescompliancestatus',
      '3': 187471858,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceStatus',
      '10': 'availablesecurityupdatescompliancestatus'
    },
    {
      '1': 'globalfilters',
      '3': 263302754,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'globalfilters'
    },
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {
      '1': 'rejectedpatches',
      '3': 309657116,
      '4': 3,
      '5': 9,
      '10': 'rejectedpatches'
    },
    {
      '1': 'rejectedpatchesaction',
      '3': 356538330,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchAction',
      '10': 'rejectedpatchesaction'
    },
    {
      '1': 'sources',
      '3': 46625746,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchSource',
      '10': 'sources'
    },
  ],
};

/// Descriptor for `BaselineOverride`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List baselineOverrideDescriptor = $convert.base64Decode(
    'ChBCYXNlbGluZU92ZXJyaWRlEjwKDWFwcHJvdmFscnVsZXMY6vqJIiABKAsyEy5zc20uUGF0Y2'
    'hSdWxlR3JvdXBSDWFwcHJvdmFscnVsZXMSKwoPYXBwcm92ZWRwYXRjaGVzGIW9iV8gAygJUg9h'
    'cHByb3ZlZHBhdGNoZXMSZAoeYXBwcm92ZWRwYXRjaGVzY29tcGxpYW5jZWxldmVsGNDRvR4gAS'
    'gOMhkuc3NtLlBhdGNoQ29tcGxpYW5jZUxldmVsUh5hcHByb3ZlZHBhdGNoZXNjb21wbGlhbmNl'
    'bGV2ZWwSTgogYXBwcm92ZWRwYXRjaGVzZW5hYmxlbm9uc2VjdXJpdHkYvab3jAEgASgIUiBhcH'
    'Byb3ZlZHBhdGNoZXNlbmFibGVub25zZWN1cml0eRJ5CihhdmFpbGFibGVzZWN1cml0eXVwZGF0'
    'ZXNjb21wbGlhbmNlc3RhdHVzGPKvslkgASgOMhouc3NtLlBhdGNoQ29tcGxpYW5jZVN0YXR1c1'
    'IoYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVzY29tcGxpYW5jZXN0YXR1cxI+Cg1nbG9iYWxmaWx0'
    'ZXJzGOLcxn0gASgLMhUuc3NtLlBhdGNoRmlsdGVyR3JvdXBSDWdsb2JhbGZpbHRlcnMSQQoPb3'
    'BlcmF0aW5nc3lzdGVtGOr9wRIgASgOMhQuc3NtLk9wZXJhdGluZ1N5c3RlbVIPb3BlcmF0aW5n'
    'c3lzdGVtEiwKD3JlamVjdGVkcGF0Y2hlcxic/NOTASADKAlSD3JlamVjdGVkcGF0Y2hlcxJKCh'
    'VyZWplY3RlZHBhdGNoZXNhY3Rpb24Y2q+BqgEgASgOMhAuc3NtLlBhdGNoQWN0aW9uUhVyZWpl'
    'Y3RlZHBhdGNoZXNhY3Rpb24SLQoHc291cmNlcxjS550WIAMoCzIQLnNzbS5QYXRjaFNvdXJjZV'
    'IHc291cmNlcw==');

@$core.Deprecated('Use cancelCommandRequestDescriptor instead')
const CancelCommandRequest$json = {
  '1': 'CancelCommandRequest',
  '2': [
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {'1': 'instanceids', '3': 312792453, '4': 3, '5': 9, '10': 'instanceids'},
  ],
};

/// Descriptor for `CancelCommandRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelCommandRequestDescriptor = $convert.base64Decode(
    'ChRDYW5jZWxDb21tYW5kUmVxdWVzdBIfCgljb21tYW5kaWQYgNuATCABKAlSCWNvbW1hbmRpZB'
    'IkCgtpbnN0YW5jZWlkcxiFq5OVASADKAlSC2luc3RhbmNlaWRz');

@$core.Deprecated('Use cancelCommandResultDescriptor instead')
const CancelCommandResult$json = {
  '1': 'CancelCommandResult',
};

/// Descriptor for `CancelCommandResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelCommandResultDescriptor =
    $convert.base64Decode('ChNDYW5jZWxDb21tYW5kUmVzdWx0');

@$core
    .Deprecated('Use cancelMaintenanceWindowExecutionRequestDescriptor instead')
const CancelMaintenanceWindowExecutionRequest$json = {
  '1': 'CancelMaintenanceWindowExecutionRequest',
  '2': [
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `CancelMaintenanceWindowExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelMaintenanceWindowExecutionRequestDescriptor =
    $convert.base64Decode(
        'CidDYW5jZWxNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblJlcXVlc3QSMAoRd2luZG93ZXhlY3'
        'V0aW9uaWQYieunqgEgASgJUhF3aW5kb3dleGVjdXRpb25pZA==');

@$core
    .Deprecated('Use cancelMaintenanceWindowExecutionResultDescriptor instead')
const CancelMaintenanceWindowExecutionResult$json = {
  '1': 'CancelMaintenanceWindowExecutionResult',
  '2': [
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `CancelMaintenanceWindowExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelMaintenanceWindowExecutionResultDescriptor =
    $convert.base64Decode(
        'CiZDYW5jZWxNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblJlc3VsdBIwChF3aW5kb3dleGVjdX'
        'Rpb25pZBiJ66eqASABKAlSEXdpbmRvd2V4ZWN1dGlvbmlk');

@$core.Deprecated('Use cloudWatchOutputConfigDescriptor instead')
const CloudWatchOutputConfig$json = {
  '1': 'CloudWatchOutputConfig',
  '2': [
    {
      '1': 'cloudwatchloggroupname',
      '3': 2322068,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchloggroupname'
    },
    {
      '1': 'cloudwatchoutputenabled',
      '3': 30430974,
      '4': 1,
      '5': 8,
      '10': 'cloudwatchoutputenabled'
    },
  ],
};

/// Descriptor for `CloudWatchOutputConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchOutputConfigDescriptor = $convert.base64Decode(
    'ChZDbG91ZFdhdGNoT3V0cHV0Q29uZmlnEjkKFmNsb3Vkd2F0Y2hsb2dncm91cG5hbWUYlN2NAS'
    'ABKAlSFmNsb3Vkd2F0Y2hsb2dncm91cG5hbWUSOwoXY2xvdWR3YXRjaG91dHB1dGVuYWJsZWQY'
    '/q3BDiABKAhSF2Nsb3Vkd2F0Y2hvdXRwdXRlbmFibGVk');

@$core.Deprecated('Use commandDescriptor instead')
const Command$json = {
  '1': 'Command',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cloudwatchoutputconfig',
      '3': 21186555,
      '4': 1,
      '5': 11,
      '6': '.ssm.CloudWatchOutputConfig',
      '10': 'cloudwatchoutputconfig'
    },
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'completedcount',
      '3': 11380122,
      '4': 1,
      '5': 5,
      '10': 'completedcount'
    },
    {
      '1': 'deliverytimedoutcount',
      '3': 176157036,
      '4': 1,
      '5': 5,
      '10': 'deliverytimedoutcount'
    },
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'errorcount', '3': 311137001, '4': 1, '5': 5, '10': 'errorcount'},
    {'1': 'expiresafter', '3': 350882700, '4': 1, '5': 9, '10': 'expiresafter'},
    {'1': 'instanceids', '3': 312792453, '4': 3, '5': 9, '10': 'instanceids'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'notificationconfig',
      '3': 346074145,
      '4': 1,
      '5': 11,
      '6': '.ssm.NotificationConfig',
      '10': 'notificationconfig'
    },
    {
      '1': 'outputs3bucketname',
      '3': 186756480,
      '4': 1,
      '5': 9,
      '10': 'outputs3bucketname'
    },
    {
      '1': 'outputs3keyprefix',
      '3': 17840974,
      '4': 1,
      '5': 9,
      '10': 'outputs3keyprefix'
    },
    {
      '1': 'outputs3region',
      '3': 398718159,
      '4': 1,
      '5': 9,
      '10': 'outputs3region'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.Command.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'requesteddatetime',
      '3': 111054087,
      '4': 1,
      '5': 9,
      '10': 'requesteddatetime'
    },
    {'1': 'servicerole', '3': 47807725, '4': 1, '5': 9, '10': 'servicerole'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.CommandStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {'1': 'targetcount', '3': 525114440, '4': 1, '5': 5, '10': 'targetcount'},
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'timeoutseconds',
      '3': 336148022,
      '4': 1,
      '5': 5,
      '10': 'timeoutseconds'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
  ],
  '3': [Command_ParametersEntry$json],
};

@$core.Deprecated('Use commandDescriptor instead')
const Command_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Command`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commandDescriptor = $convert.base64Decode(
    'CgdDb21tYW5kEkoKEmFsYXJtY29uZmlndXJhdGlvbhiJmbkhIAEoCzIXLnNzbS5BbGFybUNvbm'
    'ZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlvbhJWChZjbG91ZHdhdGNob3V0cHV0Y29uZmln'
    'GPuPjQogASgLMhsuc3NtLkNsb3VkV2F0Y2hPdXRwdXRDb25maWdSFmNsb3Vkd2F0Y2hvdXRwdX'
    'Rjb25maWcSHwoJY29tbWFuZGlkGIDbgEwgASgJUgljb21tYW5kaWQSHAoHY29tbWVudBj/v77C'
    'ASABKAlSB2NvbW1lbnQSKQoOY29tcGxldGVkY291bnQYmsu2BSABKAVSDmNvbXBsZXRlZGNvdW'
    '50EjcKFWRlbGl2ZXJ5dGltZWRvdXRjb3VudBjs4v9TIAEoBVIVZGVsaXZlcnl0aW1lZG91dGNv'
    'dW50EiUKDGRvY3VtZW50bmFtZRjQo8c5IAEoCVIMZG9jdW1lbnRuYW1lEisKD2RvY3VtZW50dm'
    'Vyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEiIKCmVycm9yY291bnQY6aWulAEgASgF'
    'UgplcnJvcmNvdW50EiYKDGV4cGlyZXNhZnRlchiMl6inASABKAlSDGV4cGlyZXNhZnRlchIkCg'
    'tpbnN0YW5jZWlkcxiFq5OVASADKAlSC2luc3RhbmNlaWRzEikKDm1heGNvbmN1cnJlbmN5GP3B'
    'jg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcnMYq8L1PSABKAlSCW1heGVycm9ycx'
    'JLChJub3RpZmljYXRpb25jb25maWcYodiCpQEgASgLMhcuc3NtLk5vdGlmaWNhdGlvbkNvbmZp'
    'Z1ISbm90aWZpY2F0aW9uY29uZmlnEjEKEm91dHB1dHMzYnVja2V0bmFtZRiA24ZZIAEoCVISb3'
    'V0cHV0czNidWNrZXRuYW1lEi8KEW91dHB1dHMza2V5cHJlZml4GM72wAggASgJUhFvdXRwdXRz'
    'M2tleXByZWZpeBIqCg5vdXRwdXRzM3JlZ2lvbhjP6Y++ASABKAlSDm91dHB1dHMzcmVnaW9uEk'
    'AKCnBhcmFtZXRlcnMY+qf+6wEgAygLMhwuc3NtLkNvbW1hbmQuUGFyYW1ldGVyc0VudHJ5Ugpw'
    'YXJhbWV0ZXJzEi8KEXJlcXVlc3RlZGRhdGV0aW1lGIea+jQgASgJUhFyZXF1ZXN0ZWRkYXRldG'
    'ltZRIjCgtzZXJ2aWNlcm9sZRjt+eUWIAEoCVILc2VydmljZXJvbGUSLQoGc3RhdHVzGJDk+wIg'
    'ASgOMhIuc3NtLkNvbW1hbmRTdGF0dXNSBnN0YXR1cxIoCg1zdGF0dXNkZXRhaWxzGKiSwbEBIA'
    'EoCVINc3RhdHVzZGV0YWlscxIkCgt0YXJnZXRjb3VudBjIuLL6ASABKAVSC3RhcmdldGNvdW50'
    'EigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEioKDnRpbWVvdXRzZW'
    'NvbmRzGLbspKABIAEoBVIOdGltZW91dHNlY29uZHMSRwoPdHJpZ2dlcmVkYWxhcm1zGIXtwX0g'
    'AygLMhouc3NtLkFsYXJtU3RhdGVJbmZvcm1hdGlvblIPdHJpZ2dlcmVkYWxhcm1zGj0KD1Bhcm'
    'FtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use commandFilterDescriptor instead')
const CommandFilter$json = {
  '1': 'CommandFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.CommandFilterKey',
      '10': 'key'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `CommandFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commandFilterDescriptor = $convert.base64Decode(
    'Cg1Db21tYW5kRmlsdGVyEioKA2tleRjtkNdAIAEoDjIVLnNzbS5Db21tYW5kRmlsdGVyS2V5Ug'
    'NrZXkSFwoFdmFsdWUYy6f7EiABKAlSBXZhbHVl');

@$core.Deprecated('Use commandInvocationDescriptor instead')
const CommandInvocation$json = {
  '1': 'CommandInvocation',
  '2': [
    {
      '1': 'cloudwatchoutputconfig',
      '3': 21186555,
      '4': 1,
      '5': 11,
      '6': '.ssm.CloudWatchOutputConfig',
      '10': 'cloudwatchoutputconfig'
    },
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {
      '1': 'commandplugins',
      '3': 508176155,
      '4': 3,
      '5': 11,
      '6': '.ssm.CommandPlugin',
      '10': 'commandplugins'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'instancename', '3': 281444442, '4': 1, '5': 9, '10': 'instancename'},
    {
      '1': 'notificationconfig',
      '3': 346074145,
      '4': 1,
      '5': 11,
      '6': '.ssm.NotificationConfig',
      '10': 'notificationconfig'
    },
    {
      '1': 'requesteddatetime',
      '3': 111054087,
      '4': 1,
      '5': 9,
      '10': 'requesteddatetime'
    },
    {'1': 'servicerole', '3': 47807725, '4': 1, '5': 9, '10': 'servicerole'},
    {
      '1': 'standarderrorurl',
      '3': 403407680,
      '4': 1,
      '5': 9,
      '10': 'standarderrorurl'
    },
    {
      '1': 'standardoutputurl',
      '3': 347642271,
      '4': 1,
      '5': 9,
      '10': 'standardoutputurl'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.CommandInvocationStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {'1': 'traceoutput', '3': 519666568, '4': 1, '5': 9, '10': 'traceoutput'},
  ],
};

/// Descriptor for `CommandInvocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commandInvocationDescriptor = $convert.base64Decode(
    'ChFDb21tYW5kSW52b2NhdGlvbhJWChZjbG91ZHdhdGNob3V0cHV0Y29uZmlnGPuPjQogASgLMh'
    'suc3NtLkNsb3VkV2F0Y2hPdXRwdXRDb25maWdSFmNsb3Vkd2F0Y2hvdXRwdXRjb25maWcSHwoJ'
    'Y29tbWFuZGlkGIDbgEwgASgJUgljb21tYW5kaWQSPgoOY29tbWFuZHBsdWdpbnMYm86o8gEgAy'
    'gLMhIuc3NtLkNvbW1hbmRQbHVnaW5SDmNvbW1hbmRwbHVnaW5zEhwKB2NvbW1lbnQY/7++wgEg'
    'ASgJUgdjb21tZW50EiUKDGRvY3VtZW50bmFtZRjQo8c5IAEoCVIMZG9jdW1lbnRuYW1lEisKD2'
    'RvY3VtZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEiEKCmluc3RhbmNlaWQY'
    'oK3RFyABKAlSCmluc3RhbmNlaWQSJgoMaW5zdGFuY2VuYW1lGNqAmoYBIAEoCVIMaW5zdGFuY2'
    'VuYW1lEksKEm5vdGlmaWNhdGlvbmNvbmZpZxih2IKlASABKAsyFy5zc20uTm90aWZpY2F0aW9u'
    'Q29uZmlnUhJub3RpZmljYXRpb25jb25maWcSLwoRcmVxdWVzdGVkZGF0ZXRpbWUYh5r6NCABKA'
    'lSEXJlcXVlc3RlZGRhdGV0aW1lEiMKC3NlcnZpY2Vyb2xlGO355RYgASgJUgtzZXJ2aWNlcm9s'
    'ZRIuChBzdGFuZGFyZGVycm9ydXJsGMCGrsABIAEoCVIQc3RhbmRhcmRlcnJvcnVybBIwChFzdG'
    'FuZGFyZG91dHB1dHVybBifs+KlASABKAlSEXN0YW5kYXJkb3V0cHV0dXJsEjcKBnN0YXR1cxiQ'
    '5PsCIAEoDjIcLnNzbS5Db21tYW5kSW52b2NhdGlvblN0YXR1c1IGc3RhdHVzEigKDXN0YXR1c2'
    'RldGFpbHMYqJLBsQEgASgJUg1zdGF0dXNkZXRhaWxzEiQKC3RyYWNlb3V0cHV0GIj35fcBIAEo'
    'CVILdHJhY2VvdXRwdXQ=');

@$core.Deprecated('Use commandPluginDescriptor instead')
const CommandPlugin$json = {
  '1': 'CommandPlugin',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'output', '3': 242631461, '4': 1, '5': 9, '10': 'output'},
    {
      '1': 'outputs3bucketname',
      '3': 186756480,
      '4': 1,
      '5': 9,
      '10': 'outputs3bucketname'
    },
    {
      '1': 'outputs3keyprefix',
      '3': 17840974,
      '4': 1,
      '5': 9,
      '10': 'outputs3keyprefix'
    },
    {
      '1': 'outputs3region',
      '3': 398718159,
      '4': 1,
      '5': 9,
      '10': 'outputs3region'
    },
    {'1': 'responsecode', '3': 447553700, '4': 1, '5': 5, '10': 'responsecode'},
    {
      '1': 'responsefinishdatetime',
      '3': 259049989,
      '4': 1,
      '5': 9,
      '10': 'responsefinishdatetime'
    },
    {
      '1': 'responsestartdatetime',
      '3': 382241276,
      '4': 1,
      '5': 9,
      '10': 'responsestartdatetime'
    },
    {
      '1': 'standarderrorurl',
      '3': 403407680,
      '4': 1,
      '5': 9,
      '10': 'standarderrorurl'
    },
    {
      '1': 'standardoutputurl',
      '3': 347642271,
      '4': 1,
      '5': 9,
      '10': 'standardoutputurl'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.CommandPluginStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
  ],
};

/// Descriptor for `CommandPlugin`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commandPluginDescriptor = $convert.base64Decode(
    'Cg1Db21tYW5kUGx1Z2luEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGQoGb3V0cHV0GKWG2XMgAS'
    'gJUgZvdXRwdXQSMQoSb3V0cHV0czNidWNrZXRuYW1lGIDbhlkgASgJUhJvdXRwdXRzM2J1Y2tl'
    'dG5hbWUSLwoRb3V0cHV0czNrZXlwcmVmaXgYzvbACCABKAlSEW91dHB1dHMza2V5cHJlZml4Ei'
    'oKDm91dHB1dHMzcmVnaW9uGM/pj74BIAEoCVIOb3V0cHV0czNyZWdpb24SJgoMcmVzcG9uc2Vj'
    'b2RlGKTBtNUBIAEoBVIMcmVzcG9uc2Vjb2RlEjkKFnJlc3BvbnNlZmluaXNoZGF0ZXRpbWUYhZ'
    'TDeyABKAlSFnJlc3BvbnNlZmluaXNoZGF0ZXRpbWUSOAoVcmVzcG9uc2VzdGFydGRhdGV0aW1l'
    'GPyTorYBIAEoCVIVcmVzcG9uc2VzdGFydGRhdGV0aW1lEi4KEHN0YW5kYXJkZXJyb3J1cmwYwI'
    'auwAEgASgJUhBzdGFuZGFyZGVycm9ydXJsEjAKEXN0YW5kYXJkb3V0cHV0dXJsGJ+z4qUBIAEo'
    'CVIRc3RhbmRhcmRvdXRwdXR1cmwSMwoGc3RhdHVzGJDk+wIgASgOMhguc3NtLkNvbW1hbmRQbH'
    'VnaW5TdGF0dXNSBnN0YXR1cxIoCg1zdGF0dXNkZXRhaWxzGKiSwbEBIAEoCVINc3RhdHVzZGV0'
    'YWlscw==');

@$core.Deprecated('Use complianceExecutionSummaryDescriptor instead')
const ComplianceExecutionSummary$json = {
  '1': 'ComplianceExecutionSummary',
  '2': [
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {
      '1': 'executiontime',
      '3': 379716053,
      '4': 1,
      '5': 9,
      '10': 'executiontime'
    },
    {
      '1': 'executiontype',
      '3': 224869374,
      '4': 1,
      '5': 9,
      '10': 'executiontype'
    },
  ],
};

/// Descriptor for `ComplianceExecutionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complianceExecutionSummaryDescriptor =
    $convert.base64Decode(
        'ChpDb21wbGlhbmNlRXhlY3V0aW9uU3VtbWFyeRIjCgtleGVjdXRpb25pZBixz69GIAEoCVILZX'
        'hlY3V0aW9uaWQSKAoNZXhlY3V0aW9udGltZRjVg4i1ASABKAlSDWV4ZWN1dGlvbnRpbWUSJwoN'
        'ZXhlY3V0aW9udHlwZRj+95xrIAEoCVINZXhlY3V0aW9udHlwZQ==');

@$core.Deprecated('Use complianceItemDescriptor instead')
const ComplianceItem$json = {
  '1': 'ComplianceItem',
  '2': [
    {
      '1': 'compliancetype',
      '3': 451385291,
      '4': 1,
      '5': 9,
      '10': 'compliancetype'
    },
    {
      '1': 'details',
      '3': 247611974,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceItem.DetailsEntry',
      '10': 'details'
    },
    {
      '1': 'executionsummary',
      '3': 71480964,
      '4': 1,
      '5': 11,
      '6': '.ssm.ComplianceExecutionSummary',
      '10': 'executionsummary'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {
      '1': 'severity',
      '3': 276886227,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceSeverity',
      '10': 'severity'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
  '3': [ComplianceItem_DetailsEntry$json],
};

@$core.Deprecated('Use complianceItemDescriptor instead')
const ComplianceItem_DetailsEntry$json = {
  '1': 'DetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ComplianceItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complianceItemDescriptor = $convert.base64Decode(
    'Cg5Db21wbGlhbmNlSXRlbRIqCg5jb21wbGlhbmNldHlwZRjLr57XASABKAlSDmNvbXBsaWFuY2'
    'V0eXBlEj0KB2RldGFpbHMYxoSJdiADKAsyIC5zc20uQ29tcGxpYW5jZUl0ZW0uRGV0YWlsc0Vu'
    'dHJ5UgdkZXRhaWxzEk4KEGV4ZWN1dGlvbnN1bW1hcnkYhO2KIiABKAsyHy5zc20uQ29tcGxpYW'
    '5jZUV4ZWN1dGlvblN1bW1hcnlSEGV4ZWN1dGlvbnN1bW1hcnkSEgoCaWQYgfKitwEgASgJUgJp'
    'ZBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBImCgxyZXNvdXJjZXR5cGUY3r'
    '7YjwEgASgJUgxyZXNvdXJjZXR5cGUSNwoIc2V2ZXJpdHkY0+WDhAEgASgOMhcuc3NtLkNvbXBs'
    'aWFuY2VTZXZlcml0eVIIc2V2ZXJpdHkSMAoGc3RhdHVzGJDk+wIgASgOMhUuc3NtLkNvbXBsaW'
    'FuY2VTdGF0dXNSBnN0YXR1cxIXCgV0aXRsZRiq49EmIAEoCVIFdGl0bGUaOgoMRGV0YWlsc0Vu'
    'dHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use complianceItemEntryDescriptor instead')
const ComplianceItemEntry$json = {
  '1': 'ComplianceItemEntry',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceItemEntry.DetailsEntry',
      '10': 'details'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'severity',
      '3': 276886227,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceSeverity',
      '10': 'severity'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
  '3': [ComplianceItemEntry_DetailsEntry$json],
};

@$core.Deprecated('Use complianceItemEntryDescriptor instead')
const ComplianceItemEntry_DetailsEntry$json = {
  '1': 'DetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ComplianceItemEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complianceItemEntryDescriptor = $convert.base64Decode(
    'ChNDb21wbGlhbmNlSXRlbUVudHJ5EkIKB2RldGFpbHMYxoSJdiADKAsyJS5zc20uQ29tcGxpYW'
    '5jZUl0ZW1FbnRyeS5EZXRhaWxzRW50cnlSB2RldGFpbHMSEgoCaWQYgfKitwEgASgJUgJpZBI3'
    'CghzZXZlcml0eRjT5YOEASABKA4yFy5zc20uQ29tcGxpYW5jZVNldmVyaXR5UghzZXZlcml0eR'
    'IwCgZzdGF0dXMYkOT7AiABKA4yFS5zc20uQ29tcGxpYW5jZVN0YXR1c1IGc3RhdHVzEhcKBXRp'
    'dGxlGKrj0SYgASgJUgV0aXRsZRo6CgxEZXRhaWxzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFA'
    'oFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use complianceStringFilterDescriptor instead')
const ComplianceStringFilter$json = {
  '1': 'ComplianceStringFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceQueryOperatorType',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `ComplianceStringFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complianceStringFilterDescriptor = $convert.base64Decode(
    'ChZDb21wbGlhbmNlU3RyaW5nRmlsdGVyEhMKA2tleRiNkutoIAEoCVIDa2V5EjgKBHR5cGUY7q'
    'DXigEgASgOMiAuc3NtLkNvbXBsaWFuY2VRdWVyeU9wZXJhdG9yVHlwZVIEdHlwZRIZCgZ2YWx1'
    'ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use complianceSummaryItemDescriptor instead')
const ComplianceSummaryItem$json = {
  '1': 'ComplianceSummaryItem',
  '2': [
    {
      '1': 'compliancetype',
      '3': 451385291,
      '4': 1,
      '5': 9,
      '10': 'compliancetype'
    },
    {
      '1': 'compliantsummary',
      '3': 133218055,
      '4': 1,
      '5': 11,
      '6': '.ssm.CompliantSummary',
      '10': 'compliantsummary'
    },
    {
      '1': 'noncompliantsummary',
      '3': 294594444,
      '4': 1,
      '5': 11,
      '6': '.ssm.NonCompliantSummary',
      '10': 'noncompliantsummary'
    },
  ],
};

/// Descriptor for `ComplianceSummaryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complianceSummaryItemDescriptor = $convert.base64Decode(
    'ChVDb21wbGlhbmNlU3VtbWFyeUl0ZW0SKgoOY29tcGxpYW5jZXR5cGUYy6+e1wEgASgJUg5jb2'
    '1wbGlhbmNldHlwZRJEChBjb21wbGlhbnRzdW1tYXJ5GIf+wj8gASgLMhUuc3NtLkNvbXBsaWFu'
    'dFN1bW1hcnlSEGNvbXBsaWFudHN1bW1hcnkSTgoTbm9uY29tcGxpYW50c3VtbWFyeRiMz7yMAS'
    'ABKAsyGC5zc20uTm9uQ29tcGxpYW50U3VtbWFyeVITbm9uY29tcGxpYW50c3VtbWFyeQ==');

@$core.Deprecated(
    'Use complianceTypeCountLimitExceededExceptionDescriptor instead')
const ComplianceTypeCountLimitExceededException$json = {
  '1': 'ComplianceTypeCountLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ComplianceTypeCountLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    complianceTypeCountLimitExceededExceptionDescriptor = $convert.base64Decode(
        'CilDb21wbGlhbmNlVHlwZUNvdW50TGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use compliantSummaryDescriptor instead')
const CompliantSummary$json = {
  '1': 'CompliantSummary',
  '2': [
    {
      '1': 'compliantcount',
      '3': 480617586,
      '4': 1,
      '5': 5,
      '10': 'compliantcount'
    },
    {
      '1': 'severitysummary',
      '3': 20421279,
      '4': 1,
      '5': 11,
      '6': '.ssm.SeveritySummary',
      '10': 'severitysummary'
    },
  ],
};

/// Descriptor for `CompliantSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List compliantSummaryDescriptor = $convert.base64Decode(
    'ChBDb21wbGlhbnRTdW1tYXJ5EioKDmNvbXBsaWFudGNvdW50GPLIluUBIAEoBVIOY29tcGxpYW'
    '50Y291bnQSQQoPc2V2ZXJpdHlzdW1tYXJ5GJ+13gkgASgLMhQuc3NtLlNldmVyaXR5U3VtbWFy'
    'eVIPc2V2ZXJpdHlzdW1tYXJ5');

@$core.Deprecated('Use createActivationRequestDescriptor instead')
const CreateActivationRequest$json = {
  '1': 'CreateActivationRequest',
  '2': [
    {
      '1': 'defaultinstancename',
      '3': 27070947,
      '4': 1,
      '5': 9,
      '10': 'defaultinstancename'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'expirationdate',
      '3': 17577725,
      '4': 1,
      '5': 9,
      '10': 'expirationdate'
    },
    {'1': 'iamrole', '3': 242424351, '4': 1, '5': 9, '10': 'iamrole'},
    {
      '1': 'registrationlimit',
      '3': 312393964,
      '4': 1,
      '5': 5,
      '10': 'registrationlimit'
    },
    {
      '1': 'registrationmetadata',
      '3': 360246862,
      '4': 3,
      '5': 11,
      '6': '.ssm.RegistrationMetadataItem',
      '10': 'registrationmetadata'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateActivationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createActivationRequestDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVBY3RpdmF0aW9uUmVxdWVzdBIzChNkZWZhdWx0aW5zdGFuY2VuYW1lGOOj9AwgAS'
    'gJUhNkZWZhdWx0aW5zdGFuY2VuYW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlw'
    'dGlvbhIpCg5leHBpcmF0aW9uZGF0ZRj97bAIIAEoCVIOZXhwaXJhdGlvbmRhdGUSGwoHaWFtcm'
    '9sZRiftMxzIAEoCVIHaWFtcm9sZRIwChFyZWdpc3RyYXRpb25saW1pdBjsgfuUASABKAVSEXJl'
    'Z2lzdHJhdGlvbmxpbWl0ElUKFHJlZ2lzdHJhdGlvbm1ldGFkYXRhGM7c46sBIAMoCzIdLnNzbS'
    '5SZWdpc3RyYXRpb25NZXRhZGF0YUl0ZW1SFHJlZ2lzdHJhdGlvbm1ldGFkYXRhEiAKBHRhZ3MY'
    'wcH2tQEgAygLMgguc3NtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createActivationResultDescriptor instead')
const CreateActivationResult$json = {
  '1': 'CreateActivationResult',
  '2': [
    {
      '1': 'activationcode',
      '3': 290753501,
      '4': 1,
      '5': 9,
      '10': 'activationcode'
    },
    {'1': 'activationid', '3': 146858601, '4': 1, '5': 9, '10': 'activationid'},
  ],
};

/// Descriptor for `CreateActivationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createActivationResultDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVBY3RpdmF0aW9uUmVzdWx0EioKDmFjdGl2YXRpb25jb2RlGN2X0ooBIAEoCVIOYW'
    'N0aXZhdGlvbmNvZGUSJQoMYWN0aXZhdGlvbmlkGOnEg0YgASgJUgxhY3RpdmF0aW9uaWQ=');

@$core.Deprecated('Use createAssociationBatchRequestDescriptor instead')
const CreateAssociationBatchRequest$json = {
  '1': 'CreateAssociationBatchRequest',
  '2': [
    {
      '1': 'associationdispatchassumerole',
      '3': 281627425,
      '4': 1,
      '5': 9,
      '10': 'associationdispatchassumerole'
    },
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateAssociationBatchRequestEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `CreateAssociationBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAssociationBatchRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVBc3NvY2lhdGlvbkJhdGNoUmVxdWVzdBJICh1hc3NvY2lhdGlvbmRpc3BhdGNoYX'
    'NzdW1lcm9sZRihlqWGASABKAlSHWFzc29jaWF0aW9uZGlzcGF0Y2hhc3N1bWVyb2xlEkUKB2Vu'
    'dHJpZXMYlMWy5QEgAygLMicuc3NtLkNyZWF0ZUFzc29jaWF0aW9uQmF0Y2hSZXF1ZXN0RW50cn'
    'lSB2VudHJpZXM=');

@$core.Deprecated('Use createAssociationBatchRequestEntryDescriptor instead')
const CreateAssociationBatchRequestEntry$json = {
  '1': 'CreateAssociationBatchRequestEntry',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'applyonlyatcroninterval',
      '3': 285867158,
      '4': 1,
      '5': 8,
      '10': 'applyonlyatcroninterval'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'automationtargetparametername',
      '3': 348584826,
      '4': 1,
      '5': 9,
      '10': 'automationtargetparametername'
    },
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
    {
      '1': 'complianceseverity',
      '3': 278891158,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationComplianceSeverity',
      '10': 'complianceseverity'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateAssociationBatchRequestEntry.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'synccompliance',
      '3': 500469318,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationSyncCompliance',
      '10': 'synccompliance'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateAssociationBatchRequestEntry.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    CreateAssociationBatchRequestEntry_ParametersEntry$json,
    CreateAssociationBatchRequestEntry_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use createAssociationBatchRequestEntryDescriptor instead')
const CreateAssociationBatchRequestEntry_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createAssociationBatchRequestEntryDescriptor instead')
const CreateAssociationBatchRequestEntry_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateAssociationBatchRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAssociationBatchRequestEntryDescriptor = $convert.base64Decode(
    'CiJDcmVhdGVBc3NvY2lhdGlvbkJhdGNoUmVxdWVzdEVudHJ5EkoKEmFsYXJtY29uZmlndXJhdG'
    'lvbhiJmbkhIAEoCzIXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlv'
    'bhI8ChdhcHBseW9ubHlhdGNyb25pbnRlcnZhbBiW+aeIASABKAhSF2FwcGx5b25seWF0Y3Jvbm'
    'ludGVydmFsEiwKD2Fzc29jaWF0aW9ubmFtZRiYkMWVASABKAlSD2Fzc29jaWF0aW9ubmFtZRJI'
    'Ch1hdXRvbWF0aW9udGFyZ2V0cGFyYW1ldGVybmFtZRj69pumASABKAlSHWF1dG9tYXRpb250YX'
    'JnZXRwYXJhbWV0ZXJuYW1lEicKDWNhbGVuZGFybmFtZXMYvvOZESADKAlSDWNhbGVuZGFybmFt'
    'ZXMSVgoSY29tcGxpYW5jZXNldmVyaXR5GJaV/oQBIAEoDjIiLnNzbS5Bc3NvY2lhdGlvbkNvbX'
    'BsaWFuY2VTZXZlcml0eVISY29tcGxpYW5jZXNldmVyaXR5EisKD2RvY3VtZW50dmVyc2lvbhjJ'
    '76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEh4KCGR1cmF0aW9uGK6SnaYBIAEoBVIIZHVyYXRpb2'
    '4SIQoKaW5zdGFuY2VpZBigrdEXIAEoCVIKaW5zdGFuY2VpZBIpCg5tYXhjb25jdXJyZW5jeRj9'
    'wY4OIAEoCVIObWF4Y29uY3VycmVuY3kSHwoJbWF4ZXJyb3JzGKvC9T0gASgJUgltYXhlcnJvcn'
    'MSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRJRCg5vdXRwdXRsb2NhdGlvbhj067UgIAEoCzImLnNz'
    'bS5JbnN0YW5jZUFzc29jaWF0aW9uT3V0cHV0TG9jYXRpb25SDm91dHB1dGxvY2F0aW9uElsKCn'
    'BhcmFtZXRlcnMY+qf+6wEgAygLMjcuc3NtLkNyZWF0ZUFzc29jaWF0aW9uQmF0Y2hSZXF1ZXN0'
    'RW50cnkuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEjIKEnNjaGVkdWxlZXhwcmVzc2lvbh'
    'j/kdvUASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhIpCg5zY2hlZHVsZW9mZnNldBjYkepKIAEo'
    'BVIOc2NoZWR1bGVvZmZzZXQSSgoOc3luY2NvbXBsaWFuY2UYxpzS7gEgASgOMh4uc3NtLkFzc2'
    '9jaWF0aW9uU3luY0NvbXBsaWFuY2VSDnN5bmNjb21wbGlhbmNlEkEKD3RhcmdldGxvY2F0aW9u'
    'cxilu/GJASADKAsyEy5zc20uVGFyZ2V0TG9jYXRpb25SD3RhcmdldGxvY2F0aW9ucxJaCgp0YX'
    'JnZXRtYXBzGLi81SMgAygLMjcuc3NtLkNyZWF0ZUFzc29jaWF0aW9uQmF0Y2hSZXF1ZXN0RW50'
    'cnkuVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc2'
    '0uVGFyZ2V0Ugd0YXJnZXRzGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIU'
    'CgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGj0KD1RhcmdldG1hcHNFbnRyeRIQCgNrZXkYASABKA'
    'lSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createAssociationBatchResultDescriptor instead')
const CreateAssociationBatchResult$json = {
  '1': 'CreateAssociationBatchResult',
  '2': [
    {
      '1': 'failed',
      '3': 360301525,
      '4': 3,
      '5': 11,
      '6': '.ssm.FailedCreateAssociation',
      '10': 'failed'
    },
    {
      '1': 'successful',
      '3': 412818844,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationDescription',
      '10': 'successful'
    },
  ],
};

/// Descriptor for `CreateAssociationBatchResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAssociationBatchResultDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVBc3NvY2lhdGlvbkJhdGNoUmVzdWx0EjgKBmZhaWxlZBjVh+erASADKAsyHC5zc2'
        '0uRmFpbGVkQ3JlYXRlQXNzb2NpYXRpb25SBmZhaWxlZBI/CgpzdWNjZXNzZnVsGJy77MQBIAMo'
        'CzIbLnNzbS5Bc3NvY2lhdGlvbkRlc2NyaXB0aW9uUgpzdWNjZXNzZnVs');

@$core.Deprecated('Use createAssociationRequestDescriptor instead')
const CreateAssociationRequest$json = {
  '1': 'CreateAssociationRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'applyonlyatcroninterval',
      '3': 285867158,
      '4': 1,
      '5': 8,
      '10': 'applyonlyatcroninterval'
    },
    {
      '1': 'associationdispatchassumerole',
      '3': 281627425,
      '4': 1,
      '5': 9,
      '10': 'associationdispatchassumerole'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'automationtargetparametername',
      '3': 348584826,
      '4': 1,
      '5': 9,
      '10': 'automationtargetparametername'
    },
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
    {
      '1': 'complianceseverity',
      '3': 278891158,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationComplianceSeverity',
      '10': 'complianceseverity'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateAssociationRequest.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'synccompliance',
      '3': 500469318,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationSyncCompliance',
      '10': 'synccompliance'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateAssociationRequest.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    CreateAssociationRequest_ParametersEntry$json,
    CreateAssociationRequest_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use createAssociationRequestDescriptor instead')
const CreateAssociationRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createAssociationRequestDescriptor instead')
const CreateAssociationRequest_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAssociationRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVBc3NvY2lhdGlvblJlcXVlc3QSSgoSYWxhcm1jb25maWd1cmF0aW9uGImZuSEgAS'
    'gLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uEjwKF2FwcGx5'
    'b25seWF0Y3JvbmludGVydmFsGJb5p4gBIAEoCFIXYXBwbHlvbmx5YXRjcm9uaW50ZXJ2YWwSSA'
    'odYXNzb2NpYXRpb25kaXNwYXRjaGFzc3VtZXJvbGUYoZalhgEgASgJUh1hc3NvY2lhdGlvbmRp'
    'c3BhdGNoYXNzdW1lcm9sZRIsCg9hc3NvY2lhdGlvbm5hbWUYmJDFlQEgASgJUg9hc3NvY2lhdG'
    'lvbm5hbWUSSAodYXV0b21hdGlvbnRhcmdldHBhcmFtZXRlcm5hbWUY+vabpgEgASgJUh1hdXRv'
    'bWF0aW9udGFyZ2V0cGFyYW1ldGVybmFtZRInCg1jYWxlbmRhcm5hbWVzGL7zmREgAygJUg1jYW'
    'xlbmRhcm5hbWVzElYKEmNvbXBsaWFuY2VzZXZlcml0eRiWlf6EASABKA4yIi5zc20uQXNzb2Np'
    'YXRpb25Db21wbGlhbmNlU2V2ZXJpdHlSEmNvbXBsaWFuY2VzZXZlcml0eRIrCg9kb2N1bWVudH'
    'ZlcnNpb24Yye+pKCABKAlSD2RvY3VtZW50dmVyc2lvbhIeCghkdXJhdGlvbhiukp2mASABKAVS'
    'CGR1cmF0aW9uEiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQSKQoObWF4Y29uY3'
    'VycmVuY3kY/cGODiABKAlSDm1heGNvbmN1cnJlbmN5Eh8KCW1heGVycm9ycxirwvU9IAEoCVIJ'
    'bWF4ZXJyb3JzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSUQoOb3V0cHV0bG9jYXRpb24Y9Ou1IC'
    'ABKAsyJi5zc20uSW5zdGFuY2VBc3NvY2lhdGlvbk91dHB1dExvY2F0aW9uUg5vdXRwdXRsb2Nh'
    'dGlvbhJRCgpwYXJhbWV0ZXJzGPqn/usBIAMoCzItLnNzbS5DcmVhdGVBc3NvY2lhdGlvblJlcX'
    'Vlc3QuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEjIKEnNjaGVkdWxlZXhwcmVzc2lvbhj/'
    'kdvUASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhIpCg5zY2hlZHVsZW9mZnNldBjYkepKIAEoBV'
    'IOc2NoZWR1bGVvZmZzZXQSSgoOc3luY2NvbXBsaWFuY2UYxpzS7gEgASgOMh4uc3NtLkFzc29j'
    'aWF0aW9uU3luY0NvbXBsaWFuY2VSDnN5bmNjb21wbGlhbmNlEiAKBHRhZ3MYwcH2tQEgAygLMg'
    'guc3NtLlRhZ1IEdGFncxJBCg90YXJnZXRsb2NhdGlvbnMYpbvxiQEgAygLMhMuc3NtLlRhcmdl'
    'dExvY2F0aW9uUg90YXJnZXRsb2NhdGlvbnMSUAoKdGFyZ2V0bWFwcxi4vNUjIAMoCzItLnNzbS'
    '5DcmVhdGVBc3NvY2lhdGlvblJlcXVlc3QuVGFyZ2V0bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEigK'
    'B3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzGj0KD1BhcmFtZXRlcnNFbn'
    'RyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGj0KD1Rhcmdl'
    'dG1hcHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createAssociationResultDescriptor instead')
const CreateAssociationResult$json = {
  '1': 'CreateAssociationResult',
  '2': [
    {
      '1': 'associationdescription',
      '3': 344755863,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationDescription',
      '10': 'associationdescription'
    },
  ],
};

/// Descriptor for `CreateAssociationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAssociationResultDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVBc3NvY2lhdGlvblJlc3VsdBJXChZhc3NvY2lhdGlvbmRlc2NyaXB0aW9uGJedsq'
    'QBIAEoCzIbLnNzbS5Bc3NvY2lhdGlvbkRlc2NyaXB0aW9uUhZhc3NvY2lhdGlvbmRlc2NyaXB0'
    'aW9u');

@$core.Deprecated('Use createDocumentRequestDescriptor instead')
const CreateDocumentRequest$json = {
  '1': 'CreateDocumentRequest',
  '2': [
    {
      '1': 'attachments',
      '3': 498946338,
      '4': 3,
      '5': 11,
      '6': '.ssm.AttachmentsSource',
      '10': 'attachments'
    },
    {'1': 'content', '3': 23568227, '4': 1, '5': 9, '10': 'content'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documenttype',
      '3': 457084477,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentType',
      '10': 'documenttype'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'requires',
      '3': 149214838,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentRequires',
      '10': 'requires'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {'1': 'targettype', '3': 397256481, '4': 1, '5': 9, '10': 'targettype'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `CreateDocumentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDocumentRequestDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVEb2N1bWVudFJlcXVlc3QSPAoLYXR0YWNobWVudHMYoqL17QEgAygLMhYuc3NtLk'
    'F0dGFjaG1lbnRzU291cmNlUgthdHRhY2htZW50cxIbCgdjb250ZW50GOO+ngsgASgJUgdjb250'
    'ZW50EiQKC2Rpc3BsYXluYW1lGLfJsscBIAEoCVILZGlzcGxheW5hbWUSPwoOZG9jdW1lbnRmb3'
    'JtYXQYiJm/9gEgASgOMhMuc3NtLkRvY3VtZW50Rm9ybWF0Ug5kb2N1bWVudGZvcm1hdBI5Cgxk'
    'b2N1bWVudHR5cGUYvZz62QEgASgOMhEuc3NtLkRvY3VtZW50VHlwZVIMZG9jdW1lbnR0eXBlEh'
    'UKBG5hbWUYh+aBfyABKAlSBG5hbWUSNAoIcmVxdWlyZXMY9qyTRyADKAsyFS5zc20uRG9jdW1l'
    'bnRSZXF1aXJlc1IIcmVxdWlyZXMSIAoEdGFncxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdzEi'
    'IKCnRhcmdldHR5cGUYoc62vQEgASgJUgp0YXJnZXR0eXBlEiMKC3ZlcnNpb25uYW1lGNWjtGwg'
    'ASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use createDocumentResultDescriptor instead')
const CreateDocumentResult$json = {
  '1': 'CreateDocumentResult',
  '2': [
    {
      '1': 'documentdescription',
      '3': 500822655,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentDescription',
      '10': 'documentdescription'
    },
  ],
};

/// Descriptor for `CreateDocumentResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDocumentResultDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVEb2N1bWVudFJlc3VsdBJOChNkb2N1bWVudGRlc2NyaXB0aW9uGP/k5+4BIAEoCz'
    'IYLnNzbS5Eb2N1bWVudERlc2NyaXB0aW9uUhNkb2N1bWVudGRlc2NyaXB0aW9u');

@$core.Deprecated('Use createMaintenanceWindowRequestDescriptor instead')
const CreateMaintenanceWindowRequest$json = {
  '1': 'CreateMaintenanceWindowRequest',
  '2': [
    {
      '1': 'allowunassociatedtargets',
      '3': 154411300,
      '4': 1,
      '5': 8,
      '10': 'allowunassociatedtargets'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'cutoff', '3': 498433089, '4': 1, '5': 5, '10': 'cutoff'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'schedule', '3': 66697965, '4': 1, '5': 9, '10': 'schedule'},
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'scheduletimezone',
      '3': 170037696,
      '4': 1,
      '5': 9,
      '10': 'scheduletimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMaintenanceWindowRequestDescriptor = $convert.base64Decode(
    'Ch5DcmVhdGVNYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSPQoYYWxsb3d1bmFzc29jaWF0ZWR0YX'
    'JnZXRzGKTC0EkgASgIUhhhbGxvd3VuYXNzb2NpYXRlZHRhcmdldHMSIwoLY2xpZW50dG9rZW4Y'
    'zPu7QSABKAlSC2NsaWVudHRva2VuEhoKBmN1dG9mZhjB+NXtASABKAVSBmN1dG9mZhIjCgtkZX'
    'NjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SHgoIZHVyYXRpb24YrpKdpgEgASgFUghk'
    'dXJhdGlvbhIbCgdlbmRkYXRlGM+z+SQgASgJUgdlbmRkYXRlEhUKBG5hbWUYh+aBfyABKAlSBG'
    '5hbWUSHQoIc2NoZWR1bGUY7fXmHyABKAlSCHNjaGVkdWxlEikKDnNjaGVkdWxlb2Zmc2V0GNiR'
    '6kogASgFUg5zY2hlZHVsZW9mZnNldBItChBzY2hlZHVsZXRpbWV6b25lGMCjilEgASgJUhBzY2'
    'hlZHVsZXRpbWV6b25lEiAKCXN0YXJ0ZGF0ZRj8+KDUASABKAlSCXN0YXJ0ZGF0ZRIgCgR0YWdz'
    'GMHB9rUBIAMoCzIILnNzbS5UYWdSBHRhZ3M=');

@$core.Deprecated('Use createMaintenanceWindowResultDescriptor instead')
const CreateMaintenanceWindowResult$json = {
  '1': 'CreateMaintenanceWindowResult',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `CreateMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMaintenanceWindowResultDescriptor =
    $convert.base64Decode(
        'Ch1DcmVhdGVNYWludGVuYW5jZVdpbmRvd1Jlc3VsdBIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2'
        'luZG93aWQ=');

@$core.Deprecated('Use createOpsItemRequestDescriptor instead')
const CreateOpsItemRequest$json = {
  '1': 'CreateOpsItemRequest',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'actualendtime',
      '3': 452787526,
      '4': 1,
      '5': 9,
      '10': 'actualendtime'
    },
    {
      '1': 'actualstarttime',
      '3': 532853117,
      '4': 1,
      '5': 9,
      '10': 'actualstarttime'
    },
    {'1': 'category', '3': 263447954, '4': 1, '5': 9, '10': 'category'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'notifications',
      '3': 404560144,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemNotification',
      '10': 'notifications'
    },
    {
      '1': 'operationaldata',
      '3': 360357940,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateOpsItemRequest.OperationaldataEntry',
      '10': 'operationaldata'
    },
    {'1': 'opsitemtype', '3': 317873397, '4': 1, '5': 9, '10': 'opsitemtype'},
    {
      '1': 'plannedendtime',
      '3': 245727820,
      '4': 1,
      '5': 9,
      '10': 'plannedendtime'
    },
    {
      '1': 'plannedstarttime',
      '3': 478079215,
      '4': 1,
      '5': 9,
      '10': 'plannedstarttime'
    },
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'relatedopsitems',
      '3': 287082393,
      '4': 3,
      '5': 11,
      '6': '.ssm.RelatedOpsItem',
      '10': 'relatedopsitems'
    },
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
  '3': [CreateOpsItemRequest_OperationaldataEntry$json],
};

@$core.Deprecated('Use createOpsItemRequestDescriptor instead')
const CreateOpsItemRequest_OperationaldataEntry$json = {
  '1': 'OperationaldataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemDataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateOpsItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpsItemRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVPcHNJdGVtUmVxdWVzdBIfCglhY2NvdW50aWQY0sG5HyABKAlSCWFjY291bnRpZB'
    'IoCg1hY3R1YWxlbmR0aW1lGMb689cBIAEoCVINYWN0dWFsZW5kdGltZRIsCg9hY3R1YWxzdGFy'
    'dHRpbWUY/eKK/gEgASgJUg9hY3R1YWxzdGFydHRpbWUSHQoIY2F0ZWdvcnkYksvPfSABKAlSCG'
    'NhdGVnb3J5EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhJCCg1ub3RpZmlj'
    'YXRpb25zGJCy9MABIAMoCzIYLnNzbS5PcHNJdGVtTm90aWZpY2F0aW9uUg1ub3RpZmljYXRpb2'
    '5zElwKD29wZXJhdGlvbmFsZGF0YRi0wOqrASADKAsyLi5zc20uQ3JlYXRlT3BzSXRlbVJlcXVl'
    'c3QuT3BlcmF0aW9uYWxkYXRhRW50cnlSD29wZXJhdGlvbmFsZGF0YRIkCgtvcHNpdGVtdHlwZR'
    'j1ucmXASABKAlSC29wc2l0ZW10eXBlEikKDnBsYW5uZWRlbmR0aW1lGMyElnUgASgJUg5wbGFu'
    'bmVkZW5kdGltZRIuChBwbGFubmVkc3RhcnR0aW1lGO/R++MBIAEoCVIQcGxhbm5lZHN0YXJ0dG'
    'ltZRIdCghwcmlvcml0eRiqvrY0IAEoBVIIcHJpb3JpdHkSQQoPcmVsYXRlZG9wc2l0ZW1zGJmP'
    '8ogBIAMoCzITLnNzbS5SZWxhdGVkT3BzSXRlbVIPcmVsYXRlZG9wc2l0ZW1zEh4KCHNldmVyaX'
    'R5GNPlg4QBIAEoCVIIc2V2ZXJpdHkSGQoGc291cmNlGPnHig8gASgJUgZzb3VyY2USIAoEdGFn'
    'cxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdzEhcKBXRpdGxlGKrj0SYgASgJUgV0aXRsZRpZCh'
    'RPcGVyYXRpb25hbGRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIrCgV2YWx1ZRgCIAEoCzIV'
    'LnNzbS5PcHNJdGVtRGF0YVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createOpsItemResponseDescriptor instead')
const CreateOpsItemResponse$json = {
  '1': 'CreateOpsItemResponse',
  '2': [
    {'1': 'opsitemarn', '3': 222489428, '4': 1, '5': 9, '10': 'opsitemarn'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `CreateOpsItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpsItemResponseDescriptor = $convert.base64Decode(
    'ChVDcmVhdGVPcHNJdGVtUmVzcG9uc2USIQoKb3BzaXRlbWFybhjU1otqIAEoCVIKb3BzaXRlbW'
    'FybhIfCglvcHNpdGVtaWQY0tKVDCABKAlSCW9wc2l0ZW1pZA==');

@$core.Deprecated('Use createOpsMetadataRequestDescriptor instead')
const CreateOpsMetadataRequest$json = {
  '1': 'CreateOpsMetadataRequest',
  '2': [
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.ssm.CreateOpsMetadataRequest.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
  '3': [CreateOpsMetadataRequest_MetadataEntry$json],
};

@$core.Deprecated('Use createOpsMetadataRequestDescriptor instead')
const CreateOpsMetadataRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MetadataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateOpsMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpsMetadataRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVPcHNNZXRhZGF0YVJlcXVlc3QSSwoIbWV0YWRhdGEY4eKP4AEgAygLMisuc3NtLk'
    'NyZWF0ZU9wc01ldGFkYXRhUmVxdWVzdC5NZXRhZGF0YUVudHJ5UghtZXRhZGF0YRIiCgpyZXNv'
    'dXJjZWlkGJG68foBIAEoCVIKcmVzb3VyY2VpZBIgCgR0YWdzGMHB9rUBIAMoCzIILnNzbS5UYW'
    'dSBHRhZ3MaTwoNTWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIoCgV2YWx1ZRgCIAEo'
    'CzISLnNzbS5NZXRhZGF0YVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createOpsMetadataResultDescriptor instead')
const CreateOpsMetadataResult$json = {
  '1': 'CreateOpsMetadataResult',
  '2': [
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
  ],
};

/// Descriptor for `CreateOpsMetadataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOpsMetadataResultDescriptor =
    $convert.base64Decode(
        'ChdDcmVhdGVPcHNNZXRhZGF0YVJlc3VsdBIqCg5vcHNtZXRhZGF0YWFybhiivoLmASABKAlSDm'
        '9wc21ldGFkYXRhYXJu');

@$core.Deprecated('Use createPatchBaselineRequestDescriptor instead')
const CreatePatchBaselineRequest$json = {
  '1': 'CreatePatchBaselineRequest',
  '2': [
    {
      '1': 'approvalrules',
      '3': 71466346,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchRuleGroup',
      '10': 'approvalrules'
    },
    {
      '1': 'approvedpatches',
      '3': 199384709,
      '4': 3,
      '5': 9,
      '10': 'approvedpatches'
    },
    {
      '1': 'approvedpatchescompliancelevel',
      '3': 63924432,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'approvedpatchescompliancelevel'
    },
    {
      '1': 'approvedpatchesenablenonsecurity',
      '3': 295555901,
      '4': 1,
      '5': 8,
      '10': 'approvedpatchesenablenonsecurity'
    },
    {
      '1': 'availablesecurityupdatescompliancestatus',
      '3': 187471858,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceStatus',
      '10': 'availablesecurityupdatescompliancestatus'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'globalfilters',
      '3': 263302754,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'globalfilters'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {
      '1': 'rejectedpatches',
      '3': 309657116,
      '4': 3,
      '5': 9,
      '10': 'rejectedpatches'
    },
    {
      '1': 'rejectedpatchesaction',
      '3': 356538330,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchAction',
      '10': 'rejectedpatchesaction'
    },
    {
      '1': 'sources',
      '3': 46625746,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchSource',
      '10': 'sources'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreatePatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPatchBaselineRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVQYXRjaEJhc2VsaW5lUmVxdWVzdBI8Cg1hcHByb3ZhbHJ1bGVzGOr6iSIgASgLMh'
    'Muc3NtLlBhdGNoUnVsZUdyb3VwUg1hcHByb3ZhbHJ1bGVzEisKD2FwcHJvdmVkcGF0Y2hlcxiF'
    'vYlfIAMoCVIPYXBwcm92ZWRwYXRjaGVzEmQKHmFwcHJvdmVkcGF0Y2hlc2NvbXBsaWFuY2VsZX'
    'ZlbBjQ0b0eIAEoDjIZLnNzbS5QYXRjaENvbXBsaWFuY2VMZXZlbFIeYXBwcm92ZWRwYXRjaGVz'
    'Y29tcGxpYW5jZWxldmVsEk4KIGFwcHJvdmVkcGF0Y2hlc2VuYWJsZW5vbnNlY3VyaXR5GL2m94'
    'wBIAEoCFIgYXBwcm92ZWRwYXRjaGVzZW5hYmxlbm9uc2VjdXJpdHkSeQooYXZhaWxhYmxlc2Vj'
    'dXJpdHl1cGRhdGVzY29tcGxpYW5jZXN0YXR1cxjyr7JZIAEoDjIaLnNzbS5QYXRjaENvbXBsaW'
    'FuY2VTdGF0dXNSKGF2YWlsYWJsZXNlY3VyaXR5dXBkYXRlc2NvbXBsaWFuY2VzdGF0dXMSIwoL'
    'Y2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaWVudHRva2VuEiMKC2Rlc2NyaXB0aW9uGIr0+TYgAS'
    'gJUgtkZXNjcmlwdGlvbhI+Cg1nbG9iYWxmaWx0ZXJzGOLcxn0gASgLMhUuc3NtLlBhdGNoRmls'
    'dGVyR3JvdXBSDWdsb2JhbGZpbHRlcnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRJBCg9vcGVyYX'
    'RpbmdzeXN0ZW0Y6v3BEiABKA4yFC5zc20uT3BlcmF0aW5nU3lzdGVtUg9vcGVyYXRpbmdzeXN0'
    'ZW0SLAoPcmVqZWN0ZWRwYXRjaGVzGJz805MBIAMoCVIPcmVqZWN0ZWRwYXRjaGVzEkoKFXJlam'
    'VjdGVkcGF0Y2hlc2FjdGlvbhjar4GqASABKA4yEC5zc20uUGF0Y2hBY3Rpb25SFXJlamVjdGVk'
    'cGF0Y2hlc2FjdGlvbhItCgdzb3VyY2VzGNLnnRYgAygLMhAuc3NtLlBhdGNoU291cmNlUgdzb3'
    'VyY2VzEiAKBHRhZ3MYwcH2tQEgAygLMgguc3NtLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createPatchBaselineResultDescriptor instead')
const CreatePatchBaselineResult$json = {
  '1': 'CreatePatchBaselineResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `CreatePatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPatchBaselineResultDescriptor =
    $convert.base64Decode(
        'ChlDcmVhdGVQYXRjaEJhc2VsaW5lUmVzdWx0EiEKCmJhc2VsaW5laWQY0OTbKCABKAlSCmJhc2'
        'VsaW5laWQ=');

@$core.Deprecated('Use createResourceDataSyncRequestDescriptor instead')
const CreateResourceDataSyncRequest$json = {
  '1': 'CreateResourceDataSyncRequest',
  '2': [
    {
      '1': 's3destination',
      '3': 526234522,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncS3Destination',
      '10': 's3destination'
    },
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
    {
      '1': 'syncsource',
      '3': 286486824,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncSource',
      '10': 'syncsource'
    },
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `CreateResourceDataSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResourceDataSyncRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVSZXNvdXJjZURhdGFTeW5jUmVxdWVzdBJMCg1zM2Rlc3RpbmF0aW9uGJrn9voBIA'
    'EoCzIiLnNzbS5SZXNvdXJjZURhdGFTeW5jUzNEZXN0aW5hdGlvblINczNkZXN0aW5hdGlvbhIe'
    'CghzeW5jbmFtZRiilrKwASABKAlSCHN5bmNuYW1lEj8KCnN5bmNzb3VyY2UYqOLNiAEgASgLMh'
    'suc3NtLlJlc291cmNlRGF0YVN5bmNTb3VyY2VSCnN5bmNzb3VyY2USHQoIc3luY3R5cGUY2+aq'
    'OiABKAlSCHN5bmN0eXBl');

@$core.Deprecated('Use createResourceDataSyncResultDescriptor instead')
const CreateResourceDataSyncResult$json = {
  '1': 'CreateResourceDataSyncResult',
};

/// Descriptor for `CreateResourceDataSyncResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createResourceDataSyncResultDescriptor =
    $convert.base64Decode('ChxDcmVhdGVSZXNvdXJjZURhdGFTeW5jUmVzdWx0');

@$core.Deprecated('Use credentialsDescriptor instead')
const Credentials$json = {
  '1': 'Credentials',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {
      '1': 'expirationtime',
      '3': 93473378,
      '4': 1,
      '5': 9,
      '10': 'expirationtime'
    },
    {
      '1': 'secretaccesskey',
      '3': 172288487,
      '4': 1,
      '5': 9,
      '10': 'secretaccesskey'
    },
    {'1': 'sessiontoken', '3': 211161069, '4': 1, '5': 9, '10': 'sessiontoken'},
  ],
};

/// Descriptor for `Credentials`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialsDescriptor = $convert.base64Decode(
    'CgtDcmVkZW50aWFscxIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2Vzc2tleWlkEikKDm'
    'V4cGlyYXRpb250aW1lGOKUySwgASgJUg5leHBpcmF0aW9udGltZRIrCg9zZWNyZXRhY2Nlc3Nr'
    'ZXkY59OTUiABKAlSD3NlY3JldGFjY2Vzc2tleRIlCgxzZXNzaW9udG9rZW4Y7Z/YZCABKAlSDH'
    'Nlc3Npb250b2tlbg==');

@$core
    .Deprecated('Use customSchemaCountLimitExceededExceptionDescriptor instead')
const CustomSchemaCountLimitExceededException$json = {
  '1': 'CustomSchemaCountLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CustomSchemaCountLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customSchemaCountLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CidDdXN0b21TY2hlbWFDb3VudExpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use deleteActivationRequestDescriptor instead')
const DeleteActivationRequest$json = {
  '1': 'DeleteActivationRequest',
  '2': [
    {'1': 'activationid', '3': 146858601, '4': 1, '5': 9, '10': 'activationid'},
  ],
};

/// Descriptor for `DeleteActivationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteActivationRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVBY3RpdmF0aW9uUmVxdWVzdBIlCgxhY3RpdmF0aW9uaWQY6cSDRiABKAlSDGFjdG'
        'l2YXRpb25pZA==');

@$core.Deprecated('Use deleteActivationResultDescriptor instead')
const DeleteActivationResult$json = {
  '1': 'DeleteActivationResult',
};

/// Descriptor for `DeleteActivationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteActivationResultDescriptor =
    $convert.base64Decode('ChZEZWxldGVBY3RpdmF0aW9uUmVzdWx0');

@$core.Deprecated('Use deleteAssociationRequestDescriptor instead')
const DeleteAssociationRequest$json = {
  '1': 'DeleteAssociationRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAssociationRequestDescriptor = $convert.base64Decode(
    'ChhEZWxldGVBc3NvY2lhdGlvblJlcXVlc3QSJwoNYXNzb2NpYXRpb25pZBiS/JVCIAEoCVINYX'
    'Nzb2NpYXRpb25pZBIhCgppbnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEhUKBG5hbWUY'
    'h+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteAssociationResultDescriptor instead')
const DeleteAssociationResult$json = {
  '1': 'DeleteAssociationResult',
};

/// Descriptor for `DeleteAssociationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAssociationResultDescriptor =
    $convert.base64Decode('ChdEZWxldGVBc3NvY2lhdGlvblJlc3VsdA==');

@$core.Deprecated('Use deleteDocumentRequestDescriptor instead')
const DeleteDocumentRequest$json = {
  '1': 'DeleteDocumentRequest',
  '2': [
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'force', '3': 526711333, '4': 1, '5': 8, '10': 'force'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DeleteDocumentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDocumentRequestDescriptor = $convert.base64Decode(
    'ChVEZWxldGVEb2N1bWVudFJlcXVlc3QSKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2'
    'N1bWVudHZlcnNpb24SGAoFZm9yY2UYpfST+wEgASgIUgVmb3JjZRIVCgRuYW1lGIfmgX8gASgJ'
    'UgRuYW1lEiMKC3ZlcnNpb25uYW1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use deleteDocumentResultDescriptor instead')
const DeleteDocumentResult$json = {
  '1': 'DeleteDocumentResult',
};

/// Descriptor for `DeleteDocumentResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDocumentResultDescriptor =
    $convert.base64Decode('ChREZWxldGVEb2N1bWVudFJlc3VsdA==');

@$core.Deprecated('Use deleteInventoryRequestDescriptor instead')
const DeleteInventoryRequest$json = {
  '1': 'DeleteInventoryRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'schemadeleteoption',
      '3': 105711169,
      '4': 1,
      '5': 14,
      '6': '.ssm.InventorySchemaDeleteOption',
      '10': 'schemadeleteoption'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `DeleteInventoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteInventoryRequestDescriptor = $convert.base64Decode(
    'ChZEZWxldGVJbnZlbnRvcnlSZXF1ZXN0EiMKC2NsaWVudHRva2VuGMz7u0EgASgJUgtjbGllbn'
    'R0b2tlbhIZCgZkcnlydW4YuN/7KyABKAhSBmRyeXJ1bhJTChJzY2hlbWFkZWxldGVvcHRpb24Y'
    'wYy0MiABKA4yIC5zc20uSW52ZW50b3J5U2NoZW1hRGVsZXRlT3B0aW9uUhJzY2hlbWFkZWxldG'
    'VvcHRpb24SHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ==');

@$core.Deprecated('Use deleteInventoryResultDescriptor instead')
const DeleteInventoryResult$json = {
  '1': 'DeleteInventoryResult',
  '2': [
    {'1': 'deletionid', '3': 126693587, '4': 1, '5': 9, '10': 'deletionid'},
    {
      '1': 'deletionsummary',
      '3': 189954074,
      '4': 1,
      '5': 11,
      '6': '.ssm.InventoryDeletionSummary',
      '10': 'deletionsummary'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `DeleteInventoryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteInventoryResultDescriptor = $convert.base64Decode(
    'ChVEZWxldGVJbnZlbnRvcnlSZXN1bHQSIQoKZGVsZXRpb25pZBjT4bQ8IAEoCVIKZGVsZXRpb2'
    '5pZBJKCg9kZWxldGlvbnN1bW1hcnkYmvDJWiABKAsyHS5zc20uSW52ZW50b3J5RGVsZXRpb25T'
    'dW1tYXJ5Ug9kZWxldGlvbnN1bW1hcnkSHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ'
    '==');

@$core.Deprecated('Use deleteMaintenanceWindowRequestDescriptor instead')
const DeleteMaintenanceWindowRequest$json = {
  '1': 'DeleteMaintenanceWindowRequest',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DeleteMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMaintenanceWindowRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVNYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSHQoId2luZG93aWQYqeSHCSABKAlSCH'
        'dpbmRvd2lk');

@$core.Deprecated('Use deleteMaintenanceWindowResultDescriptor instead')
const DeleteMaintenanceWindowResult$json = {
  '1': 'DeleteMaintenanceWindowResult',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DeleteMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMaintenanceWindowResultDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVNYWludGVuYW5jZVdpbmRvd1Jlc3VsdBIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2'
        'luZG93aWQ=');

@$core.Deprecated('Use deleteOpsItemRequestDescriptor instead')
const DeleteOpsItemRequest$json = {
  '1': 'DeleteOpsItemRequest',
  '2': [
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `DeleteOpsItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOpsItemRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVPcHNJdGVtUmVxdWVzdBIfCglvcHNpdGVtaWQY0tKVDCABKAlSCW9wc2l0ZW1pZA'
    '==');

@$core.Deprecated('Use deleteOpsItemResponseDescriptor instead')
const DeleteOpsItemResponse$json = {
  '1': 'DeleteOpsItemResponse',
};

/// Descriptor for `DeleteOpsItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOpsItemResponseDescriptor =
    $convert.base64Decode('ChVEZWxldGVPcHNJdGVtUmVzcG9uc2U=');

@$core.Deprecated('Use deleteOpsMetadataRequestDescriptor instead')
const DeleteOpsMetadataRequest$json = {
  '1': 'DeleteOpsMetadataRequest',
  '2': [
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
  ],
};

/// Descriptor for `DeleteOpsMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOpsMetadataRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVPcHNNZXRhZGF0YVJlcXVlc3QSKgoOb3BzbWV0YWRhdGFhcm4Yor6C5gEgASgJUg'
        '5vcHNtZXRhZGF0YWFybg==');

@$core.Deprecated('Use deleteOpsMetadataResultDescriptor instead')
const DeleteOpsMetadataResult$json = {
  '1': 'DeleteOpsMetadataResult',
};

/// Descriptor for `DeleteOpsMetadataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteOpsMetadataResultDescriptor =
    $convert.base64Decode('ChdEZWxldGVPcHNNZXRhZGF0YVJlc3VsdA==');

@$core.Deprecated('Use deleteParameterRequestDescriptor instead')
const DeleteParameterRequest$json = {
  '1': 'DeleteParameterRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteParameterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteParameterRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWxldGVQYXJhbWV0ZXJSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteParameterResultDescriptor instead')
const DeleteParameterResult$json = {
  '1': 'DeleteParameterResult',
};

/// Descriptor for `DeleteParameterResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteParameterResultDescriptor =
    $convert.base64Decode('ChVEZWxldGVQYXJhbWV0ZXJSZXN1bHQ=');

@$core.Deprecated('Use deleteParametersRequestDescriptor instead')
const DeleteParametersRequest$json = {
  '1': 'DeleteParametersRequest',
  '2': [
    {'1': 'names', '3': 324387120, '4': 3, '5': 9, '10': 'names'},
  ],
};

/// Descriptor for `DeleteParametersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteParametersRequestDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVQYXJhbWV0ZXJzUmVxdWVzdBIYCgVuYW1lcxiwgteaASADKAlSBW5hbWVz');

@$core.Deprecated('Use deleteParametersResultDescriptor instead')
const DeleteParametersResult$json = {
  '1': 'DeleteParametersResult',
  '2': [
    {
      '1': 'deletedparameters',
      '3': 479939243,
      '4': 3,
      '5': 9,
      '10': 'deletedparameters'
    },
    {
      '1': 'invalidparameters',
      '3': 329741157,
      '4': 3,
      '5': 9,
      '10': 'invalidparameters'
    },
  ],
};

/// Descriptor for `DeleteParametersResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteParametersResultDescriptor = $convert.base64Decode(
    'ChZEZWxldGVQYXJhbWV0ZXJzUmVzdWx0EjAKEWRlbGV0ZWRwYXJhbWV0ZXJzGKuV7eQBIAMoCV'
    'IRZGVsZXRlZHBhcmFtZXRlcnMSMAoRaW52YWxpZHBhcmFtZXRlcnMY5eadnQEgAygJUhFpbnZh'
    'bGlkcGFyYW1ldGVycw==');

@$core.Deprecated('Use deletePatchBaselineRequestDescriptor instead')
const DeletePatchBaselineRequest$json = {
  '1': 'DeletePatchBaselineRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `DeletePatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePatchBaselineRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVQYXRjaEJhc2VsaW5lUmVxdWVzdBIhCgpiYXNlbGluZWlkGNDk2yggASgJUgpiYX'
        'NlbGluZWlk');

@$core.Deprecated('Use deletePatchBaselineResultDescriptor instead')
const DeletePatchBaselineResult$json = {
  '1': 'DeletePatchBaselineResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `DeletePatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePatchBaselineResultDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVQYXRjaEJhc2VsaW5lUmVzdWx0EiEKCmJhc2VsaW5laWQY0OTbKCABKAlSCmJhc2'
        'VsaW5laWQ=');

@$core.Deprecated('Use deleteResourceDataSyncRequestDescriptor instead')
const DeleteResourceDataSyncRequest$json = {
  '1': 'DeleteResourceDataSyncRequest',
  '2': [
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `DeleteResourceDataSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourceDataSyncRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVSZXNvdXJjZURhdGFTeW5jUmVxdWVzdBIeCghzeW5jbmFtZRiilrKwASABKAlSCH'
        'N5bmNuYW1lEh0KCHN5bmN0eXBlGNvmqjogASgJUghzeW5jdHlwZQ==');

@$core.Deprecated('Use deleteResourceDataSyncResultDescriptor instead')
const DeleteResourceDataSyncResult$json = {
  '1': 'DeleteResourceDataSyncResult',
};

/// Descriptor for `DeleteResourceDataSyncResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourceDataSyncResultDescriptor =
    $convert.base64Decode('ChxEZWxldGVSZXNvdXJjZURhdGFTeW5jUmVzdWx0');

@$core.Deprecated('Use deleteResourcePolicyRequestDescriptor instead')
const DeleteResourcePolicyRequest$json = {
  '1': 'DeleteResourcePolicyRequest',
  '2': [
    {'1': 'policyhash', '3': 203037360, '4': 1, '5': 9, '10': 'policyhash'},
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `DeleteResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVSZXNvdXJjZVBvbGljeVJlcXVlc3QSIQoKcG9saWN5aGFzaBiwtehgIAEoCVIKcG'
        '9saWN5aGFzaBIeCghwb2xpY3lpZBjzo+mOASABKAlSCHBvbGljeWlkEiQKC3Jlc291cmNlYXJu'
        'GK342a0BIAEoCVILcmVzb3VyY2Vhcm4=');

@$core.Deprecated('Use deleteResourcePolicyResponseDescriptor instead')
const DeleteResourcePolicyResponse$json = {
  '1': 'DeleteResourcePolicyResponse',
};

/// Descriptor for `DeleteResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVSZXNvdXJjZVBvbGljeVJlc3BvbnNl');

@$core.Deprecated('Use deregisterManagedInstanceRequestDescriptor instead')
const DeregisterManagedInstanceRequest$json = {
  '1': 'DeregisterManagedInstanceRequest',
  '2': [
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
  ],
};

/// Descriptor for `DeregisterManagedInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deregisterManagedInstanceRequestDescriptor =
    $convert.base64Decode(
        'CiBEZXJlZ2lzdGVyTWFuYWdlZEluc3RhbmNlUmVxdWVzdBIhCgppbnN0YW5jZWlkGKCt0RcgAS'
        'gJUgppbnN0YW5jZWlk');

@$core.Deprecated('Use deregisterManagedInstanceResultDescriptor instead')
const DeregisterManagedInstanceResult$json = {
  '1': 'DeregisterManagedInstanceResult',
};

/// Descriptor for `DeregisterManagedInstanceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deregisterManagedInstanceResultDescriptor =
    $convert.base64Decode('Ch9EZXJlZ2lzdGVyTWFuYWdlZEluc3RhbmNlUmVzdWx0');

@$core.Deprecated(
    'Use deregisterPatchBaselineForPatchGroupRequestDescriptor instead')
const DeregisterPatchBaselineForPatchGroupRequest$json = {
  '1': 'DeregisterPatchBaselineForPatchGroupRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `DeregisterPatchBaselineForPatchGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterPatchBaselineForPatchGroupRequestDescriptor =
    $convert.base64Decode(
        'CitEZXJlZ2lzdGVyUGF0Y2hCYXNlbGluZUZvclBhdGNoR3JvdXBSZXF1ZXN0EiEKCmJhc2VsaW'
        '5laWQY0OTbKCABKAlSCmJhc2VsaW5laWQSIgoKcGF0Y2hncm91cBjht7H3ASABKAlSCnBhdGNo'
        'Z3JvdXA=');

@$core.Deprecated(
    'Use deregisterPatchBaselineForPatchGroupResultDescriptor instead')
const DeregisterPatchBaselineForPatchGroupResult$json = {
  '1': 'DeregisterPatchBaselineForPatchGroupResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `DeregisterPatchBaselineForPatchGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterPatchBaselineForPatchGroupResultDescriptor =
    $convert.base64Decode(
        'CipEZXJlZ2lzdGVyUGF0Y2hCYXNlbGluZUZvclBhdGNoR3JvdXBSZXN1bHQSIQoKYmFzZWxpbm'
        'VpZBjQ5NsoIAEoCVIKYmFzZWxpbmVpZBIiCgpwYXRjaGdyb3VwGOG3sfcBIAEoCVIKcGF0Y2hn'
        'cm91cA==');

@$core.Deprecated(
    'Use deregisterTargetFromMaintenanceWindowRequestDescriptor instead')
const DeregisterTargetFromMaintenanceWindowRequest$json = {
  '1': 'DeregisterTargetFromMaintenanceWindowRequest',
  '2': [
    {'1': 'safe', '3': 223623801, '4': 1, '5': 8, '10': 'safe'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `DeregisterTargetFromMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterTargetFromMaintenanceWindowRequestDescriptor =
    $convert.base64Decode(
        'CixEZXJlZ2lzdGVyVGFyZ2V0RnJvbU1haW50ZW5hbmNlV2luZG93UmVxdWVzdBIVCgRzYWZlGP'
        'n00GogASgIUgRzYWZlEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZBIpCg53aW5kb3d0'
        'YXJnZXRpZBjezup7IAEoCVIOd2luZG93dGFyZ2V0aWQ=');

@$core.Deprecated(
    'Use deregisterTargetFromMaintenanceWindowResultDescriptor instead')
const DeregisterTargetFromMaintenanceWindowResult$json = {
  '1': 'DeregisterTargetFromMaintenanceWindowResult',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `DeregisterTargetFromMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterTargetFromMaintenanceWindowResultDescriptor =
    $convert.base64Decode(
        'CitEZXJlZ2lzdGVyVGFyZ2V0RnJvbU1haW50ZW5hbmNlV2luZG93UmVzdWx0Eh0KCHdpbmRvd2'
        'lkGKnkhwkgASgJUgh3aW5kb3dpZBIpCg53aW5kb3d0YXJnZXRpZBjezup7IAEoCVIOd2luZG93'
        'dGFyZ2V0aWQ=');

@$core.Deprecated(
    'Use deregisterTaskFromMaintenanceWindowRequestDescriptor instead')
const DeregisterTaskFromMaintenanceWindowRequest$json = {
  '1': 'DeregisterTaskFromMaintenanceWindowRequest',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
};

/// Descriptor for `DeregisterTaskFromMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterTaskFromMaintenanceWindowRequestDescriptor =
    $convert.base64Decode(
        'CipEZXJlZ2lzdGVyVGFza0Zyb21NYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSHQoId2luZG93aW'
        'QYqeSHCSABKAlSCHdpbmRvd2lkEiYKDHdpbmRvd3Rhc2tpZBjajbGaASABKAlSDHdpbmRvd3Rh'
        'c2tpZA==');

@$core.Deprecated(
    'Use deregisterTaskFromMaintenanceWindowResultDescriptor instead')
const DeregisterTaskFromMaintenanceWindowResult$json = {
  '1': 'DeregisterTaskFromMaintenanceWindowResult',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
};

/// Descriptor for `DeregisterTaskFromMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deregisterTaskFromMaintenanceWindowResultDescriptor = $convert.base64Decode(
        'CilEZXJlZ2lzdGVyVGFza0Zyb21NYWludGVuYW5jZVdpbmRvd1Jlc3VsdBIdCgh3aW5kb3dpZB'
        'ip5IcJIAEoCVIId2luZG93aWQSJgoMd2luZG93dGFza2lkGNqNsZoBIAEoCVIMd2luZG93dGFz'
        'a2lk');

@$core.Deprecated('Use describeActivationsFilterDescriptor instead')
const DescribeActivationsFilter$json = {
  '1': 'DescribeActivationsFilter',
  '2': [
    {
      '1': 'filterkey',
      '3': 300347055,
      '4': 1,
      '5': 14,
      '6': '.ssm.DescribeActivationsFilterKeys',
      '10': 'filterkey'
    },
    {'1': 'filtervalues', '3': 471484302, '4': 3, '5': 9, '10': 'filtervalues'},
  ],
};

/// Descriptor for `DescribeActivationsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeActivationsFilterDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUFjdGl2YXRpb25zRmlsdGVyEkQKCWZpbHRlcmtleRiv3ZuPASABKA4yIi5zc2'
    '0uRGVzY3JpYmVBY3RpdmF0aW9uc0ZpbHRlcktleXNSCWZpbHRlcmtleRImCgxmaWx0ZXJ2YWx1'
    'ZXMYjo/p4AEgAygJUgxmaWx0ZXJ2YWx1ZXM=');

@$core.Deprecated('Use describeActivationsRequestDescriptor instead')
const DescribeActivationsRequest$json = {
  '1': 'DescribeActivationsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.DescribeActivationsFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeActivationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeActivationsRequestDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZUFjdGl2YXRpb25zUmVxdWVzdBI7CgdmaWx0ZXJzGO3N6lkgAygLMh4uc3NtLk'
        'Rlc2NyaWJlQWN0aXZhdGlvbnNGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASAB'
        'KAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeActivationsResultDescriptor instead')
const DescribeActivationsResult$json = {
  '1': 'DescribeActivationsResult',
  '2': [
    {
      '1': 'activationlist',
      '3': 341191290,
      '4': 3,
      '5': 11,
      '6': '.ssm.Activation',
      '10': 'activationlist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeActivationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeActivationsResultDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUFjdGl2YXRpb25zUmVzdWx0EjsKDmFjdGl2YXRpb25saXN0GPrU2KIBIAMoCz'
    'IPLnNzbS5BY3RpdmF0aW9uUg5hY3RpdmF0aW9ubGlzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlS'
    'CW5leHR0b2tlbg==');

@$core.Deprecated(
    'Use describeAssociationExecutionTargetsRequestDescriptor instead')
const DescribeAssociationExecutionTargetsRequest$json = {
  '1': 'DescribeAssociationExecutionTargetsRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationExecutionTargetsFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAssociationExecutionTargetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeAssociationExecutionTargetsRequestDescriptor =
    $convert.base64Decode(
        'CipEZXNjcmliZUFzc29jaWF0aW9uRXhlY3V0aW9uVGFyZ2V0c1JlcXVlc3QSJwoNYXNzb2NpYX'
        'Rpb25pZBiS/JVCIAEoCVINYXNzb2NpYXRpb25pZBIjCgtleGVjdXRpb25pZBixz69GIAEoCVIL'
        'ZXhlY3V0aW9uaWQSQwoHZmlsdGVycxjtzepZIAMoCzImLnNzbS5Bc3NvY2lhdGlvbkV4ZWN1dG'
        'lvblRhcmdldHNGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJl'
        'c3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated(
    'Use describeAssociationExecutionTargetsResultDescriptor instead')
const DescribeAssociationExecutionTargetsResult$json = {
  '1': 'DescribeAssociationExecutionTargetsResult',
  '2': [
    {
      '1': 'associationexecutiontargets',
      '3': 124720489,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationExecutionTarget',
      '10': 'associationexecutiontargets'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAssociationExecutionTargetsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeAssociationExecutionTargetsResultDescriptor = $convert.base64Decode(
        'CilEZXNjcmliZUFzc29jaWF0aW9uRXhlY3V0aW9uVGFyZ2V0c1Jlc3VsdBJkChthc3NvY2lhdG'
        'lvbmV4ZWN1dGlvbnRhcmdldHMY6aq8OyADKAsyHy5zc20uQXNzb2NpYXRpb25FeGVjdXRpb25U'
        'YXJnZXRSG2Fzc29jaWF0aW9uZXhlY3V0aW9udGFyZ2V0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeAssociationExecutionsRequestDescriptor instead')
const DescribeAssociationExecutionsRequest$json = {
  '1': 'DescribeAssociationExecutionsRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationExecutionFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAssociationExecutionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAssociationExecutionsRequestDescriptor =
    $convert.base64Decode(
        'CiREZXNjcmliZUFzc29jaWF0aW9uRXhlY3V0aW9uc1JlcXVlc3QSJwoNYXNzb2NpYXRpb25pZB'
        'iS/JVCIAEoCVINYXNzb2NpYXRpb25pZBI8CgdmaWx0ZXJzGO3N6lkgAygLMh8uc3NtLkFzc29j'
        'aWF0aW9uRXhlY3V0aW9uRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUg'
        'ptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeAssociationExecutionsResultDescriptor instead')
const DescribeAssociationExecutionsResult$json = {
  '1': 'DescribeAssociationExecutionsResult',
  '2': [
    {
      '1': 'associationexecutions',
      '3': 406414972,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationExecution',
      '10': 'associationexecutions'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAssociationExecutionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAssociationExecutionsResultDescriptor =
    $convert.base64Decode(
        'CiNEZXNjcmliZUFzc29jaWF0aW9uRXhlY3V0aW9uc1Jlc3VsdBJTChVhc3NvY2lhdGlvbmV4ZW'
        'N1dGlvbnMY/MzlwQEgAygLMhkuc3NtLkFzc29jaWF0aW9uRXhlY3V0aW9uUhVhc3NvY2lhdGlv'
        'bmV4ZWN1dGlvbnMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeAssociationRequestDescriptor instead')
const DescribeAssociationRequest$json = {
  '1': 'DescribeAssociationRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DescribeAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAssociationRequestDescriptor = $convert.base64Decode(
    'ChpEZXNjcmliZUFzc29jaWF0aW9uUmVxdWVzdBInCg1hc3NvY2lhdGlvbmlkGJL8lUIgASgJUg'
    '1hc3NvY2lhdGlvbmlkEjIKEmFzc29jaWF0aW9udmVyc2lvbhiRisnVASABKAlSEmFzc29jaWF0'
    'aW9udmVyc2lvbhIhCgppbnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEhUKBG5hbWUYh+'
    'aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use describeAssociationResultDescriptor instead')
const DescribeAssociationResult$json = {
  '1': 'DescribeAssociationResult',
  '2': [
    {
      '1': 'associationdescription',
      '3': 344755863,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationDescription',
      '10': 'associationdescription'
    },
  ],
};

/// Descriptor for `DescribeAssociationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAssociationResultDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZUFzc29jaWF0aW9uUmVzdWx0ElcKFmFzc29jaWF0aW9uZGVzY3JpcHRpb24Yl5'
    '2ypAEgASgLMhsuc3NtLkFzc29jaWF0aW9uRGVzY3JpcHRpb25SFmFzc29jaWF0aW9uZGVzY3Jp'
    'cHRpb24=');

@$core.Deprecated('Use describeAutomationExecutionsRequestDescriptor instead')
const DescribeAutomationExecutionsRequest$json = {
  '1': 'DescribeAutomationExecutionsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAutomationExecutionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAutomationExecutionsRequestDescriptor =
    $convert.base64Decode(
        'CiNEZXNjcmliZUF1dG9tYXRpb25FeGVjdXRpb25zUmVxdWVzdBI7CgdmaWx0ZXJzGO3N6lkgAy'
        'gLMh4uc3NtLkF1dG9tYXRpb25FeGVjdXRpb25GaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0'
        'cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
        '4=');

@$core.Deprecated('Use describeAutomationExecutionsResultDescriptor instead')
const DescribeAutomationExecutionsResult$json = {
  '1': 'DescribeAutomationExecutionsResult',
  '2': [
    {
      '1': 'automationexecutionmetadatalist',
      '3': 202240522,
      '4': 3,
      '5': 11,
      '6': '.ssm.AutomationExecutionMetadata',
      '10': 'automationexecutionmetadatalist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAutomationExecutionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAutomationExecutionsResultDescriptor =
    $convert.base64Decode(
        'CiJEZXNjcmliZUF1dG9tYXRpb25FeGVjdXRpb25zUmVzdWx0Em0KH2F1dG9tYXRpb25leGVjdX'
        'Rpb25tZXRhZGF0YWxpc3QYiuS3YCADKAsyIC5zc20uQXV0b21hdGlvbkV4ZWN1dGlvbk1ldGFk'
        'YXRhUh9hdXRvbWF0aW9uZXhlY3V0aW9ubWV0YWRhdGFsaXN0Eh8KCW5leHR0b2tlbhj+hLpnIA'
        'EoCVIJbmV4dHRva2Vu');

@$core
    .Deprecated('Use describeAutomationStepExecutionsRequestDescriptor instead')
const DescribeAutomationStepExecutionsRequest$json = {
  '1': 'DescribeAutomationStepExecutionsRequest',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecutionFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'reverseorder', '3': 427445336, '4': 1, '5': 8, '10': 'reverseorder'},
  ],
};

/// Descriptor for `DescribeAutomationStepExecutionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAutomationStepExecutionsRequestDescriptor =
    $convert.base64Decode(
        'CidEZXNjcmliZUF1dG9tYXRpb25TdGVwRXhlY3V0aW9uc1JlcXVlc3QSNwoVYXV0b21hdGlvbm'
        'V4ZWN1dGlvbmlkGObv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQSNQoHZmlsdGVycxjt'
        'zepZIAMoCzIYLnNzbS5TdGVwRXhlY3V0aW9uRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdH'
        'MYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu'
        'EiYKDHJldmVyc2VvcmRlchjYmOnLASABKAhSDHJldmVyc2VvcmRlcg==');

@$core
    .Deprecated('Use describeAutomationStepExecutionsResultDescriptor instead')
const DescribeAutomationStepExecutionsResult$json = {
  '1': 'DescribeAutomationStepExecutionsResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'stepexecutions',
      '3': 164157981,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecution',
      '10': 'stepexecutions'
    },
  ],
};

/// Descriptor for `DescribeAutomationStepExecutionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAutomationStepExecutionsResultDescriptor =
    $convert.base64Decode(
        'CiZEZXNjcmliZUF1dG9tYXRpb25TdGVwRXhlY3V0aW9uc1Jlc3VsdBIfCgluZXh0dG9rZW4Y/o'
        'S6ZyABKAlSCW5leHR0b2tlbhI9Cg5zdGVwZXhlY3V0aW9ucxidtKNOIAMoCzISLnNzbS5TdGVw'
        'RXhlY3V0aW9uUg5zdGVwZXhlY3V0aW9ucw==');

@$core.Deprecated('Use describeAvailablePatchesRequestDescriptor instead')
const DescribeAvailablePatchesRequest$json = {
  '1': 'DescribeAvailablePatchesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchOrchestratorFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeAvailablePatchesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAvailablePatchesRequestDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUF2YWlsYWJsZVBhdGNoZXNSZXF1ZXN0EjkKB2ZpbHRlcnMY7c3qWSADKAsyHC'
        '5zc20uUGF0Y2hPcmNoZXN0cmF0b3JGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuD'
        'ASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeAvailablePatchesResultDescriptor instead')
const DescribeAvailablePatchesResult$json = {
  '1': 'DescribeAvailablePatchesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'patches',
      '3': 66243222,
      '4': 3,
      '5': 11,
      '6': '.ssm.Patch',
      '10': 'patches'
    },
  ],
};

/// Descriptor for `DescribeAvailablePatchesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAvailablePatchesResultDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUF2YWlsYWJsZVBhdGNoZXNSZXN1bHQSHwoJbmV4dHRva2VuGP6EumcgASgJUg'
        'luZXh0dG9rZW4SJwoHcGF0Y2hlcxiWlcsfIAMoCzIKLnNzbS5QYXRjaFIHcGF0Y2hlcw==');

@$core.Deprecated('Use describeDocumentPermissionRequestDescriptor instead')
const DescribeDocumentPermissionRequest$json = {
  '1': 'DescribeDocumentPermissionRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'permissiontype',
      '3': 479110739,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentPermissionType',
      '10': 'permissiontype'
    },
  ],
};

/// Descriptor for `DescribeDocumentPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDocumentPermissionRequestDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZURvY3VtZW50UGVybWlzc2lvblJlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDAS'
        'ABKAVSCm1heHJlc3VsdHMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIfCgluZXh0dG9rZW4Y/oS6'
        'ZyABKAlSCW5leHR0b2tlbhJHCg5wZXJtaXNzaW9udHlwZRjTzLrkASABKA4yGy5zc20uRG9jdW'
        '1lbnRQZXJtaXNzaW9uVHlwZVIOcGVybWlzc2lvbnR5cGU=');

@$core.Deprecated('Use describeDocumentPermissionResponseDescriptor instead')
const DescribeDocumentPermissionResponse$json = {
  '1': 'DescribeDocumentPermissionResponse',
  '2': [
    {'1': 'accountids', '3': 306323207, '4': 3, '5': 9, '10': 'accountids'},
    {
      '1': 'accountsharinginfolist',
      '3': 250360901,
      '4': 3,
      '5': 11,
      '6': '.ssm.AccountSharingInfo',
      '10': 'accountsharinginfolist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeDocumentPermissionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDocumentPermissionResponseDescriptor =
    $convert.base64Decode(
        'CiJEZXNjcmliZURvY3VtZW50UGVybWlzc2lvblJlc3BvbnNlEiIKCmFjY291bnRpZHMYh76Ikg'
        'EgAygJUgphY2NvdW50aWRzElIKFmFjY291bnRzaGFyaW5naW5mb2xpc3QYxeiwdyADKAsyFy5z'
        'c20uQWNjb3VudFNoYXJpbmdJbmZvUhZhY2NvdW50c2hhcmluZ2luZm9saXN0Eh8KCW5leHR0b2'
        'tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeDocumentRequestDescriptor instead')
const DescribeDocumentRequest$json = {
  '1': 'DescribeDocumentRequest',
  '2': [
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DescribeDocumentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDocumentRequestDescriptor = $convert.base64Decode(
    'ChdEZXNjcmliZURvY3VtZW50UmVxdWVzdBIrCg9kb2N1bWVudHZlcnNpb24Yye+pKCABKAlSD2'
    'RvY3VtZW50dmVyc2lvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiMKC3ZlcnNpb25uYW1lGNWj'
    'tGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use describeDocumentResultDescriptor instead')
const DescribeDocumentResult$json = {
  '1': 'DescribeDocumentResult',
  '2': [
    {
      '1': 'document',
      '3': 407108341,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentDescription',
      '10': 'document'
    },
  ],
};

/// Descriptor for `DescribeDocumentResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeDocumentResultDescriptor =
    $convert.base64Decode(
        'ChZEZXNjcmliZURvY3VtZW50UmVzdWx0EjgKCGRvY3VtZW50GPX1j8IBIAEoCzIYLnNzbS5Eb2'
        'N1bWVudERlc2NyaXB0aW9uUghkb2N1bWVudA==');

@$core.Deprecated(
    'Use describeEffectiveInstanceAssociationsRequestDescriptor instead')
const DescribeEffectiveInstanceAssociationsRequest$json = {
  '1': 'DescribeEffectiveInstanceAssociationsRequest',
  '2': [
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeEffectiveInstanceAssociationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeEffectiveInstanceAssociationsRequestDescriptor =
    $convert.base64Decode(
        'CixEZXNjcmliZUVmZmVjdGl2ZUluc3RhbmNlQXNzb2NpYXRpb25zUmVxdWVzdBIhCgppbnN0YW'
        '5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhy'
        'ZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use describeEffectiveInstanceAssociationsResultDescriptor instead')
const DescribeEffectiveInstanceAssociationsResult$json = {
  '1': 'DescribeEffectiveInstanceAssociationsResult',
  '2': [
    {
      '1': 'associations',
      '3': 149658718,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceAssociation',
      '10': 'associations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeEffectiveInstanceAssociationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeEffectiveInstanceAssociationsResultDescriptor =
    $convert.base64Decode(
        'CitEZXNjcmliZUVmZmVjdGl2ZUluc3RhbmNlQXNzb2NpYXRpb25zUmVzdWx0Ej8KDGFzc29jaW'
        'F0aW9ucxjeuK5HIAMoCzIYLnNzbS5JbnN0YW5jZUFzc29jaWF0aW9uUgxhc3NvY2lhdGlvbnMS'
        'HwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated(
    'Use describeEffectivePatchesForPatchBaselineRequestDescriptor instead')
const DescribeEffectivePatchesForPatchBaselineRequest$json = {
  '1': 'DescribeEffectivePatchesForPatchBaselineRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeEffectivePatchesForPatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeEffectivePatchesForPatchBaselineRequestDescriptor =
    $convert.base64Decode(
        'Ci9EZXNjcmliZUVmZmVjdGl2ZVBhdGNoZXNGb3JQYXRjaEJhc2VsaW5lUmVxdWVzdBIhCgpiYX'
        'NlbGluZWlkGNDk2yggASgJUgpiYXNlbGluZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgpt'
        'YXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use describeEffectivePatchesForPatchBaselineResultDescriptor instead')
const DescribeEffectivePatchesForPatchBaselineResult$json = {
  '1': 'DescribeEffectivePatchesForPatchBaselineResult',
  '2': [
    {
      '1': 'effectivepatches',
      '3': 494972541,
      '4': 3,
      '5': 11,
      '6': '.ssm.EffectivePatch',
      '10': 'effectivepatches'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeEffectivePatchesForPatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeEffectivePatchesForPatchBaselineResultDescriptor =
    $convert.base64Decode(
        'Ci5EZXNjcmliZUVmZmVjdGl2ZVBhdGNoZXNGb3JQYXRjaEJhc2VsaW5lUmVzdWx0EkMKEGVmZm'
        'VjdGl2ZXBhdGNoZXMY/dyC7AEgAygLMhMuc3NtLkVmZmVjdGl2ZVBhdGNoUhBlZmZlY3RpdmVw'
        'YXRjaGVzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use describeInstanceAssociationsStatusRequestDescriptor instead')
const DescribeInstanceAssociationsStatusRequest$json = {
  '1': 'DescribeInstanceAssociationsStatusRequest',
  '2': [
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstanceAssociationsStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeInstanceAssociationsStatusRequestDescriptor = $convert.base64Decode(
        'CilEZXNjcmliZUluc3RhbmNlQXNzb2NpYXRpb25zU3RhdHVzUmVxdWVzdBIhCgppbnN0YW5jZW'
        'lkGKCt0RcgASgJUgppbnN0YW5jZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1'
        'bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use describeInstanceAssociationsStatusResultDescriptor instead')
const DescribeInstanceAssociationsStatusResult$json = {
  '1': 'DescribeInstanceAssociationsStatusResult',
  '2': [
    {
      '1': 'instanceassociationstatusinfos',
      '3': 261960249,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceAssociationStatusInfo',
      '10': 'instanceassociationstatusinfos'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstanceAssociationsStatusResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstanceAssociationsStatusResultDescriptor =
    $convert.base64Decode(
        'CihEZXNjcmliZUluc3RhbmNlQXNzb2NpYXRpb25zU3RhdHVzUmVzdWx0Em0KHmluc3RhbmNlYX'
        'Nzb2NpYXRpb25zdGF0dXNpbmZvcxi55PR8IAMoCzIiLnNzbS5JbnN0YW5jZUFzc29jaWF0aW9u'
        'U3RhdHVzSW5mb1IeaW5zdGFuY2Vhc3NvY2lhdGlvbnN0YXR1c2luZm9zEh8KCW5leHR0b2tlbh'
        'j+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeInstanceInformationRequestDescriptor instead')
const DescribeInstanceInformationRequest$json = {
  '1': 'DescribeInstanceInformationRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceInformationStringFilter',
      '10': 'filters'
    },
    {
      '1': 'instanceinformationfilterlist',
      '3': 364026675,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceInformationFilter',
      '10': 'instanceinformationfilterlist'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstanceInformationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstanceInformationRequestDescriptor = $convert.base64Decode(
    'CiJEZXNjcmliZUluc3RhbmNlSW5mb3JtYXRpb25SZXF1ZXN0EkEKB2ZpbHRlcnMY7c3qWSADKA'
    'syJC5zc20uSW5zdGFuY2VJbmZvcm1hdGlvblN0cmluZ0ZpbHRlclIHZmlsdGVycxJoCh1pbnN0'
    'YW5jZWluZm9ybWF0aW9uZmlsdGVybGlzdBiztsqtASADKAsyHi5zc20uSW5zdGFuY2VJbmZvcm'
    '1hdGlvbkZpbHRlclIdaW5zdGFuY2VpbmZvcm1hdGlvbmZpbHRlcmxpc3QSIgoKbWF4cmVzdWx0'
    'cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
    '4=');

@$core.Deprecated('Use describeInstanceInformationResultDescriptor instead')
const DescribeInstanceInformationResult$json = {
  '1': 'DescribeInstanceInformationResult',
  '2': [
    {
      '1': 'instanceinformationlist',
      '3': 170547885,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceInformation',
      '10': 'instanceinformationlist'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstanceInformationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstanceInformationResultDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUluc3RhbmNlSW5mb3JtYXRpb25SZXN1bHQSVQoXaW5zdGFuY2VpbmZvcm1hdG'
        'lvbmxpc3QYrbWpUSADKAsyGC5zc20uSW5zdGFuY2VJbmZvcm1hdGlvblIXaW5zdGFuY2VpbmZv'
        'cm1hdGlvbmxpc3QSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated(
    'Use describeInstancePatchStatesForPatchGroupRequestDescriptor instead')
const DescribeInstancePatchStatesForPatchGroupRequest$json = {
  '1': 'DescribeInstancePatchStatesForPatchGroupRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstancePatchStateFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `DescribeInstancePatchStatesForPatchGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeInstancePatchStatesForPatchGroupRequestDescriptor =
    $convert.base64Decode(
        'Ci9EZXNjcmliZUluc3RhbmNlUGF0Y2hTdGF0ZXNGb3JQYXRjaEdyb3VwUmVxdWVzdBI6CgdmaW'
        'x0ZXJzGO3N6lkgAygLMh0uc3NtLkluc3RhbmNlUGF0Y2hTdGF0ZUZpbHRlclIHZmlsdGVycxIi'
        'CgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhIiCgpwYXRjaGdyb3VwGOG3sfcBIAEoCVIKcGF0Y2hncm91cA==');

@$core.Deprecated(
    'Use describeInstancePatchStatesForPatchGroupResultDescriptor instead')
const DescribeInstancePatchStatesForPatchGroupResult$json = {
  '1': 'DescribeInstancePatchStatesForPatchGroupResult',
  '2': [
    {
      '1': 'instancepatchstates',
      '3': 373960971,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstancePatchState',
      '10': 'instancepatchstates'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePatchStatesForPatchGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeInstancePatchStatesForPatchGroupResultDescriptor =
    $convert.base64Decode(
        'Ci5EZXNjcmliZUluc3RhbmNlUGF0Y2hTdGF0ZXNGb3JQYXRjaEdyb3VwUmVzdWx0Ek0KE2luc3'
        'RhbmNlcGF0Y2hzdGF0ZXMYi+KosgEgAygLMhcuc3NtLkluc3RhbmNlUGF0Y2hTdGF0ZVITaW5z'
        'dGFuY2VwYXRjaHN0YXRlcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeInstancePatchStatesRequestDescriptor instead')
const DescribeInstancePatchStatesRequest$json = {
  '1': 'DescribeInstancePatchStatesRequest',
  '2': [
    {'1': 'instanceids', '3': 312792453, '4': 3, '5': 9, '10': 'instanceids'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePatchStatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePatchStatesRequestDescriptor =
    $convert.base64Decode(
        'CiJEZXNjcmliZUluc3RhbmNlUGF0Y2hTdGF0ZXNSZXF1ZXN0EiQKC2luc3RhbmNlaWRzGIWrk5'
        'UBIAMoCVILaW5zdGFuY2VpZHMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMS'
        'HwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describeInstancePatchStatesResultDescriptor instead')
const DescribeInstancePatchStatesResult$json = {
  '1': 'DescribeInstancePatchStatesResult',
  '2': [
    {
      '1': 'instancepatchstates',
      '3': 373960971,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstancePatchState',
      '10': 'instancepatchstates'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePatchStatesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePatchStatesResultDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUluc3RhbmNlUGF0Y2hTdGF0ZXNSZXN1bHQSTQoTaW5zdGFuY2VwYXRjaHN0YX'
        'RlcxiL4qiyASADKAsyFy5zc20uSW5zdGFuY2VQYXRjaFN0YXRlUhNpbnN0YW5jZXBhdGNoc3Rh'
        'dGVzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeInstancePatchesRequestDescriptor instead')
const DescribeInstancePatchesRequest$json = {
  '1': 'DescribeInstancePatchesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchOrchestratorFilter',
      '10': 'filters'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePatchesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePatchesRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUluc3RhbmNlUGF0Y2hlc1JlcXVlc3QSOQoHZmlsdGVycxjtzepZIAMoCzIcLn'
        'NzbS5QYXRjaE9yY2hlc3RyYXRvckZpbHRlclIHZmlsdGVycxIhCgppbnN0YW5jZWlkGKCt0Rcg'
        'ASgJUgppbnN0YW5jZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW'
        '5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describeInstancePatchesResultDescriptor instead')
const DescribeInstancePatchesResult$json = {
  '1': 'DescribeInstancePatchesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'patches',
      '3': 66243222,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchComplianceData',
      '10': 'patches'
    },
  ],
};

/// Descriptor for `DescribeInstancePatchesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePatchesResultDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZUluc3RhbmNlUGF0Y2hlc1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
        '5leHR0b2tlbhI1CgdwYXRjaGVzGJaVyx8gAygLMhguc3NtLlBhdGNoQ29tcGxpYW5jZURhdGFS'
        'B3BhdGNoZXM=');

@$core.Deprecated('Use describeInstancePropertiesRequestDescriptor instead')
const DescribeInstancePropertiesRequest$json = {
  '1': 'DescribeInstancePropertiesRequest',
  '2': [
    {
      '1': 'filterswithoperator',
      '3': 500781557,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstancePropertyStringFilter',
      '10': 'filterswithoperator'
    },
    {
      '1': 'instancepropertyfilterlist',
      '3': 72826574,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstancePropertyFilter',
      '10': 'instancepropertyfilterlist'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePropertiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePropertiesRequestDescriptor = $convert.base64Decode(
    'CiFEZXNjcmliZUluc3RhbmNlUHJvcGVydGllc1JlcXVlc3QSVwoTZmlsdGVyc3dpdGhvcGVyYX'
    'Rvchj1o+XuASADKAsyIS5zc20uSW5zdGFuY2VQcm9wZXJ0eVN0cmluZ0ZpbHRlclITZmlsdGVy'
    'c3dpdGhvcGVyYXRvchJeChppbnN0YW5jZXByb3BlcnR5ZmlsdGVybGlzdBjO/dwiIAMoCzIbLn'
    'NzbS5JbnN0YW5jZVByb3BlcnR5RmlsdGVyUhppbnN0YW5jZXByb3BlcnR5ZmlsdGVybGlzdBIi'
    'CgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
    'lSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeInstancePropertiesResultDescriptor instead')
const DescribeInstancePropertiesResult$json = {
  '1': 'DescribeInstancePropertiesResult',
  '2': [
    {
      '1': 'instanceproperties',
      '3': 346214368,
      '4': 3,
      '5': 11,
      '6': '.ssm.InstanceProperty',
      '10': 'instanceproperties'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInstancePropertiesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInstancePropertiesResultDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUluc3RhbmNlUHJvcGVydGllc1Jlc3VsdBJJChJpbnN0YW5jZXByb3BlcnRpZX'
        'MY4J+LpQEgAygLMhUuc3NtLkluc3RhbmNlUHJvcGVydHlSEmluc3RhbmNlcHJvcGVydGllcxIf'
        'CgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeInventoryDeletionsRequestDescriptor instead')
const DescribeInventoryDeletionsRequest$json = {
  '1': 'DescribeInventoryDeletionsRequest',
  '2': [
    {'1': 'deletionid', '3': 126693587, '4': 1, '5': 9, '10': 'deletionid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInventoryDeletionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInventoryDeletionsRequestDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZUludmVudG9yeURlbGV0aW9uc1JlcXVlc3QSIQoKZGVsZXRpb25pZBjT4bQ8IA'
        'EoCVIKZGVsZXRpb25pZBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCglu'
        'ZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeInventoryDeletionsResultDescriptor instead')
const DescribeInventoryDeletionsResult$json = {
  '1': 'DescribeInventoryDeletionsResult',
  '2': [
    {
      '1': 'inventorydeletions',
      '3': 492028263,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryDeletionStatusItem',
      '10': 'inventorydeletions'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeInventoryDeletionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeInventoryDeletionsResultDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZUludmVudG9yeURlbGV0aW9uc1Jlc3VsdBJUChJpbnZlbnRvcnlkZWxldGlvbn'
        'MY54LP6gEgAygLMiAuc3NtLkludmVudG9yeURlbGV0aW9uU3RhdHVzSXRlbVISaW52ZW50b3J5'
        'ZGVsZXRpb25zEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionTaskInvocationsRequestDescriptor instead')
const DescribeMaintenanceWindowExecutionTaskInvocationsRequest$json = {
  '1': 'DescribeMaintenanceWindowExecutionTaskInvocationsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionTaskInvocationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionTaskInvocationsRequestDescriptor =
    $convert.base64Decode(
        'CjhEZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uVGFza0ludm9jYXRpb25zUmVxdW'
        'VzdBI5CgdmaWx0ZXJzGO3N6lkgAygLMhwuc3NtLk1haW50ZW5hbmNlV2luZG93RmlsdGVyUgdm'
        'aWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbh'
        'j+hLpnIAEoCVIJbmV4dHRva2VuEhkKBnRhc2tpZBiikesIIAEoCVIGdGFza2lkEjAKEXdpbmRv'
        'd2V4ZWN1dGlvbmlkGInrp6oBIAEoCVIRd2luZG93ZXhlY3V0aW9uaWQ=');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionTaskInvocationsResultDescriptor instead')
const DescribeMaintenanceWindowExecutionTaskInvocationsResult$json = {
  '1': 'DescribeMaintenanceWindowExecutionTaskInvocationsResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowexecutiontaskinvocationidentities',
      '3': 497972341,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowExecutionTaskInvocationIdentity',
      '10': 'windowexecutiontaskinvocationidentities'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionTaskInvocationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionTaskInvocationsResultDescriptor =
    $convert.base64Decode(
        'CjdEZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uVGFza0ludm9jYXRpb25zUmVzdW'
        'x0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEpMBCid3aW5kb3dleGVjdXRpb250'
        'YXNraW52b2NhdGlvbmlkZW50aXRpZXMY9ei57QEgAygLMjUuc3NtLk1haW50ZW5hbmNlV2luZG'
        '93RXhlY3V0aW9uVGFza0ludm9jYXRpb25JZGVudGl0eVInd2luZG93ZXhlY3V0aW9udGFza2lu'
        'dm9jYXRpb25pZGVudGl0aWVz');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionTasksRequestDescriptor instead')
const DescribeMaintenanceWindowExecutionTasksRequest$json = {
  '1': 'DescribeMaintenanceWindowExecutionTasksRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionTasksRequestDescriptor =
    $convert.base64Decode(
        'Ci5EZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uVGFza3NSZXF1ZXN0EjkKB2ZpbH'
        'RlcnMY7c3qWSADKAsyHC5zc20uTWFpbnRlbmFuY2VXaW5kb3dGaWx0ZXJSB2ZpbHRlcnMSIgoK'
        'bWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUg'
        'luZXh0dG9rZW4SMAoRd2luZG93ZXhlY3V0aW9uaWQYieunqgEgASgJUhF3aW5kb3dleGVjdXRp'
        'b25pZA==');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionTasksResultDescriptor instead')
const DescribeMaintenanceWindowExecutionTasksResult$json = {
  '1': 'DescribeMaintenanceWindowExecutionTasksResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowexecutiontaskidentities',
      '3': 497824727,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowExecutionTaskIdentity',
      '10': 'windowexecutiontaskidentities'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionTasksResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionTasksResultDescriptor =
    $convert.base64Decode(
        'Ci1EZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uVGFza3NSZXN1bHQSHwoJbmV4dH'
        'Rva2VuGP6EumcgASgJUgluZXh0dG9rZW4SdQodd2luZG93ZXhlY3V0aW9udGFza2lkZW50aXRp'
        'ZXMY1+ew7QEgAygLMisuc3NtLk1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uVGFza0lkZW50aX'
        'R5Uh13aW5kb3dleGVjdXRpb250YXNraWRlbnRpdGllcw==');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionsRequestDescriptor instead')
const DescribeMaintenanceWindowExecutionsRequest$json = {
  '1': 'DescribeMaintenanceWindowExecutionsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionsRequestDescriptor =
    $convert.base64Decode(
        'CipEZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uc1JlcXVlc3QSOQoHZmlsdGVycx'
        'jtzepZIAMoCzIcLnNzbS5NYWludGVuYW5jZVdpbmRvd0ZpbHRlclIHZmlsdGVycxIiCgptYXhy'
        'ZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2luZG93aWQ=');

@$core.Deprecated(
    'Use describeMaintenanceWindowExecutionsResultDescriptor instead')
const DescribeMaintenanceWindowExecutionsResult$json = {
  '1': 'DescribeMaintenanceWindowExecutionsResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowexecutions',
      '3': 165671467,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowExecution',
      '10': 'windowexecutions'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowExecutionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowExecutionsResultDescriptor = $convert.base64Decode(
        'CilEZXNjcmliZU1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uc1Jlc3VsdBIfCgluZXh0dG9rZW'
        '4Y/oS6ZyABKAlSCW5leHR0b2tlbhJOChB3aW5kb3dleGVjdXRpb25zGKvk/04gAygLMh8uc3Nt'
        'Lk1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uUhB3aW5kb3dleGVjdXRpb25z');

@$core.Deprecated(
    'Use describeMaintenanceWindowScheduleRequestDescriptor instead')
const DescribeMaintenanceWindowScheduleRequest$json = {
  '1': 'DescribeMaintenanceWindowScheduleRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchOrchestratorFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowResourceType',
      '10': 'resourcetype'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DescribeMaintenanceWindowScheduleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowScheduleRequestDescriptor =
    $convert.base64Decode(
        'CihEZXNjcmliZU1haW50ZW5hbmNlV2luZG93U2NoZWR1bGVSZXF1ZXN0EjkKB2ZpbHRlcnMY7c'
        '3qWSADKAsyHC5zc20uUGF0Y2hPcmNoZXN0cmF0b3JGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVz'
        'dWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4SSgoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoDjIiLnNzbS5NYWludGVuYW5jZVdpbmRvd1Jl'
        'c291cmNlVHlwZVIMcmVzb3VyY2V0eXBlEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2'
        'V0Ugd0YXJnZXRzEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZA==');

@$core
    .Deprecated('Use describeMaintenanceWindowScheduleResultDescriptor instead')
const DescribeMaintenanceWindowScheduleResult$json = {
  '1': 'DescribeMaintenanceWindowScheduleResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'scheduledwindowexecutions',
      '3': 493214756,
      '4': 3,
      '5': 11,
      '6': '.ssm.ScheduledWindowExecution',
      '10': 'scheduledwindowexecutions'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowScheduleResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowScheduleResultDescriptor =
    $convert.base64Decode(
        'CidEZXNjcmliZU1haW50ZW5hbmNlV2luZG93U2NoZWR1bGVSZXN1bHQSHwoJbmV4dHRva2VuGP'
        '6EumcgASgJUgluZXh0dG9rZW4SXwoZc2NoZWR1bGVkd2luZG93ZXhlY3V0aW9ucxikuJfrASAD'
        'KAsyHS5zc20uU2NoZWR1bGVkV2luZG93RXhlY3V0aW9uUhlzY2hlZHVsZWR3aW5kb3dleGVjdX'
        'Rpb25z');

@$core
    .Deprecated('Use describeMaintenanceWindowTargetsRequestDescriptor instead')
const DescribeMaintenanceWindowTargetsRequest$json = {
  '1': 'DescribeMaintenanceWindowTargetsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DescribeMaintenanceWindowTargetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowTargetsRequestDescriptor =
    $convert.base64Decode(
        'CidEZXNjcmliZU1haW50ZW5hbmNlV2luZG93VGFyZ2V0c1JlcXVlc3QSOQoHZmlsdGVycxjtze'
        'pZIAMoCzIcLnNzbS5NYWludGVuYW5jZVdpbmRvd0ZpbHRlclIHZmlsdGVycxIiCgptYXhyZXN1'
        'bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2'
        'tlbhIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2luZG93aWQ=');

@$core
    .Deprecated('Use describeMaintenanceWindowTargetsResultDescriptor instead')
const DescribeMaintenanceWindowTargetsResult$json = {
  '1': 'DescribeMaintenanceWindowTargetsResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTarget',
      '10': 'targets'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowTargetsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowTargetsResultDescriptor =
    $convert.base64Decode(
        'CiZEZXNjcmliZU1haW50ZW5hbmNlV2luZG93VGFyZ2V0c1Jlc3VsdBIfCgluZXh0dG9rZW4Y/o'
        'S6ZyABKAlSCW5leHR0b2tlbhI5Cgd0YXJnZXRzGIKbgn0gAygLMhwuc3NtLk1haW50ZW5hbmNl'
        'V2luZG93VGFyZ2V0Ugd0YXJnZXRz');

@$core.Deprecated('Use describeMaintenanceWindowTasksRequestDescriptor instead')
const DescribeMaintenanceWindowTasksRequest$json = {
  '1': 'DescribeMaintenanceWindowTasksRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `DescribeMaintenanceWindowTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowTasksRequestDescriptor =
    $convert.base64Decode(
        'CiVEZXNjcmliZU1haW50ZW5hbmNlV2luZG93VGFza3NSZXF1ZXN0EjkKB2ZpbHRlcnMY7c3qWS'
        'ADKAsyHC5zc20uTWFpbnRlbmFuY2VXaW5kb3dGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0'
        'cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
        '4SHQoId2luZG93aWQYqeSHCSABKAlSCHdpbmRvd2lk');

@$core.Deprecated('Use describeMaintenanceWindowTasksResultDescriptor instead')
const DescribeMaintenanceWindowTasksResult$json = {
  '1': 'DescribeMaintenanceWindowTasksResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tasks',
      '3': 429824110,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTask',
      '10': 'tasks'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowTasksResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowTasksResultDescriptor =
    $convert.base64Decode(
        'CiREZXNjcmliZU1haW50ZW5hbmNlV2luZG93VGFza3NSZXN1bHQSHwoJbmV4dHRva2VuGP6Eum'
        'cgASgJUgluZXh0dG9rZW4SNAoFdGFza3MY7rD6zAEgAygLMhouc3NtLk1haW50ZW5hbmNlV2lu'
        'ZG93VGFza1IFdGFza3M=');

@$core.Deprecated(
    'Use describeMaintenanceWindowsForTargetRequestDescriptor instead')
const DescribeMaintenanceWindowsForTargetRequest$json = {
  '1': 'DescribeMaintenanceWindowsForTargetRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowResourceType',
      '10': 'resourcetype'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowsForTargetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowsForTargetRequestDescriptor =
    $convert.base64Decode(
        'CipEZXNjcmliZU1haW50ZW5hbmNlV2luZG93c0ZvclRhcmdldFJlcXVlc3QSIgoKbWF4cmVzdW'
        'x0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9r'
        'ZW4SSgoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoDjIiLnNzbS5NYWludGVuYW5jZVdpbmRvd1Jlc2'
        '91cmNlVHlwZVIMcmVzb3VyY2V0eXBlEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0'
        'Ugd0YXJnZXRz');

@$core.Deprecated(
    'Use describeMaintenanceWindowsForTargetResultDescriptor instead')
const DescribeMaintenanceWindowsForTargetResult$json = {
  '1': 'DescribeMaintenanceWindowsForTargetResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowidentities',
      '3': 223310036,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowIdentityForTarget',
      '10': 'windowidentities'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowsForTargetResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    describeMaintenanceWindowsForTargetResultDescriptor = $convert.base64Decode(
        'CilEZXNjcmliZU1haW50ZW5hbmNlV2luZG93c0ZvclRhcmdldFJlc3VsdBIfCgluZXh0dG9rZW'
        '4Y/oS6ZyABKAlSCW5leHR0b2tlbhJWChB3aW5kb3dpZGVudGl0aWVzGNThvWogAygLMicuc3Nt'
        'Lk1haW50ZW5hbmNlV2luZG93SWRlbnRpdHlGb3JUYXJnZXRSEHdpbmRvd2lkZW50aXRpZXM=');

@$core.Deprecated('Use describeMaintenanceWindowsRequestDescriptor instead')
const DescribeMaintenanceWindowsRequest$json = {
  '1': 'DescribeMaintenanceWindowsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribeMaintenanceWindowsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowsRequestDescriptor =
    $convert.base64Decode(
        'CiFEZXNjcmliZU1haW50ZW5hbmNlV2luZG93c1JlcXVlc3QSOQoHZmlsdGVycxjtzepZIAMoCz'
        'IcLnNzbS5NYWludGVuYW5jZVdpbmRvd0ZpbHRlclIHZmlsdGVycxIiCgptYXhyZXN1bHRzGLKo'
        'm4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use describeMaintenanceWindowsResultDescriptor instead')
const DescribeMaintenanceWindowsResult$json = {
  '1': 'DescribeMaintenanceWindowsResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'windowidentities',
      '3': 223310036,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowIdentity',
      '10': 'windowidentities'
    },
  ],
};

/// Descriptor for `DescribeMaintenanceWindowsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeMaintenanceWindowsResultDescriptor =
    $convert.base64Decode(
        'CiBEZXNjcmliZU1haW50ZW5hbmNlV2luZG93c1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhJNChB3aW5kb3dpZGVudGl0aWVzGNThvWogAygLMh4uc3NtLk1haW50ZW5h'
        'bmNlV2luZG93SWRlbnRpdHlSEHdpbmRvd2lkZW50aXRpZXM=');

@$core.Deprecated('Use describeOpsItemsRequestDescriptor instead')
const DescribeOpsItemsRequest$json = {
  '1': 'DescribeOpsItemsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'opsitemfilters',
      '3': 81740152,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemFilter',
      '10': 'opsitemfilters'
    },
  ],
};

/// Descriptor for `DescribeOpsItemsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeOpsItemsRequestDescriptor = $convert.base64Decode(
    'ChdEZXNjcmliZU9wc0l0ZW1zUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cm'
    'VzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhI9Cg5vcHNpdGVtZmlsdGVy'
    'cxj4gv0mIAMoCzISLnNzbS5PcHNJdGVtRmlsdGVyUg5vcHNpdGVtZmlsdGVycw==');

@$core.Deprecated('Use describeOpsItemsResponseDescriptor instead')
const DescribeOpsItemsResponse$json = {
  '1': 'DescribeOpsItemsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'opsitemsummaries',
      '3': 250098577,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemSummary',
      '10': 'opsitemsummaries'
    },
  ],
};

/// Descriptor for `DescribeOpsItemsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeOpsItemsResponseDescriptor = $convert.base64Decode(
    'ChhEZXNjcmliZU9wc0l0ZW1zUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
    '9rZW4SQgoQb3BzaXRlbXN1bW1hcmllcxiR56B3IAMoCzITLnNzbS5PcHNJdGVtU3VtbWFyeVIQ'
    'b3BzaXRlbXN1bW1hcmllcw==');

@$core.Deprecated('Use describeParametersRequestDescriptor instead')
const DescribeParametersRequest$json = {
  '1': 'DescribeParametersRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParametersFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameterfilters',
      '3': 511502980,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterStringFilter',
      '10': 'parameterfilters'
    },
    {'1': 'shared', '3': 185319157, '4': 1, '5': 8, '10': 'shared'},
  ],
};

/// Descriptor for `DescribeParametersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeParametersRequestDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZVBhcmFtZXRlcnNSZXF1ZXN0EjIKB2ZpbHRlcnMY7c3qWSADKAsyFS5zc20uUG'
    'FyYW1ldGVyc0ZpbHRlclIHZmlsdGVycxIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVz'
    'dWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhJKChBwYXJhbWV0ZXJmaWx0ZX'
    'JzGITV8/MBIAMoCzIaLnNzbS5QYXJhbWV0ZXJTdHJpbmdGaWx0ZXJSEHBhcmFtZXRlcmZpbHRl'
    'cnMSGQoGc2hhcmVkGPX9rlggASgIUgZzaGFyZWQ=');

@$core.Deprecated('Use describeParametersResultDescriptor instead')
const DescribeParametersResult$json = {
  '1': 'DescribeParametersResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterMetadata',
      '10': 'parameters'
    },
  ],
};

/// Descriptor for `DescribeParametersResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeParametersResultDescriptor = $convert.base64Decode(
    'ChhEZXNjcmliZVBhcmFtZXRlcnNSZXN1bHQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
    '9rZW4SOgoKcGFyYW1ldGVycxj6p/7rASADKAsyFi5zc20uUGFyYW1ldGVyTWV0YWRhdGFSCnBh'
    'cmFtZXRlcnM=');

@$core.Deprecated('Use describePatchBaselinesRequestDescriptor instead')
const DescribePatchBaselinesRequest$json = {
  '1': 'DescribePatchBaselinesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchOrchestratorFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribePatchBaselinesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchBaselinesRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZXNjcmliZVBhdGNoQmFzZWxpbmVzUmVxdWVzdBI5CgdmaWx0ZXJzGO3N6lkgAygLMhwuc3'
        'NtLlBhdGNoT3JjaGVzdHJhdG9yRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEg'
        'ASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describePatchBaselinesResultDescriptor instead')
const DescribePatchBaselinesResult$json = {
  '1': 'DescribePatchBaselinesResult',
  '2': [
    {
      '1': 'baselineidentities',
      '3': 340730797,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchBaselineIdentity',
      '10': 'baselineidentities'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribePatchBaselinesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchBaselinesResultDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZVBhdGNoQmFzZWxpbmVzUmVzdWx0Ek4KEmJhc2VsaW5laWRlbnRpdGllcxitx7'
        'yiASADKAsyGi5zc20uUGF0Y2hCYXNlbGluZUlkZW50aXR5UhJiYXNlbGluZWlkZW50aXRpZXMS'
        'HwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describePatchGroupStateRequestDescriptor instead')
const DescribePatchGroupStateRequest$json = {
  '1': 'DescribePatchGroupStateRequest',
  '2': [
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `DescribePatchGroupStateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchGroupStateRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZVBhdGNoR3JvdXBTdGF0ZVJlcXVlc3QSIgoKcGF0Y2hncm91cBjht7H3ASABKA'
        'lSCnBhdGNoZ3JvdXA=');

@$core.Deprecated('Use describePatchGroupStateResultDescriptor instead')
const DescribePatchGroupStateResult$json = {
  '1': 'DescribePatchGroupStateResult',
  '2': [
    {'1': 'instances', '3': 181818068, '4': 1, '5': 5, '10': 'instances'},
    {
      '1': 'instanceswithavailablesecurityupdates',
      '3': 120023659,
      '4': 1,
      '5': 5,
      '10': 'instanceswithavailablesecurityupdates'
    },
    {
      '1': 'instanceswithcriticalnoncompliantpatches',
      '3': 503903877,
      '4': 1,
      '5': 5,
      '10': 'instanceswithcriticalnoncompliantpatches'
    },
    {
      '1': 'instanceswithfailedpatches',
      '3': 279342697,
      '4': 1,
      '5': 5,
      '10': 'instanceswithfailedpatches'
    },
    {
      '1': 'instanceswithinstalledotherpatches',
      '3': 517595180,
      '4': 1,
      '5': 5,
      '10': 'instanceswithinstalledotherpatches'
    },
    {
      '1': 'instanceswithinstalledpatches',
      '3': 260114088,
      '4': 1,
      '5': 5,
      '10': 'instanceswithinstalledpatches'
    },
    {
      '1': 'instanceswithinstalledpendingrebootpatches',
      '3': 366426344,
      '4': 1,
      '5': 5,
      '10': 'instanceswithinstalledpendingrebootpatches'
    },
    {
      '1': 'instanceswithinstalledrejectedpatches',
      '3': 278897530,
      '4': 1,
      '5': 5,
      '10': 'instanceswithinstalledrejectedpatches'
    },
    {
      '1': 'instanceswithmissingpatches',
      '3': 63879570,
      '4': 1,
      '5': 5,
      '10': 'instanceswithmissingpatches'
    },
    {
      '1': 'instanceswithnotapplicablepatches',
      '3': 391027218,
      '4': 1,
      '5': 5,
      '10': 'instanceswithnotapplicablepatches'
    },
    {
      '1': 'instanceswithothernoncompliantpatches',
      '3': 315419780,
      '4': 1,
      '5': 5,
      '10': 'instanceswithothernoncompliantpatches'
    },
    {
      '1': 'instanceswithsecuritynoncompliantpatches',
      '3': 136821766,
      '4': 1,
      '5': 5,
      '10': 'instanceswithsecuritynoncompliantpatches'
    },
    {
      '1': 'instanceswithunreportednotapplicablepatches',
      '3': 79380150,
      '4': 1,
      '5': 5,
      '10': 'instanceswithunreportednotapplicablepatches'
    },
  ],
};

/// Descriptor for `DescribePatchGroupStateResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchGroupStateResultDescriptor = $convert.base64Decode(
    'Ch1EZXNjcmliZVBhdGNoR3JvdXBTdGF0ZVJlc3VsdBIfCglpbnN0YW5jZXMY1KXZViABKAVSCW'
    'luc3RhbmNlcxJXCiVpbnN0YW5jZXN3aXRoYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVzGOvUnTkg'
    'ASgFUiVpbnN0YW5jZXN3aXRoYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVzEl4KKGluc3RhbmNlc3'
    'dpdGhjcml0aWNhbG5vbmNvbXBsaWFudHBhdGNoZXMYhe2j8AEgASgFUihpbnN0YW5jZXN3aXRo'
    'Y3JpdGljYWxub25jb21wbGlhbnRwYXRjaGVzEkIKGmluc3RhbmNlc3dpdGhmYWlsZWRwYXRjaG'
    'VzGOncmYUBIAEoBVIaaW5zdGFuY2Vzd2l0aGZhaWxlZHBhdGNoZXMSUgoiaW5zdGFuY2Vzd2l0'
    'aGluc3RhbGxlZG90aGVycGF0Y2hlcxiswOf2ASABKAVSImluc3RhbmNlc3dpdGhpbnN0YWxsZW'
    'RvdGhlcnBhdGNoZXMSRwodaW5zdGFuY2Vzd2l0aGluc3RhbGxlZHBhdGNoZXMYqI2EfCABKAVS'
    'HWluc3RhbmNlc3dpdGhpbnN0YWxsZWRwYXRjaGVzEmIKKmluc3RhbmNlc3dpdGhpbnN0YWxsZW'
    'RwZW5kaW5ncmVib290cGF0Y2hlcxjo8dyuASABKAVSKmluc3RhbmNlc3dpdGhpbnN0YWxsZWRw'
    'ZW5kaW5ncmVib290cGF0Y2hlcxJYCiVpbnN0YW5jZXN3aXRoaW5zdGFsbGVkcmVqZWN0ZWRwYX'
    'RjaGVzGPrG/oQBIAEoBVIlaW5zdGFuY2Vzd2l0aGluc3RhbGxlZHJlamVjdGVkcGF0Y2hlcxJD'
    'ChtpbnN0YW5jZXN3aXRobWlzc2luZ3BhdGNoZXMYkvO6HiABKAVSG2luc3RhbmNlc3dpdGhtaX'
    'NzaW5ncGF0Y2hlcxJQCiFpbnN0YW5jZXN3aXRobm90YXBwbGljYWJsZXBhdGNoZXMYkrS6ugEg'
    'ASgFUiFpbnN0YW5jZXN3aXRobm90YXBwbGljYWJsZXBhdGNoZXMSWAolaW5zdGFuY2Vzd2l0aG'
    '90aGVybm9uY29tcGxpYW50cGF0Y2hlcxiE2bOWASABKAVSJWluc3RhbmNlc3dpdGhvdGhlcm5v'
    'bmNvbXBsaWFudHBhdGNoZXMSXQooaW5zdGFuY2Vzd2l0aHNlY3VyaXR5bm9uY29tcGxpYW50cG'
    'F0Y2hlcxiG+J5BIAEoBVIoaW5zdGFuY2Vzd2l0aHNlY3VyaXR5bm9uY29tcGxpYW50cGF0Y2hl'
    'cxJjCitpbnN0YW5jZXN3aXRodW5yZXBvcnRlZG5vdGFwcGxpY2FibGVwYXRjaGVzGLb97CUgAS'
    'gFUitpbnN0YW5jZXN3aXRodW5yZXBvcnRlZG5vdGFwcGxpY2FibGVwYXRjaGVz');

@$core.Deprecated('Use describePatchGroupsRequestDescriptor instead')
const DescribePatchGroupsRequest$json = {
  '1': 'DescribePatchGroupsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchOrchestratorFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribePatchGroupsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchGroupsRequestDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZVBhdGNoR3JvdXBzUmVxdWVzdBI5CgdmaWx0ZXJzGO3N6lkgAygLMhwuc3NtLl'
        'BhdGNoT3JjaGVzdHJhdG9yRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgF'
        'UgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use describePatchGroupsResultDescriptor instead')
const DescribePatchGroupsResult$json = {
  '1': 'DescribePatchGroupsResult',
  '2': [
    {
      '1': 'mappings',
      '3': 8955491,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchGroupPatchBaselineMapping',
      '10': 'mappings'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `DescribePatchGroupsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchGroupsResultDescriptor = $convert.base64Decode(
    'ChlEZXNjcmliZVBhdGNoR3JvdXBzUmVzdWx0EkIKCG1hcHBpbmdzGOPMogQgAygLMiMuc3NtLl'
    'BhdGNoR3JvdXBQYXRjaEJhc2VsaW5lTWFwcGluZ1IIbWFwcGluZ3MSHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use describePatchPropertiesRequestDescriptor instead')
const DescribePatchPropertiesRequest$json = {
  '1': 'DescribePatchPropertiesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {
      '1': 'patchset',
      '3': 254208192,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchSet',
      '10': 'patchset'
    },
    {
      '1': 'property',
      '3': 304216553,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchProperty',
      '10': 'property'
    },
  ],
};

/// Descriptor for `DescribePatchPropertiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchPropertiesRequestDescriptor = $convert.base64Decode(
    'Ch5EZXNjcmliZVBhdGNoUHJvcGVydGllc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKA'
    'VSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SQQoPb3BlcmF0'
    'aW5nc3lzdGVtGOr9wRIgASgOMhQuc3NtLk9wZXJhdGluZ1N5c3RlbVIPb3BlcmF0aW5nc3lzdG'
    'VtEiwKCHBhdGNoc2V0GMDRm3kgASgOMg0uc3NtLlBhdGNoU2V0UghwYXRjaHNldBIyCghwcm9w'
    'ZXJ0eRjp84eRASABKA4yEi5zc20uUGF0Y2hQcm9wZXJ0eVIIcHJvcGVydHk=');

@$core.Deprecated('Use describePatchPropertiesResultDescriptor instead')
const DescribePatchPropertiesResult$json = {
  '1': 'DescribePatchPropertiesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'properties',
      '3': 29886973,
      '4': 3,
      '5': 11,
      '6': '.ssm.DescribePatchPropertiesResult.PropertiesEntry',
      '10': 'properties'
    },
  ],
  '3': [DescribePatchPropertiesResult_PropertiesEntry$json],
};

@$core.Deprecated('Use describePatchPropertiesResultDescriptor instead')
const DescribePatchPropertiesResult_PropertiesEntry$json = {
  '1': 'PropertiesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DescribePatchPropertiesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describePatchPropertiesResultDescriptor = $convert.base64Decode(
    'Ch1EZXNjcmliZVBhdGNoUHJvcGVydGllc1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
    '5leHR0b2tlbhJVCgpwcm9wZXJ0aWVzGP2ToA4gAygLMjIuc3NtLkRlc2NyaWJlUGF0Y2hQcm9w'
    'ZXJ0aWVzUmVzdWx0LlByb3BlcnRpZXNFbnRyeVIKcHJvcGVydGllcxo9Cg9Qcm9wZXJ0aWVzRW'
    '50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use describeSessionsRequestDescriptor instead')
const DescribeSessionsRequest$json = {
  '1': 'DescribeSessionsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.SessionFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.ssm.SessionState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `DescribeSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSessionsRequestDescriptor = $convert.base64Decode(
    'ChdEZXNjcmliZVNlc3Npb25zUmVxdWVzdBIvCgdmaWx0ZXJzGO3N6lkgAygLMhIuc3NtLlNlc3'
    'Npb25GaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMS'
    'HwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SKwoFc3RhdGUYl8my7wEgASgOMhEuc3'
    'NtLlNlc3Npb25TdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use describeSessionsResponseDescriptor instead')
const DescribeSessionsResponse$json = {
  '1': 'DescribeSessionsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'sessions',
      '3': 379226861,
      '4': 3,
      '5': 11,
      '6': '.ssm.Session',
      '10': 'sessions'
    },
  ],
};

/// Descriptor for `DescribeSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeSessionsResponseDescriptor =
    $convert.base64Decode(
        'ChhEZXNjcmliZVNlc3Npb25zUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4SLAoIc2Vzc2lvbnMY7ZXqtAEgAygLMgwuc3NtLlNlc3Npb25SCHNlc3Npb25z');

@$core.Deprecated('Use disassociateOpsItemRelatedItemRequestDescriptor instead')
const DisassociateOpsItemRelatedItemRequest$json = {
  '1': 'DisassociateOpsItemRelatedItemRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `DisassociateOpsItemRelatedItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateOpsItemRelatedItemRequestDescriptor =
    $convert.base64Decode(
        'CiVEaXNhc3NvY2lhdGVPcHNJdGVtUmVsYXRlZEl0ZW1SZXF1ZXN0EicKDWFzc29jaWF0aW9uaW'
        'QYkvyVQiABKAlSDWFzc29jaWF0aW9uaWQSHwoJb3BzaXRlbWlkGNLSlQwgASgJUglvcHNpdGVt'
        'aWQ=');

@$core
    .Deprecated('Use disassociateOpsItemRelatedItemResponseDescriptor instead')
const DisassociateOpsItemRelatedItemResponse$json = {
  '1': 'DisassociateOpsItemRelatedItemResponse',
};

/// Descriptor for `DisassociateOpsItemRelatedItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disassociateOpsItemRelatedItemResponseDescriptor =
    $convert.base64Decode(
        'CiZEaXNhc3NvY2lhdGVPcHNJdGVtUmVsYXRlZEl0ZW1SZXNwb25zZQ==');

@$core.Deprecated('Use documentAlreadyExistsDescriptor instead')
const DocumentAlreadyExists$json = {
  '1': 'DocumentAlreadyExists',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DocumentAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentAlreadyExistsDescriptor = $convert.base64Decode(
    'ChVEb2N1bWVudEFscmVhZHlFeGlzdHMSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use documentDefaultVersionDescriptionDescriptor instead')
const DocumentDefaultVersionDescription$json = {
  '1': 'DocumentDefaultVersionDescription',
  '2': [
    {
      '1': 'defaultversion',
      '3': 161172483,
      '4': 1,
      '5': 9,
      '10': 'defaultversion'
    },
    {
      '1': 'defaultversionname',
      '3': 61928938,
      '4': 1,
      '5': 9,
      '10': 'defaultversionname'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DocumentDefaultVersionDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentDefaultVersionDescriptionDescriptor =
    $convert.base64Decode(
        'CiFEb2N1bWVudERlZmF1bHRWZXJzaW9uRGVzY3JpcHRpb24SKQoOZGVmYXVsdHZlcnNpb24Yg5'
        'jtTCABKAlSDmRlZmF1bHR2ZXJzaW9uEjEKEmRlZmF1bHR2ZXJzaW9ubmFtZRjq68MdIAEoCVIS'
        'ZGVmYXVsdHZlcnNpb25uYW1lEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use documentDescriptionDescriptor instead')
const DocumentDescription$json = {
  '1': 'DocumentDescription',
  '2': [
    {
      '1': 'approvedversion',
      '3': 457739495,
      '4': 1,
      '5': 9,
      '10': 'approvedversion'
    },
    {
      '1': 'attachmentsinformation',
      '3': 327375170,
      '4': 3,
      '5': 11,
      '6': '.ssm.AttachmentInformation',
      '10': 'attachmentsinformation'
    },
    {'1': 'author', '3': 361744247, '4': 1, '5': 9, '10': 'author'},
    {'1': 'category', '3': 263447954, '4': 3, '5': 9, '10': 'category'},
    {'1': 'categoryenum', '3': 42168029, '4': 3, '5': 9, '10': 'categoryenum'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {
      '1': 'defaultversion',
      '3': 161172483,
      '4': 1,
      '5': 9,
      '10': 'defaultversion'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documenttype',
      '3': 457084477,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentType',
      '10': 'documenttype'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'hash', '3': 250828530, '4': 1, '5': 9, '10': 'hash'},
    {
      '1': 'hashtype',
      '3': 172838330,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentHashType',
      '10': 'hashtype'
    },
    {
      '1': 'latestversion',
      '3': 424864587,
      '4': 1,
      '5': 9,
      '10': 'latestversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'owner', '3': 455261813, '4': 1, '5': 9, '10': 'owner'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentParameter',
      '10': 'parameters'
    },
    {
      '1': 'pendingreviewversion',
      '3': 160998637,
      '4': 1,
      '5': 9,
      '10': 'pendingreviewversion'
    },
    {
      '1': 'platformtypes',
      '3': 447196676,
      '4': 3,
      '5': 14,
      '6': '.ssm.PlatformType',
      '10': 'platformtypes'
    },
    {
      '1': 'requires',
      '3': 149214838,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentRequires',
      '10': 'requires'
    },
    {
      '1': 'reviewinformation',
      '3': 512648418,
      '4': 3,
      '5': 11,
      '6': '.ssm.ReviewInformation',
      '10': 'reviewinformation'
    },
    {
      '1': 'reviewstatus',
      '3': 34562404,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'reviewstatus'
    },
    {
      '1': 'schemaversion',
      '3': 371681851,
      '4': 1,
      '5': 9,
      '10': 'schemaversion'
    },
    {'1': 'sha1', '3': 462897589, '4': 1, '5': 9, '10': 'sha1'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentStatus',
      '10': 'status'
    },
    {
      '1': 'statusinformation',
      '3': 14795748,
      '4': 1,
      '5': 9,
      '10': 'statusinformation'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {'1': 'targettype', '3': 397256481, '4': 1, '5': 9, '10': 'targettype'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DocumentDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentDescriptionDescriptor = $convert.base64Decode(
    'ChNEb2N1bWVudERlc2NyaXB0aW9uEiwKD2FwcHJvdmVkdmVyc2lvbhjnmaLaASABKAlSD2FwcH'
    'JvdmVkdmVyc2lvbhJWChZhdHRhY2htZW50c2luZm9ybWF0aW9uGMKyjZwBIAMoCzIaLnNzbS5B'
    'dHRhY2htZW50SW5mb3JtYXRpb25SFmF0dGFjaG1lbnRzaW5mb3JtYXRpb24SGgoGYXV0aG9yGP'
    'eOv6wBIAEoCVIGYXV0aG9yEh0KCGNhdGVnb3J5GJLLz30gAygJUghjYXRlZ29yeRIlCgxjYXRl'
    'Z29yeWVudW0Y3d2NFCADKAlSDGNhdGVnb3J5ZW51bRIkCgtjcmVhdGVkZGF0ZRiwsOfGASABKA'
    'lSC2NyZWF0ZWRkYXRlEikKDmRlZmF1bHR2ZXJzaW9uGIOY7UwgASgJUg5kZWZhdWx0dmVyc2lv'
    'bhIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJAoLZGlzcGxheW5hbWUYt8'
    'myxwEgASgJUgtkaXNwbGF5bmFtZRI/Cg5kb2N1bWVudGZvcm1hdBiImb/2ASABKA4yEy5zc20u'
    'RG9jdW1lbnRGb3JtYXRSDmRvY3VtZW50Zm9ybWF0EjkKDGRvY3VtZW50dHlwZRi9nPrZASABKA'
    '4yES5zc20uRG9jdW1lbnRUeXBlUgxkb2N1bWVudHR5cGUSKwoPZG9jdW1lbnR2ZXJzaW9uGMnv'
    'qSggASgJUg9kb2N1bWVudHZlcnNpb24SFQoEaGFzaBjyrc13IAEoCVIEaGFzaBI0CghoYXNodH'
    'lwZRi6m7VSIAEoDjIVLnNzbS5Eb2N1bWVudEhhc2hUeXBlUghoYXNodHlwZRIoCg1sYXRlc3R2'
    'ZXJzaW9uGMvWy8oBIAEoCVINbGF0ZXN0dmVyc2lvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh'
    'gKBW93bmVyGPX8itkBIAEoCVIFb3duZXISOgoKcGFyYW1ldGVycxj6p/7rASADKAsyFi5zc20u'
    'RG9jdW1lbnRQYXJhbWV0ZXJSCnBhcmFtZXRlcnMSNQoUcGVuZGluZ3Jldmlld3ZlcnNpb24Y7c'
    'niTCABKAlSFHBlbmRpbmdyZXZpZXd2ZXJzaW9uEjsKDXBsYXRmb3JtdHlwZXMYhNye1QEgAygO'
    'MhEuc3NtLlBsYXRmb3JtVHlwZVINcGxhdGZvcm10eXBlcxI0CghyZXF1aXJlcxj2rJNHIAMoCz'
    'IVLnNzbS5Eb2N1bWVudFJlcXVpcmVzUghyZXF1aXJlcxJIChFyZXZpZXdpbmZvcm1hdGlvbhji'
    'ybn0ASADKAsyFi5zc20uUmV2aWV3SW5mb3JtYXRpb25SEXJldmlld2luZm9ybWF0aW9uEjgKDH'
    'Jldmlld3N0YXR1cxjkwr0QIAEoDjIRLnNzbS5SZXZpZXdTdGF0dXNSDHJldmlld3N0YXR1cxIo'
    'Cg1zY2hlbWF2ZXJzaW9uGLvUnbEBIAEoCVINc2NoZW1hdmVyc2lvbhIWCgRzaGExGLWD3dwBIA'
    'EoCVIEc2hhMRIuCgZzdGF0dXMYkOT7AiABKA4yEy5zc20uRG9jdW1lbnRTdGF0dXNSBnN0YXR1'
    'cxIvChFzdGF0dXNpbmZvcm1hdGlvbhjkh4cHIAEoCVIRc3RhdHVzaW5mb3JtYXRpb24SIAoEdG'
    'FncxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdzEiIKCnRhcmdldHR5cGUYoc62vQEgASgJUgp0'
    'YXJnZXR0eXBlEiMKC3ZlcnNpb25uYW1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use documentFilterDescriptor instead')
const DocumentFilter$json = {
  '1': 'DocumentFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFilterKey',
      '10': 'key'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `DocumentFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentFilterDescriptor = $convert.base64Decode(
    'Cg5Eb2N1bWVudEZpbHRlchIrCgNrZXkY7ZDXQCABKA4yFi5zc20uRG9jdW1lbnRGaWx0ZXJLZX'
    'lSA2tleRIXCgV2YWx1ZRjLp/sSIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use documentIdentifierDescriptor instead')
const DocumentIdentifier$json = {
  '1': 'DocumentIdentifier',
  '2': [
    {'1': 'author', '3': 361744247, '4': 1, '5': 9, '10': 'author'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documenttype',
      '3': 457084477,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentType',
      '10': 'documenttype'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'owner', '3': 455261813, '4': 1, '5': 9, '10': 'owner'},
    {
      '1': 'platformtypes',
      '3': 447196676,
      '4': 3,
      '5': 14,
      '6': '.ssm.PlatformType',
      '10': 'platformtypes'
    },
    {
      '1': 'requires',
      '3': 149214838,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentRequires',
      '10': 'requires'
    },
    {
      '1': 'reviewstatus',
      '3': 34562404,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'reviewstatus'
    },
    {
      '1': 'schemaversion',
      '3': 371681851,
      '4': 1,
      '5': 9,
      '10': 'schemaversion'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {'1': 'targettype', '3': 397256481, '4': 1, '5': 9, '10': 'targettype'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DocumentIdentifier`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentIdentifierDescriptor = $convert.base64Decode(
    'ChJEb2N1bWVudElkZW50aWZpZXISGgoGYXV0aG9yGPeOv6wBIAEoCVIGYXV0aG9yEiQKC2NyZW'
    'F0ZWRkYXRlGLCw58YBIAEoCVILY3JlYXRlZGRhdGUSJAoLZGlzcGxheW5hbWUYt8myxwEgASgJ'
    'UgtkaXNwbGF5bmFtZRI/Cg5kb2N1bWVudGZvcm1hdBiImb/2ASABKA4yEy5zc20uRG9jdW1lbn'
    'RGb3JtYXRSDmRvY3VtZW50Zm9ybWF0EjkKDGRvY3VtZW50dHlwZRi9nPrZASABKA4yES5zc20u'
    'RG9jdW1lbnRUeXBlUgxkb2N1bWVudHR5cGUSKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg'
    '9kb2N1bWVudHZlcnNpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIYCgVvd25lchj1/IrZASAB'
    'KAlSBW93bmVyEjsKDXBsYXRmb3JtdHlwZXMYhNye1QEgAygOMhEuc3NtLlBsYXRmb3JtVHlwZV'
    'INcGxhdGZvcm10eXBlcxI0CghyZXF1aXJlcxj2rJNHIAMoCzIVLnNzbS5Eb2N1bWVudFJlcXVp'
    'cmVzUghyZXF1aXJlcxI4CgxyZXZpZXdzdGF0dXMY5MK9ECABKA4yES5zc20uUmV2aWV3U3RhdH'
    'VzUgxyZXZpZXdzdGF0dXMSKAoNc2NoZW1hdmVyc2lvbhi71J2xASABKAlSDXNjaGVtYXZlcnNp'
    'b24SIAoEdGFncxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdzEiIKCnRhcmdldHR5cGUYoc62vQ'
    'EgASgJUgp0YXJnZXR0eXBlEiMKC3ZlcnNpb25uYW1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use documentKeyValuesFilterDescriptor instead')
const DocumentKeyValuesFilter$json = {
  '1': 'DocumentKeyValuesFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `DocumentKeyValuesFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentKeyValuesFilterDescriptor =
    $convert.base64Decode(
        'ChdEb2N1bWVudEtleVZhbHVlc0ZpbHRlchITCgNrZXkYjZLraCABKAlSA2tleRIZCgZ2YWx1ZX'
        'MY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use documentLimitExceededDescriptor instead')
const DocumentLimitExceeded$json = {
  '1': 'DocumentLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DocumentLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentLimitExceededDescriptor = $convert.base64Decode(
    'ChVEb2N1bWVudExpbWl0RXhjZWVkZWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use documentMetadataResponseInfoDescriptor instead')
const DocumentMetadataResponseInfo$json = {
  '1': 'DocumentMetadataResponseInfo',
  '2': [
    {
      '1': 'reviewerresponse',
      '3': 336575822,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentReviewerResponseSource',
      '10': 'reviewerresponse'
    },
  ],
};

/// Descriptor for `DocumentMetadataResponseInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentMetadataResponseInfoDescriptor =
    $convert.base64Decode(
        'ChxEb2N1bWVudE1ldGFkYXRhUmVzcG9uc2VJbmZvElMKEHJldmlld2VycmVzcG9uc2UYzvq+oA'
        'EgAygLMiMuc3NtLkRvY3VtZW50UmV2aWV3ZXJSZXNwb25zZVNvdXJjZVIQcmV2aWV3ZXJyZXNw'
        'b25zZQ==');

@$core.Deprecated('Use documentParameterDescriptor instead')
const DocumentParameter$json = {
  '1': 'DocumentParameter',
  '2': [
    {'1': 'defaultvalue', '3': 218709920, '4': 1, '5': 9, '10': 'defaultvalue'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentParameterType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `DocumentParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentParameterDescriptor = $convert.base64Decode(
    'ChFEb2N1bWVudFBhcmFtZXRlchIlCgxkZWZhdWx0dmFsdWUYoP+kaCABKAlSDGRlZmF1bHR2YW'
    'x1ZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SFQoEbmFtZRiH5oF/IAEo'
    'CVIEbmFtZRIyCgR0eXBlGO6g14oBIAEoDjIaLnNzbS5Eb2N1bWVudFBhcmFtZXRlclR5cGVSBH'
    'R5cGU=');

@$core.Deprecated('Use documentPermissionLimitDescriptor instead')
const DocumentPermissionLimit$json = {
  '1': 'DocumentPermissionLimit',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DocumentPermissionLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentPermissionLimitDescriptor =
    $convert.base64Decode(
        'ChdEb2N1bWVudFBlcm1pc3Npb25MaW1pdBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use documentRequiresDescriptor instead')
const DocumentRequires$json = {
  '1': 'DocumentRequires',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'requiretype', '3': 56574093, '4': 1, '5': 9, '10': 'requiretype'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DocumentRequires`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentRequiresDescriptor = $convert.base64Decode(
    'ChBEb2N1bWVudFJlcXVpcmVzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSIwoLcmVxdWlyZXR5cG'
    'UYjYH9GiABKAlSC3JlcXVpcmV0eXBlEhwKB3ZlcnNpb24YuKq37gEgASgJUgd2ZXJzaW9uEiMK'
    'C3ZlcnNpb25uYW1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use documentReviewCommentSourceDescriptor instead')
const DocumentReviewCommentSource$json = {
  '1': 'DocumentReviewCommentSource',
  '2': [
    {'1': 'content', '3': 23568227, '4': 1, '5': 9, '10': 'content'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentReviewCommentType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `DocumentReviewCommentSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentReviewCommentSourceDescriptor =
    $convert.base64Decode(
        'ChtEb2N1bWVudFJldmlld0NvbW1lbnRTb3VyY2USGwoHY29udGVudBjjvp4LIAEoCVIHY29udG'
        'VudBI2CgR0eXBlGO6g14oBIAEoDjIeLnNzbS5Eb2N1bWVudFJldmlld0NvbW1lbnRUeXBlUgR0'
        'eXBl');

@$core.Deprecated('Use documentReviewerResponseSourceDescriptor instead')
const DocumentReviewerResponseSource$json = {
  '1': 'DocumentReviewerResponseSource',
  '2': [
    {
      '1': 'comment',
      '3': 407871487,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentReviewCommentSource',
      '10': 'comment'
    },
    {'1': 'createtime', '3': 490895933, '4': 1, '5': 9, '10': 'createtime'},
    {
      '1': 'reviewstatus',
      '3': 34562404,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'reviewstatus'
    },
    {'1': 'reviewer', '3': 436444219, '4': 1, '5': 9, '10': 'reviewer'},
    {'1': 'updatedtime', '3': 156274570, '4': 1, '5': 9, '10': 'updatedtime'},
  ],
};

/// Descriptor for `DocumentReviewerResponseSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentReviewerResponseSourceDescriptor = $convert.base64Decode(
    'Ch5Eb2N1bWVudFJldmlld2VyUmVzcG9uc2VTb3VyY2USPgoHY29tbWVudBj/v77CASADKAsyIC'
    '5zc20uRG9jdW1lbnRSZXZpZXdDb21tZW50U291cmNlUgdjb21tZW50EiIKCmNyZWF0ZXRpbWUY'
    'vfSJ6gEgASgJUgpjcmVhdGV0aW1lEjgKDHJldmlld3N0YXR1cxjkwr0QIAEoDjIRLnNzbS5SZX'
    'ZpZXdTdGF0dXNSDHJldmlld3N0YXR1cxIeCghyZXZpZXdlchi7uI7QASABKAlSCHJldmlld2Vy'
    'EiMKC3VwZGF0ZWR0aW1lGIqfwkogASgJUgt1cGRhdGVkdGltZQ==');

@$core.Deprecated('Use documentReviewsDescriptor instead')
const DocumentReviews$json = {
  '1': 'DocumentReviews',
  '2': [
    {
      '1': 'action',
      '3': 175614240,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentReviewAction',
      '10': 'action'
    },
    {
      '1': 'comment',
      '3': 407871487,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentReviewCommentSource',
      '10': 'comment'
    },
  ],
};

/// Descriptor for `DocumentReviews`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentReviewsDescriptor = $convert.base64Decode(
    'Cg9Eb2N1bWVudFJldmlld3MSNAoGYWN0aW9uGKDS3lMgASgOMhkuc3NtLkRvY3VtZW50UmV2aW'
    'V3QWN0aW9uUgZhY3Rpb24SPgoHY29tbWVudBj/v77CASADKAsyIC5zc20uRG9jdW1lbnRSZXZp'
    'ZXdDb21tZW50U291cmNlUgdjb21tZW50');

@$core.Deprecated('Use documentVersionInfoDescriptor instead')
const DocumentVersionInfo$json = {
  '1': 'DocumentVersionInfo',
  '2': [
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'isdefaultversion',
      '3': 465655635,
      '4': 1,
      '5': 8,
      '10': 'isdefaultversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'reviewstatus',
      '3': 34562404,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'reviewstatus'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentStatus',
      '10': 'status'
    },
    {
      '1': 'statusinformation',
      '3': 14795748,
      '4': 1,
      '5': 9,
      '10': 'statusinformation'
    },
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `DocumentVersionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentVersionInfoDescriptor = $convert.base64Decode(
    'ChNEb2N1bWVudFZlcnNpb25JbmZvEiQKC2NyZWF0ZWRkYXRlGLCw58YBIAEoCVILY3JlYXRlZG'
    'RhdGUSJAoLZGlzcGxheW5hbWUYt8myxwEgASgJUgtkaXNwbGF5bmFtZRI/Cg5kb2N1bWVudGZv'
    'cm1hdBiImb/2ASABKA4yEy5zc20uRG9jdW1lbnRGb3JtYXRSDmRvY3VtZW50Zm9ybWF0EisKD2'
    'RvY3VtZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEi4KEGlzZGVmYXVsdHZl'
    'cnNpb24Y066F3gEgASgIUhBpc2RlZmF1bHR2ZXJzaW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USOAoMcmV2aWV3c3RhdHVzGOTCvRAgASgOMhEuc3NtLlJldmlld1N0YXR1c1IMcmV2aWV3c3Rh'
    'dHVzEi4KBnN0YXR1cxiQ5PsCIAEoDjITLnNzbS5Eb2N1bWVudFN0YXR1c1IGc3RhdHVzEi8KEX'
    'N0YXR1c2luZm9ybWF0aW9uGOSHhwcgASgJUhFzdGF0dXNpbmZvcm1hdGlvbhIjCgt2ZXJzaW9u'
    'bmFtZRjVo7RsIAEoCVILdmVyc2lvbm5hbWU=');

@$core.Deprecated('Use documentVersionLimitExceededDescriptor instead')
const DocumentVersionLimitExceeded$json = {
  '1': 'DocumentVersionLimitExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DocumentVersionLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List documentVersionLimitExceededDescriptor =
    $convert.base64Decode(
        'ChxEb2N1bWVudFZlcnNpb25MaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use doesNotExistExceptionDescriptor instead')
const DoesNotExistException$json = {
  '1': 'DoesNotExistException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DoesNotExistException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List doesNotExistExceptionDescriptor = $convert.base64Decode(
    'ChVEb2VzTm90RXhpc3RFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use duplicateDocumentContentDescriptor instead')
const DuplicateDocumentContent$json = {
  '1': 'DuplicateDocumentContent',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateDocumentContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateDocumentContentDescriptor =
    $convert.base64Decode(
        'ChhEdXBsaWNhdGVEb2N1bWVudENvbnRlbnQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use duplicateDocumentVersionNameDescriptor instead')
const DuplicateDocumentVersionName$json = {
  '1': 'DuplicateDocumentVersionName',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DuplicateDocumentVersionName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateDocumentVersionNameDescriptor =
    $convert.base64Decode(
        'ChxEdXBsaWNhdGVEb2N1bWVudFZlcnNpb25OYW1lEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use duplicateInstanceIdDescriptor instead')
const DuplicateInstanceId$json = {
  '1': 'DuplicateInstanceId',
};

/// Descriptor for `DuplicateInstanceId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List duplicateInstanceIdDescriptor =
    $convert.base64Decode('ChNEdXBsaWNhdGVJbnN0YW5jZUlk');

@$core.Deprecated('Use effectivePatchDescriptor instead')
const EffectivePatch$json = {
  '1': 'EffectivePatch',
  '2': [
    {
      '1': 'patch',
      '3': 185240906,
      '4': 1,
      '5': 11,
      '6': '.ssm.Patch',
      '10': 'patch'
    },
    {
      '1': 'patchstatus',
      '3': 454549580,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchStatus',
      '10': 'patchstatus'
    },
  ],
};

/// Descriptor for `EffectivePatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List effectivePatchDescriptor = $convert.base64Decode(
    'Cg5FZmZlY3RpdmVQYXRjaBIjCgVwYXRjaBjKmqpYIAEoCzIKLnNzbS5QYXRjaFIFcGF0Y2gSNg'
    'oLcGF0Y2hzdGF0dXMYzMDf2AEgASgLMhAuc3NtLlBhdGNoU3RhdHVzUgtwYXRjaHN0YXR1cw==');

@$core.Deprecated('Use executionInputsDescriptor instead')
const ExecutionInputs$json = {
  '1': 'ExecutionInputs',
  '2': [
    {
      '1': 'automation',
      '3': 73629771,
      '4': 1,
      '5': 11,
      '6': '.ssm.AutomationExecutionInputs',
      '10': 'automation'
    },
  ],
};

/// Descriptor for `ExecutionInputs`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionInputsDescriptor = $convert.base64Decode(
    'Cg9FeGVjdXRpb25JbnB1dHMSQQoKYXV0b21hdGlvbhjLgI4jIAEoCzIeLnNzbS5BdXRvbWF0aW'
    '9uRXhlY3V0aW9uSW5wdXRzUgphdXRvbWF0aW9u');

@$core.Deprecated('Use executionPreviewDescriptor instead')
const ExecutionPreview$json = {
  '1': 'ExecutionPreview',
  '2': [
    {
      '1': 'automation',
      '3': 73629771,
      '4': 1,
      '5': 11,
      '6': '.ssm.AutomationExecutionPreview',
      '10': 'automation'
    },
  ],
};

/// Descriptor for `ExecutionPreview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List executionPreviewDescriptor = $convert.base64Decode(
    'ChBFeGVjdXRpb25QcmV2aWV3EkIKCmF1dG9tYXRpb24Yy4COIyABKAsyHy5zc20uQXV0b21hdG'
    'lvbkV4ZWN1dGlvblByZXZpZXdSCmF1dG9tYXRpb24=');

@$core.Deprecated('Use failedCreateAssociationDescriptor instead')
const FailedCreateAssociation$json = {
  '1': 'FailedCreateAssociation',
  '2': [
    {
      '1': 'entry',
      '3': 482340148,
      '4': 1,
      '5': 11,
      '6': '.ssm.CreateAssociationBatchRequestEntry',
      '10': 'entry'
    },
    {
      '1': 'fault',
      '3': 460759072,
      '4': 1,
      '5': 14,
      '6': '.ssm.Fault',
      '10': 'fault'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FailedCreateAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List failedCreateAssociationDescriptor = $convert.base64Decode(
    'ChdGYWlsZWRDcmVhdGVBc3NvY2lhdGlvbhJBCgVlbnRyeRi02v/lASABKAsyJy5zc20uQ3JlYX'
    'RlQXNzb2NpYXRpb25CYXRjaFJlcXVlc3RFbnRyeVIFZW50cnkSJAoFZmF1bHQYoMDa2wEgASgO'
    'Mgouc3NtLkZhdWx0UgVmYXVsdBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use failureDetailsDescriptor instead')
const FailureDetails$json = {
  '1': 'FailureDetails',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 3,
      '5': 11,
      '6': '.ssm.FailureDetails.DetailsEntry',
      '10': 'details'
    },
    {'1': 'failurestage', '3': 27770026, '4': 1, '5': 9, '10': 'failurestage'},
    {'1': 'failuretype', '3': 448893658, '4': 1, '5': 9, '10': 'failuretype'},
  ],
  '3': [FailureDetails_DetailsEntry$json],
};

@$core.Deprecated('Use failureDetailsDescriptor instead')
const FailureDetails_DetailsEntry$json = {
  '1': 'DetailsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `FailureDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List failureDetailsDescriptor = $convert.base64Decode(
    'Cg5GYWlsdXJlRGV0YWlscxI9CgdkZXRhaWxzGMaEiXYgAygLMiAuc3NtLkZhaWx1cmVEZXRhaW'
    'xzLkRldGFpbHNFbnRyeVIHZGV0YWlscxIlCgxmYWlsdXJlc3RhZ2UYqvmeDSABKAlSDGZhaWx1'
    'cmVzdGFnZRIkCgtmYWlsdXJldHlwZRjapYbWASABKAlSC2ZhaWx1cmV0eXBlGjoKDERldGFpbH'
    'NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use featureNotAvailableExceptionDescriptor instead')
const FeatureNotAvailableException$json = {
  '1': 'FeatureNotAvailableException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FeatureNotAvailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List featureNotAvailableExceptionDescriptor =
    $convert.base64Decode(
        'ChxGZWF0dXJlTm90QXZhaWxhYmxlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use getAccessTokenRequestDescriptor instead')
const GetAccessTokenRequest$json = {
  '1': 'GetAccessTokenRequest',
  '2': [
    {
      '1': 'accessrequestid',
      '3': 377097210,
      '4': 1,
      '5': 9,
      '10': 'accessrequestid'
    },
  ],
};

/// Descriptor for `GetAccessTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessTokenRequestDescriptor = $convert.base64Decode(
    'ChVHZXRBY2Nlc3NUb2tlblJlcXVlc3QSLAoPYWNjZXNzcmVxdWVzdGlkGPqX6LMBIAEoCVIPYW'
    'NjZXNzcmVxdWVzdGlk');

@$core.Deprecated('Use getAccessTokenResponseDescriptor instead')
const GetAccessTokenResponse$json = {
  '1': 'GetAccessTokenResponse',
  '2': [
    {
      '1': 'accessrequeststatus',
      '3': 425567007,
      '4': 1,
      '5': 14,
      '6': '.ssm.AccessRequestStatus',
      '10': 'accessrequeststatus'
    },
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.ssm.Credentials',
      '10': 'credentials'
    },
  ],
};

/// Descriptor for `GetAccessTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessTokenResponseDescriptor = $convert.base64Decode(
    'ChZHZXRBY2Nlc3NUb2tlblJlc3BvbnNlEk4KE2FjY2Vzc3JlcXVlc3RzdGF0dXMYn8b2ygEgAS'
    'gOMhguc3NtLkFjY2Vzc1JlcXVlc3RTdGF0dXNSE2FjY2Vzc3JlcXVlc3RzdGF0dXMSNgoLY3Jl'
    'ZGVudGlhbHMY8pqOtgEgASgLMhAuc3NtLkNyZWRlbnRpYWxzUgtjcmVkZW50aWFscw==');

@$core.Deprecated('Use getAutomationExecutionRequestDescriptor instead')
const GetAutomationExecutionRequest$json = {
  '1': 'GetAutomationExecutionRequest',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
  ],
};

/// Descriptor for `GetAutomationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAutomationExecutionRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRBdXRvbWF0aW9uRXhlY3V0aW9uUmVxdWVzdBI3ChVhdXRvbWF0aW9uZXhlY3V0aW9uaW'
        'QY5u/3BSABKAlSFWF1dG9tYXRpb25leGVjdXRpb25pZA==');

@$core.Deprecated('Use getAutomationExecutionResultDescriptor instead')
const GetAutomationExecutionResult$json = {
  '1': 'GetAutomationExecutionResult',
  '2': [
    {
      '1': 'automationexecution',
      '3': 127241049,
      '4': 1,
      '5': 11,
      '6': '.ssm.AutomationExecution',
      '10': 'automationexecution'
    },
  ],
};

/// Descriptor for `GetAutomationExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAutomationExecutionResultDescriptor =
    $convert.base64Decode(
        'ChxHZXRBdXRvbWF0aW9uRXhlY3V0aW9uUmVzdWx0Ek0KE2F1dG9tYXRpb25leGVjdXRpb24Y2Z'
        'bWPCABKAsyGC5zc20uQXV0b21hdGlvbkV4ZWN1dGlvblITYXV0b21hdGlvbmV4ZWN1dGlvbg==');

@$core.Deprecated('Use getCalendarStateRequestDescriptor instead')
const GetCalendarStateRequest$json = {
  '1': 'GetCalendarStateRequest',
  '2': [
    {'1': 'attime', '3': 404952130, '4': 1, '5': 9, '10': 'attime'},
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
  ],
};

/// Descriptor for `GetCalendarStateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalendarStateRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRDYWxlbmRhclN0YXRlUmVxdWVzdBIaCgZhdHRpbWUYwqiMwQEgASgJUgZhdHRpbWUSJw'
        'oNY2FsZW5kYXJuYW1lcxi+85kRIAMoCVINY2FsZW5kYXJuYW1lcw==');

@$core.Deprecated('Use getCalendarStateResponseDescriptor instead')
const GetCalendarStateResponse$json = {
  '1': 'GetCalendarStateResponse',
  '2': [
    {'1': 'attime', '3': 404952130, '4': 1, '5': 9, '10': 'attime'},
    {
      '1': 'nexttransitiontime',
      '3': 426659123,
      '4': 1,
      '5': 9,
      '10': 'nexttransitiontime'
    },
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.ssm.CalendarState',
      '10': 'state'
    },
  ],
};

/// Descriptor for `GetCalendarStateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalendarStateResponseDescriptor = $convert.base64Decode(
    'ChhHZXRDYWxlbmRhclN0YXRlUmVzcG9uc2USGgoGYXR0aW1lGMKojMEBIAEoCVIGYXR0aW1lEj'
    'IKEm5leHR0cmFuc2l0aW9udGltZRizmrnLASABKAlSEm5leHR0cmFuc2l0aW9udGltZRIsCgVz'
    'dGF0ZRiXybLvASABKA4yEi5zc20uQ2FsZW5kYXJTdGF0ZVIFc3RhdGU=');

@$core.Deprecated('Use getCommandInvocationRequestDescriptor instead')
const GetCommandInvocationRequest$json = {
  '1': 'GetCommandInvocationRequest',
  '2': [
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'pluginname', '3': 115989314, '4': 1, '5': 9, '10': 'pluginname'},
  ],
};

/// Descriptor for `GetCommandInvocationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCommandInvocationRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRDb21tYW5kSW52b2NhdGlvblJlcXVlc3QSHwoJY29tbWFuZGlkGIDbgEwgASgJUgljb2'
        '1tYW5kaWQSIQoKaW5zdGFuY2VpZBigrdEXIAEoCVIKaW5zdGFuY2VpZBIhCgpwbHVnaW5uYW1l'
        'GMK2pzcgASgJUgpwbHVnaW5uYW1l');

@$core.Deprecated('Use getCommandInvocationResultDescriptor instead')
const GetCommandInvocationResult$json = {
  '1': 'GetCommandInvocationResult',
  '2': [
    {
      '1': 'cloudwatchoutputconfig',
      '3': 21186555,
      '4': 1,
      '5': 11,
      '6': '.ssm.CloudWatchOutputConfig',
      '10': 'cloudwatchoutputconfig'
    },
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'executionelapsedtime',
      '3': 71300547,
      '4': 1,
      '5': 9,
      '10': 'executionelapsedtime'
    },
    {
      '1': 'executionenddatetime',
      '3': 521894840,
      '4': 1,
      '5': 9,
      '10': 'executionenddatetime'
    },
    {
      '1': 'executionstartdatetime',
      '3': 309261539,
      '4': 1,
      '5': 9,
      '10': 'executionstartdatetime'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'pluginname', '3': 115989314, '4': 1, '5': 9, '10': 'pluginname'},
    {'1': 'responsecode', '3': 447553700, '4': 1, '5': 5, '10': 'responsecode'},
    {
      '1': 'standarderrorcontent',
      '3': 473735696,
      '4': 1,
      '5': 9,
      '10': 'standarderrorcontent'
    },
    {
      '1': 'standarderrorurl',
      '3': 403407680,
      '4': 1,
      '5': 9,
      '10': 'standarderrorurl'
    },
    {
      '1': 'standardoutputcontent',
      '3': 294611267,
      '4': 1,
      '5': 9,
      '10': 'standardoutputcontent'
    },
    {
      '1': 'standardoutputurl',
      '3': 347642271,
      '4': 1,
      '5': 9,
      '10': 'standardoutputurl'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.CommandInvocationStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
  ],
};

/// Descriptor for `GetCommandInvocationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCommandInvocationResultDescriptor = $convert.base64Decode(
    'ChpHZXRDb21tYW5kSW52b2NhdGlvblJlc3VsdBJWChZjbG91ZHdhdGNob3V0cHV0Y29uZmlnGP'
    'uPjQogASgLMhsuc3NtLkNsb3VkV2F0Y2hPdXRwdXRDb25maWdSFmNsb3Vkd2F0Y2hvdXRwdXRj'
    'b25maWcSHwoJY29tbWFuZGlkGIDbgEwgASgJUgljb21tYW5kaWQSHAoHY29tbWVudBj/v77CAS'
    'ABKAlSB2NvbW1lbnQSJQoMZG9jdW1lbnRuYW1lGNCjxzkgASgJUgxkb2N1bWVudG5hbWUSKwoP'
    'ZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SNQoUZXhlY3V0aW9uZW'
    'xhcHNlZHRpbWUYw+v/ISABKAlSFGV4ZWN1dGlvbmVsYXBzZWR0aW1lEjYKFGV4ZWN1dGlvbmVu'
    'ZGRhdGV0aW1lGLj37fgBIAEoCVIUZXhlY3V0aW9uZW5kZGF0ZXRpbWUSOgoWZXhlY3V0aW9uc3'
    'RhcnRkYXRldGltZRjj6buTASABKAlSFmV4ZWN1dGlvbnN0YXJ0ZGF0ZXRpbWUSIQoKaW5zdGFu'
    'Y2VpZBigrdEXIAEoCVIKaW5zdGFuY2VpZBIhCgpwbHVnaW5uYW1lGMK2pzcgASgJUgpwbHVnaW'
    '5uYW1lEiYKDHJlc3BvbnNlY29kZRikwbTVASABKAVSDHJlc3BvbnNlY29kZRI2ChRzdGFuZGFy'
    'ZGVycm9yY29udGVudBiQxPLhASABKAlSFHN0YW5kYXJkZXJyb3Jjb250ZW50Ei4KEHN0YW5kYX'
    'JkZXJyb3J1cmwYwIauwAEgASgJUhBzdGFuZGFyZGVycm9ydXJsEjgKFXN0YW5kYXJkb3V0cHV0'
    'Y29udGVudBjD0r2MASABKAlSFXN0YW5kYXJkb3V0cHV0Y29udGVudBIwChFzdGFuZGFyZG91dH'
    'B1dHVybBifs+KlASABKAlSEXN0YW5kYXJkb3V0cHV0dXJsEjcKBnN0YXR1cxiQ5PsCIAEoDjIc'
    'LnNzbS5Db21tYW5kSW52b2NhdGlvblN0YXR1c1IGc3RhdHVzEigKDXN0YXR1c2RldGFpbHMYqJ'
    'LBsQEgASgJUg1zdGF0dXNkZXRhaWxz');

@$core.Deprecated('Use getConnectionStatusRequestDescriptor instead')
const GetConnectionStatusRequest$json = {
  '1': 'GetConnectionStatusRequest',
  '2': [
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
  ],
};

/// Descriptor for `GetConnectionStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionStatusRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDb25uZWN0aW9uU3RhdHVzUmVxdWVzdBIZCgZ0YXJnZXQY6eKfWyABKAlSBnRhcmdldA'
        '==');

@$core.Deprecated('Use getConnectionStatusResponseDescriptor instead')
const GetConnectionStatusResponse$json = {
  '1': 'GetConnectionStatusResponse',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ConnectionStatus',
      '10': 'status'
    },
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
  ],
};

/// Descriptor for `GetConnectionStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConnectionStatusResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRDb25uZWN0aW9uU3RhdHVzUmVzcG9uc2USMAoGc3RhdHVzGJDk+wIgASgOMhUuc3NtLk'
        'Nvbm5lY3Rpb25TdGF0dXNSBnN0YXR1cxIZCgZ0YXJnZXQY6eKfWyABKAlSBnRhcmdldA==');

@$core.Deprecated('Use getDefaultPatchBaselineRequestDescriptor instead')
const GetDefaultPatchBaselineRequest$json = {
  '1': 'GetDefaultPatchBaselineRequest',
  '2': [
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
  ],
};

/// Descriptor for `GetDefaultPatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDefaultPatchBaselineRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXREZWZhdWx0UGF0Y2hCYXNlbGluZVJlcXVlc3QSQQoPb3BlcmF0aW5nc3lzdGVtGOr9wR'
        'IgASgOMhQuc3NtLk9wZXJhdGluZ1N5c3RlbVIPb3BlcmF0aW5nc3lzdGVt');

@$core.Deprecated('Use getDefaultPatchBaselineResultDescriptor instead')
const GetDefaultPatchBaselineResult$json = {
  '1': 'GetDefaultPatchBaselineResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
  ],
};

/// Descriptor for `GetDefaultPatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDefaultPatchBaselineResultDescriptor =
    $convert.base64Decode(
        'Ch1HZXREZWZhdWx0UGF0Y2hCYXNlbGluZVJlc3VsdBIhCgpiYXNlbGluZWlkGNDk2yggASgJUg'
        'piYXNlbGluZWlkEkEKD29wZXJhdGluZ3N5c3RlbRjq/cESIAEoDjIULnNzbS5PcGVyYXRpbmdT'
        'eXN0ZW1SD29wZXJhdGluZ3N5c3RlbQ==');

@$core.Deprecated(
    'Use getDeployablePatchSnapshotForInstanceRequestDescriptor instead')
const GetDeployablePatchSnapshotForInstanceRequest$json = {
  '1': 'GetDeployablePatchSnapshotForInstanceRequest',
  '2': [
    {
      '1': 'baselineoverride',
      '3': 360801147,
      '4': 1,
      '5': 11,
      '6': '.ssm.BaselineOverride',
      '10': 'baselineoverride'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'snapshotid', '3': 99585817, '4': 1, '5': 9, '10': 'snapshotid'},
    {
      '1': 'uses3dualstackendpoint',
      '3': 398619818,
      '4': 1,
      '5': 8,
      '10': 'uses3dualstackendpoint'
    },
  ],
};

/// Descriptor for `GetDeployablePatchSnapshotForInstanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getDeployablePatchSnapshotForInstanceRequestDescriptor =
    $convert.base64Decode(
        'CixHZXREZXBsb3lhYmxlUGF0Y2hTbmFwc2hvdEZvckluc3RhbmNlUmVxdWVzdBJFChBiYXNlbG'
        'luZW92ZXJyaWRlGPvGhawBIAEoCzIVLnNzbS5CYXNlbGluZU92ZXJyaWRlUhBiYXNlbGluZW92'
        'ZXJyaWRlEiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQSIQoKc25hcHNob3RpZB'
        'iZnr4vIAEoCVIKc25hcHNob3RpZBI6ChZ1c2VzM2R1YWxzdGFja2VuZHBvaW50GKrpib4BIAEo'
        'CFIWdXNlczNkdWFsc3RhY2tlbmRwb2ludA==');

@$core.Deprecated(
    'Use getDeployablePatchSnapshotForInstanceResultDescriptor instead')
const GetDeployablePatchSnapshotForInstanceResult$json = {
  '1': 'GetDeployablePatchSnapshotForInstanceResult',
  '2': [
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'product', '3': 312238285, '4': 1, '5': 9, '10': 'product'},
    {
      '1': 'snapshotdownloadurl',
      '3': 419550867,
      '4': 1,
      '5': 9,
      '10': 'snapshotdownloadurl'
    },
    {'1': 'snapshotid', '3': 99585817, '4': 1, '5': 9, '10': 'snapshotid'},
  ],
};

/// Descriptor for `GetDeployablePatchSnapshotForInstanceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getDeployablePatchSnapshotForInstanceResultDescriptor =
    $convert.base64Decode(
        'CitHZXREZXBsb3lhYmxlUGF0Y2hTbmFwc2hvdEZvckluc3RhbmNlUmVzdWx0EiEKCmluc3Rhbm'
        'NlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQSHAoHcHJvZHVjdBjNwfGUASABKAlSB3Byb2R1Y3QS'
        'NAoTc25hcHNob3Rkb3dubG9hZHVybBiTrYfIASABKAlSE3NuYXBzaG90ZG93bmxvYWR1cmwSIQ'
        'oKc25hcHNob3RpZBiZnr4vIAEoCVIKc25hcHNob3RpZA==');

@$core.Deprecated('Use getDocumentRequestDescriptor instead')
const GetDocumentRequest$json = {
  '1': 'GetDocumentRequest',
  '2': [
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `GetDocumentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentRequestDescriptor = $convert.base64Decode(
    'ChJHZXREb2N1bWVudFJlcXVlc3QSPwoOZG9jdW1lbnRmb3JtYXQYiJm/9gEgASgOMhMuc3NtLk'
    'RvY3VtZW50Rm9ybWF0Ug5kb2N1bWVudGZvcm1hdBIrCg9kb2N1bWVudHZlcnNpb24Yye+pKCAB'
    'KAlSD2RvY3VtZW50dmVyc2lvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEiMKC3ZlcnNpb25uYW'
    '1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use getDocumentResultDescriptor instead')
const GetDocumentResult$json = {
  '1': 'GetDocumentResult',
  '2': [
    {
      '1': 'attachmentscontent',
      '3': 467878119,
      '4': 3,
      '5': 11,
      '6': '.ssm.AttachmentContent',
      '10': 'attachmentscontent'
    },
    {'1': 'content', '3': 23568227, '4': 1, '5': 9, '10': 'content'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documenttype',
      '3': 457084477,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentType',
      '10': 'documenttype'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'requires',
      '3': 149214838,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentRequires',
      '10': 'requires'
    },
    {
      '1': 'reviewstatus',
      '3': 34562404,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'reviewstatus'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentStatus',
      '10': 'status'
    },
    {
      '1': 'statusinformation',
      '3': 14795748,
      '4': 1,
      '5': 9,
      '10': 'statusinformation'
    },
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `GetDocumentResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDocumentResultDescriptor = $convert.base64Decode(
    'ChFHZXREb2N1bWVudFJlc3VsdBJKChJhdHRhY2htZW50c2NvbnRlbnQY54GN3wEgAygLMhYuc3'
    'NtLkF0dGFjaG1lbnRDb250ZW50UhJhdHRhY2htZW50c2NvbnRlbnQSGwoHY29udGVudBjjvp4L'
    'IAEoCVIHY29udGVudBIkCgtjcmVhdGVkZGF0ZRiwsOfGASABKAlSC2NyZWF0ZWRkYXRlEiQKC2'
    'Rpc3BsYXluYW1lGLfJsscBIAEoCVILZGlzcGxheW5hbWUSPwoOZG9jdW1lbnRmb3JtYXQYiJm/'
    '9gEgASgOMhMuc3NtLkRvY3VtZW50Rm9ybWF0Ug5kb2N1bWVudGZvcm1hdBI5Cgxkb2N1bWVudH'
    'R5cGUYvZz62QEgASgOMhEuc3NtLkRvY3VtZW50VHlwZVIMZG9jdW1lbnR0eXBlEisKD2RvY3Vt'
    'ZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEhUKBG5hbWUYh+aBfyABKAlSBG'
    '5hbWUSNAoIcmVxdWlyZXMY9qyTRyADKAsyFS5zc20uRG9jdW1lbnRSZXF1aXJlc1IIcmVxdWly'
    'ZXMSOAoMcmV2aWV3c3RhdHVzGOTCvRAgASgOMhEuc3NtLlJldmlld1N0YXR1c1IMcmV2aWV3c3'
    'RhdHVzEi4KBnN0YXR1cxiQ5PsCIAEoDjITLnNzbS5Eb2N1bWVudFN0YXR1c1IGc3RhdHVzEi8K'
    'EXN0YXR1c2luZm9ybWF0aW9uGOSHhwcgASgJUhFzdGF0dXNpbmZvcm1hdGlvbhIjCgt2ZXJzaW'
    '9ubmFtZRjVo7RsIAEoCVILdmVyc2lvbm5hbWU=');

@$core.Deprecated('Use getExecutionPreviewRequestDescriptor instead')
const GetExecutionPreviewRequest$json = {
  '1': 'GetExecutionPreviewRequest',
  '2': [
    {
      '1': 'executionpreviewid',
      '3': 35285163,
      '4': 1,
      '5': 9,
      '10': 'executionpreviewid'
    },
  ],
};

/// Descriptor for `GetExecutionPreviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExecutionPreviewRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRFeGVjdXRpb25QcmV2aWV3UmVxdWVzdBIxChJleGVjdXRpb25wcmV2aWV3aWQYq9HpEC'
        'ABKAlSEmV4ZWN1dGlvbnByZXZpZXdpZA==');

@$core.Deprecated('Use getExecutionPreviewResponseDescriptor instead')
const GetExecutionPreviewResponse$json = {
  '1': 'GetExecutionPreviewResponse',
  '2': [
    {'1': 'endedat', '3': 104122351, '4': 1, '5': 9, '10': 'endedat'},
    {
      '1': 'executionpreview',
      '3': 108797384,
      '4': 1,
      '5': 11,
      '6': '.ssm.ExecutionPreview',
      '10': 'executionpreview'
    },
    {
      '1': 'executionpreviewid',
      '3': 35285163,
      '4': 1,
      '5': 9,
      '10': 'executionpreviewid'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ExecutionPreviewStatus',
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

/// Descriptor for `GetExecutionPreviewResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExecutionPreviewResponseDescriptor = $convert.base64Decode(
    'ChtHZXRFeGVjdXRpb25QcmV2aWV3UmVzcG9uc2USGwoHZW5kZWRhdBjvj9MxIAEoCVIHZW5kZW'
    'RhdBJEChBleGVjdXRpb25wcmV2aWV3GMi78DMgASgLMhUuc3NtLkV4ZWN1dGlvblByZXZpZXdS'
    'EGV4ZWN1dGlvbnByZXZpZXcSMQoSZXhlY3V0aW9ucHJldmlld2lkGKvR6RAgASgJUhJleGVjdX'
    'Rpb25wcmV2aWV3aWQSNgoGc3RhdHVzGJDk+wIgASgOMhsuc3NtLkV4ZWN1dGlvblByZXZpZXdT'
    'dGF0dXNSBnN0YXR1cxInCg1zdGF0dXNtZXNzYWdlGI/GziIgASgJUg1zdGF0dXNtZXNzYWdl');

@$core.Deprecated('Use getInventoryRequestDescriptor instead')
const GetInventoryRequest$json = {
  '1': 'GetInventoryRequest',
  '2': [
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryAggregator',
      '10': 'aggregators'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resultattributes',
      '3': 430325814,
      '4': 3,
      '5': 11,
      '6': '.ssm.ResultAttribute',
      '10': 'resultattributes'
    },
  ],
};

/// Descriptor for `GetInventoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventoryRequestDescriptor = $convert.base64Decode(
    'ChNHZXRJbnZlbnRvcnlSZXF1ZXN0Ej0KC2FnZ3JlZ2F0b3JzGI79g00gAygLMhguc3NtLkludm'
    'VudG9yeUFnZ3JlZ2F0b3JSC2FnZ3JlZ2F0b3JzEjEKB2ZpbHRlcnMY7c3qWSADKAsyFC5zc20u'
    'SW52ZW50b3J5RmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZX'
    'N1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEkQKEHJlc3VsdGF0dHJpYnV0'
    'ZXMYtoCZzQEgAygLMhQuc3NtLlJlc3VsdEF0dHJpYnV0ZVIQcmVzdWx0YXR0cmlidXRlcw==');

@$core.Deprecated('Use getInventoryResultDescriptor instead')
const GetInventoryResult$json = {
  '1': 'GetInventoryResult',
  '2': [
    {
      '1': 'entities',
      '3': 406339595,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryResultEntity',
      '10': 'entities'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `GetInventoryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventoryResultDescriptor = $convert.base64Decode(
    'ChJHZXRJbnZlbnRvcnlSZXN1bHQSOgoIZW50aXRpZXMYi4DhwQEgAygLMhouc3NtLkludmVudG'
    '9yeVJlc3VsdEVudGl0eVIIZW50aXRpZXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9r'
    'ZW4=');

@$core.Deprecated('Use getInventorySchemaRequestDescriptor instead')
const GetInventorySchemaRequest$json = {
  '1': 'GetInventorySchemaRequest',
  '2': [
    {'1': 'aggregator', '3': 117674797, '4': 1, '5': 8, '10': 'aggregator'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'subtype', '3': 152988350, '4': 1, '5': 8, '10': 'subtype'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `GetInventorySchemaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventorySchemaRequestDescriptor = $convert.base64Decode(
    'ChlHZXRJbnZlbnRvcnlTY2hlbWFSZXF1ZXN0EiEKCmFnZ3JlZ2F0b3IYraaOOCABKAhSCmFnZ3'
    'JlZ2F0b3ISIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2Vu'
    'GP6EumcgASgJUgluZXh0dG9rZW4SGwoHc3VidHlwZRi+1flIIAEoCFIHc3VidHlwZRIeCgh0eX'
    'BlbmFtZRjPztnUASABKAlSCHR5cGVuYW1l');

@$core.Deprecated('Use getInventorySchemaResultDescriptor instead')
const GetInventorySchemaResult$json = {
  '1': 'GetInventorySchemaResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'schemas',
      '3': 276103488,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryItemSchema',
      '10': 'schemas'
    },
  ],
};

/// Descriptor for `GetInventorySchemaResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventorySchemaResultDescriptor = $convert.base64Decode(
    'ChhHZXRJbnZlbnRvcnlTY2hlbWFSZXN1bHQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
    '9rZW4SNgoHc2NoZW1hcxjAgtSDASADKAsyGC5zc20uSW52ZW50b3J5SXRlbVNjaGVtYVIHc2No'
    'ZW1hcw==');

@$core.Deprecated('Use getMaintenanceWindowExecutionRequestDescriptor instead')
const GetMaintenanceWindowExecutionRequest$json = {
  '1': 'GetMaintenanceWindowExecutionRequest',
  '2': [
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `GetMaintenanceWindowExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowExecutionRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblJlcXVlc3QSMAoRd2luZG93ZXhlY3V0aW'
        '9uaWQYieunqgEgASgJUhF3aW5kb3dleGVjdXRpb25pZA==');

@$core.Deprecated('Use getMaintenanceWindowExecutionResultDescriptor instead')
const GetMaintenanceWindowExecutionResult$json = {
  '1': 'GetMaintenanceWindowExecutionResult',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {'1': 'taskids', '3': 254419607, '4': 3, '5': 9, '10': 'taskids'},
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `GetMaintenanceWindowExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowExecutionResultDescriptor = $convert.base64Decode(
    'CiNHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblJlc3VsdBIbCgdlbmR0aW1lGMzvvB4gAS'
    'gJUgdlbmR0aW1lEiAKCXN0YXJ0dGltZRjvtOWwASABKAlSCXN0YXJ0dGltZRJACgZzdGF0dXMY'
    'kOT7AiABKA4yJS5zc20uTWFpbnRlbmFuY2VXaW5kb3dFeGVjdXRpb25TdGF0dXNSBnN0YXR1cx'
    'IoCg1zdGF0dXNkZXRhaWxzGKiSwbEBIAEoCVINc3RhdHVzZGV0YWlscxIbCgd0YXNraWRzGJfF'
    'qHkgAygJUgd0YXNraWRzEjAKEXdpbmRvd2V4ZWN1dGlvbmlkGInrp6oBIAEoCVIRd2luZG93ZX'
    'hlY3V0aW9uaWQ=');

@$core.Deprecated(
    'Use getMaintenanceWindowExecutionTaskInvocationRequestDescriptor instead')
const GetMaintenanceWindowExecutionTaskInvocationRequest$json = {
  '1': 'GetMaintenanceWindowExecutionTaskInvocationRequest',
  '2': [
    {'1': 'invocationid', '3': 116064639, '4': 1, '5': 9, '10': 'invocationid'},
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `GetMaintenanceWindowExecutionTaskInvocationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getMaintenanceWindowExecutionTaskInvocationRequestDescriptor =
    $convert.base64Decode(
        'CjJHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tJbnZvY2F0aW9uUmVxdWVzdBIlCg'
        'xpbnZvY2F0aW9uaWQY/4KsNyABKAlSDGludm9jYXRpb25pZBIZCgZ0YXNraWQYopHrCCABKAlS'
        'BnRhc2tpZBIwChF3aW5kb3dleGVjdXRpb25pZBiJ66eqASABKAlSEXdpbmRvd2V4ZWN1dGlvbm'
        'lk');

@$core.Deprecated(
    'Use getMaintenanceWindowExecutionTaskInvocationResultDescriptor instead')
const GetMaintenanceWindowExecutionTaskInvocationResult$json = {
  '1': 'GetMaintenanceWindowExecutionTaskInvocationResult',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {'1': 'invocationid', '3': 116064639, '4': 1, '5': 9, '10': 'invocationid'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {'1': 'parameters', '3': 494900218, '4': 1, '5': 9, '10': 'parameters'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {
      '1': 'taskexecutionid',
      '3': 474180408,
      '4': 1,
      '5': 9,
      '10': 'taskexecutionid'
    },
    {
      '1': 'tasktype',
      '3': 370407909,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'tasktype'
    },
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `GetMaintenanceWindowExecutionTaskInvocationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowExecutionTaskInvocationResultDescriptor = $convert.base64Decode(
    'CjFHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tJbnZvY2F0aW9uUmVzdWx0EhsKB2'
    'VuZHRpbWUYzO+8HiABKAlSB2VuZHRpbWUSIwoLZXhlY3V0aW9uaWQYsc+vRiABKAlSC2V4ZWN1'
    'dGlvbmlkEiUKDGludm9jYXRpb25pZBj/gqw3IAEoCVIMaW52b2NhdGlvbmlkEi0KEG93bmVyaW'
    '5mb3JtYXRpb24Yg5rFICABKAlSEG93bmVyaW5mb3JtYXRpb24SIgoKcGFyYW1ldGVycxj6p/7r'
    'ASABKAlSCnBhcmFtZXRlcnMSIAoJc3RhcnR0aW1lGO+05bABIAEoCVIJc3RhcnR0aW1lEkAKBn'
    'N0YXR1cxiQ5PsCIAEoDjIlLnNzbS5NYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblN0YXR1c1IG'
    'c3RhdHVzEigKDXN0YXR1c2RldGFpbHMYqJLBsQEgASgJUg1zdGF0dXNkZXRhaWxzEiwKD3Rhc2'
    'tleGVjdXRpb25pZBi41o3iASABKAlSD3Rhc2tleGVjdXRpb25pZBI+Cgh0YXNrdHlwZRjl88+w'
    'ASABKA4yHi5zc20uTWFpbnRlbmFuY2VXaW5kb3dUYXNrVHlwZVIIdGFza3R5cGUSMAoRd2luZG'
    '93ZXhlY3V0aW9uaWQYieunqgEgASgJUhF3aW5kb3dleGVjdXRpb25pZBIpCg53aW5kb3d0YXJn'
    'ZXRpZBjezup7IAEoCVIOd2luZG93dGFyZ2V0aWQ=');

@$core.Deprecated(
    'Use getMaintenanceWindowExecutionTaskRequestDescriptor instead')
const GetMaintenanceWindowExecutionTaskRequest$json = {
  '1': 'GetMaintenanceWindowExecutionTaskRequest',
  '2': [
    {'1': 'taskid', '3': 18532514, '4': 1, '5': 9, '10': 'taskid'},
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `GetMaintenanceWindowExecutionTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowExecutionTaskRequestDescriptor =
    $convert.base64Decode(
        'CihHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tSZXF1ZXN0EhkKBnRhc2tpZBiike'
        'sIIAEoCVIGdGFza2lkEjAKEXdpbmRvd2V4ZWN1dGlvbmlkGInrp6oBIAEoCVIRd2luZG93ZXhl'
        'Y3V0aW9uaWQ=');

@$core
    .Deprecated('Use getMaintenanceWindowExecutionTaskResultDescriptor instead')
const GetMaintenanceWindowExecutionTaskResult$json = {
  '1': 'GetMaintenanceWindowExecutionTaskResult',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {'1': 'servicerole', '3': 47807725, '4': 1, '5': 9, '10': 'servicerole'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskexecutionid',
      '3': 474180408,
      '4': 1,
      '5': 9,
      '10': 'taskexecutionid'
    },
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.GetMaintenanceWindowExecutionTaskResult.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'type'
    },
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
  '3': [GetMaintenanceWindowExecutionTaskResult_TaskparametersEntry$json],
};

@$core
    .Deprecated('Use getMaintenanceWindowExecutionTaskResultDescriptor instead')
const GetMaintenanceWindowExecutionTaskResult_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `GetMaintenanceWindowExecutionTaskResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowExecutionTaskResultDescriptor = $convert.base64Decode(
    'CidHZXRNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tSZXN1bHQSSgoSYWxhcm1jb25maW'
    'd1cmF0aW9uGImZuSEgASgLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1'
    'cmF0aW9uEhsKB2VuZHRpbWUYzO+8HiABKAlSB2VuZHRpbWUSKQoObWF4Y29uY3VycmVuY3kY/c'
    'GODiABKAlSDm1heGNvbmN1cnJlbmN5Eh8KCW1heGVycm9ycxirwvU9IAEoCVIJbWF4ZXJyb3Jz'
    'Eh0KCHByaW9yaXR5GKq+tjQgASgFUghwcmlvcml0eRIjCgtzZXJ2aWNlcm9sZRjt+eUWIAEoCV'
    'ILc2VydmljZXJvbGUSIAoJc3RhcnR0aW1lGO+05bABIAEoCVIJc3RhcnR0aW1lEkAKBnN0YXR1'
    'cxiQ5PsCIAEoDjIlLnNzbS5NYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblN0YXR1c1IGc3RhdH'
    'VzEigKDXN0YXR1c2RldGFpbHMYqJLBsQEgASgJUg1zdGF0dXNkZXRhaWxzEhwKB3Rhc2thcm4Y'
    '5Mn6lAEgASgJUgd0YXNrYXJuEiwKD3Rhc2tleGVjdXRpb25pZBi41o3iASABKAlSD3Rhc2tleG'
    'VjdXRpb25pZBJsCg50YXNrcGFyYW1ldGVycxiBj+a3ASADKAsyQC5zc20uR2V0TWFpbnRlbmFu'
    'Y2VXaW5kb3dFeGVjdXRpb25UYXNrUmVzdWx0LlRhc2twYXJhbWV0ZXJzRW50cnlSDnRhc2twYX'
    'JhbWV0ZXJzEkcKD3RyaWdnZXJlZGFsYXJtcxiF7cF9IAMoCzIaLnNzbS5BbGFybVN0YXRlSW5m'
    'b3JtYXRpb25SD3RyaWdnZXJlZGFsYXJtcxI2CgR0eXBlGO6g14oBIAEoDjIeLnNzbS5NYWludG'
    'VuYW5jZVdpbmRvd1Rhc2tUeXBlUgR0eXBlEjAKEXdpbmRvd2V4ZWN1dGlvbmlkGInrp6oBIAEo'
    'CVIRd2luZG93ZXhlY3V0aW9uaWQadQoTVGFza3BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKA'
    'lSA2tleRJICgV2YWx1ZRgCIAEoCzIyLnNzbS5NYWludGVuYW5jZVdpbmRvd1Rhc2tQYXJhbWV0'
    'ZXJWYWx1ZUV4cHJlc3Npb25SBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getMaintenanceWindowRequestDescriptor instead')
const GetMaintenanceWindowRequest$json = {
  '1': 'GetMaintenanceWindowRequest',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `GetMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRNYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSHQoId2luZG93aWQYqeSHCSABKAlSCHdpbm'
        'Rvd2lk');

@$core.Deprecated('Use getMaintenanceWindowResultDescriptor instead')
const GetMaintenanceWindowResult$json = {
  '1': 'GetMaintenanceWindowResult',
  '2': [
    {
      '1': 'allowunassociatedtargets',
      '3': 154411300,
      '4': 1,
      '5': 8,
      '10': 'allowunassociatedtargets'
    },
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'cutoff', '3': 498433089, '4': 1, '5': 5, '10': 'cutoff'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'modifieddate', '3': 210609143, '4': 1, '5': 9, '10': 'modifieddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextexecutiontime',
      '3': 528679720,
      '4': 1,
      '5': 9,
      '10': 'nextexecutiontime'
    },
    {'1': 'schedule', '3': 66697965, '4': 1, '5': 9, '10': 'schedule'},
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'scheduletimezone',
      '3': 170037696,
      '4': 1,
      '5': 9,
      '10': 'scheduletimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `GetMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowResultDescriptor = $convert.base64Decode(
    'ChpHZXRNYWludGVuYW5jZVdpbmRvd1Jlc3VsdBI9ChhhbGxvd3VuYXNzb2NpYXRlZHRhcmdldH'
    'MYpMLQSSABKAhSGGFsbG93dW5hc3NvY2lhdGVkdGFyZ2V0cxIkCgtjcmVhdGVkZGF0ZRiwsOfG'
    'ASABKAlSC2NyZWF0ZWRkYXRlEhoKBmN1dG9mZhjB+NXtASABKAVSBmN1dG9mZhIjCgtkZXNjcm'
    'lwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SHgoIZHVyYXRpb24YrpKdpgEgASgFUghkdXJh'
    'dGlvbhIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIbCgdlbmRkYXRlGM+z+SQgASgJUg'
    'dlbmRkYXRlEiUKDG1vZGlmaWVkZGF0ZRj3x7ZkIAEoCVIMbW9kaWZpZWRkYXRlEhUKBG5hbWUY'
    'h+aBfyABKAlSBG5hbWUSMAoRbmV4dGV4ZWN1dGlvbnRpbWUYqIaM/AEgASgJUhFuZXh0ZXhlY3'
    'V0aW9udGltZRIdCghzY2hlZHVsZRjt9eYfIAEoCVIIc2NoZWR1bGUSKQoOc2NoZWR1bGVvZmZz'
    'ZXQY2JHqSiABKAVSDnNjaGVkdWxlb2Zmc2V0Ei0KEHNjaGVkdWxldGltZXpvbmUYwKOKUSABKA'
    'lSEHNjaGVkdWxldGltZXpvbmUSIAoJc3RhcnRkYXRlGPz4oNQBIAEoCVIJc3RhcnRkYXRlEh0K'
    'CHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZA==');

@$core.Deprecated('Use getMaintenanceWindowTaskRequestDescriptor instead')
const GetMaintenanceWindowTaskRequest$json = {
  '1': 'GetMaintenanceWindowTaskRequest',
  '2': [
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
};

/// Descriptor for `GetMaintenanceWindowTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowTaskRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRNYWludGVuYW5jZVdpbmRvd1Rhc2tSZXF1ZXN0Eh0KCHdpbmRvd2lkGKnkhwkgASgJUg'
        'h3aW5kb3dpZBImCgx3aW5kb3d0YXNraWQY2o2xmgEgASgJUgx3aW5kb3d0YXNraWQ=');

@$core.Deprecated('Use getMaintenanceWindowTaskResultDescriptor instead')
const GetMaintenanceWindowTaskResult$json = {
  '1': 'GetMaintenanceWindowTaskResult',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cutoffbehavior',
      '3': 120608587,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskCutoffBehavior',
      '10': 'cutoffbehavior'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'logginginfo',
      '3': 448312415,
      '4': 1,
      '5': 11,
      '6': '.ssm.LoggingInfo',
      '10': 'logginginfo'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskinvocationparameters',
      '3': 347127635,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskInvocationParameters',
      '10': 'taskinvocationparameters'
    },
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.GetMaintenanceWindowTaskResult.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {
      '1': 'tasktype',
      '3': 370407909,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'tasktype'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
  '3': [GetMaintenanceWindowTaskResult_TaskparametersEntry$json],
};

@$core.Deprecated('Use getMaintenanceWindowTaskResultDescriptor instead')
const GetMaintenanceWindowTaskResult_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `GetMaintenanceWindowTaskResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMaintenanceWindowTaskResultDescriptor = $convert.base64Decode(
    'Ch5HZXRNYWludGVuYW5jZVdpbmRvd1Rhc2tSZXN1bHQSSgoSYWxhcm1jb25maWd1cmF0aW9uGI'
    'mZuSEgASgLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uElMK'
    'DmN1dG9mZmJlaGF2aW9yGMuuwTkgASgOMiguc3NtLk1haW50ZW5hbmNlV2luZG93VGFza0N1dG'
    '9mZkJlaGF2aW9yUg5jdXRvZmZiZWhhdmlvchIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVz'
    'Y3JpcHRpb24SNgoLbG9nZ2luZ2luZm8Y3+ji1QEgASgLMhAuc3NtLkxvZ2dpbmdJbmZvUgtsb2'
    'dnaW5naW5mbxIpCg5tYXhjb25jdXJyZW5jeRj9wY4OIAEoCVIObWF4Y29uY3VycmVuY3kSHwoJ'
    'bWF4ZXJyb3JzGKvC9T0gASgJUgltYXhlcnJvcnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIdCg'
    'hwcmlvcml0eRiqvrY0IAEoBVIIcHJpb3JpdHkSKgoOc2VydmljZXJvbGVhcm4YhOPatgEgASgJ'
    'Ug5zZXJ2aWNlcm9sZWFybhIoCgd0YXJnZXRzGIKbgn0gAygLMgsuc3NtLlRhcmdldFIHdGFyZ2'
    'V0cxIcCgd0YXNrYXJuGOTJ+pQBIAEoCVIHdGFza2FybhJuChh0YXNraW52b2NhdGlvbnBhcmFt'
    'ZXRlcnMY0/7CpQEgASgLMi4uc3NtLk1haW50ZW5hbmNlV2luZG93VGFza0ludm9jYXRpb25QYX'
    'JhbWV0ZXJzUhh0YXNraW52b2NhdGlvbnBhcmFtZXRlcnMSYwoOdGFza3BhcmFtZXRlcnMYgY/m'
    'twEgAygLMjcuc3NtLkdldE1haW50ZW5hbmNlV2luZG93VGFza1Jlc3VsdC5UYXNrcGFyYW1ldG'
    'Vyc0VudHJ5Ug50YXNrcGFyYW1ldGVycxI+Cgh0YXNrdHlwZRjl88+wASABKA4yHi5zc20uTWFp'
    'bnRlbmFuY2VXaW5kb3dUYXNrVHlwZVIIdGFza3R5cGUSHQoId2luZG93aWQYqeSHCSABKAlSCH'
    'dpbmRvd2lkEiYKDHdpbmRvd3Rhc2tpZBjajbGaASABKAlSDHdpbmRvd3Rhc2tpZBp1ChNUYXNr'
    'cGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EkgKBXZhbHVlGAIgASgLMjIuc3NtLk'
    '1haW50ZW5hbmNlV2luZG93VGFza1BhcmFtZXRlclZhbHVlRXhwcmVzc2lvblIFdmFsdWU6AjgB');

@$core.Deprecated('Use getOpsItemRequestDescriptor instead')
const GetOpsItemRequest$json = {
  '1': 'GetOpsItemRequest',
  '2': [
    {'1': 'opsitemarn', '3': 222489428, '4': 1, '5': 9, '10': 'opsitemarn'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `GetOpsItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsItemRequestDescriptor = $convert.base64Decode(
    'ChFHZXRPcHNJdGVtUmVxdWVzdBIhCgpvcHNpdGVtYXJuGNTWi2ogASgJUgpvcHNpdGVtYXJuEh'
    '8KCW9wc2l0ZW1pZBjS0pUMIAEoCVIJb3BzaXRlbWlk');

@$core.Deprecated('Use getOpsItemResponseDescriptor instead')
const GetOpsItemResponse$json = {
  '1': 'GetOpsItemResponse',
  '2': [
    {
      '1': 'opsitem',
      '3': 218191741,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItem',
      '10': 'opsitem'
    },
  ],
};

/// Descriptor for `GetOpsItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsItemResponseDescriptor = $convert.base64Decode(
    'ChJHZXRPcHNJdGVtUmVzcG9uc2USKQoHb3BzaXRlbRj9roVoIAEoCzIMLnNzbS5PcHNJdGVtUg'
    'dvcHNpdGVt');

@$core.Deprecated('Use getOpsMetadataRequestDescriptor instead')
const GetOpsMetadataRequest$json = {
  '1': 'GetOpsMetadataRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
  ],
};

/// Descriptor for `GetOpsMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsMetadataRequestDescriptor = $convert.base64Decode(
    'ChVHZXRPcHNNZXRhZGF0YVJlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3'
    'VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SKgoOb3BzbWV0YWRhdGFhcm4Y'
    'or6C5gEgASgJUg5vcHNtZXRhZGF0YWFybg==');

@$core.Deprecated('Use getOpsMetadataResultDescriptor instead')
const GetOpsMetadataResult$json = {
  '1': 'GetOpsMetadataResult',
  '2': [
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 3,
      '5': 11,
      '6': '.ssm.GetOpsMetadataResult.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
  ],
  '3': [GetOpsMetadataResult_MetadataEntry$json],
};

@$core.Deprecated('Use getOpsMetadataResultDescriptor instead')
const GetOpsMetadataResult_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MetadataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `GetOpsMetadataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsMetadataResultDescriptor = $convert.base64Decode(
    'ChRHZXRPcHNNZXRhZGF0YVJlc3VsdBJHCghtZXRhZGF0YRjh4o/gASADKAsyJy5zc20uR2V0T3'
    'BzTWV0YWRhdGFSZXN1bHQuTWV0YWRhdGFFbnRyeVIIbWV0YWRhdGESHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4SIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc291cmNlaWQaTw'
    'oNTWV0YWRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIoCgV2YWx1ZRgCIAEoCzISLnNzbS5N'
    'ZXRhZGF0YVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getOpsSummaryRequestDescriptor instead')
const GetOpsSummaryRequest$json = {
  '1': 'GetOpsSummaryRequest',
  '2': [
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsAggregator',
      '10': 'aggregators'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resultattributes',
      '3': 430325814,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsResultAttribute',
      '10': 'resultattributes'
    },
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
  ],
};

/// Descriptor for `GetOpsSummaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsSummaryRequestDescriptor = $convert.base64Decode(
    'ChRHZXRPcHNTdW1tYXJ5UmVxdWVzdBI3CgthZ2dyZWdhdG9ycxiO/YNNIAMoCzISLnNzbS5PcH'
    'NBZ2dyZWdhdG9yUgthZ2dyZWdhdG9ycxIrCgdmaWx0ZXJzGO3N6lkgAygLMg4uc3NtLk9wc0Zp'
    'bHRlclIHZmlsdGVycxIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZX'
    'h0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhJHChByZXN1bHRhdHRyaWJ1dGVzGLaAmc0BIAMo'
    'CzIXLnNzbS5PcHNSZXN1bHRBdHRyaWJ1dGVSEHJlc3VsdGF0dHJpYnV0ZXMSHgoIc3luY25hbW'
    'UYopaysAEgASgJUghzeW5jbmFtZQ==');

@$core.Deprecated('Use getOpsSummaryResultDescriptor instead')
const GetOpsSummaryResult$json = {
  '1': 'GetOpsSummaryResult',
  '2': [
    {
      '1': 'entities',
      '3': 406339595,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsEntity',
      '10': 'entities'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `GetOpsSummaryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpsSummaryResultDescriptor = $convert.base64Decode(
    'ChNHZXRPcHNTdW1tYXJ5UmVzdWx0Ei4KCGVudGl0aWVzGIuA4cEBIAMoCzIOLnNzbS5PcHNFbn'
    'RpdHlSCGVudGl0aWVzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use getParameterHistoryRequestDescriptor instead')
const GetParameterHistoryRequest$json = {
  '1': 'GetParameterHistoryRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'withdecryption',
      '3': 387365121,
      '4': 1,
      '5': 8,
      '10': 'withdecryption'
    },
  ],
};

/// Descriptor for `GetParameterHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParameterHistoryRequestDescriptor = $convert.base64Decode(
    'ChpHZXRQYXJhbWV0ZXJIaXN0b3J5UmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
    'F4cmVzdWx0cxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJ'
    'bmV4dHRva2VuEioKDndpdGhkZWNyeXB0aW9uGIHy2rgBIAEoCFIOd2l0aGRlY3J5cHRpb24=');

@$core.Deprecated('Use getParameterHistoryResultDescriptor instead')
const GetParameterHistoryResult$json = {
  '1': 'GetParameterHistoryResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterHistory',
      '10': 'parameters'
    },
  ],
};

/// Descriptor for `GetParameterHistoryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParameterHistoryResultDescriptor = $convert.base64Decode(
    'ChlHZXRQYXJhbWV0ZXJIaXN0b3J5UmVzdWx0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
    'Rva2VuEjkKCnBhcmFtZXRlcnMY+qf+6wEgAygLMhUuc3NtLlBhcmFtZXRlckhpc3RvcnlSCnBh'
    'cmFtZXRlcnM=');

@$core.Deprecated('Use getParameterRequestDescriptor instead')
const GetParameterRequest$json = {
  '1': 'GetParameterRequest',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'withdecryption',
      '3': 387365121,
      '4': 1,
      '5': 8,
      '10': 'withdecryption'
    },
  ],
};

/// Descriptor for `GetParameterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParameterRequestDescriptor = $convert.base64Decode(
    'ChNHZXRQYXJhbWV0ZXJSZXF1ZXN0EhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSKgoOd2l0aGRlY3'
    'J5cHRpb24YgfLauAEgASgIUg53aXRoZGVjcnlwdGlvbg==');

@$core.Deprecated('Use getParameterResultDescriptor instead')
const GetParameterResult$json = {
  '1': 'GetParameterResult',
  '2': [
    {
      '1': 'parameter',
      '3': 407419825,
      '4': 1,
      '5': 11,
      '6': '.ssm.Parameter',
      '10': 'parameter'
    },
  ],
};

/// Descriptor for `GetParameterResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParameterResultDescriptor = $convert.base64Decode(
    'ChJHZXRQYXJhbWV0ZXJSZXN1bHQSMAoJcGFyYW1ldGVyGLH3osIBIAEoCzIOLnNzbS5QYXJhbW'
    'V0ZXJSCXBhcmFtZXRlcg==');

@$core.Deprecated('Use getParametersByPathRequestDescriptor instead')
const GetParametersByPathRequest$json = {
  '1': 'GetParametersByPathRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameterfilters',
      '3': 511502980,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterStringFilter',
      '10': 'parameterfilters'
    },
    {'1': 'path', '3': 191292503, '4': 1, '5': 9, '10': 'path'},
    {'1': 'recursive', '3': 335804696, '4': 1, '5': 8, '10': 'recursive'},
    {
      '1': 'withdecryption',
      '3': 387365121,
      '4': 1,
      '5': 8,
      '10': 'withdecryption'
    },
  ],
};

/// Descriptor for `GetParametersByPathRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersByPathRequestDescriptor = $convert.base64Decode(
    'ChpHZXRQYXJhbWV0ZXJzQnlQYXRoUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
    'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhJKChBwYXJhbWV0ZXJm'
    'aWx0ZXJzGITV8/MBIAMoCzIaLnNzbS5QYXJhbWV0ZXJTdHJpbmdGaWx0ZXJSEHBhcmFtZXRlcm'
    'ZpbHRlcnMSFQoEcGF0aBjXyJtbIAEoCVIEcGF0aBIgCglyZWN1cnNpdmUYmPKPoAEgASgIUgly'
    'ZWN1cnNpdmUSKgoOd2l0aGRlY3J5cHRpb24YgfLauAEgASgIUg53aXRoZGVjcnlwdGlvbg==');

@$core.Deprecated('Use getParametersByPathResultDescriptor instead')
const GetParametersByPathResult$json = {
  '1': 'GetParametersByPathResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.Parameter',
      '10': 'parameters'
    },
  ],
};

/// Descriptor for `GetParametersByPathResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersByPathResultDescriptor = $convert.base64Decode(
    'ChlHZXRQYXJhbWV0ZXJzQnlQYXRoUmVzdWx0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
    'Rva2VuEjIKCnBhcmFtZXRlcnMY+qf+6wEgAygLMg4uc3NtLlBhcmFtZXRlclIKcGFyYW1ldGVy'
    'cw==');

@$core.Deprecated('Use getParametersRequestDescriptor instead')
const GetParametersRequest$json = {
  '1': 'GetParametersRequest',
  '2': [
    {'1': 'names', '3': 324387120, '4': 3, '5': 9, '10': 'names'},
    {
      '1': 'withdecryption',
      '3': 387365121,
      '4': 1,
      '5': 8,
      '10': 'withdecryption'
    },
  ],
};

/// Descriptor for `GetParametersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersRequestDescriptor = $convert.base64Decode(
    'ChRHZXRQYXJhbWV0ZXJzUmVxdWVzdBIYCgVuYW1lcxiwgteaASADKAlSBW5hbWVzEioKDndpdG'
    'hkZWNyeXB0aW9uGIHy2rgBIAEoCFIOd2l0aGRlY3J5cHRpb24=');

@$core.Deprecated('Use getParametersResultDescriptor instead')
const GetParametersResult$json = {
  '1': 'GetParametersResult',
  '2': [
    {
      '1': 'invalidparameters',
      '3': 329741157,
      '4': 3,
      '5': 9,
      '10': 'invalidparameters'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.Parameter',
      '10': 'parameters'
    },
  ],
};

/// Descriptor for `GetParametersResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersResultDescriptor = $convert.base64Decode(
    'ChNHZXRQYXJhbWV0ZXJzUmVzdWx0EjAKEWludmFsaWRwYXJhbWV0ZXJzGOXmnZ0BIAMoCVIRaW'
    '52YWxpZHBhcmFtZXRlcnMSMgoKcGFyYW1ldGVycxj6p/7rASADKAsyDi5zc20uUGFyYW1ldGVy'
    'UgpwYXJhbWV0ZXJz');

@$core.Deprecated('Use getPatchBaselineForPatchGroupRequestDescriptor instead')
const GetPatchBaselineForPatchGroupRequest$json = {
  '1': 'GetPatchBaselineForPatchGroupRequest',
  '2': [
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `GetPatchBaselineForPatchGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPatchBaselineForPatchGroupRequestDescriptor =
    $convert.base64Decode(
        'CiRHZXRQYXRjaEJhc2VsaW5lRm9yUGF0Y2hHcm91cFJlcXVlc3QSQQoPb3BlcmF0aW5nc3lzdG'
        'VtGOr9wRIgASgOMhQuc3NtLk9wZXJhdGluZ1N5c3RlbVIPb3BlcmF0aW5nc3lzdGVtEiIKCnBh'
        'dGNoZ3JvdXAY4bex9wEgASgJUgpwYXRjaGdyb3Vw');

@$core.Deprecated('Use getPatchBaselineForPatchGroupResultDescriptor instead')
const GetPatchBaselineForPatchGroupResult$json = {
  '1': 'GetPatchBaselineForPatchGroupResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `GetPatchBaselineForPatchGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPatchBaselineForPatchGroupResultDescriptor =
    $convert.base64Decode(
        'CiNHZXRQYXRjaEJhc2VsaW5lRm9yUGF0Y2hHcm91cFJlc3VsdBIhCgpiYXNlbGluZWlkGNDk2y'
        'ggASgJUgpiYXNlbGluZWlkEkEKD29wZXJhdGluZ3N5c3RlbRjq/cESIAEoDjIULnNzbS5PcGVy'
        'YXRpbmdTeXN0ZW1SD29wZXJhdGluZ3N5c3RlbRIiCgpwYXRjaGdyb3VwGOG3sfcBIAEoCVIKcG'
        'F0Y2hncm91cA==');

@$core.Deprecated('Use getPatchBaselineRequestDescriptor instead')
const GetPatchBaselineRequest$json = {
  '1': 'GetPatchBaselineRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `GetPatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPatchBaselineRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRQYXRjaEJhc2VsaW5lUmVxdWVzdBIhCgpiYXNlbGluZWlkGNDk2yggASgJUgpiYXNlbG'
        'luZWlk');

@$core.Deprecated('Use getPatchBaselineResultDescriptor instead')
const GetPatchBaselineResult$json = {
  '1': 'GetPatchBaselineResult',
  '2': [
    {
      '1': 'approvalrules',
      '3': 71466346,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchRuleGroup',
      '10': 'approvalrules'
    },
    {
      '1': 'approvedpatches',
      '3': 199384709,
      '4': 3,
      '5': 9,
      '10': 'approvedpatches'
    },
    {
      '1': 'approvedpatchescompliancelevel',
      '3': 63924432,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'approvedpatchescompliancelevel'
    },
    {
      '1': 'approvedpatchesenablenonsecurity',
      '3': 295555901,
      '4': 1,
      '5': 8,
      '10': 'approvedpatchesenablenonsecurity'
    },
    {
      '1': 'availablesecurityupdatescompliancestatus',
      '3': 187471858,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceStatus',
      '10': 'availablesecurityupdatescompliancestatus'
    },
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'globalfilters',
      '3': 263302754,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'globalfilters'
    },
    {'1': 'modifieddate', '3': 210609143, '4': 1, '5': 9, '10': 'modifieddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {'1': 'patchgroups', '3': 18098282, '4': 3, '5': 9, '10': 'patchgroups'},
    {
      '1': 'rejectedpatches',
      '3': 309657116,
      '4': 3,
      '5': 9,
      '10': 'rejectedpatches'
    },
    {
      '1': 'rejectedpatchesaction',
      '3': 356538330,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchAction',
      '10': 'rejectedpatchesaction'
    },
    {
      '1': 'sources',
      '3': 46625746,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchSource',
      '10': 'sources'
    },
  ],
};

/// Descriptor for `GetPatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPatchBaselineResultDescriptor = $convert.base64Decode(
    'ChZHZXRQYXRjaEJhc2VsaW5lUmVzdWx0EjwKDWFwcHJvdmFscnVsZXMY6vqJIiABKAsyEy5zc2'
    '0uUGF0Y2hSdWxlR3JvdXBSDWFwcHJvdmFscnVsZXMSKwoPYXBwcm92ZWRwYXRjaGVzGIW9iV8g'
    'AygJUg9hcHByb3ZlZHBhdGNoZXMSZAoeYXBwcm92ZWRwYXRjaGVzY29tcGxpYW5jZWxldmVsGN'
    'DRvR4gASgOMhkuc3NtLlBhdGNoQ29tcGxpYW5jZUxldmVsUh5hcHByb3ZlZHBhdGNoZXNjb21w'
    'bGlhbmNlbGV2ZWwSTgogYXBwcm92ZWRwYXRjaGVzZW5hYmxlbm9uc2VjdXJpdHkYvab3jAEgAS'
    'gIUiBhcHByb3ZlZHBhdGNoZXNlbmFibGVub25zZWN1cml0eRJ5CihhdmFpbGFibGVzZWN1cml0'
    'eXVwZGF0ZXNjb21wbGlhbmNlc3RhdHVzGPKvslkgASgOMhouc3NtLlBhdGNoQ29tcGxpYW5jZV'
    'N0YXR1c1IoYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVzY29tcGxpYW5jZXN0YXR1cxIhCgpiYXNl'
    'bGluZWlkGNDk2yggASgJUgpiYXNlbGluZWlkEiQKC2NyZWF0ZWRkYXRlGLCw58YBIAEoCVILY3'
    'JlYXRlZGRhdGUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEj4KDWdsb2Jh'
    'bGZpbHRlcnMY4tzGfSABKAsyFS5zc20uUGF0Y2hGaWx0ZXJHcm91cFINZ2xvYmFsZmlsdGVycx'
    'IlCgxtb2RpZmllZGRhdGUY98e2ZCABKAlSDG1vZGlmaWVkZGF0ZRIVCgRuYW1lGIfmgX8gASgJ'
    'UgRuYW1lEkEKD29wZXJhdGluZ3N5c3RlbRjq/cESIAEoDjIULnNzbS5PcGVyYXRpbmdTeXN0ZW'
    '1SD29wZXJhdGluZ3N5c3RlbRIjCgtwYXRjaGdyb3Vwcxjq0NAIIAMoCVILcGF0Y2hncm91cHMS'
    'LAoPcmVqZWN0ZWRwYXRjaGVzGJz805MBIAMoCVIPcmVqZWN0ZWRwYXRjaGVzEkoKFXJlamVjdG'
    'VkcGF0Y2hlc2FjdGlvbhjar4GqASABKA4yEC5zc20uUGF0Y2hBY3Rpb25SFXJlamVjdGVkcGF0'
    'Y2hlc2FjdGlvbhItCgdzb3VyY2VzGNLnnRYgAygLMhAuc3NtLlBhdGNoU291cmNlUgdzb3VyY2'
    'Vz');

@$core.Deprecated('Use getResourcePoliciesRequestDescriptor instead')
const GetResourcePoliciesRequest$json = {
  '1': 'GetResourcePoliciesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetResourcePoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePoliciesRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRSZXNvdXJjZVBvbGljaWVzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbW'
        'F4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIkCgtyZXNvdXJjZWFy'
        'bhit+NmtASABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use getResourcePoliciesResponseDescriptor instead')
const GetResourcePoliciesResponse$json = {
  '1': 'GetResourcePoliciesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.ssm.GetResourcePoliciesResponseEntry',
      '10': 'policies'
    },
  ],
};

/// Descriptor for `GetResourcePoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePoliciesResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRSZXNvdXJjZVBvbGljaWVzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SRAoIcG9saWNpZXMYmKyKEyADKAsyJS5zc20uR2V0UmVzb3VyY2VQb2xpY2llc1Jl'
        'c3BvbnNlRW50cnlSCHBvbGljaWVz');

@$core.Deprecated('Use getResourcePoliciesResponseEntryDescriptor instead')
const GetResourcePoliciesResponseEntry$json = {
  '1': 'GetResourcePoliciesResponseEntry',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyhash', '3': 203037360, '4': 1, '5': 9, '10': 'policyhash'},
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
  ],
};

/// Descriptor for `GetResourcePoliciesResponseEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePoliciesResponseEntryDescriptor =
    $convert.base64Decode(
        'CiBHZXRSZXNvdXJjZVBvbGljaWVzUmVzcG9uc2VFbnRyeRIaCgZwb2xpY3kYoO/w4AEgASgJUg'
        'Zwb2xpY3kSIQoKcG9saWN5aGFzaBiwtehgIAEoCVIKcG9saWN5aGFzaBIeCghwb2xpY3lpZBjz'
        'o+mOASABKAlSCHBvbGljeWlk');

@$core.Deprecated('Use getServiceSettingRequestDescriptor instead')
const GetServiceSettingRequest$json = {
  '1': 'GetServiceSettingRequest',
  '2': [
    {'1': 'settingid', '3': 422452711, '4': 1, '5': 9, '10': 'settingid'},
  ],
};

/// Descriptor for `GetServiceSettingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServiceSettingRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRTZXJ2aWNlU2V0dGluZ1JlcXVlc3QSIAoJc2V0dGluZ2lkGOe7uMkBIAEoCVIJc2V0dG'
        'luZ2lk');

@$core.Deprecated('Use getServiceSettingResultDescriptor instead')
const GetServiceSettingResult$json = {
  '1': 'GetServiceSettingResult',
  '2': [
    {
      '1': 'servicesetting',
      '3': 486565473,
      '4': 1,
      '5': 11,
      '6': '.ssm.ServiceSetting',
      '10': 'servicesetting'
    },
  ],
};

/// Descriptor for `GetServiceSettingResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getServiceSettingResultDescriptor =
    $convert.base64Decode(
        'ChdHZXRTZXJ2aWNlU2V0dGluZ1Jlc3VsdBI/Cg5zZXJ2aWNlc2V0dGluZxjhzIHoASABKAsyEy'
        '5zc20uU2VydmljZVNldHRpbmdSDnNlcnZpY2VzZXR0aW5n');

@$core.Deprecated('Use hierarchyLevelLimitExceededExceptionDescriptor instead')
const HierarchyLevelLimitExceededException$json = {
  '1': 'HierarchyLevelLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HierarchyLevelLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hierarchyLevelLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiRIaWVyYXJjaHlMZXZlbExpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use hierarchyTypeMismatchExceptionDescriptor instead')
const HierarchyTypeMismatchException$json = {
  '1': 'HierarchyTypeMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `HierarchyTypeMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hierarchyTypeMismatchExceptionDescriptor =
    $convert.base64Decode(
        'Ch5IaWVyYXJjaHlUeXBlTWlzbWF0Y2hFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use idempotentParameterMismatchDescriptor instead')
const IdempotentParameterMismatch$json = {
  '1': 'IdempotentParameterMismatch',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IdempotentParameterMismatch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List idempotentParameterMismatchDescriptor =
    $convert.base64Decode(
        'ChtJZGVtcG90ZW50UGFyYW1ldGVyTWlzbWF0Y2gSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use incompatiblePolicyExceptionDescriptor instead')
const IncompatiblePolicyException$json = {
  '1': 'IncompatiblePolicyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IncompatiblePolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incompatiblePolicyExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbmNvbXBhdGlibGVQb2xpY3lFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use instanceAggregatedAssociationOverviewDescriptor instead')
const InstanceAggregatedAssociationOverview$json = {
  '1': 'InstanceAggregatedAssociationOverview',
  '2': [
    {
      '1': 'detailedstatus',
      '3': 223171376,
      '4': 1,
      '5': 9,
      '10': 'detailedstatus'
    },
    {
      '1': 'instanceassociationstatusaggregatedcount',
      '3': 202072750,
      '4': 3,
      '5': 11,
      '6':
          '.ssm.InstanceAggregatedAssociationOverview.InstanceassociationstatusaggregatedcountEntry',
      '10': 'instanceassociationstatusaggregatedcount'
    },
  ],
  '3': [
    InstanceAggregatedAssociationOverview_InstanceassociationstatusaggregatedcountEntry$json
  ],
};

@$core.Deprecated('Use instanceAggregatedAssociationOverviewDescriptor instead')
const InstanceAggregatedAssociationOverview_InstanceassociationstatusaggregatedcountEntry$json =
    {
  '1': 'InstanceassociationstatusaggregatedcountEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InstanceAggregatedAssociationOverview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceAggregatedAssociationOverviewDescriptor = $convert.base64Decode(
    'CiVJbnN0YW5jZUFnZ3JlZ2F0ZWRBc3NvY2lhdGlvbk92ZXJ2aWV3EikKDmRldGFpbGVkc3RhdH'
    'VzGLCmtWogASgJUg5kZXRhaWxlZHN0YXR1cxK3AQooaW5zdGFuY2Vhc3NvY2lhdGlvbnN0YXR1'
    'c2FnZ3JlZ2F0ZWRjb3VudBiuxa1gIAMoCzJYLnNzbS5JbnN0YW5jZUFnZ3JlZ2F0ZWRBc3NvY2'
    'lhdGlvbk92ZXJ2aWV3Lkluc3RhbmNlYXNzb2NpYXRpb25zdGF0dXNhZ2dyZWdhdGVkY291bnRF'
    'bnRyeVIoaW5zdGFuY2Vhc3NvY2lhdGlvbnN0YXR1c2FnZ3JlZ2F0ZWRjb3VudBpbCi1JbnN0YW'
    '5jZWFzc29jaWF0aW9uc3RhdHVzYWdncmVnYXRlZGNvdW50RW50cnkSEAoDa2V5GAEgASgJUgNr'
    'ZXkSFAoFdmFsdWUYAiABKAVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use instanceAssociationDescriptor instead')
const InstanceAssociation$json = {
  '1': 'InstanceAssociation',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {'1': 'content', '3': 23568227, '4': 1, '5': 9, '10': 'content'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
  ],
};

/// Descriptor for `InstanceAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceAssociationDescriptor = $convert.base64Decode(
    'ChNJbnN0YW5jZUFzc29jaWF0aW9uEicKDWFzc29jaWF0aW9uaWQYkvyVQiABKAlSDWFzc29jaW'
    'F0aW9uaWQSMgoSYXNzb2NpYXRpb252ZXJzaW9uGJGKydUBIAEoCVISYXNzb2NpYXRpb252ZXJz'
    'aW9uEhsKB2NvbnRlbnQY476eCyABKAlSB2NvbnRlbnQSIQoKaW5zdGFuY2VpZBigrdEXIAEoCV'
    'IKaW5zdGFuY2VpZA==');

@$core.Deprecated('Use instanceAssociationOutputLocationDescriptor instead')
const InstanceAssociationOutputLocation$json = {
  '1': 'InstanceAssociationOutputLocation',
  '2': [
    {
      '1': 's3location',
      '3': 517005333,
      '4': 1,
      '5': 11,
      '6': '.ssm.S3OutputLocation',
      '10': 's3location'
    },
  ],
};

/// Descriptor for `InstanceAssociationOutputLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceAssociationOutputLocationDescriptor =
    $convert.base64Decode(
        'CiFJbnN0YW5jZUFzc29jaWF0aW9uT3V0cHV0TG9jYXRpb24SOQoKczNsb2NhdGlvbhiVwMP2AS'
        'ABKAsyFS5zc20uUzNPdXRwdXRMb2NhdGlvblIKczNsb2NhdGlvbg==');

@$core.Deprecated('Use instanceAssociationOutputUrlDescriptor instead')
const InstanceAssociationOutputUrl$json = {
  '1': 'InstanceAssociationOutputUrl',
  '2': [
    {
      '1': 's3outputurl',
      '3': 490515856,
      '4': 1,
      '5': 11,
      '6': '.ssm.S3OutputUrl',
      '10': 's3outputurl'
    },
  ],
};

/// Descriptor for `InstanceAssociationOutputUrl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceAssociationOutputUrlDescriptor =
    $convert.base64Decode(
        'ChxJbnN0YW5jZUFzc29jaWF0aW9uT3V0cHV0VXJsEjYKC3Mzb3V0cHV0dXJsGJDb8ukBIAEoCz'
        'IQLnNzbS5TM091dHB1dFVybFILczNvdXRwdXR1cmw=');

@$core.Deprecated('Use instanceAssociationStatusInfoDescriptor instead')
const InstanceAssociationStatusInfo$json = {
  '1': 'InstanceAssociationStatusInfo',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'detailedstatus',
      '3': 223171376,
      '4': 1,
      '5': 9,
      '10': 'detailedstatus'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {
      '1': 'executiondate',
      '3': 432181066,
      '4': 1,
      '5': 9,
      '10': 'executiondate'
    },
    {
      '1': 'executionsummary',
      '3': 71480964,
      '4': 1,
      '5': 9,
      '10': 'executionsummary'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputurl',
      '3': 42495998,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputUrl',
      '10': 'outputurl'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `InstanceAssociationStatusInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceAssociationStatusInfoDescriptor = $convert.base64Decode(
    'Ch1JbnN0YW5jZUFzc29jaWF0aW9uU3RhdHVzSW5mbxInCg1hc3NvY2lhdGlvbmlkGJL8lUIgAS'
    'gJUg1hc3NvY2lhdGlvbmlkEiwKD2Fzc29jaWF0aW9ubmFtZRiYkMWVASABKAlSD2Fzc29jaWF0'
    'aW9ubmFtZRIyChJhc3NvY2lhdGlvbnZlcnNpb24YkYrJ1QEgASgJUhJhc3NvY2lhdGlvbnZlcn'
    'Npb24SKQoOZGV0YWlsZWRzdGF0dXMYsKa1aiABKAlSDmRldGFpbGVkc3RhdHVzEisKD2RvY3Vt'
    'ZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEh8KCWVycm9yY29kZRiZ1sMQIA'
    'EoCVIJZXJyb3Jjb2RlEigKDWV4ZWN1dGlvbmRhdGUYyp6KzgEgASgJUg1leGVjdXRpb25kYXRl'
    'Ei0KEGV4ZWN1dGlvbnN1bW1hcnkYhO2KIiABKAlSEGV4ZWN1dGlvbnN1bW1hcnkSIQoKaW5zdG'
    'FuY2VpZBigrdEXIAEoCVIKaW5zdGFuY2VpZBIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkIKCW91'
    'dHB1dHVybBj+36EUIAEoCzIhLnNzbS5JbnN0YW5jZUFzc29jaWF0aW9uT3V0cHV0VXJsUglvdX'
    'RwdXR1cmwSGQoGc3RhdHVzGJDk+wIgASgJUgZzdGF0dXM=');

@$core.Deprecated('Use instanceInfoDescriptor instead')
const InstanceInfo$json = {
  '1': 'InstanceInfo',
  '2': [
    {'1': 'agenttype', '3': 408897495, '4': 1, '5': 9, '10': 'agenttype'},
    {'1': 'agentversion', '3': 173670635, '4': 1, '5': 9, '10': 'agentversion'},
    {'1': 'computername', '3': 67735292, '4': 1, '5': 9, '10': 'computername'},
    {
      '1': 'instancestatus',
      '3': 29604317,
      '4': 1,
      '5': 9,
      '10': 'instancestatus'
    },
    {'1': 'ipaddress', '3': 1800397, '4': 1, '5': 9, '10': 'ipaddress'},
    {
      '1': 'managedstatus',
      '3': 200879673,
      '4': 1,
      '5': 14,
      '6': '.ssm.ManagedStatus',
      '10': 'managedstatus'
    },
    {'1': 'platformname', '3': 188917178, '4': 1, '5': 9, '10': 'platformname'},
    {
      '1': 'platformtype',
      '3': 378742691,
      '4': 1,
      '5': 14,
      '6': '.ssm.PlatformType',
      '10': 'platformtype'
    },
    {
      '1': 'platformversion',
      '3': 139924287,
      '4': 1,
      '5': 9,
      '10': 'platformversion'
    },
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `InstanceInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceInfoDescriptor = $convert.base64Decode(
    'CgxJbnN0YW5jZUluZm8SIAoJYWdlbnR0eXBlGNeP/cIBIAEoCVIJYWdlbnR0eXBlEiUKDGFnZW'
    '50dmVyc2lvbhjrgehSIAEoCVIMYWdlbnR2ZXJzaW9uEiUKDGNvbXB1dGVybmFtZRj8naYgIAEo'
    'CVIMY29tcHV0ZXJuYW1lEikKDmluc3RhbmNlc3RhdHVzGN3zjg4gASgJUg5pbnN0YW5jZXN0YX'
    'R1cxIeCglpcGFkZHJlc3MYzfFtIAEoCVIJaXBhZGRyZXNzEjsKDW1hbmFnZWRzdGF0dXMYudzk'
    'XyABKA4yEi5zc20uTWFuYWdlZFN0YXR1c1INbWFuYWdlZHN0YXR1cxIlCgxwbGF0Zm9ybW5hbW'
    'UYusuKWiABKAlSDHBsYXRmb3JtbmFtZRI5CgxwbGF0Zm9ybXR5cGUYo8/MtAEgASgOMhEuc3Nt'
    'LlBsYXRmb3JtVHlwZVIMcGxhdGZvcm10eXBlEisKD3BsYXRmb3JtdmVyc2lvbhi/ptxCIAEoCV'
    'IPcGxhdGZvcm12ZXJzaW9uEjkKDHJlc291cmNldHlwZRjevtiPASABKA4yES5zc20uUmVzb3Vy'
    'Y2VUeXBlUgxyZXNvdXJjZXR5cGU=');

@$core.Deprecated('Use instanceInformationDescriptor instead')
const InstanceInformation$json = {
  '1': 'InstanceInformation',
  '2': [
    {'1': 'activationid', '3': 146858601, '4': 1, '5': 9, '10': 'activationid'},
    {'1': 'agentversion', '3': 173670635, '4': 1, '5': 9, '10': 'agentversion'},
    {
      '1': 'associationoverview',
      '3': 458852672,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAggregatedAssociationOverview',
      '10': 'associationoverview'
    },
    {
      '1': 'associationstatus',
      '3': 486964487,
      '4': 1,
      '5': 9,
      '10': 'associationstatus'
    },
    {'1': 'computername', '3': 67735292, '4': 1, '5': 9, '10': 'computername'},
    {'1': 'ipaddress', '3': 169333741, '4': 1, '5': 9, '10': 'ipaddress'},
    {'1': 'iamrole', '3': 242424351, '4': 1, '5': 9, '10': 'iamrole'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'islatestversion',
      '3': 27844379,
      '4': 1,
      '5': 8,
      '10': 'islatestversion'
    },
    {
      '1': 'lastassociationexecutiondate',
      '3': 37666619,
      '4': 1,
      '5': 9,
      '10': 'lastassociationexecutiondate'
    },
    {
      '1': 'lastpingdatetime',
      '3': 49842265,
      '4': 1,
      '5': 9,
      '10': 'lastpingdatetime'
    },
    {
      '1': 'lastsuccessfulassociationexecutiondate',
      '3': 162380165,
      '4': 1,
      '5': 9,
      '10': 'lastsuccessfulassociationexecutiondate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'pingstatus',
      '3': 282358380,
      '4': 1,
      '5': 14,
      '6': '.ssm.PingStatus',
      '10': 'pingstatus'
    },
    {'1': 'platformname', '3': 188917178, '4': 1, '5': 9, '10': 'platformname'},
    {
      '1': 'platformtype',
      '3': 378742691,
      '4': 1,
      '5': 14,
      '6': '.ssm.PlatformType',
      '10': 'platformtype'
    },
    {
      '1': 'platformversion',
      '3': 139924287,
      '4': 1,
      '5': 9,
      '10': 'platformversion'
    },
    {
      '1': 'registrationdate',
      '3': 370887405,
      '4': 1,
      '5': 9,
      '10': 'registrationdate'
    },
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceType',
      '10': 'resourcetype'
    },
    {'1': 'sourceid', '3': 425309766, '4': 1, '5': 9, '10': 'sourceid'},
    {
      '1': 'sourcetype',
      '3': 195731217,
      '4': 1,
      '5': 14,
      '6': '.ssm.SourceType',
      '10': 'sourcetype'
    },
  ],
};

/// Descriptor for `InstanceInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceInformationDescriptor = $convert.base64Decode(
    'ChNJbnN0YW5jZUluZm9ybWF0aW9uEiUKDGFjdGl2YXRpb25pZBjpxINGIAEoCVIMYWN0aXZhdG'
    'lvbmlkEiUKDGFnZW50dmVyc2lvbhjrgehSIAEoCVIMYWdlbnR2ZXJzaW9uEmAKE2Fzc29jaWF0'
    'aW9ub3ZlcnZpZXcYwJLm2gEgASgLMiouc3NtLkluc3RhbmNlQWdncmVnYXRlZEFzc29jaWF0aW'
    '9uT3ZlcnZpZXdSE2Fzc29jaWF0aW9ub3ZlcnZpZXcSMAoRYXNzb2NpYXRpb25zdGF0dXMYh/qZ'
    '6AEgASgJUhFhc3NvY2lhdGlvbnN0YXR1cxIlCgxjb21wdXRlcm5hbWUY/J2mICABKAlSDGNvbX'
    'B1dGVybmFtZRIfCglpcGFkZHJlc3MY7affUCABKAlSCWlwYWRkcmVzcxIbCgdpYW1yb2xlGJ+0'
    'zHMgASgJUgdpYW1yb2xlEiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQSKwoPaX'
    'NsYXRlc3R2ZXJzaW9uGJu+ow0gASgIUg9pc2xhdGVzdHZlcnNpb24SRQocbGFzdGFzc29jaWF0'
    'aW9uZXhlY3V0aW9uZGF0ZRi7/voRIAEoCVIcbGFzdGFzc29jaWF0aW9uZXhlY3V0aW9uZGF0ZR'
    'ItChBsYXN0cGluZ2RhdGV0aW1lGNmQ4hcgASgJUhBsYXN0cGluZ2RhdGV0aW1lElkKJmxhc3Rz'
    'dWNjZXNzZnVsYXNzb2NpYXRpb25leGVjdXRpb25kYXRlGIXztk0gASgJUiZsYXN0c3VjY2Vzc2'
    'Z1bGFzc29jaWF0aW9uZXhlY3V0aW9uZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjMKCnBp'
    'bmdzdGF0dXMY7OTRhgEgASgOMg8uc3NtLlBpbmdTdGF0dXNSCnBpbmdzdGF0dXMSJQoMcGxhdG'
    'Zvcm1uYW1lGLrLilogASgJUgxwbGF0Zm9ybW5hbWUSOQoMcGxhdGZvcm10eXBlGKPPzLQBIAEo'
    'DjIRLnNzbS5QbGF0Zm9ybVR5cGVSDHBsYXRmb3JtdHlwZRIrCg9wbGF0Zm9ybXZlcnNpb24Yv6'
    'bcQiABKAlSD3BsYXRmb3JtdmVyc2lvbhIuChByZWdpc3RyYXRpb25kYXRlGO2V7bABIAEoCVIQ'
    'cmVnaXN0cmF0aW9uZGF0ZRI5CgxyZXNvdXJjZXR5cGUY3r7YjwEgASgOMhEuc3NtLlJlc291cm'
    'NlVHlwZVIMcmVzb3VyY2V0eXBlEh4KCHNvdXJjZWlkGMbs5soBIAEoCVIIc291cmNlaWQSMgoK'
    'c291cmNldHlwZRiRvqpdIAEoDjIPLnNzbS5Tb3VyY2VUeXBlUgpzb3VyY2V0eXBl');

@$core.Deprecated('Use instanceInformationFilterDescriptor instead')
const InstanceInformationFilter$json = {
  '1': 'InstanceInformationFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.InstanceInformationFilterKey',
      '10': 'key'
    },
    {'1': 'valueset', '3': 250257195, '4': 3, '5': 9, '10': 'valueset'},
  ],
};

/// Descriptor for `InstanceInformationFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceInformationFilterDescriptor = $convert.base64Decode(
    'ChlJbnN0YW5jZUluZm9ybWF0aW9uRmlsdGVyEjYKA2tleRjtkNdAIAEoDjIhLnNzbS5JbnN0YW'
    '5jZUluZm9ybWF0aW9uRmlsdGVyS2V5UgNrZXkSHQoIdmFsdWVzZXQYq76qdyADKAlSCHZhbHVl'
    'c2V0');

@$core.Deprecated('Use instanceInformationStringFilterDescriptor instead')
const InstanceInformationStringFilter$json = {
  '1': 'InstanceInformationStringFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `InstanceInformationStringFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instanceInformationStringFilterDescriptor =
    $convert.base64Decode(
        'Ch9JbnN0YW5jZUluZm9ybWF0aW9uU3RyaW5nRmlsdGVyEhMKA2tleRiNkutoIAEoCVIDa2V5Eh'
        'kKBnZhbHVlcxjcxLRqIAMoCVIGdmFsdWVz');

@$core.Deprecated('Use instancePatchStateDescriptor instead')
const InstancePatchState$json = {
  '1': 'InstancePatchState',
  '2': [
    {
      '1': 'availablesecurityupdatecount',
      '3': 6211183,
      '4': 1,
      '5': 5,
      '10': 'availablesecurityupdatecount'
    },
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {
      '1': 'criticalnoncompliantcount',
      '3': 446575416,
      '4': 1,
      '5': 5,
      '10': 'criticalnoncompliantcount'
    },
    {'1': 'failedcount', '3': 121973932, '4': 1, '5': 5, '10': 'failedcount'},
    {
      '1': 'installoverridelist',
      '3': 487484089,
      '4': 1,
      '5': 9,
      '10': 'installoverridelist'
    },
    {
      '1': 'installedcount',
      '3': 300247257,
      '4': 1,
      '5': 5,
      '10': 'installedcount'
    },
    {
      '1': 'installedothercount',
      '3': 95481817,
      '4': 1,
      '5': 5,
      '10': 'installedothercount'
    },
    {
      '1': 'installedpendingrebootcount',
      '3': 401618817,
      '4': 1,
      '5': 5,
      '10': 'installedpendingrebootcount'
    },
    {
      '1': 'installedrejectedcount',
      '3': 456810059,
      '4': 1,
      '5': 5,
      '10': 'installedrejectedcount'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'lastnorebootinstalloperationtime',
      '3': 135277587,
      '4': 1,
      '5': 9,
      '10': 'lastnorebootinstalloperationtime'
    },
    {'1': 'missingcount', '3': 269209919, '4': 1, '5': 5, '10': 'missingcount'},
    {
      '1': 'notapplicablecount',
      '3': 403992231,
      '4': 1,
      '5': 5,
      '10': 'notapplicablecount'
    },
    {
      '1': 'operation',
      '3': 26084007,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchOperationType',
      '10': 'operation'
    },
    {
      '1': 'operationendtime',
      '3': 166452443,
      '4': 1,
      '5': 9,
      '10': 'operationendtime'
    },
    {
      '1': 'operationstarttime',
      '3': 120622464,
      '4': 1,
      '5': 9,
      '10': 'operationstarttime'
    },
    {
      '1': 'othernoncompliantcount',
      '3': 22394029,
      '4': 1,
      '5': 5,
      '10': 'othernoncompliantcount'
    },
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
    {
      '1': 'rebootoption',
      '3': 106110930,
      '4': 1,
      '5': 14,
      '6': '.ssm.RebootOption',
      '10': 'rebootoption'
    },
    {
      '1': 'securitynoncompliantcount',
      '3': 500346499,
      '4': 1,
      '5': 5,
      '10': 'securitynoncompliantcount'
    },
    {'1': 'snapshotid', '3': 99585817, '4': 1, '5': 9, '10': 'snapshotid'},
    {
      '1': 'unreportednotapplicablecount',
      '3': 295417459,
      '4': 1,
      '5': 5,
      '10': 'unreportednotapplicablecount'
    },
  ],
};

/// Descriptor for `InstancePatchState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instancePatchStateDescriptor = $convert.base64Decode(
    'ChJJbnN0YW5jZVBhdGNoU3RhdGUSRQocYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVjb3VudBjvjP'
    'sCIAEoBVIcYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVjb3VudBIhCgpiYXNlbGluZWlkGNDk2ygg'
    'ASgJUgpiYXNlbGluZWlkEkAKGWNyaXRpY2Fsbm9uY29tcGxpYW50Y291bnQYuOb41AEgASgFUh'
    'ljcml0aWNhbG5vbmNvbXBsaWFudGNvdW50EiMKC2ZhaWxlZGNvdW50GKzZlDogASgFUgtmYWls'
    'ZWRjb3VudBI0ChNpbnN0YWxsb3ZlcnJpZGVsaXN0GLnVuegBIAEoCVITaW5zdGFsbG92ZXJyaW'
    'RlbGlzdBIqCg5pbnN0YWxsZWRjb3VudBjZ0ZWPASABKAVSDmluc3RhbGxlZGNvdW50EjMKE2lu'
    'c3RhbGxlZG90aGVyY291bnQY2d/DLSABKAVSE2luc3RhbGxlZG90aGVyY291bnQSRAobaW5zdG'
    'FsbGVkcGVuZGluZ3JlYm9vdGNvdW50GIHvwL8BIAEoBVIbaW5zdGFsbGVkcGVuZGluZ3JlYm9v'
    'dGNvdW50EjoKFmluc3RhbGxlZHJlamVjdGVkY291bnQYy7zp2QEgASgFUhZpbnN0YWxsZWRyZW'
    'plY3RlZGNvdW50EiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQSTQogbGFzdG5v'
    'cmVib290aW5zdGFsbG9wZXJhdGlvbnRpbWUYk9jAQCABKAlSIGxhc3Rub3JlYm9vdGluc3RhbG'
    'xvcGVyYXRpb250aW1lEiYKDG1pc3Npbmdjb3VudBi/oq+AASABKAVSDG1pc3Npbmdjb3VudBIy'
    'ChJub3RhcHBsaWNhYmxlY291bnQYp93RwAEgASgFUhJub3RhcHBsaWNhYmxlY291bnQSOAoJb3'
    'BlcmF0aW9uGKeFuAwgASgOMhcuc3NtLlBhdGNoT3BlcmF0aW9uVHlwZVIJb3BlcmF0aW9uEi0K'
    'EG9wZXJhdGlvbmVuZHRpbWUY27mvTyABKAlSEG9wZXJhdGlvbmVuZHRpbWUSMQoSb3BlcmF0aW'
    '9uc3RhcnR0aW1lGICbwjkgASgJUhJvcGVyYXRpb25zdGFydHRpbWUSOQoWb3RoZXJub25jb21w'
    'bGlhbnRjb3VudBit6dYKIAEoBVIWb3RoZXJub25jb21wbGlhbnRjb3VudBItChBvd25lcmluZm'
    '9ybWF0aW9uGIOaxSAgASgJUhBvd25lcmluZm9ybWF0aW9uEiIKCnBhdGNoZ3JvdXAY4bex9wEg'
    'ASgJUgpwYXRjaGdyb3VwEjgKDHJlYm9vdG9wdGlvbhjSv8wyIAEoDjIRLnNzbS5SZWJvb3RPcH'
    'Rpb25SDHJlYm9vdG9wdGlvbhJAChlzZWN1cml0eW5vbmNvbXBsaWFudGNvdW50GIPdyu4BIAEo'
    'BVIZc2VjdXJpdHlub25jb21wbGlhbnRjb3VudBIhCgpzbmFwc2hvdGlkGJmevi8gASgJUgpzbm'
    'Fwc2hvdGlkEkYKHHVucmVwb3J0ZWRub3RhcHBsaWNhYmxlY291bnQY8+zujAEgASgFUhx1bnJl'
    'cG9ydGVkbm90YXBwbGljYWJsZWNvdW50');

@$core.Deprecated('Use instancePatchStateFilterDescriptor instead')
const InstancePatchStateFilter$json = {
  '1': 'InstancePatchStateFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.InstancePatchStateOperatorType',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `InstancePatchStateFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instancePatchStateFilterDescriptor = $convert.base64Decode(
    'ChhJbnN0YW5jZVBhdGNoU3RhdGVGaWx0ZXISEwoDa2V5GI2S62ggASgJUgNrZXkSOwoEdHlwZR'
    'juoNeKASABKA4yIy5zc20uSW5zdGFuY2VQYXRjaFN0YXRlT3BlcmF0b3JUeXBlUgR0eXBlEhkK'
    'BnZhbHVlcxjcxLRqIAMoCVIGdmFsdWVz');

@$core.Deprecated('Use instancePropertyDescriptor instead')
const InstanceProperty$json = {
  '1': 'InstanceProperty',
  '2': [
    {'1': 'activationid', '3': 146858601, '4': 1, '5': 9, '10': 'activationid'},
    {'1': 'agentversion', '3': 173670635, '4': 1, '5': 9, '10': 'agentversion'},
    {'1': 'architecture', '3': 420247267, '4': 1, '5': 9, '10': 'architecture'},
    {
      '1': 'associationoverview',
      '3': 458852672,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAggregatedAssociationOverview',
      '10': 'associationoverview'
    },
    {
      '1': 'associationstatus',
      '3': 486964487,
      '4': 1,
      '5': 9,
      '10': 'associationstatus'
    },
    {'1': 'computername', '3': 67735292, '4': 1, '5': 9, '10': 'computername'},
    {'1': 'ipaddress', '3': 169333741, '4': 1, '5': 9, '10': 'ipaddress'},
    {'1': 'iamrole', '3': 242424351, '4': 1, '5': 9, '10': 'iamrole'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'instancerole', '3': 143624115, '4': 1, '5': 9, '10': 'instancerole'},
    {
      '1': 'instancestate',
      '3': 217046808,
      '4': 1,
      '5': 9,
      '10': 'instancestate'
    },
    {'1': 'instancetype', '3': 46812483, '4': 1, '5': 9, '10': 'instancetype'},
    {'1': 'keyname', '3': 234857320, '4': 1, '5': 9, '10': 'keyname'},
    {
      '1': 'lastassociationexecutiondate',
      '3': 37666619,
      '4': 1,
      '5': 9,
      '10': 'lastassociationexecutiondate'
    },
    {
      '1': 'lastpingdatetime',
      '3': 49842265,
      '4': 1,
      '5': 9,
      '10': 'lastpingdatetime'
    },
    {
      '1': 'lastsuccessfulassociationexecutiondate',
      '3': 162380165,
      '4': 1,
      '5': 9,
      '10': 'lastsuccessfulassociationexecutiondate'
    },
    {'1': 'launchtime', '3': 338860032, '4': 1, '5': 9, '10': 'launchtime'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'pingstatus',
      '3': 282358380,
      '4': 1,
      '5': 14,
      '6': '.ssm.PingStatus',
      '10': 'pingstatus'
    },
    {'1': 'platformname', '3': 188917178, '4': 1, '5': 9, '10': 'platformname'},
    {
      '1': 'platformtype',
      '3': 378742691,
      '4': 1,
      '5': 14,
      '6': '.ssm.PlatformType',
      '10': 'platformtype'
    },
    {
      '1': 'platformversion',
      '3': 139924287,
      '4': 1,
      '5': 9,
      '10': 'platformversion'
    },
    {
      '1': 'registrationdate',
      '3': 370887405,
      '4': 1,
      '5': 9,
      '10': 'registrationdate'
    },
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'sourceid', '3': 425309766, '4': 1, '5': 9, '10': 'sourceid'},
    {
      '1': 'sourcetype',
      '3': 195731217,
      '4': 1,
      '5': 14,
      '6': '.ssm.SourceType',
      '10': 'sourcetype'
    },
  ],
};

/// Descriptor for `InstanceProperty`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instancePropertyDescriptor = $convert.base64Decode(
    'ChBJbnN0YW5jZVByb3BlcnR5EiUKDGFjdGl2YXRpb25pZBjpxINGIAEoCVIMYWN0aXZhdGlvbm'
    'lkEiUKDGFnZW50dmVyc2lvbhjrgehSIAEoCVIMYWdlbnR2ZXJzaW9uEiYKDGFyY2hpdGVjdHVy'
    'ZRjj7bHIASABKAlSDGFyY2hpdGVjdHVyZRJgChNhc3NvY2lhdGlvbm92ZXJ2aWV3GMCS5toBIA'
    'EoCzIqLnNzbS5JbnN0YW5jZUFnZ3JlZ2F0ZWRBc3NvY2lhdGlvbk92ZXJ2aWV3UhNhc3NvY2lh'
    'dGlvbm92ZXJ2aWV3EjAKEWFzc29jaWF0aW9uc3RhdHVzGIf6megBIAEoCVIRYXNzb2NpYXRpb2'
    '5zdGF0dXMSJQoMY29tcHV0ZXJuYW1lGPydpiAgASgJUgxjb21wdXRlcm5hbWUSHwoJaXBhZGRy'
    'ZXNzGO2n31AgASgJUglpcGFkZHJlc3MSGwoHaWFtcm9sZRiftMxzIAEoCVIHaWFtcm9sZRIhCg'
    'ppbnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEiUKDGluc3RhbmNlcm9sZRizj75EIAEo'
    'CVIMaW5zdGFuY2Vyb2xlEicKDWluc3RhbmNlc3RhdGUYmL6/ZyABKAlSDWluc3RhbmNlc3RhdG'
    'USJQoMaW5zdGFuY2V0eXBlGMOaqRYgASgJUgxpbnN0YW5jZXR5cGUSGwoHa2V5bmFtZRjoxv5v'
    'IAEoCVIHa2V5bmFtZRJFChxsYXN0YXNzb2NpYXRpb25leGVjdXRpb25kYXRlGLv++hEgASgJUh'
    'xsYXN0YXNzb2NpYXRpb25leGVjdXRpb25kYXRlEi0KEGxhc3RwaW5nZGF0ZXRpbWUY2ZDiFyAB'
    'KAlSEGxhc3RwaW5nZGF0ZXRpbWUSWQombGFzdHN1Y2Nlc3NmdWxhc3NvY2lhdGlvbmV4ZWN1dG'
    'lvbmRhdGUYhfO2TSABKAlSJmxhc3RzdWNjZXNzZnVsYXNzb2NpYXRpb25leGVjdXRpb25kYXRl'
    'EiIKCmxhdW5jaHRpbWUYgLDKoQEgASgJUgpsYXVuY2h0aW1lEhUKBG5hbWUYh+aBfyABKAlSBG'
    '5hbWUSMwoKcGluZ3N0YXR1cxjs5NGGASABKA4yDy5zc20uUGluZ1N0YXR1c1IKcGluZ3N0YXR1'
    'cxIlCgxwbGF0Zm9ybW5hbWUYusuKWiABKAlSDHBsYXRmb3JtbmFtZRI5CgxwbGF0Zm9ybXR5cG'
    'UYo8/MtAEgASgOMhEuc3NtLlBsYXRmb3JtVHlwZVIMcGxhdGZvcm10eXBlEisKD3BsYXRmb3Jt'
    'dmVyc2lvbhi/ptxCIAEoCVIPcGxhdGZvcm12ZXJzaW9uEi4KEHJlZ2lzdHJhdGlvbmRhdGUY7Z'
    'XtsAEgASgJUhByZWdpc3RyYXRpb25kYXRlEiYKDHJlc291cmNldHlwZRjevtiPASABKAlSDHJl'
    'c291cmNldHlwZRIeCghzb3VyY2VpZBjG7ObKASABKAlSCHNvdXJjZWlkEjIKCnNvdXJjZXR5cG'
    'UYkb6qXSABKA4yDy5zc20uU291cmNlVHlwZVIKc291cmNldHlwZQ==');

@$core.Deprecated('Use instancePropertyFilterDescriptor instead')
const InstancePropertyFilter$json = {
  '1': 'InstancePropertyFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.InstancePropertyFilterKey',
      '10': 'key'
    },
    {'1': 'valueset', '3': 250257195, '4': 3, '5': 9, '10': 'valueset'},
  ],
};

/// Descriptor for `InstancePropertyFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instancePropertyFilterDescriptor = $convert.base64Decode(
    'ChZJbnN0YW5jZVByb3BlcnR5RmlsdGVyEjMKA2tleRjtkNdAIAEoDjIeLnNzbS5JbnN0YW5jZV'
    'Byb3BlcnR5RmlsdGVyS2V5UgNrZXkSHQoIdmFsdWVzZXQYq76qdyADKAlSCHZhbHVlc2V0');

@$core.Deprecated('Use instancePropertyStringFilterDescriptor instead')
const InstancePropertyStringFilter$json = {
  '1': 'InstancePropertyStringFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'operator',
      '3': 31807518,
      '4': 1,
      '5': 14,
      '6': '.ssm.InstancePropertyFilterOperator',
      '10': 'operator'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `InstancePropertyStringFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List instancePropertyStringFilterDescriptor =
    $convert.base64Decode(
        'ChxJbnN0YW5jZVByb3BlcnR5U3RyaW5nRmlsdGVyEhMKA2tleRiNkutoIAEoCVIDa2V5EkIKCG'
        '9wZXJhdG9yGJ6wlQ8gASgOMiMuc3NtLkluc3RhbmNlUHJvcGVydHlGaWx0ZXJPcGVyYXRvclII'
        'b3BlcmF0b3ISGQoGdmFsdWVzGNzEtGogAygJUgZ2YWx1ZXM=');

@$core.Deprecated('Use internalServerErrorDescriptor instead')
const InternalServerError$json = {
  '1': 'InternalServerError',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServerError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServerErrorDescriptor =
    $convert.base64Decode(
        'ChNJbnRlcm5hbFNlcnZlckVycm9yEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidActivationDescriptor instead')
const InvalidActivation$json = {
  '1': 'InvalidActivation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidActivation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidActivationDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkQWN0aXZhdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidActivationIdDescriptor instead')
const InvalidActivationId$json = {
  '1': 'InvalidActivationId',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidActivationId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidActivationIdDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkQWN0aXZhdGlvbklkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidAggregatorExceptionDescriptor instead')
const InvalidAggregatorException$json = {
  '1': 'InvalidAggregatorException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAggregatorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAggregatorExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkQWdncmVnYXRvckV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidAllowedPatternExceptionDescriptor instead')
const InvalidAllowedPatternException$json = {
  '1': 'InvalidAllowedPatternException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAllowedPatternException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAllowedPatternExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkQWxsb3dlZFBhdHRlcm5FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidAssociationDescriptor instead')
const InvalidAssociation$json = {
  '1': 'InvalidAssociation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAssociation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAssociationDescriptor =
    $convert.base64Decode(
        'ChJJbnZhbGlkQXNzb2NpYXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidAssociationVersionDescriptor instead')
const InvalidAssociationVersion$json = {
  '1': 'InvalidAssociationVersion',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAssociationVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAssociationVersionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkQXNzb2NpYXRpb25WZXJzaW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated(
    'Use invalidAutomationExecutionParametersExceptionDescriptor instead')
const InvalidAutomationExecutionParametersException$json = {
  '1': 'InvalidAutomationExecutionParametersException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAutomationExecutionParametersException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    invalidAutomationExecutionParametersExceptionDescriptor =
    $convert.base64Decode(
        'Ci1JbnZhbGlkQXV0b21hdGlvbkV4ZWN1dGlvblBhcmFtZXRlcnNFeGNlcHRpb24SGwoHbWVzc2'
        'FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidAutomationSignalExceptionDescriptor instead')
const InvalidAutomationSignalException$json = {
  '1': 'InvalidAutomationSignalException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAutomationSignalException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAutomationSignalExceptionDescriptor =
    $convert.base64Decode(
        'CiBJbnZhbGlkQXV0b21hdGlvblNpZ25hbEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core
    .Deprecated('Use invalidAutomationStatusUpdateExceptionDescriptor instead')
const InvalidAutomationStatusUpdateException$json = {
  '1': 'InvalidAutomationStatusUpdateException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAutomationStatusUpdateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAutomationStatusUpdateExceptionDescriptor =
    $convert.base64Decode(
        'CiZJbnZhbGlkQXV0b21hdGlvblN0YXR1c1VwZGF0ZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3'
        'AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidCommandIdDescriptor instead')
const InvalidCommandId$json = {
  '1': 'InvalidCommandId',
};

/// Descriptor for `InvalidCommandId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidCommandIdDescriptor =
    $convert.base64Decode('ChBJbnZhbGlkQ29tbWFuZElk');

@$core.Deprecated(
    'Use invalidDeleteInventoryParametersExceptionDescriptor instead')
const InvalidDeleteInventoryParametersException$json = {
  '1': 'InvalidDeleteInventoryParametersException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDeleteInventoryParametersException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    invalidDeleteInventoryParametersExceptionDescriptor = $convert.base64Decode(
        'CilJbnZhbGlkRGVsZXRlSW52ZW50b3J5UGFyYW1ldGVyc0V4Y2VwdGlvbhIbCgdtZXNzYWdlGI'
        'Wzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidDeletionIdExceptionDescriptor instead')
const InvalidDeletionIdException$json = {
  '1': 'InvalidDeletionIdException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDeletionIdException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDeletionIdExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkRGVsZXRpb25JZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidDocumentDescriptor instead')
const InvalidDocument$json = {
  '1': 'InvalidDocument',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocument`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkRG9jdW1lbnQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidDocumentContentDescriptor instead')
const InvalidDocumentContent$json = {
  '1': 'InvalidDocumentContent',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocumentContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentContentDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkRG9jdW1lbnRDb250ZW50EhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidDocumentOperationDescriptor instead')
const InvalidDocumentOperation$json = {
  '1': 'InvalidDocumentOperation',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocumentOperation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentOperationDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkRG9jdW1lbnRPcGVyYXRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidDocumentSchemaVersionDescriptor instead')
const InvalidDocumentSchemaVersion$json = {
  '1': 'InvalidDocumentSchemaVersion',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocumentSchemaVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentSchemaVersionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkRG9jdW1lbnRTY2hlbWFWZXJzaW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidDocumentTypeDescriptor instead')
const InvalidDocumentType$json = {
  '1': 'InvalidDocumentType',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocumentType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentTypeDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkRG9jdW1lbnRUeXBlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidDocumentVersionDescriptor instead')
const InvalidDocumentVersion$json = {
  '1': 'InvalidDocumentVersion',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidDocumentVersion`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidDocumentVersionDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkRG9jdW1lbnRWZXJzaW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidFilterDescriptor instead')
const InvalidFilter$json = {
  '1': 'InvalidFilter',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFilterDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkRmlsdGVyEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidFilterKeyDescriptor instead')
const InvalidFilterKey$json = {
  '1': 'InvalidFilterKey',
};

/// Descriptor for `InvalidFilterKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFilterKeyDescriptor =
    $convert.base64Decode('ChBJbnZhbGlkRmlsdGVyS2V5');

@$core.Deprecated('Use invalidFilterOptionDescriptor instead')
const InvalidFilterOption$json = {
  '1': 'InvalidFilterOption',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidFilterOption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFilterOptionDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkRmlsdGVyT3B0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidFilterValueDescriptor instead')
const InvalidFilterValue$json = {
  '1': 'InvalidFilterValue',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidFilterValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidFilterValueDescriptor =
    $convert.base64Decode(
        'ChJJbnZhbGlkRmlsdGVyVmFsdWUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidInstanceIdDescriptor instead')
const InvalidInstanceId$json = {
  '1': 'InvalidInstanceId',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInstanceId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInstanceIdDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkSW5zdGFuY2VJZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidInstanceInformationFilterValueDescriptor instead')
const InvalidInstanceInformationFilterValue$json = {
  '1': 'InvalidInstanceInformationFilterValue',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInstanceInformationFilterValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInstanceInformationFilterValueDescriptor =
    $convert.base64Decode(
        'CiVJbnZhbGlkSW5zdGFuY2VJbmZvcm1hdGlvbkZpbHRlclZhbHVlEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidInstancePropertyFilterValueDescriptor instead')
const InvalidInstancePropertyFilterValue$json = {
  '1': 'InvalidInstancePropertyFilterValue',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInstancePropertyFilterValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInstancePropertyFilterValueDescriptor =
    $convert.base64Decode(
        'CiJJbnZhbGlkSW5zdGFuY2VQcm9wZXJ0eUZpbHRlclZhbHVlEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidInventoryGroupExceptionDescriptor instead')
const InvalidInventoryGroupException$json = {
  '1': 'InvalidInventoryGroupException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInventoryGroupException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInventoryGroupExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkSW52ZW50b3J5R3JvdXBFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidInventoryItemContextExceptionDescriptor instead')
const InvalidInventoryItemContextException$json = {
  '1': 'InvalidInventoryItemContextException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInventoryItemContextException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInventoryItemContextExceptionDescriptor =
    $convert.base64Decode(
        'CiRJbnZhbGlkSW52ZW50b3J5SXRlbUNvbnRleHRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidInventoryRequestExceptionDescriptor instead')
const InvalidInventoryRequestException$json = {
  '1': 'InvalidInventoryRequestException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidInventoryRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidInventoryRequestExceptionDescriptor =
    $convert.base64Decode(
        'CiBJbnZhbGlkSW52ZW50b3J5UmVxdWVzdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use invalidItemContentExceptionDescriptor instead')
const InvalidItemContentException$json = {
  '1': 'InvalidItemContentException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `InvalidItemContentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidItemContentExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkSXRlbUNvbnRlbnRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2'
        'FnZRIeCgh0eXBlbmFtZRjPztnUASABKAlSCHR5cGVuYW1l');

@$core.Deprecated('Use invalidKeyIdDescriptor instead')
const InvalidKeyId$json = {
  '1': 'InvalidKeyId',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKeyId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKeyIdDescriptor = $convert.base64Decode(
    'CgxJbnZhbGlkS2V5SWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidNextTokenDescriptor instead')
const InvalidNextToken$json = {
  '1': 'InvalidNextToken',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNextToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNextTokenDescriptor = $convert.base64Decode(
    'ChBJbnZhbGlkTmV4dFRva2VuEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidNotificationConfigDescriptor instead')
const InvalidNotificationConfig$json = {
  '1': 'InvalidNotificationConfig',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNotificationConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNotificationConfigDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkTm90aWZpY2F0aW9uQ29uZmlnEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidOptionExceptionDescriptor instead')
const InvalidOptionException$json = {
  '1': 'InvalidOptionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidOptionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOptionExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkT3B0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidOutputFolderDescriptor instead')
const InvalidOutputFolder$json = {
  '1': 'InvalidOutputFolder',
};

/// Descriptor for `InvalidOutputFolder`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOutputFolderDescriptor =
    $convert.base64Decode('ChNJbnZhbGlkT3V0cHV0Rm9sZGVy');

@$core.Deprecated('Use invalidOutputLocationDescriptor instead')
const InvalidOutputLocation$json = {
  '1': 'InvalidOutputLocation',
};

/// Descriptor for `InvalidOutputLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidOutputLocationDescriptor =
    $convert.base64Decode('ChVJbnZhbGlkT3V0cHV0TG9jYXRpb24=');

@$core.Deprecated('Use invalidParametersDescriptor instead')
const InvalidParameters$json = {
  '1': 'InvalidParameters',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParametersDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkUGFyYW1ldGVycxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidPermissionTypeDescriptor instead')
const InvalidPermissionType$json = {
  '1': 'InvalidPermissionType',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPermissionType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPermissionTypeDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkUGVybWlzc2lvblR5cGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidPluginNameDescriptor instead')
const InvalidPluginName$json = {
  '1': 'InvalidPluginName',
};

/// Descriptor for `InvalidPluginName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPluginNameDescriptor =
    $convert.base64Decode('ChFJbnZhbGlkUGx1Z2luTmFtZQ==');

@$core.Deprecated('Use invalidPolicyAttributeExceptionDescriptor instead')
const InvalidPolicyAttributeException$json = {
  '1': 'InvalidPolicyAttributeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPolicyAttributeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPolicyAttributeExceptionDescriptor =
    $convert.base64Decode(
        'Ch9JbnZhbGlkUG9saWN5QXR0cmlidXRlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use invalidPolicyTypeExceptionDescriptor instead')
const InvalidPolicyTypeException$json = {
  '1': 'InvalidPolicyTypeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidPolicyTypeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidPolicyTypeExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkUG9saWN5VHlwZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidResourceIdDescriptor instead')
const InvalidResourceId$json = {
  '1': 'InvalidResourceId',
};

/// Descriptor for `InvalidResourceId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidResourceIdDescriptor =
    $convert.base64Decode('ChFJbnZhbGlkUmVzb3VyY2VJZA==');

@$core.Deprecated('Use invalidResourceTypeDescriptor instead')
const InvalidResourceType$json = {
  '1': 'InvalidResourceType',
};

/// Descriptor for `InvalidResourceType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidResourceTypeDescriptor =
    $convert.base64Decode('ChNJbnZhbGlkUmVzb3VyY2VUeXBl');

@$core.Deprecated('Use invalidResultAttributeExceptionDescriptor instead')
const InvalidResultAttributeException$json = {
  '1': 'InvalidResultAttributeException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidResultAttributeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidResultAttributeExceptionDescriptor =
    $convert.base64Decode(
        'Ch9JbnZhbGlkUmVzdWx0QXR0cmlidXRlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use invalidRoleDescriptor instead')
const InvalidRole$json = {
  '1': 'InvalidRole',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidRole`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidRoleDescriptor = $convert
    .base64Decode('CgtJbnZhbGlkUm9sZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidScheduleDescriptor instead')
const InvalidSchedule$json = {
  '1': 'InvalidSchedule',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSchedule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidScheduleDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkU2NoZWR1bGUSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidTagDescriptor instead')
const InvalidTag$json = {
  '1': 'InvalidTag',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTagDescriptor = $convert
    .base64Decode('CgpJbnZhbGlkVGFnEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidTargetDescriptor instead')
const InvalidTarget$json = {
  '1': 'InvalidTarget',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTargetDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkVGFyZ2V0EhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidTargetMapsDescriptor instead')
const InvalidTargetMaps$json = {
  '1': 'InvalidTargetMaps',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTargetMaps`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTargetMapsDescriptor = $convert.base64Decode(
    'ChFJbnZhbGlkVGFyZ2V0TWFwcxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidTypeNameExceptionDescriptor instead')
const InvalidTypeNameException$json = {
  '1': 'InvalidTypeNameException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidTypeNameException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidTypeNameExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkVHlwZU5hbWVFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidUpdateDescriptor instead')
const InvalidUpdate$json = {
  '1': 'InvalidUpdate',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidUpdateDescriptor = $convert.base64Decode(
    'Cg1JbnZhbGlkVXBkYXRlEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use inventoryAggregatorDescriptor instead')
const InventoryAggregator$json = {
  '1': 'InventoryAggregator',
  '2': [
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryAggregator',
      '10': 'aggregators'
    },
    {'1': 'expression', '3': 193051916, '4': 1, '5': 9, '10': 'expression'},
    {
      '1': 'groups',
      '3': 360662414,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryGroup',
      '10': 'groups'
    },
  ],
};

/// Descriptor for `InventoryAggregator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryAggregatorDescriptor = $convert.base64Decode(
    'ChNJbnZlbnRvcnlBZ2dyZWdhdG9yEj0KC2FnZ3JlZ2F0b3JzGI79g00gAygLMhguc3NtLkludm'
    'VudG9yeUFnZ3JlZ2F0b3JSC2FnZ3JlZ2F0b3JzEiEKCmV4cHJlc3Npb24YjPqGXCABKAlSCmV4'
    'cHJlc3Npb24SLwoGZ3JvdXBzGI6L/asBIAMoCzITLnNzbS5JbnZlbnRvcnlHcm91cFIGZ3JvdX'
    'Bz');

@$core.Deprecated('Use inventoryDeletionStatusItemDescriptor instead')
const InventoryDeletionStatusItem$json = {
  '1': 'InventoryDeletionStatusItem',
  '2': [
    {'1': 'deletionid', '3': 126693587, '4': 1, '5': 9, '10': 'deletionid'},
    {
      '1': 'deletionstarttime',
      '3': 413142901,
      '4': 1,
      '5': 9,
      '10': 'deletionstarttime'
    },
    {
      '1': 'deletionsummary',
      '3': 189954074,
      '4': 1,
      '5': 11,
      '6': '.ssm.InventoryDeletionSummary',
      '10': 'deletionsummary'
    },
    {
      '1': 'laststatus',
      '3': 163326556,
      '4': 1,
      '5': 14,
      '6': '.ssm.InventoryDeletionStatus',
      '10': 'laststatus'
    },
    {
      '1': 'laststatusmessage',
      '3': 59961987,
      '4': 1,
      '5': 9,
      '10': 'laststatusmessage'
    },
    {
      '1': 'laststatusupdatetime',
      '3': 66342972,
      '4': 1,
      '5': 9,
      '10': 'laststatusupdatetime'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `InventoryDeletionStatusItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryDeletionStatusItemDescriptor = $convert.base64Decode(
    'ChtJbnZlbnRvcnlEZWxldGlvblN0YXR1c0l0ZW0SIQoKZGVsZXRpb25pZBjT4bQ8IAEoCVIKZG'
    'VsZXRpb25pZBIwChFkZWxldGlvbnN0YXJ0dGltZRj1noDFASABKAlSEWRlbGV0aW9uc3RhcnR0'
    'aW1lEkoKD2RlbGV0aW9uc3VtbWFyeRia8MlaIAEoCzIdLnNzbS5JbnZlbnRvcnlEZWxldGlvbl'
    'N1bW1hcnlSD2RlbGV0aW9uc3VtbWFyeRI/CgpsYXN0c3RhdHVzGNzU8E0gASgOMhwuc3NtLklu'
    'dmVudG9yeURlbGV0aW9uU3RhdHVzUgpsYXN0c3RhdHVzEi8KEWxhc3RzdGF0dXNtZXNzYWdlGI'
    'PlyxwgASgJUhFsYXN0c3RhdHVzbWVzc2FnZRI1ChRsYXN0c3RhdHVzdXBkYXRldGltZRi8oNEf'
    'IAEoCVIUbGFzdHN0YXR1c3VwZGF0ZXRpbWUSHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbm'
    'FtZQ==');

@$core.Deprecated('Use inventoryDeletionSummaryDescriptor instead')
const InventoryDeletionSummary$json = {
  '1': 'InventoryDeletionSummary',
  '2': [
    {
      '1': 'remainingcount',
      '3': 230680675,
      '4': 1,
      '5': 5,
      '10': 'remainingcount'
    },
    {
      '1': 'summaryitems',
      '3': 103452326,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryDeletionSummaryItem',
      '10': 'summaryitems'
    },
    {'1': 'totalcount', '3': 502519869, '4': 1, '5': 5, '10': 'totalcount'},
  ],
};

/// Descriptor for `InventoryDeletionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryDeletionSummaryDescriptor = $convert.base64Decode(
    'ChhJbnZlbnRvcnlEZWxldGlvblN1bW1hcnkSKQoOcmVtYWluaW5nY291bnQY49D/bSABKAVSDn'
    'JlbWFpbmluZ2NvdW50EkgKDHN1bW1hcnlpdGVtcximnaoxIAMoCzIhLnNzbS5JbnZlbnRvcnlE'
    'ZWxldGlvblN1bW1hcnlJdGVtUgxzdW1tYXJ5aXRlbXMSIgoKdG90YWxjb3VudBi9sM/vASABKA'
    'VSCnRvdGFsY291bnQ=');

@$core.Deprecated('Use inventoryDeletionSummaryItemDescriptor instead')
const InventoryDeletionSummaryItem$json = {
  '1': 'InventoryDeletionSummaryItem',
  '2': [
    {'1': 'count', '3': 31963285, '4': 1, '5': 5, '10': 'count'},
    {
      '1': 'remainingcount',
      '3': 230680675,
      '4': 1,
      '5': 5,
      '10': 'remainingcount'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `InventoryDeletionSummaryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryDeletionSummaryItemDescriptor =
    $convert.base64Decode(
        'ChxJbnZlbnRvcnlEZWxldGlvblN1bW1hcnlJdGVtEhcKBWNvdW50GJXxng8gASgFUgVjb3VudB'
        'IpCg5yZW1haW5pbmdjb3VudBjj0P9tIAEoBVIOcmVtYWluaW5nY291bnQSHAoHdmVyc2lvbhi4'
        'qrfuASABKAlSB3ZlcnNpb24=');

@$core.Deprecated('Use inventoryFilterDescriptor instead')
const InventoryFilter$json = {
  '1': 'InventoryFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.InventoryQueryOperatorType',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `InventoryFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryFilterDescriptor = $convert.base64Decode(
    'Cg9JbnZlbnRvcnlGaWx0ZXISEwoDa2V5GI2S62ggASgJUgNrZXkSNwoEdHlwZRjuoNeKASABKA'
    '4yHy5zc20uSW52ZW50b3J5UXVlcnlPcGVyYXRvclR5cGVSBHR5cGUSGQoGdmFsdWVzGNzEtGog'
    'AygJUgZ2YWx1ZXM=');

@$core.Deprecated('Use inventoryGroupDescriptor instead')
const InventoryGroup$json = {
  '1': 'InventoryGroup',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryFilter',
      '10': 'filters'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `InventoryGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryGroupDescriptor = $convert.base64Decode(
    'Cg5JbnZlbnRvcnlHcm91cBIxCgdmaWx0ZXJzGO3N6lkgAygLMhQuc3NtLkludmVudG9yeUZpbH'
    'RlclIHZmlsdGVycxIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use inventoryItemDescriptor instead')
const InventoryItem$json = {
  '1': 'InventoryItem',
  '2': [
    {'1': 'capturetime', '3': 72276089, '4': 1, '5': 9, '10': 'capturetime'},
    {
      '1': 'content',
      '3': 23568227,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryItem.ContentEntry',
      '10': 'content'
    },
    {'1': 'contenthash', '3': 99301015, '4': 1, '5': 9, '10': 'contenthash'},
    {
      '1': 'context',
      '3': 489783069,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryItem.ContextEntry',
      '10': 'context'
    },
    {
      '1': 'schemaversion',
      '3': 371681851,
      '4': 1,
      '5': 9,
      '10': 'schemaversion'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
  '3': [InventoryItem_ContentEntry$json, InventoryItem_ContextEntry$json],
};

@$core.Deprecated('Use inventoryItemDescriptor instead')
const InventoryItem_ContentEntry$json = {
  '1': 'ContentEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use inventoryItemDescriptor instead')
const InventoryItem_ContextEntry$json = {
  '1': 'ContextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InventoryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryItemDescriptor = $convert.base64Decode(
    'Cg1JbnZlbnRvcnlJdGVtEiMKC2NhcHR1cmV0aW1lGPmwuyIgASgJUgtjYXB0dXJldGltZRI8Cg'
    'djb250ZW50GOO+ngsgAygLMh8uc3NtLkludmVudG9yeUl0ZW0uQ29udGVudEVudHJ5Ugdjb250'
    'ZW50EiMKC2NvbnRlbnRoYXNoGJftrC8gASgJUgtjb250ZW50aGFzaBI9Cgdjb250ZXh0GJ3+xe'
    'kBIAMoCzIfLnNzbS5JbnZlbnRvcnlJdGVtLkNvbnRleHRFbnRyeVIHY29udGV4dBIoCg1zY2hl'
    'bWF2ZXJzaW9uGLvUnbEBIAEoCVINc2NoZW1hdmVyc2lvbhIeCgh0eXBlbmFtZRjPztnUASABKA'
    'lSCHR5cGVuYW1lGjoKDENvbnRlbnRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgC'
    'IAEoCVIFdmFsdWU6AjgBGjoKDENvbnRleHRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YW'
    'x1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use inventoryItemAttributeDescriptor instead')
const InventoryItemAttribute$json = {
  '1': 'InventoryItemAttribute',
  '2': [
    {
      '1': 'datatype',
      '3': 67988590,
      '4': 1,
      '5': 14,
      '6': '.ssm.InventoryAttributeDataType',
      '10': 'datatype'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `InventoryItemAttribute`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryItemAttributeDescriptor = $convert.base64Decode(
    'ChZJbnZlbnRvcnlJdGVtQXR0cmlidXRlEj4KCGRhdGF0eXBlGO7YtSAgASgOMh8uc3NtLkludm'
    'VudG9yeUF0dHJpYnV0ZURhdGFUeXBlUghkYXRhdHlwZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use inventoryItemSchemaDescriptor instead')
const InventoryItemSchema$json = {
  '1': 'InventoryItemSchema',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryItemAttribute',
      '10': 'attributes'
    },
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `InventoryItemSchema`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryItemSchemaDescriptor = $convert.base64Decode(
    'ChNJbnZlbnRvcnlJdGVtU2NoZW1hEj4KCmF0dHJpYnV0ZXMYtan7YyADKAsyGy5zc20uSW52ZW'
    '50b3J5SXRlbUF0dHJpYnV0ZVIKYXR0cmlidXRlcxIkCgtkaXNwbGF5bmFtZRi3ybLHASABKAlS'
    'C2Rpc3BsYXluYW1lEh4KCHR5cGVuYW1lGM/O2dQBIAEoCVIIdHlwZW5hbWUSHAoHdmVyc2lvbh'
    'i4qrfuASABKAlSB3ZlcnNpb24=');

@$core.Deprecated('Use inventoryResultEntityDescriptor instead')
const InventoryResultEntity$json = {
  '1': 'InventoryResultEntity',
  '2': [
    {
      '1': 'data',
      '3': 525498822,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryResultEntity.DataEntry',
      '10': 'data'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
  '3': [InventoryResultEntity_DataEntry$json],
};

@$core.Deprecated('Use inventoryResultEntityDescriptor instead')
const InventoryResultEntity_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.InventoryResultItem',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `InventoryResultEntity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryResultEntityDescriptor = $convert.base64Decode(
    'ChVJbnZlbnRvcnlSZXN1bHRFbnRpdHkSPAoEZGF0YRjG88n6ASADKAsyJC5zc20uSW52ZW50b3'
    'J5UmVzdWx0RW50aXR5LkRhdGFFbnRyeVIEZGF0YRISCgJpZBiB8qK3ASABKAlSAmlkGlEKCURh'
    'dGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIuCgV2YWx1ZRgCIAEoCzIYLnNzbS5JbnZlbnRvcn'
    'lSZXN1bHRJdGVtUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use inventoryResultItemDescriptor instead')
const InventoryResultItem$json = {
  '1': 'InventoryResultItem',
  '2': [
    {'1': 'capturetime', '3': 72276089, '4': 1, '5': 9, '10': 'capturetime'},
    {
      '1': 'content',
      '3': 23568227,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryResultItem.ContentEntry',
      '10': 'content'
    },
    {'1': 'contenthash', '3': 99301015, '4': 1, '5': 9, '10': 'contenthash'},
    {
      '1': 'schemaversion',
      '3': 371681851,
      '4': 1,
      '5': 9,
      '10': 'schemaversion'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
  '3': [InventoryResultItem_ContentEntry$json],
};

@$core.Deprecated('Use inventoryResultItemDescriptor instead')
const InventoryResultItem_ContentEntry$json = {
  '1': 'ContentEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `InventoryResultItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryResultItemDescriptor = $convert.base64Decode(
    'ChNJbnZlbnRvcnlSZXN1bHRJdGVtEiMKC2NhcHR1cmV0aW1lGPmwuyIgASgJUgtjYXB0dXJldG'
    'ltZRJCCgdjb250ZW50GOO+ngsgAygLMiUuc3NtLkludmVudG9yeVJlc3VsdEl0ZW0uQ29udGVu'
    'dEVudHJ5Ugdjb250ZW50EiMKC2NvbnRlbnRoYXNoGJftrC8gASgJUgtjb250ZW50aGFzaBIoCg'
    '1zY2hlbWF2ZXJzaW9uGLvUnbEBIAEoCVINc2NoZW1hdmVyc2lvbhIeCgh0eXBlbmFtZRjPztnU'
    'ASABKAlSCHR5cGVuYW1lGjoKDENvbnRlbnRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YW'
    'x1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use invocationDoesNotExistDescriptor instead')
const InvocationDoesNotExist$json = {
  '1': 'InvocationDoesNotExist',
};

/// Descriptor for `InvocationDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invocationDoesNotExistDescriptor =
    $convert.base64Decode('ChZJbnZvY2F0aW9uRG9lc05vdEV4aXN0');

@$core.Deprecated('Use itemContentMismatchExceptionDescriptor instead')
const ItemContentMismatchException$json = {
  '1': 'ItemContentMismatchException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `ItemContentMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List itemContentMismatchExceptionDescriptor =
    $convert.base64Decode(
        'ChxJdGVtQ29udGVudE1pc21hdGNoRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2USHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ==');

@$core.Deprecated('Use itemSizeLimitExceededExceptionDescriptor instead')
const ItemSizeLimitExceededException$json = {
  '1': 'ItemSizeLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `ItemSizeLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List itemSizeLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JdGVtU2l6ZUxpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZRIeCgh0eXBlbmFtZRjPztnUASABKAlSCHR5cGVuYW1l');

@$core.Deprecated('Use labelParameterVersionRequestDescriptor instead')
const LabelParameterVersionRequest$json = {
  '1': 'LabelParameterVersionRequest',
  '2': [
    {'1': 'labels', '3': 178416811, '4': 3, '5': 9, '10': 'labels'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameterversion',
      '3': 141162077,
      '4': 1,
      '5': 3,
      '10': 'parameterversion'
    },
  ],
};

/// Descriptor for `LabelParameterVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List labelParameterVersionRequestDescriptor =
    $convert.base64Decode(
        'ChxMYWJlbFBhcmFtZXRlclZlcnNpb25SZXF1ZXN0EhkKBmxhYmVscxir2YlVIAMoCVIGbGFiZW'
        'xzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLQoQcGFyYW1ldGVydmVyc2lvbhjd7KdDIAEoA1IQ'
        'cGFyYW1ldGVydmVyc2lvbg==');

@$core.Deprecated('Use labelParameterVersionResultDescriptor instead')
const LabelParameterVersionResult$json = {
  '1': 'LabelParameterVersionResult',
  '2': [
    {
      '1': 'invalidlabels',
      '3': 281240648,
      '4': 3,
      '5': 9,
      '10': 'invalidlabels'
    },
    {
      '1': 'parameterversion',
      '3': 141162077,
      '4': 1,
      '5': 3,
      '10': 'parameterversion'
    },
  ],
};

/// Descriptor for `LabelParameterVersionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List labelParameterVersionResultDescriptor =
    $convert.base64Decode(
        'ChtMYWJlbFBhcmFtZXRlclZlcnNpb25SZXN1bHQSKAoNaW52YWxpZGxhYmVscxjIyI2GASADKA'
        'lSDWludmFsaWRsYWJlbHMSLQoQcGFyYW1ldGVydmVyc2lvbhjd7KdDIAEoA1IQcGFyYW1ldGVy'
        'dmVyc2lvbg==');

@$core.Deprecated('Use listAssociationVersionsRequestDescriptor instead')
const ListAssociationVersionsRequest$json = {
  '1': 'ListAssociationVersionsRequest',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAssociationVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssociationVersionsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0QXNzb2NpYXRpb25WZXJzaW9uc1JlcXVlc3QSJwoNYXNzb2NpYXRpb25pZBiS/JVCIA'
        'EoCVINYXNzb2NpYXRpb25pZBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIf'
        'CgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listAssociationVersionsResultDescriptor instead')
const ListAssociationVersionsResult$json = {
  '1': 'ListAssociationVersionsResult',
  '2': [
    {
      '1': 'associationversions',
      '3': 161692378,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationVersionInfo',
      '10': 'associationversions'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAssociationVersionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssociationVersionsResultDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0QXNzb2NpYXRpb25WZXJzaW9uc1Jlc3VsdBJQChNhc3NvY2lhdGlvbnZlcnNpb25zGN'
        'r1jE0gAygLMhsuc3NtLkFzc29jaWF0aW9uVmVyc2lvbkluZm9SE2Fzc29jaWF0aW9udmVyc2lv'
        'bnMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listAssociationsRequestDescriptor instead')
const ListAssociationsRequest$json = {
  '1': 'ListAssociationsRequest',
  '2': [
    {
      '1': 'associationfilterlist',
      '3': 464032811,
      '4': 3,
      '5': 11,
      '6': '.ssm.AssociationFilter',
      '10': 'associationfilterlist'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAssociationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssociationsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0QXNzb2NpYXRpb25zUmVxdWVzdBJQChVhc3NvY2lhdGlvbmZpbHRlcmxpc3QYq6ii3Q'
    'EgAygLMhYuc3NtLkFzc29jaWF0aW9uRmlsdGVyUhVhc3NvY2lhdGlvbmZpbHRlcmxpc3QSIgoK'
    'bWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUg'
    'luZXh0dG9rZW4=');

@$core.Deprecated('Use listAssociationsResultDescriptor instead')
const ListAssociationsResult$json = {
  '1': 'ListAssociationsResult',
  '2': [
    {
      '1': 'associations',
      '3': 149658718,
      '4': 3,
      '5': 11,
      '6': '.ssm.Association',
      '10': 'associations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListAssociationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssociationsResultDescriptor = $convert.base64Decode(
    'ChZMaXN0QXNzb2NpYXRpb25zUmVzdWx0EjcKDGFzc29jaWF0aW9ucxjeuK5HIAMoCzIQLnNzbS'
    '5Bc3NvY2lhdGlvblIMYXNzb2NpYXRpb25zEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRv'
    'a2Vu');

@$core.Deprecated('Use listCommandInvocationsRequestDescriptor instead')
const ListCommandInvocationsRequest$json = {
  '1': 'ListCommandInvocationsRequest',
  '2': [
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {'1': 'details', '3': 247611974, '4': 1, '5': 8, '10': 'details'},
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.CommandFilter',
      '10': 'filters'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCommandInvocationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommandInvocationsRequestDescriptor = $convert.base64Decode(
    'Ch1MaXN0Q29tbWFuZEludm9jYXRpb25zUmVxdWVzdBIfCgljb21tYW5kaWQYgNuATCABKAlSCW'
    'NvbW1hbmRpZBIbCgdkZXRhaWxzGMaEiXYgASgIUgdkZXRhaWxzEi8KB2ZpbHRlcnMY7c3qWSAD'
    'KAsyEi5zc20uQ29tbWFuZEZpbHRlclIHZmlsdGVycxIhCgppbnN0YW5jZWlkGKCt0RcgASgJUg'
    'ppbnN0YW5jZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0'
    'b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listCommandInvocationsResultDescriptor instead')
const ListCommandInvocationsResult$json = {
  '1': 'ListCommandInvocationsResult',
  '2': [
    {
      '1': 'commandinvocations',
      '3': 19955226,
      '4': 3,
      '5': 11,
      '6': '.ssm.CommandInvocation',
      '10': 'commandinvocations'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCommandInvocationsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommandInvocationsResultDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q29tbWFuZEludm9jYXRpb25zUmVzdWx0EkkKEmNvbW1hbmRpbnZvY2F0aW9ucxia/M'
        'EJIAMoCzIWLnNzbS5Db21tYW5kSW52b2NhdGlvblISY29tbWFuZGludm9jYXRpb25zEh8KCW5l'
        'eHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listCommandsRequestDescriptor instead')
const ListCommandsRequest$json = {
  '1': 'ListCommandsRequest',
  '2': [
    {'1': 'commandid', '3': 159395200, '4': 1, '5': 9, '10': 'commandid'},
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.CommandFilter',
      '10': 'filters'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCommandsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommandsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Q29tbWFuZHNSZXF1ZXN0Eh8KCWNvbW1hbmRpZBiA24BMIAEoCVIJY29tbWFuZGlkEi'
    '8KB2ZpbHRlcnMY7c3qWSADKAsyEi5zc20uQ29tbWFuZEZpbHRlclIHZmlsdGVycxIhCgppbnN0'
    'YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYX'
    'hyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listCommandsResultDescriptor instead')
const ListCommandsResult$json = {
  '1': 'ListCommandsResult',
  '2': [
    {
      '1': 'commands',
      '3': 387272948,
      '4': 3,
      '5': 11,
      '6': '.ssm.Command',
      '10': 'commands'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCommandsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommandsResultDescriptor = $convert.base64Decode(
    'ChJMaXN0Q29tbWFuZHNSZXN1bHQSLAoIY29tbWFuZHMY9KHVuAEgAygLMgwuc3NtLkNvbW1hbm'
    'RSCGNvbW1hbmRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listComplianceItemsRequestDescriptor instead')
const ListComplianceItemsRequest$json = {
  '1': 'ListComplianceItemsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceStringFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'resourceids', '3': 23528154, '4': 3, '5': 9, '10': 'resourceids'},
    {
      '1': 'resourcetypes',
      '3': 343086443,
      '4': 3,
      '5': 9,
      '10': 'resourcetypes'
    },
  ],
};

/// Descriptor for `ListComplianceItemsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listComplianceItemsRequestDescriptor = $convert.base64Decode(
    'ChpMaXN0Q29tcGxpYW5jZUl0ZW1zUmVxdWVzdBI4CgdmaWx0ZXJzGO3N6lkgAygLMhsuc3NtLk'
    'NvbXBsaWFuY2VTdHJpbmdGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVS'
    'Cm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIwoLcmVzb3VyY2'
    'VpZHMY2oWcCyADKAlSC3Jlc291cmNlaWRzEigKDXJlc291cmNldHlwZXMY66rMowEgAygJUg1y'
    'ZXNvdXJjZXR5cGVz');

@$core.Deprecated('Use listComplianceItemsResultDescriptor instead')
const ListComplianceItemsResult$json = {
  '1': 'ListComplianceItemsResult',
  '2': [
    {
      '1': 'complianceitems',
      '3': 129846551,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceItem',
      '10': 'complianceitems'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListComplianceItemsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listComplianceItemsResultDescriptor = $convert.base64Decode(
    'ChlMaXN0Q29tcGxpYW5jZUl0ZW1zUmVzdWx0EkAKD2NvbXBsaWFuY2VpdGVtcxiXmvU9IAMoCz'
    'ITLnNzbS5Db21wbGlhbmNlSXRlbVIPY29tcGxpYW5jZWl0ZW1zEh8KCW5leHR0b2tlbhj+hLpn'
    'IAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listComplianceSummariesRequestDescriptor instead')
const ListComplianceSummariesRequest$json = {
  '1': 'ListComplianceSummariesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceStringFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListComplianceSummariesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listComplianceSummariesRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0Q29tcGxpYW5jZVN1bW1hcmllc1JlcXVlc3QSOAoHZmlsdGVycxjtzepZIAMoCzIbLn'
        'NzbS5Db21wbGlhbmNlU3RyaW5nRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEg'
        'ASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listComplianceSummariesResultDescriptor instead')
const ListComplianceSummariesResult$json = {
  '1': 'ListComplianceSummariesResult',
  '2': [
    {
      '1': 'compliancesummaryitems',
      '3': 47794643,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceSummaryItem',
      '10': 'compliancesummaryitems'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListComplianceSummariesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listComplianceSummariesResultDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q29tcGxpYW5jZVN1bW1hcmllc1Jlc3VsdBJVChZjb21wbGlhbmNlc3VtbWFyeWl0ZW'
        '1zGNOT5RYgAygLMhouc3NtLkNvbXBsaWFuY2VTdW1tYXJ5SXRlbVIWY29tcGxpYW5jZXN1bW1h'
        'cnlpdGVtcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listDocumentMetadataHistoryRequestDescriptor instead')
const ListDocumentMetadataHistoryRequest$json = {
  '1': 'ListDocumentMetadataHistoryRequest',
  '2': [
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentMetadataEnum',
      '10': 'metadata'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentMetadataHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentMetadataHistoryRequestDescriptor =
    $convert.base64Decode(
        'CiJMaXN0RG9jdW1lbnRNZXRhZGF0YUhpc3RvcnlSZXF1ZXN0EisKD2RvY3VtZW50dmVyc2lvbh'
        'jJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhy'
        'ZXN1bHRzEjkKCG1ldGFkYXRhGOHij+ABIAEoDjIZLnNzbS5Eb2N1bWVudE1ldGFkYXRhRW51bV'
        'IIbWV0YWRhdGESFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlS'
        'CW5leHR0b2tlbg==');

@$core.Deprecated('Use listDocumentMetadataHistoryResponseDescriptor instead')
const ListDocumentMetadataHistoryResponse$json = {
  '1': 'ListDocumentMetadataHistoryResponse',
  '2': [
    {'1': 'author', '3': 361744247, '4': 1, '5': 9, '10': 'author'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'metadata',
      '3': 470020449,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentMetadataResponseInfo',
      '10': 'metadata'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentMetadataHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentMetadataHistoryResponseDescriptor =
    $convert.base64Decode(
        'CiNMaXN0RG9jdW1lbnRNZXRhZGF0YUhpc3RvcnlSZXNwb25zZRIaCgZhdXRob3IY946/rAEgAS'
        'gJUgZhdXRob3ISKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24S'
        'QQoIbWV0YWRhdGEY4eKP4AEgASgLMiEuc3NtLkRvY3VtZW50TWV0YWRhdGFSZXNwb25zZUluZm'
        '9SCG1ldGFkYXRhEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSHwoJbmV4dHRva2VuGP6EumcgASgJ'
        'UgluZXh0dG9rZW4=');

@$core.Deprecated('Use listDocumentVersionsRequestDescriptor instead')
const ListDocumentVersionsRequest$json = {
  '1': 'ListDocumentVersionsRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentVersionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentVersionsRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0RG9jdW1lbnRWZXJzaW9uc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlS'
        'CW5leHR0b2tlbg==');

@$core.Deprecated('Use listDocumentVersionsResultDescriptor instead')
const ListDocumentVersionsResult$json = {
  '1': 'ListDocumentVersionsResult',
  '2': [
    {
      '1': 'documentversions',
      '3': 175799138,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentVersionInfo',
      '10': 'documentversions'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentVersionsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentVersionsResultDescriptor =
    $convert.base64Decode(
        'ChpMaXN0RG9jdW1lbnRWZXJzaW9uc1Jlc3VsdBJHChBkb2N1bWVudHZlcnNpb25zGOL26VMgAy'
        'gLMhguc3NtLkRvY3VtZW50VmVyc2lvbkluZm9SEGRvY3VtZW50dmVyc2lvbnMSHwoJbmV4dHRv'
        'a2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listDocumentsRequestDescriptor instead')
const ListDocumentsRequest$json = {
  '1': 'ListDocumentsRequest',
  '2': [
    {
      '1': 'documentfilterlist',
      '3': 102169731,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentFilter',
      '10': 'documentfilterlist'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentKeyValuesFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0RG9jdW1lbnRzUmVxdWVzdBJGChJkb2N1bWVudGZpbHRlcmxpc3QYg/nbMCADKAsyEy'
    '5zc20uRG9jdW1lbnRGaWx0ZXJSEmRvY3VtZW50ZmlsdGVybGlzdBI5CgdmaWx0ZXJzGO3N6lkg'
    'AygLMhwuc3NtLkRvY3VtZW50S2V5VmFsdWVzRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdH'
    'MYsqibgwEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listDocumentsResultDescriptor instead')
const ListDocumentsResult$json = {
  '1': 'ListDocumentsResult',
  '2': [
    {
      '1': 'documentidentifiers',
      '3': 367445629,
      '4': 3,
      '5': 11,
      '6': '.ssm.DocumentIdentifier',
      '10': 'documentidentifiers'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDocumentsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDocumentsResultDescriptor = $convert.base64Decode(
    'ChNMaXN0RG9jdW1lbnRzUmVzdWx0Ek0KE2RvY3VtZW50aWRlbnRpZmllcnMY/YybrwEgAygLMh'
    'cuc3NtLkRvY3VtZW50SWRlbnRpZmllclITZG9jdW1lbnRpZGVudGlmaWVycxIfCgluZXh0dG9r'
    'ZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listInventoryEntriesRequestDescriptor instead')
const ListInventoryEntriesRequest$json = {
  '1': 'ListInventoryEntriesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryFilter',
      '10': 'filters'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `ListInventoryEntriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInventoryEntriesRequestDescriptor = $convert.base64Decode(
    'ChtMaXN0SW52ZW50b3J5RW50cmllc1JlcXVlc3QSMQoHZmlsdGVycxjtzepZIAMoCzIULnNzbS'
    '5JbnZlbnRvcnlGaWx0ZXJSB2ZpbHRlcnMSIQoKaW5zdGFuY2VpZBigrdEXIAEoCVIKaW5zdGFu'
    'Y2VpZBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/o'
    'S6ZyABKAlSCW5leHR0b2tlbhIeCgh0eXBlbmFtZRjPztnUASABKAlSCHR5cGVuYW1l');

@$core.Deprecated('Use listInventoryEntriesResultDescriptor instead')
const ListInventoryEntriesResult$json = {
  '1': 'ListInventoryEntriesResult',
  '2': [
    {'1': 'capturetime', '3': 72276089, '4': 1, '5': 9, '10': 'capturetime'},
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.ssm.ListInventoryEntriesResult.EntriesEntry',
      '10': 'entries'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'schemaversion',
      '3': 371681851,
      '4': 1,
      '5': 9,
      '10': 'schemaversion'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
  '3': [ListInventoryEntriesResult_EntriesEntry$json],
};

@$core.Deprecated('Use listInventoryEntriesResultDescriptor instead')
const ListInventoryEntriesResult_EntriesEntry$json = {
  '1': 'EntriesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListInventoryEntriesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listInventoryEntriesResultDescriptor = $convert.base64Decode(
    'ChpMaXN0SW52ZW50b3J5RW50cmllc1Jlc3VsdBIjCgtjYXB0dXJldGltZRj5sLsiIAEoCVILY2'
    'FwdHVyZXRpbWUSSgoHZW50cmllcxiUxbLlASADKAsyLC5zc20uTGlzdEludmVudG9yeUVudHJp'
    'ZXNSZXN1bHQuRW50cmllc0VudHJ5UgdlbnRyaWVzEiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCm'
    'luc3RhbmNlaWQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SKAoNc2NoZW1hdmVy'
    'c2lvbhi71J2xASABKAlSDXNjaGVtYXZlcnNpb24SHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eX'
    'BlbmFtZRo6CgxFbnRyaWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4AQ==');

@$core.Deprecated('Use listNodesRequestDescriptor instead')
const ListNodesRequest$json = {
  '1': 'ListNodesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.NodeFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
  ],
};

/// Descriptor for `ListNodesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNodesRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0Tm9kZXNSZXF1ZXN0EiwKB2ZpbHRlcnMY7c3qWSADKAsyDy5zc20uTm9kZUZpbHRlcl'
    'IHZmlsdGVycxIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9r'
    'ZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghzeW5jbmFtZRiilrKwASABKAlSCHN5bmNuYW1l');

@$core.Deprecated('Use listNodesResultDescriptor instead')
const ListNodesResult$json = {
  '1': 'ListNodesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'nodes',
      '3': 498945275,
      '4': 3,
      '5': 11,
      '6': '.ssm.Node',
      '10': 'nodes'
    },
  ],
};

/// Descriptor for `ListNodesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNodesResultDescriptor = $convert.base64Decode(
    'Cg9MaXN0Tm9kZXNSZXN1bHQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIwoFbm'
    '9kZXMY+5n17QEgAygLMgkuc3NtLk5vZGVSBW5vZGVz');

@$core.Deprecated('Use listNodesSummaryRequestDescriptor instead')
const ListNodesSummaryRequest$json = {
  '1': 'ListNodesSummaryRequest',
  '2': [
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.NodeAggregator',
      '10': 'aggregators'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.NodeFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
  ],
};

/// Descriptor for `ListNodesSummaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNodesSummaryRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0Tm9kZXNTdW1tYXJ5UmVxdWVzdBI4CgthZ2dyZWdhdG9ycxiO/YNNIAMoCzITLnNzbS'
    '5Ob2RlQWdncmVnYXRvclILYWdncmVnYXRvcnMSLAoHZmlsdGVycxjtzepZIAMoCzIPLnNzbS5O'
    'b2RlRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEh'
    '8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh4KCHN5bmNuYW1lGKKWsrABIAEoCVII'
    'c3luY25hbWU=');

@$core.Deprecated('Use listNodesSummaryResultDescriptor instead')
const ListNodesSummaryResult$json = {
  '1': 'ListNodesSummaryResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'summary',
      '3': 21935540,
      '4': 3,
      '5': 11,
      '6': '.ssm.ListNodesSummaryResult.SummaryEntry',
      '10': 'summary'
    },
  ],
  '3': [ListNodesSummaryResult_SummaryEntry$json],
};

@$core.Deprecated('Use listNodesSummaryResultDescriptor instead')
const ListNodesSummaryResult_SummaryEntry$json = {
  '1': 'SummaryEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListNodesSummaryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNodesSummaryResultDescriptor = $convert.base64Decode(
    'ChZMaXN0Tm9kZXNTdW1tYXJ5UmVzdWx0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2'
    'VuEkUKB3N1bW1hcnkYtOu6CiADKAsyKC5zc20uTGlzdE5vZGVzU3VtbWFyeVJlc3VsdC5TdW1t'
    'YXJ5RW50cnlSB3N1bW1hcnkaOgoMU3VtbWFyeUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBX'
    'ZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listOpsItemEventsRequestDescriptor instead')
const ListOpsItemEventsRequest$json = {
  '1': 'ListOpsItemEventsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemEventFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListOpsItemEventsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsItemEventsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0T3BzSXRlbUV2ZW50c1JlcXVlc3QSNAoHZmlsdGVycxjtzepZIAMoCzIXLnNzbS5PcH'
    'NJdGVtRXZlbnRGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJl'
    'c3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listOpsItemEventsResponseDescriptor instead')
const ListOpsItemEventsResponse$json = {
  '1': 'ListOpsItemEventsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'summaries',
      '3': 399105812,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemEventSummary',
      '10': 'summaries'
    },
  ],
};

/// Descriptor for `ListOpsItemEventsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsItemEventsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0T3BzSXRlbUV2ZW50c1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
    'Rva2VuEjoKCXN1bW1hcmllcxiUvqe+ASADKAsyGC5zc20uT3BzSXRlbUV2ZW50U3VtbWFyeVIJ'
    'c3VtbWFyaWVz');

@$core.Deprecated('Use listOpsItemRelatedItemsRequestDescriptor instead')
const ListOpsItemRelatedItemsRequest$json = {
  '1': 'ListOpsItemRelatedItemsRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemRelatedItemsFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `ListOpsItemRelatedItemsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsItemRelatedItemsRequestDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0T3BzSXRlbVJlbGF0ZWRJdGVtc1JlcXVlc3QSOwoHZmlsdGVycxjtzepZIAMoCzIeLn'
        'NzbS5PcHNJdGVtUmVsYXRlZEl0ZW1zRmlsdGVyUgdmaWx0ZXJzEiIKCm1heHJlc3VsdHMYsqib'
        'gwEgASgFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh8KCW'
        '9wc2l0ZW1pZBjS0pUMIAEoCVIJb3BzaXRlbWlk');

@$core.Deprecated('Use listOpsItemRelatedItemsResponseDescriptor instead')
const ListOpsItemRelatedItemsResponse$json = {
  '1': 'ListOpsItemRelatedItemsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'summaries',
      '3': 399105812,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemRelatedItemSummary',
      '10': 'summaries'
    },
  ],
};

/// Descriptor for `ListOpsItemRelatedItemsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsItemRelatedItemsResponseDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0T3BzSXRlbVJlbGF0ZWRJdGVtc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCV'
        'IJbmV4dHRva2VuEkAKCXN1bW1hcmllcxiUvqe+ASADKAsyHi5zc20uT3BzSXRlbVJlbGF0ZWRJ'
        'dGVtU3VtbWFyeVIJc3VtbWFyaWVz');

@$core.Deprecated('Use listOpsMetadataRequestDescriptor instead')
const ListOpsMetadataRequest$json = {
  '1': 'ListOpsMetadataRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsMetadataFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListOpsMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsMetadataRequestDescriptor = $convert.base64Decode(
    'ChZMaXN0T3BzTWV0YWRhdGFSZXF1ZXN0EjMKB2ZpbHRlcnMY7c3qWSADKAsyFi5zc20uT3BzTW'
    'V0YWRhdGFGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3Vs'
    'dHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listOpsMetadataResultDescriptor instead')
const ListOpsMetadataResult$json = {
  '1': 'ListOpsMetadataResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'opsmetadatalist',
      '3': 196633523,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsMetadata',
      '10': 'opsmetadatalist'
    },
  ],
};

/// Descriptor for `ListOpsMetadataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOpsMetadataResultDescriptor = $convert.base64Decode(
    'ChVMaXN0T3BzTWV0YWRhdGFSZXN1bHQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
    '4SPQoPb3BzbWV0YWRhdGFsaXN0GLPH4V0gAygLMhAuc3NtLk9wc01ldGFkYXRhUg9vcHNtZXRh'
    'ZGF0YWxpc3Q=');

@$core
    .Deprecated('Use listResourceComplianceSummariesRequestDescriptor instead')
const ListResourceComplianceSummariesRequest$json = {
  '1': 'ListResourceComplianceSummariesRequest',
  '2': [
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceStringFilter',
      '10': 'filters'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListResourceComplianceSummariesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceComplianceSummariesRequestDescriptor =
    $convert.base64Decode(
        'CiZMaXN0UmVzb3VyY2VDb21wbGlhbmNlU3VtbWFyaWVzUmVxdWVzdBI4CgdmaWx0ZXJzGO3N6l'
        'kgAygLMhsuc3NtLkNvbXBsaWFuY2VTdHJpbmdGaWx0ZXJSB2ZpbHRlcnMSIgoKbWF4cmVzdWx0'
        'cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW'
        '4=');

@$core.Deprecated('Use listResourceComplianceSummariesResultDescriptor instead')
const ListResourceComplianceSummariesResult$json = {
  '1': 'ListResourceComplianceSummariesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcecompliancesummaryitems',
      '3': 286761155,
      '4': 3,
      '5': 11,
      '6': '.ssm.ResourceComplianceSummaryItem',
      '10': 'resourcecompliancesummaryitems'
    },
  ],
};

/// Descriptor for `ListResourceComplianceSummariesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceComplianceSummariesResultDescriptor =
    $convert.base64Decode(
        'CiVMaXN0UmVzb3VyY2VDb21wbGlhbmNlU3VtbWFyaWVzUmVzdWx0Eh8KCW5leHR0b2tlbhj+hL'
        'pnIAEoCVIJbmV4dHRva2VuEm4KHnJlc291cmNlY29tcGxpYW5jZXN1bW1hcnlpdGVtcxjDwd6I'
        'ASADKAsyIi5zc20uUmVzb3VyY2VDb21wbGlhbmNlU3VtbWFyeUl0ZW1SHnJlc291cmNlY29tcG'
        'xpYW5jZXN1bW1hcnlpdGVtcw==');

@$core.Deprecated('Use listResourceDataSyncRequestDescriptor instead')
const ListResourceDataSyncRequest$json = {
  '1': 'ListResourceDataSyncRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `ListResourceDataSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceDataSyncRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0UmVzb3VyY2VEYXRhU3luY1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SHQoIc3luY3R5cGUY'
        '2+aqOiABKAlSCHN5bmN0eXBl');

@$core.Deprecated('Use listResourceDataSyncResultDescriptor instead')
const ListResourceDataSyncResult$json = {
  '1': 'ListResourceDataSyncResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcedatasyncitems',
      '3': 406717751,
      '4': 3,
      '5': 11,
      '6': '.ssm.ResourceDataSyncItem',
      '10': 'resourcedatasyncitems'
    },
  ],
};

/// Descriptor for `ListResourceDataSyncResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceDataSyncResultDescriptor =
    $convert.base64Decode(
        'ChpMaXN0UmVzb3VyY2VEYXRhU3luY1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhJTChVyZXNvdXJjZWRhdGFzeW5jaXRlbXMYt4r4wQEgAygLMhkuc3NtLlJlc291cmNl'
        'RGF0YVN5bmNJdGVtUhVyZXNvdXJjZWRhdGFzeW5jaXRlbXM=');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceTypeForTagging',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCVIKcm'
        'Vzb3VyY2VpZBJDCgxyZXNvdXJjZXR5cGUY3r7YjwEgASgOMhsuc3NtLlJlc291cmNlVHlwZUZv'
        'clRhZ2dpbmdSDHJlc291cmNldHlwZQ==');

@$core.Deprecated('Use listTagsForResourceResultDescriptor instead')
const ListTagsForResourceResult$json = {
  '1': 'ListTagsForResourceResult',
  '2': [
    {
      '1': 'taglist',
      '3': 429416860,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'taglist'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResultDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlUmVzdWx0EiYKB3RhZ2xpc3QYnMPhzAEgAygLMgguc3NtLl'
        'RhZ1IHdGFnbGlzdA==');

@$core.Deprecated('Use loggingInfoDescriptor instead')
const LoggingInfo$json = {
  '1': 'LoggingInfo',
  '2': [
    {'1': 's3bucketname', '3': 320495427, '4': 1, '5': 9, '10': 's3bucketname'},
    {'1': 's3keyprefix', '3': 206015359, '4': 1, '5': 9, '10': 's3keyprefix'},
    {'1': 's3region', '3': 354623468, '4': 1, '5': 9, '10': 's3region'},
  ],
};

/// Descriptor for `LoggingInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loggingInfoDescriptor = $convert.base64Decode(
    'CgtMb2dnaW5nSW5mbxImCgxzM2J1Y2tldG5hbWUYw77pmAEgASgJUgxzM2J1Y2tldG5hbWUSIw'
    'oLczNrZXlwcmVmaXgY/5aeYiABKAlSC3Mza2V5cHJlZml4Eh4KCHMzcmVnaW9uGOy/jKkBIAEo'
    'CVIIczNyZWdpb24=');

@$core.Deprecated('Use maintenanceWindowAutomationParametersDescriptor instead')
const MaintenanceWindowAutomationParameters$json = {
  '1': 'MaintenanceWindowAutomationParameters',
  '2': [
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowAutomationParameters.ParametersEntry',
      '10': 'parameters'
    },
  ],
  '3': [MaintenanceWindowAutomationParameters_ParametersEntry$json],
};

@$core.Deprecated('Use maintenanceWindowAutomationParametersDescriptor instead')
const MaintenanceWindowAutomationParameters_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MaintenanceWindowAutomationParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowAutomationParametersDescriptor =
    $convert.base64Decode(
        'CiVNYWludGVuYW5jZVdpbmRvd0F1dG9tYXRpb25QYXJhbWV0ZXJzEisKD2RvY3VtZW50dmVyc2'
        'lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEl4KCnBhcmFtZXRlcnMY+qf+6wEgAygLMjou'
        'c3NtLk1haW50ZW5hbmNlV2luZG93QXV0b21hdGlvblBhcmFtZXRlcnMuUGFyYW1ldGVyc0VudH'
        'J5UgpwYXJhbWV0ZXJzGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2'
        'YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use maintenanceWindowExecutionDescriptor instead')
const MaintenanceWindowExecution$json = {
  '1': 'MaintenanceWindowExecution',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `MaintenanceWindowExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowExecutionDescriptor = $convert.base64Decode(
    'ChpNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvbhIbCgdlbmR0aW1lGMzvvB4gASgJUgdlbmR0aW'
    '1lEiAKCXN0YXJ0dGltZRjvtOWwASABKAlSCXN0YXJ0dGltZRJACgZzdGF0dXMYkOT7AiABKA4y'
    'JS5zc20uTWFpbnRlbmFuY2VXaW5kb3dFeGVjdXRpb25TdGF0dXNSBnN0YXR1cxIoCg1zdGF0dX'
    'NkZXRhaWxzGKiSwbEBIAEoCVINc3RhdHVzZGV0YWlscxIwChF3aW5kb3dleGVjdXRpb25pZBiJ'
    '66eqASABKAlSEXdpbmRvd2V4ZWN1dGlvbmlkEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3'
    'dpZA==');

@$core
    .Deprecated('Use maintenanceWindowExecutionTaskIdentityDescriptor instead')
const MaintenanceWindowExecutionTaskIdentity$json = {
  '1': 'MaintenanceWindowExecutionTaskIdentity',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskexecutionid',
      '3': 474180408,
      '4': 1,
      '5': 9,
      '10': 'taskexecutionid'
    },
    {
      '1': 'tasktype',
      '3': 370407909,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'tasktype'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
  ],
};

/// Descriptor for `MaintenanceWindowExecutionTaskIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowExecutionTaskIdentityDescriptor = $convert.base64Decode(
    'CiZNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tJZGVudGl0eRJKChJhbGFybWNvbmZpZ3'
    'VyYXRpb24YiZm5ISABKAsyFy5zc20uQWxhcm1Db25maWd1cmF0aW9uUhJhbGFybWNvbmZpZ3Vy'
    'YXRpb24SGwoHZW5kdGltZRjM77weIAEoCVIHZW5kdGltZRIgCglzdGFydHRpbWUY77TlsAEgAS'
    'gJUglzdGFydHRpbWUSQAoGc3RhdHVzGJDk+wIgASgOMiUuc3NtLk1haW50ZW5hbmNlV2luZG93'
    'RXhlY3V0aW9uU3RhdHVzUgZzdGF0dXMSKAoNc3RhdHVzZGV0YWlscxioksGxASABKAlSDXN0YX'
    'R1c2RldGFpbHMSHAoHdGFza2FybhjkyfqUASABKAlSB3Rhc2thcm4SLAoPdGFza2V4ZWN1dGlv'
    'bmlkGLjWjeIBIAEoCVIPdGFza2V4ZWN1dGlvbmlkEj4KCHRhc2t0eXBlGOXzz7ABIAEoDjIeLn'
    'NzbS5NYWludGVuYW5jZVdpbmRvd1Rhc2tUeXBlUgh0YXNrdHlwZRJHCg90cmlnZ2VyZWRhbGFy'
    'bXMYhe3BfSADKAsyGi5zc20uQWxhcm1TdGF0ZUluZm9ybWF0aW9uUg90cmlnZ2VyZWRhbGFybX'
    'MSMAoRd2luZG93ZXhlY3V0aW9uaWQYieunqgEgASgJUhF3aW5kb3dleGVjdXRpb25pZA==');

@$core.Deprecated(
    'Use maintenanceWindowExecutionTaskInvocationIdentityDescriptor instead')
const MaintenanceWindowExecutionTaskInvocationIdentity$json = {
  '1': 'MaintenanceWindowExecutionTaskInvocationIdentity',
  '2': [
    {'1': 'endtime', '3': 63911884, '4': 1, '5': 9, '10': 'endtime'},
    {'1': 'executionid', '3': 147580849, '4': 1, '5': 9, '10': 'executionid'},
    {'1': 'invocationid', '3': 116064639, '4': 1, '5': 9, '10': 'invocationid'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {'1': 'parameters', '3': 494900218, '4': 1, '5': 9, '10': 'parameters'},
    {'1': 'starttime', '3': 370760303, '4': 1, '5': 9, '10': 'starttime'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowExecutionStatus',
      '10': 'status'
    },
    {
      '1': 'statusdetails',
      '3': 372263208,
      '4': 1,
      '5': 9,
      '10': 'statusdetails'
    },
    {
      '1': 'taskexecutionid',
      '3': 474180408,
      '4': 1,
      '5': 9,
      '10': 'taskexecutionid'
    },
    {
      '1': 'tasktype',
      '3': 370407909,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'tasktype'
    },
    {
      '1': 'windowexecutionid',
      '3': 357168521,
      '4': 1,
      '5': 9,
      '10': 'windowexecutionid'
    },
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `MaintenanceWindowExecutionTaskInvocationIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowExecutionTaskInvocationIdentityDescriptor = $convert.base64Decode(
    'CjBNYWludGVuYW5jZVdpbmRvd0V4ZWN1dGlvblRhc2tJbnZvY2F0aW9uSWRlbnRpdHkSGwoHZW'
    '5kdGltZRjM77weIAEoCVIHZW5kdGltZRIjCgtleGVjdXRpb25pZBixz69GIAEoCVILZXhlY3V0'
    'aW9uaWQSJQoMaW52b2NhdGlvbmlkGP+CrDcgASgJUgxpbnZvY2F0aW9uaWQSLQoQb3duZXJpbm'
    'Zvcm1hdGlvbhiDmsUgIAEoCVIQb3duZXJpbmZvcm1hdGlvbhIiCgpwYXJhbWV0ZXJzGPqn/usB'
    'IAEoCVIKcGFyYW1ldGVycxIgCglzdGFydHRpbWUY77TlsAEgASgJUglzdGFydHRpbWUSQAoGc3'
    'RhdHVzGJDk+wIgASgOMiUuc3NtLk1haW50ZW5hbmNlV2luZG93RXhlY3V0aW9uU3RhdHVzUgZz'
    'dGF0dXMSKAoNc3RhdHVzZGV0YWlscxioksGxASABKAlSDXN0YXR1c2RldGFpbHMSLAoPdGFza2'
    'V4ZWN1dGlvbmlkGLjWjeIBIAEoCVIPdGFza2V4ZWN1dGlvbmlkEj4KCHRhc2t0eXBlGOXzz7AB'
    'IAEoDjIeLnNzbS5NYWludGVuYW5jZVdpbmRvd1Rhc2tUeXBlUgh0YXNrdHlwZRIwChF3aW5kb3'
    'dleGVjdXRpb25pZBiJ66eqASABKAlSEXdpbmRvd2V4ZWN1dGlvbmlkEikKDndpbmRvd3Rhcmdl'
    'dGlkGN7O6nsgASgJUg53aW5kb3d0YXJnZXRpZA==');

@$core.Deprecated('Use maintenanceWindowFilterDescriptor instead')
const MaintenanceWindowFilter$json = {
  '1': 'MaintenanceWindowFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `MaintenanceWindowFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowFilterDescriptor =
    $convert.base64Decode(
        'ChdNYWludGVuYW5jZVdpbmRvd0ZpbHRlchITCgNrZXkYjZLraCABKAlSA2tleRIZCgZ2YWx1ZX'
        'MY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use maintenanceWindowIdentityDescriptor instead')
const MaintenanceWindowIdentity$json = {
  '1': 'MaintenanceWindowIdentity',
  '2': [
    {'1': 'cutoff', '3': 498433089, '4': 1, '5': 5, '10': 'cutoff'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'nextexecutiontime',
      '3': 528679720,
      '4': 1,
      '5': 9,
      '10': 'nextexecutiontime'
    },
    {'1': 'schedule', '3': 66697965, '4': 1, '5': 9, '10': 'schedule'},
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'scheduletimezone',
      '3': 170037696,
      '4': 1,
      '5': 9,
      '10': 'scheduletimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `MaintenanceWindowIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowIdentityDescriptor = $convert.base64Decode(
    'ChlNYWludGVuYW5jZVdpbmRvd0lkZW50aXR5EhoKBmN1dG9mZhjB+NXtASABKAVSBmN1dG9mZh'
    'IjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SHgoIZHVyYXRpb24YrpKdpgEg'
    'ASgFUghkdXJhdGlvbhIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBIbCgdlbmRkYXRlGM'
    '+z+SQgASgJUgdlbmRkYXRlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSMAoRbmV4dGV4ZWN1dGlv'
    'bnRpbWUYqIaM/AEgASgJUhFuZXh0ZXhlY3V0aW9udGltZRIdCghzY2hlZHVsZRjt9eYfIAEoCV'
    'IIc2NoZWR1bGUSKQoOc2NoZWR1bGVvZmZzZXQY2JHqSiABKAVSDnNjaGVkdWxlb2Zmc2V0Ei0K'
    'EHNjaGVkdWxldGltZXpvbmUYwKOKUSABKAlSEHNjaGVkdWxldGltZXpvbmUSIAoJc3RhcnRkYX'
    'RlGPz4oNQBIAEoCVIJc3RhcnRkYXRlEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZA==');

@$core.Deprecated('Use maintenanceWindowIdentityForTargetDescriptor instead')
const MaintenanceWindowIdentityForTarget$json = {
  '1': 'MaintenanceWindowIdentityForTarget',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `MaintenanceWindowIdentityForTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowIdentityForTargetDescriptor =
    $convert.base64Decode(
        'CiJNYWludGVuYW5jZVdpbmRvd0lkZW50aXR5Rm9yVGFyZ2V0EhUKBG5hbWUYh+aBfyABKAlSBG'
        '5hbWUSHQoId2luZG93aWQYqeSHCSABKAlSCHdpbmRvd2lk');

@$core.Deprecated('Use maintenanceWindowLambdaParametersDescriptor instead')
const MaintenanceWindowLambdaParameters$json = {
  '1': 'MaintenanceWindowLambdaParameters',
  '2': [
    {
      '1': 'clientcontext',
      '3': 354405962,
      '4': 1,
      '5': 9,
      '10': 'clientcontext'
    },
    {'1': 'payload', '3': 6526790, '4': 1, '5': 12, '10': 'payload'},
    {'1': 'qualifier', '3': 526670560, '4': 1, '5': 9, '10': 'qualifier'},
  ],
};

/// Descriptor for `MaintenanceWindowLambdaParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowLambdaParametersDescriptor =
    $convert.base64Decode(
        'CiFNYWludGVuYW5jZVdpbmRvd0xhbWJkYVBhcmFtZXRlcnMSKAoNY2xpZW50Y29udGV4dBjKnP'
        '+oASABKAlSDWNsaWVudGNvbnRleHQSGwoHcGF5bG9hZBjGro4DIAEoDFIHcGF5bG9hZBIgCglx'
        'dWFsaWZpZXIY4LWR+wEgASgJUglxdWFsaWZpZXI=');

@$core.Deprecated('Use maintenanceWindowRunCommandParametersDescriptor instead')
const MaintenanceWindowRunCommandParameters$json = {
  '1': 'MaintenanceWindowRunCommandParameters',
  '2': [
    {
      '1': 'cloudwatchoutputconfig',
      '3': 21186555,
      '4': 1,
      '5': 11,
      '6': '.ssm.CloudWatchOutputConfig',
      '10': 'cloudwatchoutputconfig'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'documenthash', '3': 251248469, '4': 1, '5': 9, '10': 'documenthash'},
    {
      '1': 'documenthashtype',
      '3': 93041117,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentHashType',
      '10': 'documenthashtype'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'notificationconfig',
      '3': 346074145,
      '4': 1,
      '5': 11,
      '6': '.ssm.NotificationConfig',
      '10': 'notificationconfig'
    },
    {
      '1': 'outputs3bucketname',
      '3': 186756480,
      '4': 1,
      '5': 9,
      '10': 'outputs3bucketname'
    },
    {
      '1': 'outputs3keyprefix',
      '3': 17840974,
      '4': 1,
      '5': 9,
      '10': 'outputs3keyprefix'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowRunCommandParameters.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'timeoutseconds',
      '3': 336148022,
      '4': 1,
      '5': 5,
      '10': 'timeoutseconds'
    },
  ],
  '3': [MaintenanceWindowRunCommandParameters_ParametersEntry$json],
};

@$core.Deprecated('Use maintenanceWindowRunCommandParametersDescriptor instead')
const MaintenanceWindowRunCommandParameters_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MaintenanceWindowRunCommandParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowRunCommandParametersDescriptor = $convert.base64Decode(
    'CiVNYWludGVuYW5jZVdpbmRvd1J1bkNvbW1hbmRQYXJhbWV0ZXJzElYKFmNsb3Vkd2F0Y2hvdX'
    'RwdXRjb25maWcY+4+NCiABKAsyGy5zc20uQ2xvdWRXYXRjaE91dHB1dENvbmZpZ1IWY2xvdWR3'
    'YXRjaG91dHB1dGNvbmZpZxIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBIlCgxkb2N1bW'
    'VudGhhc2gY1f7mdyABKAlSDGRvY3VtZW50aGFzaBJEChBkb2N1bWVudGhhc2h0eXBlGN3jriwg'
    'ASgOMhUuc3NtLkRvY3VtZW50SGFzaFR5cGVSEGRvY3VtZW50aGFzaHR5cGUSKwoPZG9jdW1lbn'
    'R2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVudHZlcnNpb24SSwoSbm90aWZpY2F0aW9uY29uZmln'
    'GKHYgqUBIAEoCzIXLnNzbS5Ob3RpZmljYXRpb25Db25maWdSEm5vdGlmaWNhdGlvbmNvbmZpZx'
    'IxChJvdXRwdXRzM2J1Y2tldG5hbWUYgNuGWSABKAlSEm91dHB1dHMzYnVja2V0bmFtZRIvChFv'
    'dXRwdXRzM2tleXByZWZpeBjO9sAIIAEoCVIRb3V0cHV0czNrZXlwcmVmaXgSXgoKcGFyYW1ldG'
    'Vycxj6p/7rASADKAsyOi5zc20uTWFpbnRlbmFuY2VXaW5kb3dSdW5Db21tYW5kUGFyYW1ldGVy'
    'cy5QYXJhbWV0ZXJzRW50cnlSCnBhcmFtZXRlcnMSKgoOc2VydmljZXJvbGVhcm4YhOPatgEgAS'
    'gJUg5zZXJ2aWNlcm9sZWFybhIqCg50aW1lb3V0c2Vjb25kcxi27KSgASABKAVSDnRpbWVvdXRz'
    'ZWNvbmRzGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIA'
    'EoCVIFdmFsdWU6AjgB');

@$core.Deprecated(
    'Use maintenanceWindowStepFunctionsParametersDescriptor instead')
const MaintenanceWindowStepFunctionsParameters$json = {
  '1': 'MaintenanceWindowStepFunctionsParameters',
  '2': [
    {'1': 'input', '3': 529785116, '4': 1, '5': 9, '10': 'input'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `MaintenanceWindowStepFunctionsParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowStepFunctionsParametersDescriptor =
    $convert.base64Decode(
        'CihNYWludGVuYW5jZVdpbmRvd1N0ZXBGdW5jdGlvbnNQYXJhbWV0ZXJzEhgKBWlucHV0GJzCz/'
        'wBIAEoCVIFaW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use maintenanceWindowTargetDescriptor instead')
const MaintenanceWindowTarget$json = {
  '1': 'MaintenanceWindowTarget',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowResourceType',
      '10': 'resourcetype'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `MaintenanceWindowTarget`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowTargetDescriptor = $convert.base64Decode(
    'ChdNYWludGVuYW5jZVdpbmRvd1RhcmdldBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3'
    'JpcHRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRItChBvd25lcmluZm9ybWF0aW9uGIOaxSAg'
    'ASgJUhBvd25lcmluZm9ybWF0aW9uEkoKDHJlc291cmNldHlwZRjevtiPASABKA4yIi5zc20uTW'
    'FpbnRlbmFuY2VXaW5kb3dSZXNvdXJjZVR5cGVSDHJlc291cmNldHlwZRIoCgd0YXJnZXRzGIKb'
    'gn0gAygLMgsuc3NtLlRhcmdldFIHdGFyZ2V0cxIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2luZG'
    '93aWQSKQoOd2luZG93dGFyZ2V0aWQY3s7qeyABKAlSDndpbmRvd3RhcmdldGlk');

@$core.Deprecated('Use maintenanceWindowTaskDescriptor instead')
const MaintenanceWindowTask$json = {
  '1': 'MaintenanceWindowTask',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cutoffbehavior',
      '3': 120608587,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskCutoffBehavior',
      '10': 'cutoffbehavior'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'logginginfo',
      '3': 448312415,
      '4': 1,
      '5': 11,
      '6': '.ssm.LoggingInfo',
      '10': 'logginginfo'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTask.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'type'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
  '3': [MaintenanceWindowTask_TaskparametersEntry$json],
};

@$core.Deprecated('Use maintenanceWindowTaskDescriptor instead')
const MaintenanceWindowTask_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `MaintenanceWindowTask`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowTaskDescriptor = $convert.base64Decode(
    'ChVNYWludGVuYW5jZVdpbmRvd1Rhc2sSSgoSYWxhcm1jb25maWd1cmF0aW9uGImZuSEgASgLMh'
    'cuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uElMKDmN1dG9mZmJl'
    'aGF2aW9yGMuuwTkgASgOMiguc3NtLk1haW50ZW5hbmNlV2luZG93VGFza0N1dG9mZkJlaGF2aW'
    '9yUg5jdXRvZmZiZWhhdmlvchIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24S'
    'NgoLbG9nZ2luZ2luZm8Y3+ji1QEgASgLMhAuc3NtLkxvZ2dpbmdJbmZvUgtsb2dnaW5naW5mbx'
    'IpCg5tYXhjb25jdXJyZW5jeRj9wY4OIAEoCVIObWF4Y29uY3VycmVuY3kSHwoJbWF4ZXJyb3Jz'
    'GKvC9T0gASgJUgltYXhlcnJvcnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIdCghwcmlvcml0eR'
    'iqvrY0IAEoBVIIcHJpb3JpdHkSKgoOc2VydmljZXJvbGVhcm4YhOPatgEgASgJUg5zZXJ2aWNl'
    'cm9sZWFybhIoCgd0YXJnZXRzGIKbgn0gAygLMgsuc3NtLlRhcmdldFIHdGFyZ2V0cxIcCgd0YX'
    'NrYXJuGOTJ+pQBIAEoCVIHdGFza2FybhJaCg50YXNrcGFyYW1ldGVycxiBj+a3ASADKAsyLi5z'
    'c20uTWFpbnRlbmFuY2VXaW5kb3dUYXNrLlRhc2twYXJhbWV0ZXJzRW50cnlSDnRhc2twYXJhbW'
    'V0ZXJzEjYKBHR5cGUY7qDXigEgASgOMh4uc3NtLk1haW50ZW5hbmNlV2luZG93VGFza1R5cGVS'
    'BHR5cGUSHQoId2luZG93aWQYqeSHCSABKAlSCHdpbmRvd2lkEiYKDHdpbmRvd3Rhc2tpZBjajb'
    'GaASABKAlSDHdpbmRvd3Rhc2tpZBp1ChNUYXNrcGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEo'
    'CVIDa2V5EkgKBXZhbHVlGAIgASgLMjIuc3NtLk1haW50ZW5hbmNlV2luZG93VGFza1BhcmFtZX'
    'RlclZhbHVlRXhwcmVzc2lvblIFdmFsdWU6AjgB');

@$core.Deprecated(
    'Use maintenanceWindowTaskInvocationParametersDescriptor instead')
const MaintenanceWindowTaskInvocationParameters$json = {
  '1': 'MaintenanceWindowTaskInvocationParameters',
  '2': [
    {
      '1': 'automation',
      '3': 73629771,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowAutomationParameters',
      '10': 'automation'
    },
    {
      '1': 'lambda',
      '3': 85519371,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowLambdaParameters',
      '10': 'lambda'
    },
    {
      '1': 'runcommand',
      '3': 289316476,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowRunCommandParameters',
      '10': 'runcommand'
    },
    {
      '1': 'stepfunctions',
      '3': 361131221,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowStepFunctionsParameters',
      '10': 'stepfunctions'
    },
  ],
};

/// Descriptor for `MaintenanceWindowTaskInvocationParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maintenanceWindowTaskInvocationParametersDescriptor = $convert.base64Decode(
    'CilNYWludGVuYW5jZVdpbmRvd1Rhc2tJbnZvY2F0aW9uUGFyYW1ldGVycxJNCgphdXRvbWF0aW'
    '9uGMuAjiMgASgLMiouc3NtLk1haW50ZW5hbmNlV2luZG93QXV0b21hdGlvblBhcmFtZXRlcnNS'
    'CmF1dG9tYXRpb24SQQoGbGFtYmRhGIvY4yggASgLMiYuc3NtLk1haW50ZW5hbmNlV2luZG93TG'
    'FtYmRhUGFyYW1ldGVyc1IGbGFtYmRhEk4KCnJ1bmNvbW1hbmQY/Lz6iQEgASgLMiouc3NtLk1h'
    'aW50ZW5hbmNlV2luZG93UnVuQ29tbWFuZFBhcmFtZXRlcnNSCnJ1bmNvbW1hbmQSVwoNc3RlcG'
    'Z1bmN0aW9ucxjV2ZmsASABKAsyLS5zc20uTWFpbnRlbmFuY2VXaW5kb3dTdGVwRnVuY3Rpb25z'
    'UGFyYW1ldGVyc1INc3RlcGZ1bmN0aW9ucw==');

@$core.Deprecated(
    'Use maintenanceWindowTaskParameterValueExpressionDescriptor instead')
const MaintenanceWindowTaskParameterValueExpression$json = {
  '1': 'MaintenanceWindowTaskParameterValueExpression',
  '2': [
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `MaintenanceWindowTaskParameterValueExpression`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    maintenanceWindowTaskParameterValueExpressionDescriptor =
    $convert.base64Decode(
        'Ci1NYWludGVuYW5jZVdpbmRvd1Rhc2tQYXJhbWV0ZXJWYWx1ZUV4cHJlc3Npb24SGQoGdmFsdW'
        'VzGNzEtGogAygJUgZ2YWx1ZXM=');

@$core.Deprecated(
    'Use malformedResourcePolicyDocumentExceptionDescriptor instead')
const MalformedResourcePolicyDocumentException$json = {
  '1': 'MalformedResourcePolicyDocumentException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MalformedResourcePolicyDocumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedResourcePolicyDocumentExceptionDescriptor =
    $convert.base64Decode(
        'CihNYWxmb3JtZWRSZXNvdXJjZVBvbGljeURvY3VtZW50RXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use maxDocumentSizeExceededDescriptor instead')
const MaxDocumentSizeExceeded$json = {
  '1': 'MaxDocumentSizeExceeded',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MaxDocumentSizeExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List maxDocumentSizeExceededDescriptor =
    $convert.base64Decode(
        'ChdNYXhEb2N1bWVudFNpemVFeGNlZWRlZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use metadataValueDescriptor instead')
const MetadataValue$json = {
  '1': 'MetadataValue',
  '2': [
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MetadataValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metadataValueDescriptor = $convert
    .base64Decode('Cg1NZXRhZGF0YVZhbHVlEhgKBXZhbHVlGOvyn4oBIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use modifyDocumentPermissionRequestDescriptor instead')
const ModifyDocumentPermissionRequest$json = {
  '1': 'ModifyDocumentPermissionRequest',
  '2': [
    {
      '1': 'accountidstoadd',
      '3': 96194831,
      '4': 3,
      '5': 9,
      '10': 'accountidstoadd'
    },
    {
      '1': 'accountidstoremove',
      '3': 482842600,
      '4': 3,
      '5': 9,
      '10': 'accountidstoremove'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'permissiontype',
      '3': 479110739,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentPermissionType',
      '10': 'permissiontype'
    },
    {
      '1': 'shareddocumentversion',
      '3': 240155144,
      '4': 1,
      '5': 9,
      '10': 'shareddocumentversion'
    },
  ],
};

/// Descriptor for `ModifyDocumentPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modifyDocumentPermissionRequestDescriptor = $convert.base64Decode(
    'Ch9Nb2RpZnlEb2N1bWVudFBlcm1pc3Npb25SZXF1ZXN0EisKD2FjY291bnRpZHN0b2FkZBiPou'
    '8tIAMoCVIPYWNjb3VudGlkc3RvYWRkEjIKEmFjY291bnRpZHN0b3JlbW92ZRjor57mASADKAlS'
    'EmFjY291bnRpZHN0b3JlbW92ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEkcKDnBlcm1pc3Npb2'
    '50eXBlGNPMuuQBIAEoDjIbLnNzbS5Eb2N1bWVudFBlcm1pc3Npb25UeXBlUg5wZXJtaXNzaW9u'
    'dHlwZRI3ChVzaGFyZWRkb2N1bWVudHZlcnNpb24YiPTBciABKAlSFXNoYXJlZGRvY3VtZW50dm'
    'Vyc2lvbg==');

@$core.Deprecated('Use modifyDocumentPermissionResponseDescriptor instead')
const ModifyDocumentPermissionResponse$json = {
  '1': 'ModifyDocumentPermissionResponse',
};

/// Descriptor for `ModifyDocumentPermissionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modifyDocumentPermissionResponseDescriptor =
    $convert.base64Decode('CiBNb2RpZnlEb2N1bWVudFBlcm1pc3Npb25SZXNwb25zZQ==');

@$core.Deprecated('Use noLongerSupportedExceptionDescriptor instead')
const NoLongerSupportedException$json = {
  '1': 'NoLongerSupportedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NoLongerSupportedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List noLongerSupportedExceptionDescriptor =
    $convert.base64Decode(
        'ChpOb0xvbmdlclN1cHBvcnRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use nodeDescriptor instead')
const Node$json = {
  '1': 'Node',
  '2': [
    {'1': 'capturetime', '3': 72276089, '4': 1, '5': 9, '10': 'capturetime'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'nodetype',
      '3': 117090310,
      '4': 1,
      '5': 11,
      '6': '.ssm.NodeType',
      '10': 'nodetype'
    },
    {
      '1': 'owner',
      '3': 455261813,
      '4': 1,
      '5': 11,
      '6': '.ssm.NodeOwnerInfo',
      '10': 'owner'
    },
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `Node`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nodeDescriptor = $convert.base64Decode(
    'CgROb2RlEiMKC2NhcHR1cmV0aW1lGPmwuyIgASgJUgtjYXB0dXJldGltZRISCgJpZBiB8qK3AS'
    'ABKAlSAmlkEiwKCG5vZGV0eXBlGIbQ6jcgASgLMg0uc3NtLk5vZGVUeXBlUghub2RldHlwZRIs'
    'CgVvd25lchj1/IrZASABKAsyEi5zc20uTm9kZU93bmVySW5mb1IFb3duZXISGQoGcmVnaW9uGJ'
    '7xuUkgASgJUgZyZWdpb24=');

@$core.Deprecated('Use nodeAggregatorDescriptor instead')
const NodeAggregator$json = {
  '1': 'NodeAggregator',
  '2': [
    {
      '1': 'aggregatortype',
      '3': 328765509,
      '4': 1,
      '5': 14,
      '6': '.ssm.NodeAggregatorType',
      '10': 'aggregatortype'
    },
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.NodeAggregator',
      '10': 'aggregators'
    },
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 14,
      '6': '.ssm.NodeAttributeName',
      '10': 'attributename'
    },
    {
      '1': 'typename',
      '3': 446064463,
      '4': 1,
      '5': 14,
      '6': '.ssm.NodeTypeName',
      '10': 'typename'
    },
  ],
};

/// Descriptor for `NodeAggregator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nodeAggregatorDescriptor = $convert.base64Decode(
    'Cg5Ob2RlQWdncmVnYXRvchJDCg5hZ2dyZWdhdG9ydHlwZRjFoOKcASABKA4yFy5zc20uTm9kZU'
    'FnZ3JlZ2F0b3JUeXBlUg5hZ2dyZWdhdG9ydHlwZRI4CgthZ2dyZWdhdG9ycxiO/YNNIAMoCzIT'
    'LnNzbS5Ob2RlQWdncmVnYXRvclILYWdncmVnYXRvcnMSQAoNYXR0cmlidXRlbmFtZRitlZioAS'
    'ABKA4yFi5zc20uTm9kZUF0dHJpYnV0ZU5hbWVSDWF0dHJpYnV0ZW5hbWUSMQoIdHlwZW5hbWUY'
    'z87Z1AEgASgOMhEuc3NtLk5vZGVUeXBlTmFtZVIIdHlwZW5hbWU=');

@$core.Deprecated('Use nodeFilterDescriptor instead')
const NodeFilter$json = {
  '1': 'NodeFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.NodeFilterKey',
      '10': 'key'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.NodeFilterOperatorType',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `NodeFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nodeFilterDescriptor = $convert.base64Decode(
    'CgpOb2RlRmlsdGVyEicKA2tleRiNkutoIAEoDjISLnNzbS5Ob2RlRmlsdGVyS2V5UgNrZXkSMw'
    'oEdHlwZRjuoNeKASABKA4yGy5zc20uTm9kZUZpbHRlck9wZXJhdG9yVHlwZVIEdHlwZRIZCgZ2'
    'YWx1ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use nodeOwnerInfoDescriptor instead')
const NodeOwnerInfo$json = {
  '1': 'NodeOwnerInfo',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'organizationalunitid',
      '3': 201174977,
      '4': 1,
      '5': 9,
      '10': 'organizationalunitid'
    },
    {
      '1': 'organizationalunitpath',
      '3': 338158359,
      '4': 1,
      '5': 9,
      '10': 'organizationalunitpath'
    },
  ],
};

/// Descriptor for `NodeOwnerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nodeOwnerInfoDescriptor = $convert.base64Decode(
    'Cg1Ob2RlT3duZXJJbmZvEh8KCWFjY291bnRpZBjSwbkfIAEoCVIJYWNjb3VudGlkEjUKFG9yZ2'
    'FuaXphdGlvbmFsdW5pdGlkGMHf9l8gASgJUhRvcmdhbml6YXRpb25hbHVuaXRpZBI6ChZvcmdh'
    'bml6YXRpb25hbHVuaXRwYXRoGJfGn6EBIAEoCVIWb3JnYW5pemF0aW9uYWx1bml0cGF0aA==');

@$core.Deprecated('Use nodeTypeDescriptor instead')
const NodeType$json = {
  '1': 'NodeType',
  '2': [
    {
      '1': 'instance',
      '3': 89665971,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceInfo',
      '10': 'instance'
    },
  ],
};

/// Descriptor for `NodeType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nodeTypeDescriptor = $convert.base64Decode(
    'CghOb2RlVHlwZRIwCghpbnN0YW5jZRiz4+AqIAEoCzIRLnNzbS5JbnN0YW5jZUluZm9SCGluc3'
    'RhbmNl');

@$core.Deprecated('Use nonCompliantSummaryDescriptor instead')
const NonCompliantSummary$json = {
  '1': 'NonCompliantSummary',
  '2': [
    {
      '1': 'noncompliantcount',
      '3': 528024125,
      '4': 1,
      '5': 5,
      '10': 'noncompliantcount'
    },
    {
      '1': 'severitysummary',
      '3': 20421279,
      '4': 1,
      '5': 11,
      '6': '.ssm.SeveritySummary',
      '10': 'severitysummary'
    },
  ],
};

/// Descriptor for `NonCompliantSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List nonCompliantSummaryDescriptor = $convert.base64Decode(
    'ChNOb25Db21wbGlhbnRTdW1tYXJ5EjAKEW5vbmNvbXBsaWFudGNvdW50GL2E5PsBIAEoBVIRbm'
    '9uY29tcGxpYW50Y291bnQSQQoPc2V2ZXJpdHlzdW1tYXJ5GJ+13gkgASgLMhQuc3NtLlNldmVy'
    'aXR5U3VtbWFyeVIPc2V2ZXJpdHlzdW1tYXJ5');

@$core.Deprecated('Use notificationConfigDescriptor instead')
const NotificationConfig$json = {
  '1': 'NotificationConfig',
  '2': [
    {
      '1': 'notificationarn',
      '3': 192280642,
      '4': 1,
      '5': 9,
      '10': 'notificationarn'
    },
    {
      '1': 'notificationevents',
      '3': 16584116,
      '4': 3,
      '5': 14,
      '6': '.ssm.NotificationEvent',
      '10': 'notificationevents'
    },
    {
      '1': 'notificationtype',
      '3': 43601335,
      '4': 1,
      '5': 14,
      '6': '.ssm.NotificationType',
      '10': 'notificationtype'
    },
  ],
};

/// Descriptor for `NotificationConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationConfigDescriptor = $convert.base64Decode(
    'ChJOb3RpZmljYXRpb25Db25maWcSKwoPbm90aWZpY2F0aW9uYXJuGMLw11sgASgJUg9ub3RpZm'
    'ljYXRpb25hcm4SSQoSbm90aWZpY2F0aW9uZXZlbnRzGLSb9AcgAygOMhYuc3NtLk5vdGlmaWNh'
    'dGlvbkV2ZW50UhJub3RpZmljYXRpb25ldmVudHMSRAoQbm90aWZpY2F0aW9udHlwZRi3m+UUIA'
    'EoDjIVLnNzbS5Ob3RpZmljYXRpb25UeXBlUhBub3RpZmljYXRpb250eXBl');

@$core.Deprecated('Use opsAggregatorDescriptor instead')
const OpsAggregator$json = {
  '1': 'OpsAggregator',
  '2': [
    {
      '1': 'aggregatortype',
      '3': 328765509,
      '4': 1,
      '5': 9,
      '10': 'aggregatortype'
    },
    {
      '1': 'aggregators',
      '3': 161545870,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsAggregator',
      '10': 'aggregators'
    },
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'filters',
      '3': 188393197,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsFilter',
      '10': 'filters'
    },
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
    {
      '1': 'values',
      '3': 223158876,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsAggregator.ValuesEntry',
      '10': 'values'
    },
  ],
  '3': [OpsAggregator_ValuesEntry$json],
};

@$core.Deprecated('Use opsAggregatorDescriptor instead')
const OpsAggregator_ValuesEntry$json = {
  '1': 'ValuesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `OpsAggregator`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsAggregatorDescriptor = $convert.base64Decode(
    'Cg1PcHNBZ2dyZWdhdG9yEioKDmFnZ3JlZ2F0b3J0eXBlGMWg4pwBIAEoCVIOYWdncmVnYXRvcn'
    'R5cGUSNwoLYWdncmVnYXRvcnMYjv2DTSADKAsyEi5zc20uT3BzQWdncmVnYXRvclILYWdncmVn'
    'YXRvcnMSKAoNYXR0cmlidXRlbmFtZRitlZioASABKAlSDWF0dHJpYnV0ZW5hbWUSKwoHZmlsdG'
    'VycxjtzepZIAMoCzIOLnNzbS5PcHNGaWx0ZXJSB2ZpbHRlcnMSHgoIdHlwZW5hbWUYz87Z1AEg'
    'ASgJUgh0eXBlbmFtZRI5CgZ2YWx1ZXMY3MS0aiADKAsyHi5zc20uT3BzQWdncmVnYXRvci5WYW'
    'x1ZXNFbnRyeVIGdmFsdWVzGjkKC1ZhbHVlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZh'
    'bHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use opsEntityDescriptor instead')
const OpsEntity$json = {
  '1': 'OpsEntity',
  '2': [
    {
      '1': 'data',
      '3': 525498822,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsEntity.DataEntry',
      '10': 'data'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
  '3': [OpsEntity_DataEntry$json],
};

@$core.Deprecated('Use opsEntityDescriptor instead')
const OpsEntity_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsEntityItem',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `OpsEntity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsEntityDescriptor = $convert.base64Decode(
    'CglPcHNFbnRpdHkSMAoEZGF0YRjG88n6ASADKAsyGC5zc20uT3BzRW50aXR5LkRhdGFFbnRyeV'
    'IEZGF0YRISCgJpZBiB8qK3ASABKAlSAmlkGksKCURhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIoCgV2YWx1ZRgCIAEoCzISLnNzbS5PcHNFbnRpdHlJdGVtUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use opsEntityItemDescriptor instead')
const OpsEntityItem$json = {
  '1': 'OpsEntityItem',
  '2': [
    {'1': 'capturetime', '3': 72276089, '4': 1, '5': 9, '10': 'capturetime'},
    {
      '1': 'content',
      '3': 23568227,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsEntityItem.ContentEntry',
      '10': 'content'
    },
  ],
  '3': [OpsEntityItem_ContentEntry$json],
};

@$core.Deprecated('Use opsEntityItemDescriptor instead')
const OpsEntityItem_ContentEntry$json = {
  '1': 'ContentEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `OpsEntityItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsEntityItemDescriptor = $convert.base64Decode(
    'Cg1PcHNFbnRpdHlJdGVtEiMKC2NhcHR1cmV0aW1lGPmwuyIgASgJUgtjYXB0dXJldGltZRI8Cg'
    'djb250ZW50GOO+ngsgAygLMh8uc3NtLk9wc0VudGl0eUl0ZW0uQ29udGVudEVudHJ5Ugdjb250'
    'ZW50GjoKDENvbnRlbnRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdm'
    'FsdWU6AjgB');

@$core.Deprecated('Use opsFilterDescriptor instead')
const OpsFilter$json = {
  '1': 'OpsFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsFilterOperatorType',
      '10': 'type'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `OpsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsFilterDescriptor = $convert.base64Decode(
    'CglPcHNGaWx0ZXISEwoDa2V5GI2S62ggASgJUgNrZXkSMgoEdHlwZRjuoNeKASABKA4yGi5zc2'
    '0uT3BzRmlsdGVyT3BlcmF0b3JUeXBlUgR0eXBlEhkKBnZhbHVlcxjcxLRqIAMoCVIGdmFsdWVz');

@$core.Deprecated('Use opsItemDescriptor instead')
const OpsItem$json = {
  '1': 'OpsItem',
  '2': [
    {
      '1': 'actualendtime',
      '3': 452787526,
      '4': 1,
      '5': 9,
      '10': 'actualendtime'
    },
    {
      '1': 'actualstarttime',
      '3': 532853117,
      '4': 1,
      '5': 9,
      '10': 'actualstarttime'
    },
    {'1': 'category', '3': 263447954, '4': 1, '5': 9, '10': 'category'},
    {'1': 'createdby', '3': 73491847, '4': 1, '5': 9, '10': 'createdby'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'lastmodifiedby',
      '3': 119539216,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedby'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'notifications',
      '3': 404560144,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemNotification',
      '10': 'notifications'
    },
    {
      '1': 'operationaldata',
      '3': 360357940,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItem.OperationaldataEntry',
      '10': 'operationaldata'
    },
    {'1': 'opsitemarn', '3': 222489428, '4': 1, '5': 9, '10': 'opsitemarn'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'opsitemtype', '3': 317873397, '4': 1, '5': 9, '10': 'opsitemtype'},
    {
      '1': 'plannedendtime',
      '3': 245727820,
      '4': 1,
      '5': 9,
      '10': 'plannedendtime'
    },
    {
      '1': 'plannedstarttime',
      '3': 478079215,
      '4': 1,
      '5': 9,
      '10': 'plannedstarttime'
    },
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'relatedopsitems',
      '3': 287082393,
      '4': 3,
      '5': 11,
      '6': '.ssm.RelatedOpsItem',
      '10': 'relatedopsitems'
    },
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
  '3': [OpsItem_OperationaldataEntry$json],
};

@$core.Deprecated('Use opsItemDescriptor instead')
const OpsItem_OperationaldataEntry$json = {
  '1': 'OperationaldataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemDataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `OpsItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemDescriptor = $convert.base64Decode(
    'CgdPcHNJdGVtEigKDWFjdHVhbGVuZHRpbWUYxvrz1wEgASgJUg1hY3R1YWxlbmR0aW1lEiwKD2'
    'FjdHVhbHN0YXJ0dGltZRj94or+ASABKAlSD2FjdHVhbHN0YXJ0dGltZRIdCghjYXRlZ29yeRiS'
    'y899IAEoCVIIY2F0ZWdvcnkSHwoJY3JlYXRlZGJ5GIfLhSMgASgJUgljcmVhdGVkYnkSIwoLY3'
    'JlYXRlZHRpbWUY8+vzOSABKAlSC2NyZWF0ZWR0aW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJ'
    'UgtkZXNjcmlwdGlvbhIpCg5sYXN0bW9kaWZpZWRieRiQjIA5IAEoCVIObGFzdG1vZGlmaWVkYn'
    'kSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCVIQbGFzdG1vZGlmaWVkdGltZRJCCg1ub3Rp'
    'ZmljYXRpb25zGJCy9MABIAMoCzIYLnNzbS5PcHNJdGVtTm90aWZpY2F0aW9uUg1ub3RpZmljYX'
    'Rpb25zEk8KD29wZXJhdGlvbmFsZGF0YRi0wOqrASADKAsyIS5zc20uT3BzSXRlbS5PcGVyYXRp'
    'b25hbGRhdGFFbnRyeVIPb3BlcmF0aW9uYWxkYXRhEiEKCm9wc2l0ZW1hcm4Y1NaLaiABKAlSCm'
    '9wc2l0ZW1hcm4SHwoJb3BzaXRlbWlkGNLSlQwgASgJUglvcHNpdGVtaWQSJAoLb3BzaXRlbXR5'
    'cGUY9bnJlwEgASgJUgtvcHNpdGVtdHlwZRIpCg5wbGFubmVkZW5kdGltZRjMhJZ1IAEoCVIOcG'
    'xhbm5lZGVuZHRpbWUSLgoQcGxhbm5lZHN0YXJ0dGltZRjv0fvjASABKAlSEHBsYW5uZWRzdGFy'
    'dHRpbWUSHQoIcHJpb3JpdHkYqr62NCABKAVSCHByaW9yaXR5EkEKD3JlbGF0ZWRvcHNpdGVtcx'
    'iZj/KIASADKAsyEy5zc20uUmVsYXRlZE9wc0l0ZW1SD3JlbGF0ZWRvcHNpdGVtcxIeCghzZXZl'
    'cml0eRjT5YOEASABKAlSCHNldmVyaXR5EhkKBnNvdXJjZRj5x4oPIAEoCVIGc291cmNlEi0KBn'
    'N0YXR1cxiQ5PsCIAEoDjISLnNzbS5PcHNJdGVtU3RhdHVzUgZzdGF0dXMSFwoFdGl0bGUYquPR'
    'JiABKAlSBXRpdGxlEhwKB3ZlcnNpb24YuKq37gEgASgJUgd2ZXJzaW9uGlkKFE9wZXJhdGlvbm'
    'FsZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EisKBXZhbHVlGAIgASgLMhUuc3NtLk9wc0l0'
    'ZW1EYXRhVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use opsItemAccessDeniedExceptionDescriptor instead')
const OpsItemAccessDeniedException$json = {
  '1': 'OpsItemAccessDeniedException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsItemAccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemAccessDeniedExceptionDescriptor =
    $convert.base64Decode(
        'ChxPcHNJdGVtQWNjZXNzRGVuaWVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use opsItemAlreadyExistsExceptionDescriptor instead')
const OpsItemAlreadyExistsException$json = {
  '1': 'OpsItemAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `OpsItemAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'Ch1PcHNJdGVtQWxyZWFkeUV4aXN0c0V4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdlEh8KCW9wc2l0ZW1pZBjS0pUMIAEoCVIJb3BzaXRlbWlk');

@$core.Deprecated('Use opsItemConflictExceptionDescriptor instead')
const OpsItemConflictException$json = {
  '1': 'OpsItemConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsItemConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChhPcHNJdGVtQ29uZmxpY3RFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use opsItemDataValueDescriptor instead')
const OpsItemDataValue$json = {
  '1': 'OpsItemDataValue',
  '2': [
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemDataType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `OpsItemDataValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemDataValueDescriptor = $convert.base64Decode(
    'ChBPcHNJdGVtRGF0YVZhbHVlEiwKBHR5cGUY7qDXigEgASgOMhQuc3NtLk9wc0l0ZW1EYXRhVH'
    'lwZVIEdHlwZRIYCgV2YWx1ZRjr8p+KASABKAlSBXZhbHVl');

@$core.Deprecated('Use opsItemEventFilterDescriptor instead')
const OpsItemEventFilter$json = {
  '1': 'OpsItemEventFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemEventFilterKey',
      '10': 'key'
    },
    {
      '1': 'operator',
      '3': 31807518,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemEventFilterOperator',
      '10': 'operator'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `OpsItemEventFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemEventFilterDescriptor = $convert.base64Decode(
    'ChJPcHNJdGVtRXZlbnRGaWx0ZXISLwoDa2V5GI2S62ggASgOMhouc3NtLk9wc0l0ZW1FdmVudE'
    'ZpbHRlcktleVIDa2V5Ej4KCG9wZXJhdG9yGJ6wlQ8gASgOMh8uc3NtLk9wc0l0ZW1FdmVudEZp'
    'bHRlck9wZXJhdG9yUghvcGVyYXRvchIZCgZ2YWx1ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use opsItemEventSummaryDescriptor instead')
const OpsItemEventSummary$json = {
  '1': 'OpsItemEventSummary',
  '2': [
    {
      '1': 'createdby',
      '3': 73491847,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemIdentity',
      '10': 'createdby'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {'1': 'detail', '3': 340030389, '4': 1, '5': 9, '10': 'detail'},
    {'1': 'detailtype', '3': 11758589, '4': 1, '5': 9, '10': 'detailtype'},
    {'1': 'eventid', '3': 376916819, '4': 1, '5': 9, '10': 'eventid'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `OpsItemEventSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemEventSummaryDescriptor = $convert.base64Decode(
    'ChNPcHNJdGVtRXZlbnRTdW1tYXJ5EjUKCWNyZWF0ZWRieRiHy4UjIAEoCzIULnNzbS5PcHNJdG'
    'VtSWRlbnRpdHlSCWNyZWF0ZWRieRIjCgtjcmVhdGVkdGltZRjz6/M5IAEoCVILY3JlYXRlZHRp'
    'bWUSGgoGZGV0YWlsGLXnkaIBIAEoCVIGZGV0YWlsEiEKCmRldGFpbHR5cGUY/dfNBSABKAlSCm'
    'RldGFpbHR5cGUSHAoHZXZlbnRpZBjTlt2zASABKAlSB2V2ZW50aWQSHwoJb3BzaXRlbWlkGNLS'
    'lQwgASgJUglvcHNpdGVtaWQSGQoGc291cmNlGPnHig8gASgJUgZzb3VyY2U=');

@$core.Deprecated('Use opsItemFilterDescriptor instead')
const OpsItemFilter$json = {
  '1': 'OpsItemFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemFilterKey',
      '10': 'key'
    },
    {
      '1': 'operator',
      '3': 31807518,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemFilterOperator',
      '10': 'operator'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `OpsItemFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemFilterDescriptor = $convert.base64Decode(
    'Cg1PcHNJdGVtRmlsdGVyEioKA2tleRiNkutoIAEoDjIVLnNzbS5PcHNJdGVtRmlsdGVyS2V5Ug'
    'NrZXkSOQoIb3BlcmF0b3IYnrCVDyABKA4yGi5zc20uT3BzSXRlbUZpbHRlck9wZXJhdG9yUghv'
    'cGVyYXRvchIZCgZ2YWx1ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use opsItemIdentityDescriptor instead')
const OpsItemIdentity$json = {
  '1': 'OpsItemIdentity',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `OpsItemIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemIdentityDescriptor = $convert
    .base64Decode('Cg9PcHNJdGVtSWRlbnRpdHkSFAoDYXJuGJ2b7b8BIAEoCVIDYXJu');

@$core.Deprecated('Use opsItemInvalidParameterExceptionDescriptor instead')
const OpsItemInvalidParameterException$json = {
  '1': 'OpsItemInvalidParameterException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'parameternames',
      '3': 72522785,
      '4': 3,
      '5': 9,
      '10': 'parameternames'
    },
  ],
};

/// Descriptor for `OpsItemInvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemInvalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'CiBPcHNJdGVtSW52YWxpZFBhcmFtZXRlckV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUg'
        'dtZXNzYWdlEikKDnBhcmFtZXRlcm5hbWVzGKG4yiIgAygJUg5wYXJhbWV0ZXJuYW1lcw==');

@$core.Deprecated('Use opsItemLimitExceededExceptionDescriptor instead')
const OpsItemLimitExceededException$json = {
  '1': 'OpsItemLimitExceededException',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'limittype', '3': 155546973, '4': 1, '5': 9, '10': 'limittype'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'resourcetypes',
      '3': 343086443,
      '4': 3,
      '5': 9,
      '10': 'resourcetypes'
    },
  ],
};

/// Descriptor for `OpsItemLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1PcHNJdGVtTGltaXRFeGNlZWRlZEV4Y2VwdGlvbhIYCgVsaW1pdBjVldnEASABKAVSBWxpbW'
        'l0Eh8KCWxpbWl0dHlwZRjd6pVKIAEoCVIJbGltaXR0eXBlEhsKB21lc3NhZ2UYhbO7cCABKAlS'
        'B21lc3NhZ2USKAoNcmVzb3VyY2V0eXBlcxjrqsyjASADKAlSDXJlc291cmNldHlwZXM=');

@$core.Deprecated('Use opsItemNotFoundExceptionDescriptor instead')
const OpsItemNotFoundException$json = {
  '1': 'OpsItemNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsItemNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChhPcHNJdGVtTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use opsItemNotificationDescriptor instead')
const OpsItemNotification$json = {
  '1': 'OpsItemNotification',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `OpsItemNotification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemNotificationDescriptor =
    $convert.base64Decode(
        'ChNPcHNJdGVtTm90aWZpY2F0aW9uEhQKA2Fybhidm+2/ASABKAlSA2Fybg==');

@$core.Deprecated(
    'Use opsItemRelatedItemAlreadyExistsExceptionDescriptor instead')
const OpsItemRelatedItemAlreadyExistsException$json = {
  '1': 'OpsItemRelatedItemAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'resourceuri', '3': 442608014, '4': 1, '5': 9, '10': 'resourceuri'},
  ],
};

/// Descriptor for `OpsItemRelatedItemAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemRelatedItemAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CihPcHNJdGVtUmVsYXRlZEl0ZW1BbHJlYWR5RXhpc3RzRXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2USHwoJb3BzaXRlbWlkGNLSlQwgASgJUglvcHNpdGVtaWQSJAoLcmVz'
        'b3VyY2V1cmkYjtOG0wEgASgJUgtyZXNvdXJjZXVyaQ==');

@$core.Deprecated(
    'Use opsItemRelatedItemAssociationNotFoundExceptionDescriptor instead')
const OpsItemRelatedItemAssociationNotFoundException$json = {
  '1': 'OpsItemRelatedItemAssociationNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsItemRelatedItemAssociationNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    opsItemRelatedItemAssociationNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ci5PcHNJdGVtUmVsYXRlZEl0ZW1Bc3NvY2lhdGlvbk5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3'
        'NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use opsItemRelatedItemSummaryDescriptor instead')
const OpsItemRelatedItemSummary$json = {
  '1': 'OpsItemRelatedItemSummary',
  '2': [
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationtype',
      '3': 165803189,
      '4': 1,
      '5': 9,
      '10': 'associationtype'
    },
    {
      '1': 'createdby',
      '3': 73491847,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemIdentity',
      '10': 'createdby'
    },
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {
      '1': 'lastmodifiedby',
      '3': 119539216,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemIdentity',
      '10': 'lastmodifiedby'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'resourceuri', '3': 442608014, '4': 1, '5': 9, '10': 'resourceuri'},
  ],
};

/// Descriptor for `OpsItemRelatedItemSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemRelatedItemSummaryDescriptor = $convert.base64Decode(
    'ChlPcHNJdGVtUmVsYXRlZEl0ZW1TdW1tYXJ5EicKDWFzc29jaWF0aW9uaWQYkvyVQiABKAlSDW'
    'Fzc29jaWF0aW9uaWQSKwoPYXNzb2NpYXRpb250eXBlGLXph08gASgJUg9hc3NvY2lhdGlvbnR5'
    'cGUSNQoJY3JlYXRlZGJ5GIfLhSMgASgLMhQuc3NtLk9wc0l0ZW1JZGVudGl0eVIJY3JlYXRlZG'
    'J5EiMKC2NyZWF0ZWR0aW1lGPPr8zkgASgJUgtjcmVhdGVkdGltZRI/Cg5sYXN0bW9kaWZpZWRi'
    'eRiQjIA5IAEoCzIULnNzbS5PcHNJdGVtSWRlbnRpdHlSDmxhc3Rtb2RpZmllZGJ5Ei0KEGxhc3'
    'Rtb2RpZmllZHRpbWUY4IL8cCABKAlSEGxhc3Rtb2RpZmllZHRpbWUSHwoJb3BzaXRlbWlkGNLS'
    'lQwgASgJUglvcHNpdGVtaWQSJgoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoCVIMcmVzb3VyY2V0eX'
    'BlEiQKC3Jlc291cmNldXJpGI7ThtMBIAEoCVILcmVzb3VyY2V1cmk=');

@$core.Deprecated('Use opsItemRelatedItemsFilterDescriptor instead')
const OpsItemRelatedItemsFilter$json = {
  '1': 'OpsItemRelatedItemsFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemRelatedItemsFilterKey',
      '10': 'key'
    },
    {
      '1': 'operator',
      '3': 31807518,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemRelatedItemsFilterOperator',
      '10': 'operator'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `OpsItemRelatedItemsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemRelatedItemsFilterDescriptor = $convert.base64Decode(
    'ChlPcHNJdGVtUmVsYXRlZEl0ZW1zRmlsdGVyEjYKA2tleRiNkutoIAEoDjIhLnNzbS5PcHNJdG'
    'VtUmVsYXRlZEl0ZW1zRmlsdGVyS2V5UgNrZXkSRQoIb3BlcmF0b3IYnrCVDyABKA4yJi5zc20u'
    'T3BzSXRlbVJlbGF0ZWRJdGVtc0ZpbHRlck9wZXJhdG9yUghvcGVyYXRvchIZCgZ2YWx1ZXMY3M'
    'S0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use opsItemSummaryDescriptor instead')
const OpsItemSummary$json = {
  '1': 'OpsItemSummary',
  '2': [
    {
      '1': 'actualendtime',
      '3': 452787526,
      '4': 1,
      '5': 9,
      '10': 'actualendtime'
    },
    {
      '1': 'actualstarttime',
      '3': 532853117,
      '4': 1,
      '5': 9,
      '10': 'actualstarttime'
    },
    {'1': 'category', '3': 263447954, '4': 1, '5': 9, '10': 'category'},
    {'1': 'createdby', '3': 73491847, '4': 1, '5': 9, '10': 'createdby'},
    {'1': 'createdtime', '3': 121435635, '4': 1, '5': 9, '10': 'createdtime'},
    {
      '1': 'lastmodifiedby',
      '3': 119539216,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedby'
    },
    {
      '1': 'lastmodifiedtime',
      '3': 236912992,
      '4': 1,
      '5': 9,
      '10': 'lastmodifiedtime'
    },
    {
      '1': 'operationaldata',
      '3': 360357940,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemSummary.OperationaldataEntry',
      '10': 'operationaldata'
    },
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {'1': 'opsitemtype', '3': 317873397, '4': 1, '5': 9, '10': 'opsitemtype'},
    {
      '1': 'plannedendtime',
      '3': 245727820,
      '4': 1,
      '5': 9,
      '10': 'plannedendtime'
    },
    {
      '1': 'plannedstarttime',
      '3': 478079215,
      '4': 1,
      '5': 9,
      '10': 'plannedstarttime'
    },
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
  '3': [OpsItemSummary_OperationaldataEntry$json],
};

@$core.Deprecated('Use opsItemSummaryDescriptor instead')
const OpsItemSummary_OperationaldataEntry$json = {
  '1': 'OperationaldataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemDataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `OpsItemSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsItemSummaryDescriptor = $convert.base64Decode(
    'Cg5PcHNJdGVtU3VtbWFyeRIoCg1hY3R1YWxlbmR0aW1lGMb689cBIAEoCVINYWN0dWFsZW5kdG'
    'ltZRIsCg9hY3R1YWxzdGFydHRpbWUY/eKK/gEgASgJUg9hY3R1YWxzdGFydHRpbWUSHQoIY2F0'
    'ZWdvcnkYksvPfSABKAlSCGNhdGVnb3J5Eh8KCWNyZWF0ZWRieRiHy4UjIAEoCVIJY3JlYXRlZG'
    'J5EiMKC2NyZWF0ZWR0aW1lGPPr8zkgASgJUgtjcmVhdGVkdGltZRIpCg5sYXN0bW9kaWZpZWRi'
    'eRiQjIA5IAEoCVIObGFzdG1vZGlmaWVkYnkSLQoQbGFzdG1vZGlmaWVkdGltZRjggvxwIAEoCV'
    'IQbGFzdG1vZGlmaWVkdGltZRJWCg9vcGVyYXRpb25hbGRhdGEYtMDqqwEgAygLMiguc3NtLk9w'
    'c0l0ZW1TdW1tYXJ5Lk9wZXJhdGlvbmFsZGF0YUVudHJ5Ug9vcGVyYXRpb25hbGRhdGESHwoJb3'
    'BzaXRlbWlkGNLSlQwgASgJUglvcHNpdGVtaWQSJAoLb3BzaXRlbXR5cGUY9bnJlwEgASgJUgtv'
    'cHNpdGVtdHlwZRIpCg5wbGFubmVkZW5kdGltZRjMhJZ1IAEoCVIOcGxhbm5lZGVuZHRpbWUSLg'
    'oQcGxhbm5lZHN0YXJ0dGltZRjv0fvjASABKAlSEHBsYW5uZWRzdGFydHRpbWUSHQoIcHJpb3Jp'
    'dHkYqr62NCABKAVSCHByaW9yaXR5Eh4KCHNldmVyaXR5GNPlg4QBIAEoCVIIc2V2ZXJpdHkSGQ'
    'oGc291cmNlGPnHig8gASgJUgZzb3VyY2USLQoGc3RhdHVzGJDk+wIgASgOMhIuc3NtLk9wc0l0'
    'ZW1TdGF0dXNSBnN0YXR1cxIXCgV0aXRsZRiq49EmIAEoCVIFdGl0bGUaWQoUT3BlcmF0aW9uYW'
    'xkYXRhRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSKwoFdmFsdWUYAiABKAsyFS5zc20uT3BzSXRl'
    'bURhdGFWYWx1ZVIFdmFsdWU6AjgB');

@$core.Deprecated('Use opsMetadataDescriptor instead')
const OpsMetadata$json = {
  '1': 'OpsMetadata',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'lastmodifieduser',
      '3': 215822778,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieduser'
    },
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
  ],
};

/// Descriptor for `OpsMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataDescriptor = $convert.base64Decode(
    'CgtPcHNNZXRhZGF0YRImCgxjcmVhdGlvbmRhdGUY4di3iQEgASgJUgxjcmVhdGlvbmRhdGUSLQ'
    'oQbGFzdG1vZGlmaWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVkZGF0ZRItChBsYXN0bW9k'
    'aWZpZWR1c2VyGLrj9GYgASgJUhBsYXN0bW9kaWZpZWR1c2VyEioKDm9wc21ldGFkYXRhYXJuGK'
    'K+guYBIAEoCVIOb3BzbWV0YWRhdGFhcm4SIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc291'
    'cmNlaWQ=');

@$core.Deprecated('Use opsMetadataAlreadyExistsExceptionDescriptor instead')
const OpsMetadataAlreadyExistsException$json = {
  '1': 'OpsMetadataAlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiFPcHNNZXRhZGF0YUFscmVhZHlFeGlzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use opsMetadataFilterDescriptor instead')
const OpsMetadataFilter$json = {
  '1': 'OpsMetadataFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `OpsMetadataFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataFilterDescriptor = $convert.base64Decode(
    'ChFPcHNNZXRhZGF0YUZpbHRlchITCgNrZXkYjZLraCABKAlSA2tleRIZCgZ2YWx1ZXMY3MS0ai'
    'ADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use opsMetadataInvalidArgumentExceptionDescriptor instead')
const OpsMetadataInvalidArgumentException$json = {
  '1': 'OpsMetadataInvalidArgumentException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataInvalidArgumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataInvalidArgumentExceptionDescriptor =
    $convert.base64Decode(
        'CiNPcHNNZXRhZGF0YUludmFsaWRBcmd1bWVudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use opsMetadataKeyLimitExceededExceptionDescriptor instead')
const OpsMetadataKeyLimitExceededException$json = {
  '1': 'OpsMetadataKeyLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataKeyLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataKeyLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiRPcHNNZXRhZGF0YUtleUxpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use opsMetadataLimitExceededExceptionDescriptor instead')
const OpsMetadataLimitExceededException$json = {
  '1': 'OpsMetadataLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiFPcHNNZXRhZGF0YUxpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use opsMetadataNotFoundExceptionDescriptor instead')
const OpsMetadataNotFoundException$json = {
  '1': 'OpsMetadataNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChxPcHNNZXRhZGF0YU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use opsMetadataTooManyUpdatesExceptionDescriptor instead')
const OpsMetadataTooManyUpdatesException$json = {
  '1': 'OpsMetadataTooManyUpdatesException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OpsMetadataTooManyUpdatesException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsMetadataTooManyUpdatesExceptionDescriptor =
    $convert.base64Decode(
        'CiJPcHNNZXRhZGF0YVRvb01hbnlVcGRhdGVzRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use opsResultAttributeDescriptor instead')
const OpsResultAttribute$json = {
  '1': 'OpsResultAttribute',
  '2': [
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `OpsResultAttribute`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List opsResultAttributeDescriptor = $convert.base64Decode(
    'ChJPcHNSZXN1bHRBdHRyaWJ1dGUSHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ==');

@$core.Deprecated('Use outputSourceDescriptor instead')
const OutputSource$json = {
  '1': 'OutputSource',
  '2': [
    {
      '1': 'outputsourceid',
      '3': 305073205,
      '4': 1,
      '5': 9,
      '10': 'outputsourceid'
    },
    {
      '1': 'outputsourcetype',
      '3': 491822634,
      '4': 1,
      '5': 9,
      '10': 'outputsourcetype'
    },
  ],
};

/// Descriptor for `OutputSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List outputSourceDescriptor = $convert.base64Decode(
    'CgxPdXRwdXRTb3VyY2USKgoOb3V0cHV0c291cmNlaWQYtZi8kQEgASgJUg5vdXRwdXRzb3VyY2'
    'VpZBIuChBvdXRwdXRzb3VyY2V0eXBlGKq8wuoBIAEoCVIQb3V0cHV0c291cmNldHlwZQ==');

@$core.Deprecated('Use parameterDescriptor instead')
const Parameter$json = {
  '1': 'Parameter',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'selector', '3': 353735197, '4': 1, '5': 9, '10': 'selector'},
    {'1': 'sourceresult', '3': 320048386, '4': 1, '5': 9, '10': 'sourceresult'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `Parameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterDescriptor = $convert.base64Decode(
    'CglQYXJhbWV0ZXISFAoDYXJuGJ2cr70BIAEoCVIDYXJuEh0KCGRhdGF0eXBlGO7YtSAgASgJUg'
    'hkYXRhdHlwZRItChBsYXN0bW9kaWZpZWRkYXRlGNOIyAsgASgJUhBsYXN0bW9kaWZpZWRkYXRl'
    'EhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSHgoIc2VsZWN0b3IYnaTWqAEgASgJUghzZWxlY3Rvch'
    'ImCgxzb3VyY2VyZXN1bHQYgprOmAEgASgJUgxzb3VyY2VyZXN1bHQSKgoEdHlwZRjuoNeKASAB'
    'KA4yEi5zc20uUGFyYW1ldGVyVHlwZVIEdHlwZRIYCgV2YWx1ZRjr8p+KASABKAlSBXZhbHVlEh'
    'wKB3ZlcnNpb24YuKq37gEgASgDUgd2ZXJzaW9u');

@$core.Deprecated('Use parameterAlreadyExistsDescriptor instead')
const ParameterAlreadyExists$json = {
  '1': 'ParameterAlreadyExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterAlreadyExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterAlreadyExistsDescriptor =
    $convert.base64Decode(
        'ChZQYXJhbWV0ZXJBbHJlYWR5RXhpc3RzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use parameterHistoryDescriptor instead')
const ParameterHistory$json = {
  '1': 'ParameterHistory',
  '2': [
    {
      '1': 'allowedpattern',
      '3': 290284364,
      '4': 1,
      '5': 9,
      '10': 'allowedpattern'
    },
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'labels', '3': 178416811, '4': 3, '5': 9, '10': 'labels'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'lastmodifieduser',
      '3': 215822778,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieduser'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterInlinePolicy',
      '10': 'policies'
    },
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterTier',
      '10': 'tier'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `ParameterHistory`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterHistoryDescriptor = $convert.base64Decode(
    'ChBQYXJhbWV0ZXJIaXN0b3J5EioKDmFsbG93ZWRwYXR0ZXJuGMzGtYoBIAEoCVIOYWxsb3dlZH'
    'BhdHRlcm4SHQoIZGF0YXR5cGUY7ti1ICABKAlSCGRhdGF0eXBlEiMKC2Rlc2NyaXB0aW9uGIr0'
    '+TYgASgJUgtkZXNjcmlwdGlvbhIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhkKBmxhYmVscx'
    'ir2YlVIAMoCVIGbGFiZWxzEi0KEGxhc3Rtb2RpZmllZGRhdGUY04jICyABKAlSEGxhc3Rtb2Rp'
    'ZmllZGRhdGUSLQoQbGFzdG1vZGlmaWVkdXNlchi64/RmIAEoCVIQbGFzdG1vZGlmaWVkdXNlch'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjkKCHBvbGljaWVzGJisihMgAygLMhouc3NtLlBhcmFt'
    'ZXRlcklubGluZVBvbGljeVIIcG9saWNpZXMSKgoEdGllchiq1OH3ASABKA4yEi5zc20uUGFyYW'
    '1ldGVyVGllclIEdGllchIqCgR0eXBlGO6g14oBIAEoDjISLnNzbS5QYXJhbWV0ZXJUeXBlUgR0'
    'eXBlEhgKBXZhbHVlGOvyn4oBIAEoCVIFdmFsdWUSHAoHdmVyc2lvbhi4qrfuASABKANSB3Zlcn'
    'Npb24=');

@$core.Deprecated('Use parameterInlinePolicyDescriptor instead')
const ParameterInlinePolicy$json = {
  '1': 'ParameterInlinePolicy',
  '2': [
    {'1': 'policystatus', '3': 256036138, '4': 1, '5': 9, '10': 'policystatus'},
    {'1': 'policytext', '3': 107473097, '4': 1, '5': 9, '10': 'policytext'},
    {'1': 'policytype', '3': 470619144, '4': 1, '5': 9, '10': 'policytype'},
  ],
};

/// Descriptor for `ParameterInlinePolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterInlinePolicyDescriptor = $convert.base64Decode(
    'ChVQYXJhbWV0ZXJJbmxpbmVQb2xpY3kSJQoMcG9saWN5c3RhdHVzGKqai3ogASgJUgxwb2xpY3'
    'lzdGF0dXMSIQoKcG9saWN5dGV4dBjJ0Z8zIAEoCVIKcG9saWN5dGV4dBIiCgpwb2xpY3l0eXBl'
    'GIiotOABIAEoCVIKcG9saWN5dHlwZQ==');

@$core.Deprecated('Use parameterLimitExceededDescriptor instead')
const ParameterLimitExceeded$json = {
  '1': 'ParameterLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterLimitExceededDescriptor =
    $convert.base64Decode(
        'ChZQYXJhbWV0ZXJMaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use parameterMaxVersionLimitExceededDescriptor instead')
const ParameterMaxVersionLimitExceeded$json = {
  '1': 'ParameterMaxVersionLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterMaxVersionLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterMaxVersionLimitExceededDescriptor =
    $convert.base64Decode(
        'CiBQYXJhbWV0ZXJNYXhWZXJzaW9uTGltaXRFeGNlZWRlZBIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use parameterMetadataDescriptor instead')
const ParameterMetadata$json = {
  '1': 'ParameterMetadata',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'allowedpattern',
      '3': 290284364,
      '4': 1,
      '5': 9,
      '10': 'allowedpattern'
    },
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'lastmodifieduser',
      '3': 215822778,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieduser'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.ssm.ParameterInlinePolicy',
      '10': 'policies'
    },
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterTier',
      '10': 'tier'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterType',
      '10': 'type'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `ParameterMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterMetadataDescriptor = $convert.base64Decode(
    'ChFQYXJhbWV0ZXJNZXRhZGF0YRIUCgNhcm4YnZyvvQEgASgJUgNhcm4SKgoOYWxsb3dlZHBhdH'
    'Rlcm4YzMa1igEgASgJUg5hbGxvd2VkcGF0dGVybhIdCghkYXRhdHlwZRju2LUgIAEoCVIIZGF0'
    'YXR5cGUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhgKBWtleWlkGKKAyI'
    'MBIAEoCVIFa2V5aWQSLQoQbGFzdG1vZGlmaWVkZGF0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVk'
    'ZGF0ZRItChBsYXN0bW9kaWZpZWR1c2VyGLrj9GYgASgJUhBsYXN0bW9kaWZpZWR1c2VyEhUKBG'
    '5hbWUYh+aBfyABKAlSBG5hbWUSOQoIcG9saWNpZXMYmKyKEyADKAsyGi5zc20uUGFyYW1ldGVy'
    'SW5saW5lUG9saWN5Ughwb2xpY2llcxIqCgR0aWVyGKrU4fcBIAEoDjISLnNzbS5QYXJhbWV0ZX'
    'JUaWVyUgR0aWVyEioKBHR5cGUY7qDXigEgASgOMhIuc3NtLlBhcmFtZXRlclR5cGVSBHR5cGUS'
    'HAoHdmVyc2lvbhi4qrfuASABKANSB3ZlcnNpb24=');

@$core.Deprecated('Use parameterNotFoundDescriptor instead')
const ParameterNotFound$json = {
  '1': 'ParameterNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterNotFoundDescriptor = $convert.base64Decode(
    'ChFQYXJhbWV0ZXJOb3RGb3VuZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use parameterPatternMismatchExceptionDescriptor instead')
const ParameterPatternMismatchException$json = {
  '1': 'ParameterPatternMismatchException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterPatternMismatchException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterPatternMismatchExceptionDescriptor =
    $convert.base64Decode(
        'CiFQYXJhbWV0ZXJQYXR0ZXJuTWlzbWF0Y2hFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use parameterStringFilterDescriptor instead')
const ParameterStringFilter$json = {
  '1': 'ParameterStringFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'option', '3': 388258997, '4': 1, '5': 9, '10': 'option'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `ParameterStringFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterStringFilterDescriptor = $convert.base64Decode(
    'ChVQYXJhbWV0ZXJTdHJpbmdGaWx0ZXISEwoDa2V5GI2S62ggASgJUgNrZXkSGgoGb3B0aW9uGL'
    'W5kbkBIAEoCVIGb3B0aW9uEhkKBnZhbHVlcxjcxLRqIAMoCVIGdmFsdWVz');

@$core.Deprecated('Use parameterVersionLabelLimitExceededDescriptor instead')
const ParameterVersionLabelLimitExceeded$json = {
  '1': 'ParameterVersionLabelLimitExceeded',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterVersionLabelLimitExceeded`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterVersionLabelLimitExceededDescriptor =
    $convert.base64Decode(
        'CiJQYXJhbWV0ZXJWZXJzaW9uTGFiZWxMaW1pdEV4Y2VlZGVkEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use parameterVersionNotFoundDescriptor instead')
const ParameterVersionNotFound$json = {
  '1': 'ParameterVersionNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ParameterVersionNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parameterVersionNotFoundDescriptor =
    $convert.base64Decode(
        'ChhQYXJhbWV0ZXJWZXJzaW9uTm90Rm91bmQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use parametersFilterDescriptor instead')
const ParametersFilter$json = {
  '1': 'ParametersFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParametersFilterKey',
      '10': 'key'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `ParametersFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parametersFilterDescriptor = $convert.base64Decode(
    'ChBQYXJhbWV0ZXJzRmlsdGVyEi0KA2tleRiNkutoIAEoDjIYLnNzbS5QYXJhbWV0ZXJzRmlsdG'
    'VyS2V5UgNrZXkSGQoGdmFsdWVzGNzEtGogAygJUgZ2YWx1ZXM=');

@$core.Deprecated('Use parentStepDetailsDescriptor instead')
const ParentStepDetails$json = {
  '1': 'ParentStepDetails',
  '2': [
    {'1': 'action', '3': 175614240, '4': 1, '5': 9, '10': 'action'},
    {'1': 'iteration', '3': 162626349, '4': 1, '5': 5, '10': 'iteration'},
    {
      '1': 'iteratorvalue',
      '3': 306879461,
      '4': 1,
      '5': 9,
      '10': 'iteratorvalue'
    },
    {
      '1': 'stepexecutionid',
      '3': 47252075,
      '4': 1,
      '5': 9,
      '10': 'stepexecutionid'
    },
    {'1': 'stepname', '3': 394700557, '4': 1, '5': 9, '10': 'stepname'},
  ],
};

/// Descriptor for `ParentStepDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List parentStepDetailsDescriptor = $convert.base64Decode(
    'ChFQYXJlbnRTdGVwRGV0YWlscxIZCgZhY3Rpb24YoNLeUyABKAlSBmFjdGlvbhIfCglpdGVyYX'
    'Rpb24YrfbFTSABKAVSCWl0ZXJhdGlvbhIoCg1pdGVyYXRvcnZhbHVlGOW3qpIBIAEoCVINaXRl'
    'cmF0b3J2YWx1ZRIrCg9zdGVwZXhlY3V0aW9uaWQY64TEFiABKAlSD3N0ZXBleGVjdXRpb25pZB'
    'IeCghzdGVwbmFtZRiNzpq8ASABKAlSCHN0ZXBuYW1l');

@$core.Deprecated('Use patchDescriptor instead')
const Patch$json = {
  '1': 'Patch',
  '2': [
    {'1': 'advisoryids', '3': 297948609, '4': 3, '5': 9, '10': 'advisoryids'},
    {'1': 'arch', '3': 312929514, '4': 1, '5': 9, '10': 'arch'},
    {'1': 'bugzillaids', '3': 31579172, '4': 3, '5': 9, '10': 'bugzillaids'},
    {'1': 'cveids', '3': 22965904, '4': 3, '5': 9, '10': 'cveids'},
    {
      '1': 'classification',
      '3': 259348066,
      '4': 1,
      '5': 9,
      '10': 'classification'
    },
    {'1': 'contenturl', '3': 199526508, '4': 1, '5': 9, '10': 'contenturl'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'epoch', '3': 190165127, '4': 1, '5': 5, '10': 'epoch'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'kbnumber', '3': 530478154, '4': 1, '5': 9, '10': 'kbnumber'},
    {'1': 'language', '3': 443800476, '4': 1, '5': 9, '10': 'language'},
    {'1': 'msrcnumber', '3': 387572764, '4': 1, '5': 9, '10': 'msrcnumber'},
    {'1': 'msrcseverity', '3': 364305242, '4': 1, '5': 9, '10': 'msrcseverity'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'product', '3': 312238285, '4': 1, '5': 9, '10': 'product'},
    {
      '1': 'productfamily',
      '3': 200630325,
      '4': 1,
      '5': 9,
      '10': 'productfamily'
    },
    {'1': 'release', '3': 220109599, '4': 1, '5': 9, '10': 'release'},
    {'1': 'releasedate', '3': 16273963, '4': 1, '5': 9, '10': 'releasedate'},
    {'1': 'repository', '3': 123979958, '4': 1, '5': 9, '10': 'repository'},
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
    {'1': 'vendor', '3': 499400308, '4': 1, '5': 9, '10': 'vendor'},
    {'1': 'version', '3': 500028728, '4': 1, '5': 9, '10': 'version'},
  ],
};

/// Descriptor for `Patch`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchDescriptor = $convert.base64Decode(
    'CgVQYXRjaBIkCgthZHZpc29yeWlkcxjBq4mOASADKAlSC2Fkdmlzb3J5aWRzEhYKBGFyY2gY6t'
    'mblQEgASgJUgRhcmNoEiMKC2J1Z3ppbGxhaWRzGKS4hw8gAygJUgtidWd6aWxsYWlkcxIZCgZj'
    'dmVpZHMYkN35CiADKAlSBmN2ZWlkcxIpCg5jbGFzc2lmaWNhdGlvbhjirNV7IAEoCVIOY2xhc3'
    'NpZmljYXRpb24SIQoKY29udGVudHVybBjskJJfIAEoCVIKY29udGVudHVybBIjCgtkZXNjcmlw'
    'dGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SFwoFZXBvY2gYh+HWWiABKAVSBWVwb2NoEhIKAm'
    'lkGIHyorcBIAEoCVICaWQSHgoIa2JudW1iZXIYyuj5/AEgASgJUghrYm51bWJlchIeCghsYW5n'
    'dWFnZRict8/TASABKAlSCGxhbmd1YWdlEiIKCm1zcmNudW1iZXIYnMjnuAEgASgJUgptc3Jjbn'
    'VtYmVyEiYKDG1zcmNzZXZlcml0eRjattutASABKAlSDG1zcmNzZXZlcml0eRIVCgRuYW1lGIfm'
    'gX8gASgJUgRuYW1lEhwKB3Byb2R1Y3QYzcHxlAEgASgJUgdwcm9kdWN0EicKDXByb2R1Y3RmYW'
    '1pbHkYtcDVXyABKAlSDXByb2R1Y3RmYW1pbHkSGwoHcmVsZWFzZRiftvpoIAEoCVIHcmVsZWFz'
    'ZRIjCgtyZWxlYXNlZGF0ZRirpOEHIAEoCVILcmVsZWFzZWRhdGUSIQoKcmVwb3NpdG9yeRi2kY'
    '87IAEoCVIKcmVwb3NpdG9yeRIeCghzZXZlcml0eRjT5YOEASABKAlSCHNldmVyaXR5EhcKBXRp'
    'dGxlGKrj0SYgASgJUgV0aXRsZRIaCgZ2ZW5kb3IY9PyQ7gEgASgJUgZ2ZW5kb3ISHAoHdmVyc2'
    'lvbhi4qrfuASABKAlSB3ZlcnNpb24=');

@$core.Deprecated('Use patchBaselineIdentityDescriptor instead')
const PatchBaselineIdentity$json = {
  '1': 'PatchBaselineIdentity',
  '2': [
    {
      '1': 'baselinedescription',
      '3': 29637269,
      '4': 1,
      '5': 9,
      '10': 'baselinedescription'
    },
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'baselinename', '3': 226231050, '4': 1, '5': 9, '10': 'baselinename'},
    {
      '1': 'defaultbaseline',
      '3': 511737610,
      '4': 1,
      '5': 8,
      '10': 'defaultbaseline'
    },
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
  ],
};

/// Descriptor for `PatchBaselineIdentity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchBaselineIdentityDescriptor = $convert.base64Decode(
    'ChVQYXRjaEJhc2VsaW5lSWRlbnRpdHkSMwoTYmFzZWxpbmVkZXNjcmlwdGlvbhiV9ZAOIAEoCV'
    'ITYmFzZWxpbmVkZXNjcmlwdGlvbhIhCgpiYXNlbGluZWlkGNDk2yggASgJUgpiYXNlbGluZWlk'
    'EiUKDGJhc2VsaW5lbmFtZRiKhvBrIAEoCVIMYmFzZWxpbmVuYW1lEiwKD2RlZmF1bHRiYXNlbG'
    'luZRiK/oH0ASABKAhSD2RlZmF1bHRiYXNlbGluZRJBCg9vcGVyYXRpbmdzeXN0ZW0Y6v3BEiAB'
    'KA4yFC5zc20uT3BlcmF0aW5nU3lzdGVtUg9vcGVyYXRpbmdzeXN0ZW0=');

@$core.Deprecated('Use patchComplianceDataDescriptor instead')
const PatchComplianceData$json = {
  '1': 'PatchComplianceData',
  '2': [
    {'1': 'cveids', '3': 22965904, '4': 1, '5': 9, '10': 'cveids'},
    {
      '1': 'classification',
      '3': 259348066,
      '4': 1,
      '5': 9,
      '10': 'classification'
    },
    {
      '1': 'installedtime',
      '3': 441523353,
      '4': 1,
      '5': 9,
      '10': 'installedtime'
    },
    {'1': 'kbid', '3': 121877862, '4': 1, '5': 9, '10': 'kbid'},
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {
      '1': 'state',
      '3': 502047895,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceDataState',
      '10': 'state'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
};

/// Descriptor for `PatchComplianceData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchComplianceDataDescriptor = $convert.base64Decode(
    'ChNQYXRjaENvbXBsaWFuY2VEYXRhEhkKBmN2ZWlkcxiQ3fkKIAEoCVIGY3ZlaWRzEikKDmNsYX'
    'NzaWZpY2F0aW9uGOKs1XsgASgJUg5jbGFzc2lmaWNhdGlvbhIoCg1pbnN0YWxsZWR0aW1lGJm5'
    'xNIBIAEoCVINaW5zdGFsbGVkdGltZRIVCgRrYmlkGObqjjogASgJUgRrYmlkEh4KCHNldmVyaX'
    'R5GNPlg4QBIAEoCVIIc2V2ZXJpdHkSNwoFc3RhdGUYl8my7wEgASgOMh0uc3NtLlBhdGNoQ29t'
    'cGxpYW5jZURhdGFTdGF0ZVIFc3RhdGUSFwoFdGl0bGUYquPRJiABKAlSBXRpdGxl');

@$core.Deprecated('Use patchFilterDescriptor instead')
const PatchFilter$json = {
  '1': 'PatchFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchFilterKey',
      '10': 'key'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `PatchFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchFilterDescriptor = $convert.base64Decode(
    'CgtQYXRjaEZpbHRlchIoCgNrZXkYjZLraCABKA4yEy5zc20uUGF0Y2hGaWx0ZXJLZXlSA2tleR'
    'IZCgZ2YWx1ZXMY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use patchFilterGroupDescriptor instead')
const PatchFilterGroup$json = {
  '1': 'PatchFilterGroup',
  '2': [
    {
      '1': 'patchfilters',
      '3': 148600313,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchFilter',
      '10': 'patchfilters'
    },
  ],
};

/// Descriptor for `PatchFilterGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchFilterGroupDescriptor = $convert.base64Decode(
    'ChBQYXRjaEZpbHRlckdyb3VwEjcKDHBhdGNoZmlsdGVycxj56+1GIAMoCzIQLnNzbS5QYXRjaE'
    'ZpbHRlclIMcGF0Y2hmaWx0ZXJz');

@$core.Deprecated('Use patchGroupPatchBaselineMappingDescriptor instead')
const PatchGroupPatchBaselineMapping$json = {
  '1': 'PatchGroupPatchBaselineMapping',
  '2': [
    {
      '1': 'baselineidentity',
      '3': 354904729,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchBaselineIdentity',
      '10': 'baselineidentity'
    },
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `PatchGroupPatchBaselineMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchGroupPatchBaselineMappingDescriptor =
    $convert.base64Decode(
        'Ch5QYXRjaEdyb3VwUGF0Y2hCYXNlbGluZU1hcHBpbmcSSgoQYmFzZWxpbmVpZGVudGl0eRiZ1Z'
        '2pASABKAsyGi5zc20uUGF0Y2hCYXNlbGluZUlkZW50aXR5UhBiYXNlbGluZWlkZW50aXR5EiIK'
        'CnBhdGNoZ3JvdXAY4bex9wEgASgJUgpwYXRjaGdyb3Vw');

@$core.Deprecated('Use patchOrchestratorFilterDescriptor instead')
const PatchOrchestratorFilter$json = {
  '1': 'PatchOrchestratorFilter',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `PatchOrchestratorFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchOrchestratorFilterDescriptor =
    $convert.base64Decode(
        'ChdQYXRjaE9yY2hlc3RyYXRvckZpbHRlchITCgNrZXkYjZLraCABKAlSA2tleRIZCgZ2YWx1ZX'
        'MY3MS0aiADKAlSBnZhbHVlcw==');

@$core.Deprecated('Use patchRuleDescriptor instead')
const PatchRule$json = {
  '1': 'PatchRule',
  '2': [
    {
      '1': 'approveafterdays',
      '3': 303775484,
      '4': 1,
      '5': 5,
      '10': 'approveafterdays'
    },
    {
      '1': 'approveuntildate',
      '3': 472356791,
      '4': 1,
      '5': 9,
      '10': 'approveuntildate'
    },
    {
      '1': 'compliancelevel',
      '3': 200031721,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'compliancelevel'
    },
    {
      '1': 'enablenonsecurity',
      '3': 135632920,
      '4': 1,
      '5': 8,
      '10': 'enablenonsecurity'
    },
    {
      '1': 'patchfiltergroup',
      '3': 432583347,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'patchfiltergroup'
    },
  ],
};

/// Descriptor for `PatchRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchRuleDescriptor = $convert.base64Decode(
    'CglQYXRjaFJ1bGUSLgoQYXBwcm92ZWFmdGVyZGF5cxj8/eyQASABKAVSEGFwcHJvdmVhZnRlcm'
    'RheXMSLgoQYXBwcm92ZXVudGlsZGF0ZRi3r57hASABKAlSEGFwcHJvdmV1bnRpbGRhdGUSRgoP'
    'Y29tcGxpYW5jZWxldmVsGOn7sF8gASgOMhkuc3NtLlBhdGNoQ29tcGxpYW5jZUxldmVsUg9jb2'
    '1wbGlhbmNlbGV2ZWwSLwoRZW5hYmxlbm9uc2VjdXJpdHkYmLDWQCABKAhSEWVuYWJsZW5vbnNl'
    'Y3VyaXR5EkUKEHBhdGNoZmlsdGVyZ3JvdXAYs+WizgEgASgLMhUuc3NtLlBhdGNoRmlsdGVyR3'
    'JvdXBSEHBhdGNoZmlsdGVyZ3JvdXA=');

@$core.Deprecated('Use patchRuleGroupDescriptor instead')
const PatchRuleGroup$json = {
  '1': 'PatchRuleGroup',
  '2': [
    {
      '1': 'patchrules',
      '3': 387495445,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchRule',
      '10': 'patchrules'
    },
  ],
};

/// Descriptor for `PatchRuleGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchRuleGroupDescriptor = $convert.base64Decode(
    'Cg5QYXRjaFJ1bGVHcm91cBIyCgpwYXRjaHJ1bGVzGJXs4rgBIAMoCzIOLnNzbS5QYXRjaFJ1bG'
    'VSCnBhdGNocnVsZXM=');

@$core.Deprecated('Use patchSourceDescriptor instead')
const PatchSource$json = {
  '1': 'PatchSource',
  '2': [
    {
      '1': 'configuration',
      '3': 442426458,
      '4': 1,
      '5': 9,
      '10': 'configuration'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'products', '3': 187452590, '4': 3, '5': 9, '10': 'products'},
  ],
};

/// Descriptor for `PatchSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchSourceDescriptor = $convert.base64Decode(
    'CgtQYXRjaFNvdXJjZRIoCg1jb25maWd1cmF0aW9uGNrI+9IBIAEoCVINY29uZmlndXJhdGlvbh'
    'IVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh0KCHByb2R1Y3RzGK6ZsVkgAygJUghwcm9kdWN0cw==');

@$core.Deprecated('Use patchStatusDescriptor instead')
const PatchStatus$json = {
  '1': 'PatchStatus',
  '2': [
    {'1': 'approvaldate', '3': 140150963, '4': 1, '5': 9, '10': 'approvaldate'},
    {
      '1': 'compliancelevel',
      '3': 200031721,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'compliancelevel'
    },
    {
      '1': 'deploymentstatus',
      '3': 440218913,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchDeploymentStatus',
      '10': 'deploymentstatus'
    },
  ],
};

/// Descriptor for `PatchStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List patchStatusDescriptor = $convert.base64Decode(
    'CgtQYXRjaFN0YXR1cxIlCgxhcHByb3ZhbGRhdGUYs5HqQiABKAlSDGFwcHJvdmFsZGF0ZRJGCg'
    '9jb21wbGlhbmNlbGV2ZWwY6fuwXyABKA4yGS5zc20uUGF0Y2hDb21wbGlhbmNlTGV2ZWxSD2Nv'
    'bXBsaWFuY2VsZXZlbBJKChBkZXBsb3ltZW50c3RhdHVzGKHq9NEBIAEoDjIaLnNzbS5QYXRjaE'
    'RlcGxveW1lbnRTdGF0dXNSEGRlcGxveW1lbnRzdGF0dXM=');

@$core.Deprecated('Use policiesLimitExceededExceptionDescriptor instead')
const PoliciesLimitExceededException$json = {
  '1': 'PoliciesLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PoliciesLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policiesLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch5Qb2xpY2llc0xpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use progressCountersDescriptor instead')
const ProgressCounters$json = {
  '1': 'ProgressCounters',
  '2': [
    {
      '1': 'cancelledsteps',
      '3': 265087304,
      '4': 1,
      '5': 5,
      '10': 'cancelledsteps'
    },
    {'1': 'failedsteps', '3': 517051856, '4': 1, '5': 5, '10': 'failedsteps'},
    {'1': 'successsteps', '3': 204073044, '4': 1, '5': 5, '10': 'successsteps'},
    {
      '1': 'timedoutsteps',
      '3': 85681760,
      '4': 1,
      '5': 5,
      '10': 'timedoutsteps'
    },
    {'1': 'totalsteps', '3': 82285501, '4': 1, '5': 5, '10': 'totalsteps'},
  ],
};

/// Descriptor for `ProgressCounters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List progressCountersDescriptor = $convert.base64Decode(
    'ChBQcm9ncmVzc0NvdW50ZXJzEikKDmNhbmNlbGxlZHN0ZXBzGMjSs34gASgFUg5jYW5jZWxsZW'
    'RzdGVwcxIkCgtmYWlsZWRzdGVwcxjQq8b2ASABKAVSC2ZhaWxlZHN0ZXBzEiUKDHN1Y2Nlc3Nz'
    'dGVwcxjU0KdhIAEoBVIMc3VjY2Vzc3N0ZXBzEicKDXRpbWVkb3V0c3RlcHMY4MztKCABKAVSDX'
    'RpbWVkb3V0c3RlcHMSIQoKdG90YWxzdGVwcxi9p54nIAEoBVIKdG90YWxzdGVwcw==');

@$core.Deprecated('Use putComplianceItemsRequestDescriptor instead')
const PutComplianceItemsRequest$json = {
  '1': 'PutComplianceItemsRequest',
  '2': [
    {
      '1': 'compliancetype',
      '3': 451385291,
      '4': 1,
      '5': 9,
      '10': 'compliancetype'
    },
    {
      '1': 'executionsummary',
      '3': 71480964,
      '4': 1,
      '5': 11,
      '6': '.ssm.ComplianceExecutionSummary',
      '10': 'executionsummary'
    },
    {
      '1': 'itemcontenthash',
      '3': 371661720,
      '4': 1,
      '5': 9,
      '10': 'itemcontenthash'
    },
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.ssm.ComplianceItemEntry',
      '10': 'items'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {
      '1': 'uploadtype',
      '3': 454755573,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceUploadType',
      '10': 'uploadtype'
    },
  ],
};

/// Descriptor for `PutComplianceItemsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putComplianceItemsRequestDescriptor = $convert.base64Decode(
    'ChlQdXRDb21wbGlhbmNlSXRlbXNSZXF1ZXN0EioKDmNvbXBsaWFuY2V0eXBlGMuvntcBIAEoCV'
    'IOY29tcGxpYW5jZXR5cGUSTgoQZXhlY3V0aW9uc3VtbWFyeRiE7YoiIAEoCzIfLnNzbS5Db21w'
    'bGlhbmNlRXhlY3V0aW9uU3VtbWFyeVIQZXhlY3V0aW9uc3VtbWFyeRIsCg9pdGVtY29udGVudG'
    'hhc2gYmLecsQEgASgJUg9pdGVtY29udGVudGhhc2gSMQoFaXRlbXMYsPDYASADKAsyGC5zc20u'
    'Q29tcGxpYW5jZUl0ZW1FbnRyeVIFaXRlbXMSIgoKcmVzb3VyY2VpZBiRuvH6ASABKAlSCnJlc2'
    '91cmNlaWQSJgoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoCVIMcmVzb3VyY2V0eXBlEj0KCnVwbG9h'
    'ZHR5cGUY9Yns2AEgASgOMhkuc3NtLkNvbXBsaWFuY2VVcGxvYWRUeXBlUgp1cGxvYWR0eXBl');

@$core.Deprecated('Use putComplianceItemsResultDescriptor instead')
const PutComplianceItemsResult$json = {
  '1': 'PutComplianceItemsResult',
};

/// Descriptor for `PutComplianceItemsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putComplianceItemsResultDescriptor =
    $convert.base64Decode('ChhQdXRDb21wbGlhbmNlSXRlbXNSZXN1bHQ=');

@$core.Deprecated('Use putInventoryRequestDescriptor instead')
const PutInventoryRequest$json = {
  '1': 'PutInventoryRequest',
  '2': [
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {
      '1': 'items',
      '3': 3553328,
      '4': 3,
      '5': 11,
      '6': '.ssm.InventoryItem',
      '10': 'items'
    },
  ],
};

/// Descriptor for `PutInventoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInventoryRequestDescriptor = $convert.base64Decode(
    'ChNQdXRJbnZlbnRvcnlSZXF1ZXN0EiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaW'
    'QSKwoFaXRlbXMYsPDYASADKAsyEi5zc20uSW52ZW50b3J5SXRlbVIFaXRlbXM=');

@$core.Deprecated('Use putInventoryResultDescriptor instead')
const PutInventoryResult$json = {
  '1': 'PutInventoryResult',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PutInventoryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putInventoryResultDescriptor =
    $convert.base64Decode(
        'ChJQdXRJbnZlbnRvcnlSZXN1bHQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use putParameterRequestDescriptor instead')
const PutParameterRequest$json = {
  '1': 'PutParameterRequest',
  '2': [
    {
      '1': 'allowedpattern',
      '3': 290284364,
      '4': 1,
      '5': 9,
      '10': 'allowedpattern'
    },
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'overwrite', '3': 419368595, '4': 1, '5': 8, '10': 'overwrite'},
    {'1': 'policies', '3': 40015384, '4': 1, '5': 9, '10': 'policies'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterTier',
      '10': 'tier'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterType',
      '10': 'type'
    },
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `PutParameterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putParameterRequestDescriptor = $convert.base64Decode(
    'ChNQdXRQYXJhbWV0ZXJSZXF1ZXN0EioKDmFsbG93ZWRwYXR0ZXJuGMzGtYoBIAEoCVIOYWxsb3'
    'dlZHBhdHRlcm4SHQoIZGF0YXR5cGUY7ti1ICABKAlSCGRhdGF0eXBlEiMKC2Rlc2NyaXB0aW9u'
    'GIr0+TYgASgJUgtkZXNjcmlwdGlvbhIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhUKBG5hbW'
    'UYh+aBfyABKAlSBG5hbWUSIAoJb3ZlcndyaXRlGJOd/McBIAEoCFIJb3ZlcndyaXRlEh0KCHBv'
    'bGljaWVzGJisihMgASgJUghwb2xpY2llcxIgCgR0YWdzGMHB9rUBIAMoCzIILnNzbS5UYWdSBH'
    'RhZ3MSKgoEdGllchiq1OH3ASABKA4yEi5zc20uUGFyYW1ldGVyVGllclIEdGllchIqCgR0eXBl'
    'GO6g14oBIAEoDjISLnNzbS5QYXJhbWV0ZXJUeXBlUgR0eXBlEhgKBXZhbHVlGOvyn4oBIAEoCV'
    'IFdmFsdWU=');

@$core.Deprecated('Use putParameterResultDescriptor instead')
const PutParameterResult$json = {
  '1': 'PutParameterResult',
  '2': [
    {
      '1': 'tier',
      '3': 519596586,
      '4': 1,
      '5': 14,
      '6': '.ssm.ParameterTier',
      '10': 'tier'
    },
    {'1': 'version', '3': 500028728, '4': 1, '5': 3, '10': 'version'},
  ],
};

/// Descriptor for `PutParameterResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putParameterResultDescriptor = $convert.base64Decode(
    'ChJQdXRQYXJhbWV0ZXJSZXN1bHQSKgoEdGllchiq1OH3ASABKA4yEi5zc20uUGFyYW1ldGVyVG'
    'llclIEdGllchIcCgd2ZXJzaW9uGLiqt+4BIAEoA1IHdmVyc2lvbg==');

@$core.Deprecated('Use putResourcePolicyRequestDescriptor instead')
const PutResourcePolicyRequest$json = {
  '1': 'PutResourcePolicyRequest',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyhash', '3': 203037360, '4': 1, '5': 9, '10': 'policyhash'},
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutResourcePolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyRequestDescriptor = $convert.base64Decode(
    'ChhQdXRSZXNvdXJjZVBvbGljeVJlcXVlc3QSGgoGcG9saWN5GKDv8OABIAEoCVIGcG9saWN5Ei'
    'EKCnBvbGljeWhhc2gYsLXoYCABKAlSCnBvbGljeWhhc2gSHgoIcG9saWN5aWQY86PpjgEgASgJ'
    'Ughwb2xpY3lpZBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cmNlYXJu');

@$core.Deprecated('Use putResourcePolicyResponseDescriptor instead')
const PutResourcePolicyResponse$json = {
  '1': 'PutResourcePolicyResponse',
  '2': [
    {'1': 'policyhash', '3': 203037360, '4': 1, '5': 9, '10': 'policyhash'},
    {'1': 'policyid', '3': 299520499, '4': 1, '5': 9, '10': 'policyid'},
  ],
};

/// Descriptor for `PutResourcePolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyResponseDescriptor =
    $convert.base64Decode(
        'ChlQdXRSZXNvdXJjZVBvbGljeVJlc3BvbnNlEiEKCnBvbGljeWhhc2gYsLXoYCABKAlSCnBvbG'
        'ljeWhhc2gSHgoIcG9saWN5aWQY86PpjgEgASgJUghwb2xpY3lpZA==');

@$core.Deprecated('Use registerDefaultPatchBaselineRequestDescriptor instead')
const RegisterDefaultPatchBaselineRequest$json = {
  '1': 'RegisterDefaultPatchBaselineRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `RegisterDefaultPatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerDefaultPatchBaselineRequestDescriptor =
    $convert.base64Decode(
        'CiNSZWdpc3RlckRlZmF1bHRQYXRjaEJhc2VsaW5lUmVxdWVzdBIhCgpiYXNlbGluZWlkGNDk2y'
        'ggASgJUgpiYXNlbGluZWlk');

@$core.Deprecated('Use registerDefaultPatchBaselineResultDescriptor instead')
const RegisterDefaultPatchBaselineResult$json = {
  '1': 'RegisterDefaultPatchBaselineResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
  ],
};

/// Descriptor for `RegisterDefaultPatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerDefaultPatchBaselineResultDescriptor =
    $convert.base64Decode(
        'CiJSZWdpc3RlckRlZmF1bHRQYXRjaEJhc2VsaW5lUmVzdWx0EiEKCmJhc2VsaW5laWQY0OTbKC'
        'ABKAlSCmJhc2VsaW5laWQ=');

@$core.Deprecated(
    'Use registerPatchBaselineForPatchGroupRequestDescriptor instead')
const RegisterPatchBaselineForPatchGroupRequest$json = {
  '1': 'RegisterPatchBaselineForPatchGroupRequest',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `RegisterPatchBaselineForPatchGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    registerPatchBaselineForPatchGroupRequestDescriptor = $convert.base64Decode(
        'CilSZWdpc3RlclBhdGNoQmFzZWxpbmVGb3JQYXRjaEdyb3VwUmVxdWVzdBIhCgpiYXNlbGluZW'
        'lkGNDk2yggASgJUgpiYXNlbGluZWlkEiIKCnBhdGNoZ3JvdXAY4bex9wEgASgJUgpwYXRjaGdy'
        'b3Vw');

@$core.Deprecated(
    'Use registerPatchBaselineForPatchGroupResultDescriptor instead')
const RegisterPatchBaselineForPatchGroupResult$json = {
  '1': 'RegisterPatchBaselineForPatchGroupResult',
  '2': [
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'patchgroup', '3': 518806497, '4': 1, '5': 9, '10': 'patchgroup'},
  ],
};

/// Descriptor for `RegisterPatchBaselineForPatchGroupResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerPatchBaselineForPatchGroupResultDescriptor =
    $convert.base64Decode(
        'CihSZWdpc3RlclBhdGNoQmFzZWxpbmVGb3JQYXRjaEdyb3VwUmVzdWx0EiEKCmJhc2VsaW5laW'
        'QY0OTbKCABKAlSCmJhc2VsaW5laWQSIgoKcGF0Y2hncm91cBjht7H3ASABKAlSCnBhdGNoZ3Jv'
        'dXA=');

@$core.Deprecated(
    'Use registerTargetWithMaintenanceWindowRequestDescriptor instead')
const RegisterTargetWithMaintenanceWindowRequest$json = {
  '1': 'RegisterTargetWithMaintenanceWindowRequest',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowResourceType',
      '10': 'resourcetype'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `RegisterTargetWithMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    registerTargetWithMaintenanceWindowRequestDescriptor =
    $convert.base64Decode(
        'CipSZWdpc3RlclRhcmdldFdpdGhNYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSIwoLY2xpZW50dG'
        '9rZW4YzPu7QSABKAlSC2NsaWVudHRva2VuEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNj'
        'cmlwdGlvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEi0KEG93bmVyaW5mb3JtYXRpb24Yg5rFIC'
        'ABKAlSEG93bmVyaW5mb3JtYXRpb24SSgoMcmVzb3VyY2V0eXBlGN6+2I8BIAEoDjIiLnNzbS5N'
        'YWludGVuYW5jZVdpbmRvd1Jlc291cmNlVHlwZVIMcmVzb3VyY2V0eXBlEigKB3RhcmdldHMYgp'
        'uCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5k'
        'b3dpZA==');

@$core.Deprecated(
    'Use registerTargetWithMaintenanceWindowResultDescriptor instead')
const RegisterTargetWithMaintenanceWindowResult$json = {
  '1': 'RegisterTargetWithMaintenanceWindowResult',
  '2': [
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `RegisterTargetWithMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    registerTargetWithMaintenanceWindowResultDescriptor = $convert.base64Decode(
        'CilSZWdpc3RlclRhcmdldFdpdGhNYWludGVuYW5jZVdpbmRvd1Jlc3VsdBIpCg53aW5kb3d0YX'
        'JnZXRpZBjezup7IAEoCVIOd2luZG93dGFyZ2V0aWQ=');

@$core.Deprecated(
    'Use registerTaskWithMaintenanceWindowRequestDescriptor instead')
const RegisterTaskWithMaintenanceWindowRequest$json = {
  '1': 'RegisterTaskWithMaintenanceWindowRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {
      '1': 'cutoffbehavior',
      '3': 120608587,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskCutoffBehavior',
      '10': 'cutoffbehavior'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'logginginfo',
      '3': 448312415,
      '4': 1,
      '5': 11,
      '6': '.ssm.LoggingInfo',
      '10': 'logginginfo'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskinvocationparameters',
      '3': 347127635,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskInvocationParameters',
      '10': 'taskinvocationparameters'
    },
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.RegisterTaskWithMaintenanceWindowRequest.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {
      '1': 'tasktype',
      '3': 370407909,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskType',
      '10': 'tasktype'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
  '3': [RegisterTaskWithMaintenanceWindowRequest_TaskparametersEntry$json],
};

@$core.Deprecated(
    'Use registerTaskWithMaintenanceWindowRequestDescriptor instead')
const RegisterTaskWithMaintenanceWindowRequest_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `RegisterTaskWithMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerTaskWithMaintenanceWindowRequestDescriptor = $convert.base64Decode(
    'CihSZWdpc3RlclRhc2tXaXRoTWFpbnRlbmFuY2VXaW5kb3dSZXF1ZXN0EkoKEmFsYXJtY29uZm'
    'lndXJhdGlvbhiJmbkhIAEoCzIXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmln'
    'dXJhdGlvbhIjCgtjbGllbnR0b2tlbhjM+7tBIAEoCVILY2xpZW50dG9rZW4SUwoOY3V0b2ZmYm'
    'VoYXZpb3IYy67BOSABKA4yKC5zc20uTWFpbnRlbmFuY2VXaW5kb3dUYXNrQ3V0b2ZmQmVoYXZp'
    'b3JSDmN1dG9mZmJlaGF2aW9yEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbh'
    'I2Cgtsb2dnaW5naW5mbxjf6OLVASABKAsyEC5zc20uTG9nZ2luZ0luZm9SC2xvZ2dpbmdpbmZv'
    'EikKDm1heGNvbmN1cnJlbmN5GP3Bjg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcn'
    'MYq8L1PSABKAlSCW1heGVycm9ycxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh0KCHByaW9yaXR5'
    'GKq+tjQgASgFUghwcmlvcml0eRIqCg5zZXJ2aWNlcm9sZWFybhiE49q2ASABKAlSDnNlcnZpY2'
    'Vyb2xlYXJuEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEhwKB3Rh'
    'c2thcm4Y5Mn6lAEgASgJUgd0YXNrYXJuEm4KGHRhc2tpbnZvY2F0aW9ucGFyYW1ldGVycxjT/s'
    'KlASABKAsyLi5zc20uTWFpbnRlbmFuY2VXaW5kb3dUYXNrSW52b2NhdGlvblBhcmFtZXRlcnNS'
    'GHRhc2tpbnZvY2F0aW9ucGFyYW1ldGVycxJtCg50YXNrcGFyYW1ldGVycxiBj+a3ASADKAsyQS'
    '5zc20uUmVnaXN0ZXJUYXNrV2l0aE1haW50ZW5hbmNlV2luZG93UmVxdWVzdC5UYXNrcGFyYW1l'
    'dGVyc0VudHJ5Ug50YXNrcGFyYW1ldGVycxI+Cgh0YXNrdHlwZRjl88+wASABKA4yHi5zc20uTW'
    'FpbnRlbmFuY2VXaW5kb3dUYXNrVHlwZVIIdGFza3R5cGUSHQoId2luZG93aWQYqeSHCSABKAlS'
    'CHdpbmRvd2lkGnUKE1Rhc2twYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSSAoFdm'
    'FsdWUYAiABKAsyMi5zc20uTWFpbnRlbmFuY2VXaW5kb3dUYXNrUGFyYW1ldGVyVmFsdWVFeHBy'
    'ZXNzaW9uUgV2YWx1ZToCOAE=');

@$core
    .Deprecated('Use registerTaskWithMaintenanceWindowResultDescriptor instead')
const RegisterTaskWithMaintenanceWindowResult$json = {
  '1': 'RegisterTaskWithMaintenanceWindowResult',
  '2': [
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
};

/// Descriptor for `RegisterTaskWithMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerTaskWithMaintenanceWindowResultDescriptor =
    $convert.base64Decode(
        'CidSZWdpc3RlclRhc2tXaXRoTWFpbnRlbmFuY2VXaW5kb3dSZXN1bHQSJgoMd2luZG93dGFza2'
        'lkGNqNsZoBIAEoCVIMd2luZG93dGFza2lk');

@$core.Deprecated('Use registrationMetadataItemDescriptor instead')
const RegistrationMetadataItem$json = {
  '1': 'RegistrationMetadataItem',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `RegistrationMetadataItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registrationMetadataItemDescriptor =
    $convert.base64Decode(
        'ChhSZWdpc3RyYXRpb25NZXRhZGF0YUl0ZW0SEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdW'
        'UY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use relatedOpsItemDescriptor instead')
const RelatedOpsItem$json = {
  '1': 'RelatedOpsItem',
  '2': [
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
  ],
};

/// Descriptor for `RelatedOpsItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List relatedOpsItemDescriptor = $convert.base64Decode(
    'Cg5SZWxhdGVkT3BzSXRlbRIfCglvcHNpdGVtaWQY0tKVDCABKAlSCW9wc2l0ZW1pZA==');

@$core.Deprecated('Use removeTagsFromResourceRequestDescriptor instead')
const RemoveTagsFromResourceRequest$json = {
  '1': 'RemoveTagsFromResourceRequest',
  '2': [
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceTypeForTagging',
      '10': 'resourcetype'
    },
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `RemoveTagsFromResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsFromResourceRequestDescriptor = $convert.base64Decode(
    'Ch1SZW1vdmVUYWdzRnJvbVJlc291cmNlUmVxdWVzdBIiCgpyZXNvdXJjZWlkGJG68foBIAEoCV'
    'IKcmVzb3VyY2VpZBJDCgxyZXNvdXJjZXR5cGUY3r7YjwEgASgOMhsuc3NtLlJlc291cmNlVHlw'
    'ZUZvclRhZ2dpbmdSDHJlc291cmNldHlwZRIcCgd0YWdrZXlzGPzD85gBIAMoCVIHdGFna2V5cw'
    '==');

@$core.Deprecated('Use removeTagsFromResourceResultDescriptor instead')
const RemoveTagsFromResourceResult$json = {
  '1': 'RemoveTagsFromResourceResult',
};

/// Descriptor for `RemoveTagsFromResourceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsFromResourceResultDescriptor =
    $convert.base64Decode('ChxSZW1vdmVUYWdzRnJvbVJlc291cmNlUmVzdWx0');

@$core.Deprecated('Use resetServiceSettingRequestDescriptor instead')
const ResetServiceSettingRequest$json = {
  '1': 'ResetServiceSettingRequest',
  '2': [
    {'1': 'settingid', '3': 422452711, '4': 1, '5': 9, '10': 'settingid'},
  ],
};

/// Descriptor for `ResetServiceSettingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetServiceSettingRequestDescriptor =
    $convert.base64Decode(
        'ChpSZXNldFNlcnZpY2VTZXR0aW5nUmVxdWVzdBIgCglzZXR0aW5naWQY57u4yQEgASgJUglzZX'
        'R0aW5naWQ=');

@$core.Deprecated('Use resetServiceSettingResultDescriptor instead')
const ResetServiceSettingResult$json = {
  '1': 'ResetServiceSettingResult',
  '2': [
    {
      '1': 'servicesetting',
      '3': 486565473,
      '4': 1,
      '5': 11,
      '6': '.ssm.ServiceSetting',
      '10': 'servicesetting'
    },
  ],
};

/// Descriptor for `ResetServiceSettingResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetServiceSettingResultDescriptor =
    $convert.base64Decode(
        'ChlSZXNldFNlcnZpY2VTZXR0aW5nUmVzdWx0Ej8KDnNlcnZpY2VzZXR0aW5nGOHMgegBIAEoCz'
        'ITLnNzbS5TZXJ2aWNlU2V0dGluZ1IOc2VydmljZXNldHRpbmc=');

@$core.Deprecated('Use resolvedTargetsDescriptor instead')
const ResolvedTargets$json = {
  '1': 'ResolvedTargets',
  '2': [
    {
      '1': 'parametervalues',
      '3': 210263223,
      '4': 3,
      '5': 9,
      '10': 'parametervalues'
    },
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ResolvedTargets`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resolvedTargetsDescriptor = $convert.base64Decode(
    'Cg9SZXNvbHZlZFRhcmdldHMSKwoPcGFyYW1ldGVydmFsdWVzGLe5oWQgAygJUg9wYXJhbWV0ZX'
    'J2YWx1ZXMSHwoJdHJ1bmNhdGVkGMrv2EggASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use resourceComplianceSummaryItemDescriptor instead')
const ResourceComplianceSummaryItem$json = {
  '1': 'ResourceComplianceSummaryItem',
  '2': [
    {
      '1': 'compliancetype',
      '3': 451385291,
      '4': 1,
      '5': 9,
      '10': 'compliancetype'
    },
    {
      '1': 'compliantsummary',
      '3': 133218055,
      '4': 1,
      '5': 11,
      '6': '.ssm.CompliantSummary',
      '10': 'compliantsummary'
    },
    {
      '1': 'executionsummary',
      '3': 71480964,
      '4': 1,
      '5': 11,
      '6': '.ssm.ComplianceExecutionSummary',
      '10': 'executionsummary'
    },
    {
      '1': 'noncompliantsummary',
      '3': 294594444,
      '4': 1,
      '5': 11,
      '6': '.ssm.NonCompliantSummary',
      '10': 'noncompliantsummary'
    },
    {
      '1': 'overallseverity',
      '3': 74875526,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceSeverity',
      '10': 'overallseverity'
    },
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ComplianceStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ResourceComplianceSummaryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceComplianceSummaryItemDescriptor = $convert.base64Decode(
    'Ch1SZXNvdXJjZUNvbXBsaWFuY2VTdW1tYXJ5SXRlbRIqCg5jb21wbGlhbmNldHlwZRjLr57XAS'
    'ABKAlSDmNvbXBsaWFuY2V0eXBlEkQKEGNvbXBsaWFudHN1bW1hcnkYh/7CPyABKAsyFS5zc20u'
    'Q29tcGxpYW50U3VtbWFyeVIQY29tcGxpYW50c3VtbWFyeRJOChBleGVjdXRpb25zdW1tYXJ5GI'
    'TtiiIgASgLMh8uc3NtLkNvbXBsaWFuY2VFeGVjdXRpb25TdW1tYXJ5UhBleGVjdXRpb25zdW1t'
    'YXJ5Ek4KE25vbmNvbXBsaWFudHN1bW1hcnkYjM+8jAEgASgLMhguc3NtLk5vbkNvbXBsaWFudF'
    'N1bW1hcnlSE25vbmNvbXBsaWFudHN1bW1hcnkSRAoPb3ZlcmFsbHNldmVyaXR5GIaF2iMgASgO'
    'Mhcuc3NtLkNvbXBsaWFuY2VTZXZlcml0eVIPb3ZlcmFsbHNldmVyaXR5EiIKCnJlc291cmNlaW'
    'QYkbrx+gEgASgJUgpyZXNvdXJjZWlkEiYKDHJlc291cmNldHlwZRjevtiPASABKAlSDHJlc291'
    'cmNldHlwZRIwCgZzdGF0dXMYkOT7AiABKA4yFS5zc20uQ29tcGxpYW5jZVN0YXR1c1IGc3RhdH'
    'Vz');

@$core
    .Deprecated('Use resourceDataSyncAlreadyExistsExceptionDescriptor instead')
const ResourceDataSyncAlreadyExistsException$json = {
  '1': 'ResourceDataSyncAlreadyExistsException',
  '2': [
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
  ],
};

/// Descriptor for `ResourceDataSyncAlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncAlreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'CiZSZXNvdXJjZURhdGFTeW5jQWxyZWFkeUV4aXN0c0V4Y2VwdGlvbhIeCghzeW5jbmFtZRiilr'
        'KwASABKAlSCHN5bmNuYW1l');

@$core
    .Deprecated('Use resourceDataSyncAwsOrganizationsSourceDescriptor instead')
const ResourceDataSyncAwsOrganizationsSource$json = {
  '1': 'ResourceDataSyncAwsOrganizationsSource',
  '2': [
    {
      '1': 'organizationsourcetype',
      '3': 60015472,
      '4': 1,
      '5': 9,
      '10': 'organizationsourcetype'
    },
    {
      '1': 'organizationalunits',
      '3': 138980163,
      '4': 3,
      '5': 11,
      '6': '.ssm.ResourceDataSyncOrganizationalUnit',
      '10': 'organizationalunits'
    },
  ],
};

/// Descriptor for `ResourceDataSyncAwsOrganizationsSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncAwsOrganizationsSourceDescriptor =
    $convert.base64Decode(
        'CiZSZXNvdXJjZURhdGFTeW5jQXdzT3JnYW5pemF0aW9uc1NvdXJjZRI5ChZvcmdhbml6YXRpb2'
        '5zb3VyY2V0eXBlGPCGzxwgASgJUhZvcmdhbml6YXRpb25zb3VyY2V0eXBlElwKE29yZ2FuaXph'
        'dGlvbmFsdW5pdHMYw9aiQiADKAsyJy5zc20uUmVzb3VyY2VEYXRhU3luY09yZ2FuaXphdGlvbm'
        'FsVW5pdFITb3JnYW5pemF0aW9uYWx1bml0cw==');

@$core.Deprecated('Use resourceDataSyncConflictExceptionDescriptor instead')
const ResourceDataSyncConflictException$json = {
  '1': 'ResourceDataSyncConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceDataSyncConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncConflictExceptionDescriptor =
    $convert.base64Decode(
        'CiFSZXNvdXJjZURhdGFTeW5jQ29uZmxpY3RFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZQ==');

@$core
    .Deprecated('Use resourceDataSyncCountExceededExceptionDescriptor instead')
const ResourceDataSyncCountExceededException$json = {
  '1': 'ResourceDataSyncCountExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceDataSyncCountExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncCountExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZSZXNvdXJjZURhdGFTeW5jQ291bnRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3'
        'AgASgJUgdtZXNzYWdl');

@$core
    .Deprecated('Use resourceDataSyncDestinationDataSharingDescriptor instead')
const ResourceDataSyncDestinationDataSharing$json = {
  '1': 'ResourceDataSyncDestinationDataSharing',
  '2': [
    {
      '1': 'destinationdatasharingtype',
      '3': 162313154,
      '4': 1,
      '5': 9,
      '10': 'destinationdatasharingtype'
    },
  ],
};

/// Descriptor for `ResourceDataSyncDestinationDataSharing`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncDestinationDataSharingDescriptor =
    $convert.base64Decode(
        'CiZSZXNvdXJjZURhdGFTeW5jRGVzdGluYXRpb25EYXRhU2hhcmluZxJBChpkZXN0aW5hdGlvbm'
        'RhdGFzaGFyaW5ndHlwZRjC57JNIAEoCVIaZGVzdGluYXRpb25kYXRhc2hhcmluZ3R5cGU=');

@$core.Deprecated(
    'Use resourceDataSyncInvalidConfigurationExceptionDescriptor instead')
const ResourceDataSyncInvalidConfigurationException$json = {
  '1': 'ResourceDataSyncInvalidConfigurationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceDataSyncInvalidConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    resourceDataSyncInvalidConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'Ci1SZXNvdXJjZURhdGFTeW5jSW52YWxpZENvbmZpZ3VyYXRpb25FeGNlcHRpb24SGwoHbWVzc2'
        'FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use resourceDataSyncItemDescriptor instead')
const ResourceDataSyncItem$json = {
  '1': 'ResourceDataSyncItem',
  '2': [
    {
      '1': 'laststatus',
      '3': 163326556,
      '4': 1,
      '5': 14,
      '6': '.ssm.LastResourceDataSyncStatus',
      '10': 'laststatus'
    },
    {
      '1': 'lastsuccessfulsynctime',
      '3': 94157050,
      '4': 1,
      '5': 9,
      '10': 'lastsuccessfulsynctime'
    },
    {
      '1': 'lastsyncstatusmessage',
      '3': 70064880,
      '4': 1,
      '5': 9,
      '10': 'lastsyncstatusmessage'
    },
    {'1': 'lastsynctime', '3': 125896272, '4': 1, '5': 9, '10': 'lastsynctime'},
    {
      '1': 's3destination',
      '3': 526234522,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncS3Destination',
      '10': 's3destination'
    },
    {
      '1': 'synccreatedtime',
      '3': 5143768,
      '4': 1,
      '5': 9,
      '10': 'synccreatedtime'
    },
    {
      '1': 'synclastmodifiedtime',
      '3': 176500261,
      '4': 1,
      '5': 9,
      '10': 'synclastmodifiedtime'
    },
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
    {
      '1': 'syncsource',
      '3': 286486824,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncSourceWithState',
      '10': 'syncsource'
    },
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `ResourceDataSyncItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncItemDescriptor = $convert.base64Decode(
    'ChRSZXNvdXJjZURhdGFTeW5jSXRlbRJCCgpsYXN0c3RhdHVzGNzU8E0gASgOMh8uc3NtLkxhc3'
    'RSZXNvdXJjZURhdGFTeW5jU3RhdHVzUgpsYXN0c3RhdHVzEjkKFmxhc3RzdWNjZXNzZnVsc3lu'
    'Y3RpbWUY+vHyLCABKAlSFmxhc3RzdWNjZXNzZnVsc3luY3RpbWUSNwoVbGFzdHN5bmNzdGF0dX'
    'NtZXNzYWdlGPC1tCEgASgJUhVsYXN0c3luY3N0YXR1c21lc3NhZ2USJQoMbGFzdHN5bmN0aW1l'
    'GNCMhDwgASgJUgxsYXN0c3luY3RpbWUSTAoNczNkZXN0aW5hdGlvbhia5/b6ASABKAsyIi5zc2'
    '0uUmVzb3VyY2VEYXRhU3luY1MzRGVzdGluYXRpb25SDXMzZGVzdGluYXRpb24SKwoPc3luY2Ny'
    'ZWF0ZWR0aW1lGNj5uQIgASgJUg9zeW5jY3JlYXRlZHRpbWUSNQoUc3luY2xhc3Rtb2RpZmllZH'
    'RpbWUYpdyUVCABKAlSFHN5bmNsYXN0bW9kaWZpZWR0aW1lEh4KCHN5bmNuYW1lGKKWsrABIAEo'
    'CVIIc3luY25hbWUSSAoKc3luY3NvdXJjZRio4s2IASABKAsyJC5zc20uUmVzb3VyY2VEYXRhU3'
    'luY1NvdXJjZVdpdGhTdGF0ZVIKc3luY3NvdXJjZRIdCghzeW5jdHlwZRjb5qo6IAEoCVIIc3lu'
    'Y3R5cGU=');

@$core.Deprecated('Use resourceDataSyncNotFoundExceptionDescriptor instead')
const ResourceDataSyncNotFoundException$json = {
  '1': 'ResourceDataSyncNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `ResourceDataSyncNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiFSZXNvdXJjZURhdGFTeW5jTm90Rm91bmRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCV'
        'IHbWVzc2FnZRIeCghzeW5jbmFtZRiilrKwASABKAlSCHN5bmNuYW1lEh0KCHN5bmN0eXBlGNvm'
        'qjogASgJUghzeW5jdHlwZQ==');

@$core.Deprecated('Use resourceDataSyncOrganizationalUnitDescriptor instead')
const ResourceDataSyncOrganizationalUnit$json = {
  '1': 'ResourceDataSyncOrganizationalUnit',
  '2': [
    {
      '1': 'organizationalunitid',
      '3': 201174977,
      '4': 1,
      '5': 9,
      '10': 'organizationalunitid'
    },
  ],
};

/// Descriptor for `ResourceDataSyncOrganizationalUnit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncOrganizationalUnitDescriptor =
    $convert.base64Decode(
        'CiJSZXNvdXJjZURhdGFTeW5jT3JnYW5pemF0aW9uYWxVbml0EjUKFG9yZ2FuaXphdGlvbmFsdW'
        '5pdGlkGMHf9l8gASgJUhRvcmdhbml6YXRpb25hbHVuaXRpZA==');

@$core.Deprecated('Use resourceDataSyncS3DestinationDescriptor instead')
const ResourceDataSyncS3Destination$json = {
  '1': 'ResourceDataSyncS3Destination',
  '2': [
    {'1': 'awskmskeyarn', '3': 533835690, '4': 1, '5': 9, '10': 'awskmskeyarn'},
    {'1': 'bucketname', '3': 208117045, '4': 1, '5': 9, '10': 'bucketname'},
    {
      '1': 'destinationdatasharing',
      '3': 448514634,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncDestinationDataSharing',
      '10': 'destinationdatasharing'
    },
    {'1': 'prefix', '3': 273996266, '4': 1, '5': 9, '10': 'prefix'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
    {
      '1': 'syncformat',
      '3': 373344378,
      '4': 1,
      '5': 14,
      '6': '.ssm.ResourceDataSyncS3Format',
      '10': 'syncformat'
    },
  ],
};

/// Descriptor for `ResourceDataSyncS3Destination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncS3DestinationDescriptor = $convert.base64Decode(
    'Ch1SZXNvdXJjZURhdGFTeW5jUzNEZXN0aW5hdGlvbhImCgxhd3NrbXNrZXlhcm4Yqt/G/gEgAS'
    'gJUgxhd3NrbXNrZXlhcm4SIQoKYnVja2V0bmFtZRi1up5jIAEoCVIKYnVja2V0bmFtZRJnChZk'
    'ZXN0aW5hdGlvbmRhdGFzaGFyaW5nGMqU79UBIAEoCzIrLnNzbS5SZXNvdXJjZURhdGFTeW5jRG'
    'VzdGluYXRpb25EYXRhU2hhcmluZ1IWZGVzdGluYXRpb25kYXRhc2hhcmluZxIaCgZwcmVmaXgY'
    '6rPTggEgASgJUgZwcmVmaXgSGQoGcmVnaW9uGJ7xuUkgASgJUgZyZWdpb24SQQoKc3luY2Zvcm'
    '1hdBj6kIOyASABKA4yHS5zc20uUmVzb3VyY2VEYXRhU3luY1MzRm9ybWF0UgpzeW5jZm9ybWF0');

@$core.Deprecated('Use resourceDataSyncSourceDescriptor instead')
const ResourceDataSyncSource$json = {
  '1': 'ResourceDataSyncSource',
  '2': [
    {
      '1': 'awsorganizationssource',
      '3': 116902612,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncAwsOrganizationsSource',
      '10': 'awsorganizationssource'
    },
    {
      '1': 'enableallopsdatasources',
      '3': 466051416,
      '4': 1,
      '5': 8,
      '10': 'enableallopsdatasources'
    },
    {
      '1': 'includefutureregions',
      '3': 72068934,
      '4': 1,
      '5': 8,
      '10': 'includefutureregions'
    },
    {
      '1': 'sourceregions',
      '3': 115078446,
      '4': 3,
      '5': 9,
      '10': 'sourceregions'
    },
    {'1': 'sourcetype', '3': 195731217, '4': 1, '5': 9, '10': 'sourcetype'},
  ],
};

/// Descriptor for `ResourceDataSyncSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncSourceDescriptor = $convert.base64Decode(
    'ChZSZXNvdXJjZURhdGFTeW5jU291cmNlEmYKFmF3c29yZ2FuaXphdGlvbnNzb3VyY2UY1JXfNy'
    'ABKAsyKy5zc20uUmVzb3VyY2VEYXRhU3luY0F3c09yZ2FuaXphdGlvbnNTb3VyY2VSFmF3c29y'
    'Z2FuaXphdGlvbnNzb3VyY2USPAoXZW5hYmxlYWxsb3BzZGF0YXNvdXJjZXMY2MKd3gEgASgIUh'
    'dlbmFibGVhbGxvcHNkYXRhc291cmNlcxI1ChRpbmNsdWRlZnV0dXJlcmVnaW9ucxjG3q4iIAEo'
    'CFIUaW5jbHVkZWZ1dHVyZXJlZ2lvbnMSJwoNc291cmNlcmVnaW9ucxiu6u82IAMoCVINc291cm'
    'NlcmVnaW9ucxIhCgpzb3VyY2V0eXBlGJG+ql0gASgJUgpzb3VyY2V0eXBl');

@$core.Deprecated('Use resourceDataSyncSourceWithStateDescriptor instead')
const ResourceDataSyncSourceWithState$json = {
  '1': 'ResourceDataSyncSourceWithState',
  '2': [
    {
      '1': 'awsorganizationssource',
      '3': 116902612,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncAwsOrganizationsSource',
      '10': 'awsorganizationssource'
    },
    {
      '1': 'enableallopsdatasources',
      '3': 466051416,
      '4': 1,
      '5': 8,
      '10': 'enableallopsdatasources'
    },
    {
      '1': 'includefutureregions',
      '3': 72068934,
      '4': 1,
      '5': 8,
      '10': 'includefutureregions'
    },
    {
      '1': 'sourceregions',
      '3': 115078446,
      '4': 3,
      '5': 9,
      '10': 'sourceregions'
    },
    {'1': 'sourcetype', '3': 195731217, '4': 1, '5': 9, '10': 'sourcetype'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `ResourceDataSyncSourceWithState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceDataSyncSourceWithStateDescriptor = $convert.base64Decode(
    'Ch9SZXNvdXJjZURhdGFTeW5jU291cmNlV2l0aFN0YXRlEmYKFmF3c29yZ2FuaXphdGlvbnNzb3'
    'VyY2UY1JXfNyABKAsyKy5zc20uUmVzb3VyY2VEYXRhU3luY0F3c09yZ2FuaXphdGlvbnNTb3Vy'
    'Y2VSFmF3c29yZ2FuaXphdGlvbnNzb3VyY2USPAoXZW5hYmxlYWxsb3BzZGF0YXNvdXJjZXMY2M'
    'Kd3gEgASgIUhdlbmFibGVhbGxvcHNkYXRhc291cmNlcxI1ChRpbmNsdWRlZnV0dXJlcmVnaW9u'
    'cxjG3q4iIAEoCFIUaW5jbHVkZWZ1dHVyZXJlZ2lvbnMSJwoNc291cmNlcmVnaW9ucxiu6u82IA'
    'MoCVINc291cmNlcmVnaW9ucxIhCgpzb3VyY2V0eXBlGJG+ql0gASgJUgpzb3VyY2V0eXBlEhgK'
    'BXN0YXRlGJfJsu8BIAEoCVIFc3RhdGU=');

@$core.Deprecated('Use resourceInUseExceptionDescriptor instead')
const ResourceInUseException$json = {
  '1': 'ResourceInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChZSZXNvdXJjZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceLimitExceededExceptionDescriptor instead')
const ResourceLimitExceededException$json = {
  '1': 'ResourceLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch5SZXNvdXJjZUxpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbW'
        'Vzc2FnZQ==');

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

@$core.Deprecated('Use resourcePolicyConflictExceptionDescriptor instead')
const ResourcePolicyConflictException$json = {
  '1': 'ResourcePolicyConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourcePolicyConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyConflictExceptionDescriptor =
    $convert.base64Decode(
        'Ch9SZXNvdXJjZVBvbGljeUNvbmZsaWN0RXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core
    .Deprecated('Use resourcePolicyInvalidParameterExceptionDescriptor instead')
const ResourcePolicyInvalidParameterException$json = {
  '1': 'ResourcePolicyInvalidParameterException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'parameternames',
      '3': 72522785,
      '4': 3,
      '5': 9,
      '10': 'parameternames'
    },
  ],
};

/// Descriptor for `ResourcePolicyInvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyInvalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'CidSZXNvdXJjZVBvbGljeUludmFsaWRQYXJhbWV0ZXJFeGNlcHRpb24SGwoHbWVzc2FnZRiFs7'
        'twIAEoCVIHbWVzc2FnZRIpCg5wYXJhbWV0ZXJuYW1lcxihuMoiIAMoCVIOcGFyYW1ldGVybmFt'
        'ZXM=');

@$core.Deprecated('Use resourcePolicyLimitExceededExceptionDescriptor instead')
const ResourcePolicyLimitExceededException$json = {
  '1': 'ResourcePolicyLimitExceededException',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'limittype', '3': 155546973, '4': 1, '5': 9, '10': 'limittype'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourcePolicyLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiRSZXNvdXJjZVBvbGljeUxpbWl0RXhjZWVkZWRFeGNlcHRpb24SGAoFbGltaXQY1ZXZxAEgAS'
        'gFUgVsaW1pdBIfCglsaW1pdHR5cGUY3eqVSiABKAlSCWxpbWl0dHlwZRIbCgdtZXNzYWdlGIWz'
        'u3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use resourcePolicyNotFoundExceptionDescriptor instead')
const ResourcePolicyNotFoundException$json = {
  '1': 'ResourcePolicyNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourcePolicyNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourcePolicyNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch9SZXNvdXJjZVBvbGljeU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use resultAttributeDescriptor instead')
const ResultAttribute$json = {
  '1': 'ResultAttribute',
  '2': [
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `ResultAttribute`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resultAttributeDescriptor = $convert.base64Decode(
    'Cg9SZXN1bHRBdHRyaWJ1dGUSHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ==');

@$core.Deprecated('Use resumeSessionRequestDescriptor instead')
const ResumeSessionRequest$json = {
  '1': 'ResumeSessionRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `ResumeSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resumeSessionRequestDescriptor = $convert.base64Decode(
    'ChRSZXN1bWVTZXNzaW9uUmVxdWVzdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZA'
    '==');

@$core.Deprecated('Use resumeSessionResponseDescriptor instead')
const ResumeSessionResponse$json = {
  '1': 'ResumeSessionResponse',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {'1': 'streamurl', '3': 423623039, '4': 1, '5': 9, '10': 'streamurl'},
    {'1': 'tokenvalue', '3': 97232092, '4': 1, '5': 9, '10': 'tokenvalue'},
  ],
};

/// Descriptor for `ResumeSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resumeSessionResponseDescriptor = $convert.base64Decode(
    'ChVSZXN1bWVTZXNzaW9uUmVzcG9uc2USHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW9uaW'
    'QSIAoJc3RyZWFtdXJsGP/y/8kBIAEoCVIJc3RyZWFtdXJsEiEKCnRva2VudmFsdWUY3MmuLiAB'
    'KAlSCnRva2VudmFsdWU=');

@$core.Deprecated('Use reviewInformationDescriptor instead')
const ReviewInformation$json = {
  '1': 'ReviewInformation',
  '2': [
    {'1': 'reviewedtime', '3': 277736338, '4': 1, '5': 9, '10': 'reviewedtime'},
    {'1': 'reviewer', '3': 436444219, '4': 1, '5': 9, '10': 'reviewer'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.ReviewStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ReviewInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reviewInformationDescriptor = $convert.base64Decode(
    'ChFSZXZpZXdJbmZvcm1hdGlvbhImCgxyZXZpZXdlZHRpbWUYkte3hAEgASgJUgxyZXZpZXdlZH'
    'RpbWUSHgoIcmV2aWV3ZXIYu7iO0AEgASgJUghyZXZpZXdlchIsCgZzdGF0dXMYkOT7AiABKA4y'
    'ES5zc20uUmV2aWV3U3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use runbookDescriptor instead')
const Runbook$json = {
  '1': 'Runbook',
  '2': [
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.Runbook.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.Runbook.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targetparametername',
      '3': 351056597,
      '4': 1,
      '5': 9,
      '10': 'targetparametername'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [Runbook_ParametersEntry$json, Runbook_TargetmapsEntry$json],
};

@$core.Deprecated('Use runbookDescriptor instead')
const Runbook_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use runbookDescriptor instead')
const Runbook_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Runbook`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List runbookDescriptor = $convert.base64Decode(
    'CgdSdW5ib29rEiUKDGRvY3VtZW50bmFtZRjQo8c5IAEoCVIMZG9jdW1lbnRuYW1lEisKD2RvY3'
    'VtZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEikKDm1heGNvbmN1cnJlbmN5'
    'GP3Bjg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcnMYq8L1PSABKAlSCW1heGVycm'
    '9ycxJACgpwYXJhbWV0ZXJzGPqn/usBIAMoCzIcLnNzbS5SdW5ib29rLlBhcmFtZXRlcnNFbnRy'
    'eVIKcGFyYW1ldGVycxJBCg90YXJnZXRsb2NhdGlvbnMYpbvxiQEgAygLMhMuc3NtLlRhcmdldE'
    'xvY2F0aW9uUg90YXJnZXRsb2NhdGlvbnMSPwoKdGFyZ2V0bWFwcxi4vNUjIAMoCzIcLnNzbS5S'
    'dW5ib29rLlRhcmdldG1hcHNFbnRyeVIKdGFyZ2V0bWFwcxI0ChN0YXJnZXRwYXJhbWV0ZXJuYW'
    '1lGNXlsqcBIAEoCVITdGFyZ2V0cGFyYW1ldGVybmFtZRIoCgd0YXJnZXRzGIKbgn0gAygLMgsu'
    'c3NtLlRhcmdldFIHdGFyZ2V0cxo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZX'
    'kSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARo9Cg9UYXJnZXRtYXBzRW50cnkSEAoDa2V5GAEg'
    'ASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use s3OutputLocationDescriptor instead')
const S3OutputLocation$json = {
  '1': 'S3OutputLocation',
  '2': [
    {
      '1': 'outputs3bucketname',
      '3': 186756480,
      '4': 1,
      '5': 9,
      '10': 'outputs3bucketname'
    },
    {
      '1': 'outputs3keyprefix',
      '3': 17840974,
      '4': 1,
      '5': 9,
      '10': 'outputs3keyprefix'
    },
    {
      '1': 'outputs3region',
      '3': 398718159,
      '4': 1,
      '5': 9,
      '10': 'outputs3region'
    },
  ],
};

/// Descriptor for `S3OutputLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3OutputLocationDescriptor = $convert.base64Decode(
    'ChBTM091dHB1dExvY2F0aW9uEjEKEm91dHB1dHMzYnVja2V0bmFtZRiA24ZZIAEoCVISb3V0cH'
    'V0czNidWNrZXRuYW1lEi8KEW91dHB1dHMza2V5cHJlZml4GM72wAggASgJUhFvdXRwdXRzM2tl'
    'eXByZWZpeBIqCg5vdXRwdXRzM3JlZ2lvbhjP6Y++ASABKAlSDm91dHB1dHMzcmVnaW9u');

@$core.Deprecated('Use s3OutputUrlDescriptor instead')
const S3OutputUrl$json = {
  '1': 'S3OutputUrl',
  '2': [
    {'1': 'outputurl', '3': 42495998, '4': 1, '5': 9, '10': 'outputurl'},
  ],
};

/// Descriptor for `S3OutputUrl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List s3OutputUrlDescriptor = $convert.base64Decode(
    'CgtTM091dHB1dFVybBIfCglvdXRwdXR1cmwY/t+hFCABKAlSCW91dHB1dHVybA==');

@$core.Deprecated('Use scheduledWindowExecutionDescriptor instead')
const ScheduledWindowExecution$json = {
  '1': 'ScheduledWindowExecution',
  '2': [
    {
      '1': 'executiontime',
      '3': 379716053,
      '4': 1,
      '5': 9,
      '10': 'executiontime'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `ScheduledWindowExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduledWindowExecutionDescriptor = $convert.base64Decode(
    'ChhTY2hlZHVsZWRXaW5kb3dFeGVjdXRpb24SKAoNZXhlY3V0aW9udGltZRjVg4i1ASABKAlSDW'
    'V4ZWN1dGlvbnRpbWUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIdCgh3aW5kb3dpZBip5IcJIAEo'
    'CVIId2luZG93aWQ=');

@$core.Deprecated('Use sendAutomationSignalRequestDescriptor instead')
const SendAutomationSignalRequest$json = {
  '1': 'SendAutomationSignalRequest',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
    {
      '1': 'payload',
      '3': 6526790,
      '4': 3,
      '5': 11,
      '6': '.ssm.SendAutomationSignalRequest.PayloadEntry',
      '10': 'payload'
    },
    {
      '1': 'signaltype',
      '3': 418520162,
      '4': 1,
      '5': 14,
      '6': '.ssm.SignalType',
      '10': 'signaltype'
    },
  ],
  '3': [SendAutomationSignalRequest_PayloadEntry$json],
};

@$core.Deprecated('Use sendAutomationSignalRequestDescriptor instead')
const SendAutomationSignalRequest_PayloadEntry$json = {
  '1': 'PayloadEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SendAutomationSignalRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendAutomationSignalRequestDescriptor = $convert.base64Decode(
    'ChtTZW5kQXV0b21hdGlvblNpZ25hbFJlcXVlc3QSNwoVYXV0b21hdGlvbmV4ZWN1dGlvbmlkGO'
    'bv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQSSgoHcGF5bG9hZBjGro4DIAMoCzItLnNz'
    'bS5TZW5kQXV0b21hdGlvblNpZ25hbFJlcXVlc3QuUGF5bG9hZEVudHJ5UgdwYXlsb2FkEjMKCn'
    'NpZ25hbHR5cGUY4rjIxwEgASgOMg8uc3NtLlNpZ25hbFR5cGVSCnNpZ25hbHR5cGUaOgoMUGF5'
    'bG9hZEVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use sendAutomationSignalResultDescriptor instead')
const SendAutomationSignalResult$json = {
  '1': 'SendAutomationSignalResult',
};

/// Descriptor for `SendAutomationSignalResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendAutomationSignalResultDescriptor =
    $convert.base64Decode('ChpTZW5kQXV0b21hdGlvblNpZ25hbFJlc3VsdA==');

@$core.Deprecated('Use sendCommandRequestDescriptor instead')
const SendCommandRequest$json = {
  '1': 'SendCommandRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cloudwatchoutputconfig',
      '3': 21186555,
      '4': 1,
      '5': 11,
      '6': '.ssm.CloudWatchOutputConfig',
      '10': 'cloudwatchoutputconfig'
    },
    {'1': 'comment', '3': 407871487, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'documenthash', '3': 251248469, '4': 1, '5': 9, '10': 'documenthash'},
    {
      '1': 'documenthashtype',
      '3': 93041117,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentHashType',
      '10': 'documenthashtype'
    },
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'instanceids', '3': 312792453, '4': 3, '5': 9, '10': 'instanceids'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'notificationconfig',
      '3': 346074145,
      '4': 1,
      '5': 11,
      '6': '.ssm.NotificationConfig',
      '10': 'notificationconfig'
    },
    {
      '1': 'outputs3bucketname',
      '3': 186756480,
      '4': 1,
      '5': 9,
      '10': 'outputs3bucketname'
    },
    {
      '1': 'outputs3keyprefix',
      '3': 17840974,
      '4': 1,
      '5': 9,
      '10': 'outputs3keyprefix'
    },
    {
      '1': 'outputs3region',
      '3': 398718159,
      '4': 1,
      '5': 9,
      '10': 'outputs3region'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.SendCommandRequest.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'timeoutseconds',
      '3': 336148022,
      '4': 1,
      '5': 5,
      '10': 'timeoutseconds'
    },
  ],
  '3': [SendCommandRequest_ParametersEntry$json],
};

@$core.Deprecated('Use sendCommandRequestDescriptor instead')
const SendCommandRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SendCommandRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendCommandRequestDescriptor = $convert.base64Decode(
    'ChJTZW5kQ29tbWFuZFJlcXVlc3QSSgoSYWxhcm1jb25maWd1cmF0aW9uGImZuSEgASgLMhcuc3'
    'NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uElYKFmNsb3Vkd2F0Y2hv'
    'dXRwdXRjb25maWcY+4+NCiABKAsyGy5zc20uQ2xvdWRXYXRjaE91dHB1dENvbmZpZ1IWY2xvdW'
    'R3YXRjaG91dHB1dGNvbmZpZxIcCgdjb21tZW50GP+/vsIBIAEoCVIHY29tbWVudBIlCgxkb2N1'
    'bWVudGhhc2gY1f7mdyABKAlSDGRvY3VtZW50aGFzaBJEChBkb2N1bWVudGhhc2h0eXBlGN3jri'
    'wgASgOMhUuc3NtLkRvY3VtZW50SGFzaFR5cGVSEGRvY3VtZW50aGFzaHR5cGUSJQoMZG9jdW1l'
    'bnRuYW1lGNCjxzkgASgJUgxkb2N1bWVudG5hbWUSKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggAS'
    'gJUg9kb2N1bWVudHZlcnNpb24SJAoLaW5zdGFuY2VpZHMYhauTlQEgAygJUgtpbnN0YW5jZWlk'
    'cxIpCg5tYXhjb25jdXJyZW5jeRj9wY4OIAEoCVIObWF4Y29uY3VycmVuY3kSHwoJbWF4ZXJyb3'
    'JzGKvC9T0gASgJUgltYXhlcnJvcnMSSwoSbm90aWZpY2F0aW9uY29uZmlnGKHYgqUBIAEoCzIX'
    'LnNzbS5Ob3RpZmljYXRpb25Db25maWdSEm5vdGlmaWNhdGlvbmNvbmZpZxIxChJvdXRwdXRzM2'
    'J1Y2tldG5hbWUYgNuGWSABKAlSEm91dHB1dHMzYnVja2V0bmFtZRIvChFvdXRwdXRzM2tleXBy'
    'ZWZpeBjO9sAIIAEoCVIRb3V0cHV0czNrZXlwcmVmaXgSKgoOb3V0cHV0czNyZWdpb24Yz+mPvg'
    'EgASgJUg5vdXRwdXRzM3JlZ2lvbhJLCgpwYXJhbWV0ZXJzGPqn/usBIAMoCzInLnNzbS5TZW5k'
    'Q29tbWFuZFJlcXVlc3QuUGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEioKDnNlcnZpY2Vyb2'
    'xlYXJuGITj2rYBIAEoCVIOc2VydmljZXJvbGVhcm4SKAoHdGFyZ2V0cxiCm4J9IAMoCzILLnNz'
    'bS5UYXJnZXRSB3RhcmdldHMSKgoOdGltZW91dHNlY29uZHMYtuykoAEgASgFUg50aW1lb3V0c2'
    'Vjb25kcxo9Cg9QYXJhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiAB'
    'KAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use sendCommandResultDescriptor instead')
const SendCommandResult$json = {
  '1': 'SendCommandResult',
  '2': [
    {
      '1': 'command',
      '3': 108826451,
      '4': 1,
      '5': 11,
      '6': '.ssm.Command',
      '10': 'command'
    },
  ],
};

/// Descriptor for `SendCommandResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendCommandResultDescriptor = $convert.base64Decode(
    'ChFTZW5kQ29tbWFuZFJlc3VsdBIpCgdjb21tYW5kGNOe8jMgASgLMgwuc3NtLkNvbW1hbmRSB2'
    'NvbW1hbmQ=');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'quotacode', '3': 306488915, '4': 1, '5': 9, '10': 'quotacode'},
    {'1': 'resourceid', '3': 526146833, '4': 1, '5': 9, '10': 'resourceid'},
    {'1': 'resourcetype', '3': 301342558, '4': 1, '5': 9, '10': 'resourcetype'},
    {'1': 'servicecode', '3': 474897394, '4': 1, '5': 9, '10': 'servicecode'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor = $convert.base64Decode(
    'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
    'NzYWdlEiAKCXF1b3RhY29kZRjTzJKSASABKAlSCXF1b3RhY29kZRIiCgpyZXNvdXJjZWlkGJG6'
    '8foBIAEoCVIKcmVzb3VyY2VpZBImCgxyZXNvdXJjZXR5cGUY3r7YjwEgASgJUgxyZXNvdXJjZX'
    'R5cGUSJAoLc2VydmljZWNvZGUY8re54gEgASgJUgtzZXJ2aWNlY29kZQ==');

@$core.Deprecated('Use serviceSettingDescriptor instead')
const ServiceSetting$json = {
  '1': 'ServiceSetting',
  '2': [
    {'1': 'arn', '3': 397135389, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {
      '1': 'lastmodifieduser',
      '3': 215822778,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieduser'
    },
    {'1': 'settingid', '3': 422452711, '4': 1, '5': 9, '10': 'settingid'},
    {'1': 'settingvalue', '3': 467016041, '4': 1, '5': 9, '10': 'settingvalue'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `ServiceSetting`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceSettingDescriptor = $convert.base64Decode(
    'Cg5TZXJ2aWNlU2V0dGluZxIUCgNhcm4YnZyvvQEgASgJUgNhcm4SLQoQbGFzdG1vZGlmaWVkZG'
    'F0ZRjTiMgLIAEoCVIQbGFzdG1vZGlmaWVkZGF0ZRItChBsYXN0bW9kaWZpZWR1c2VyGLrj9GYg'
    'ASgJUhBsYXN0bW9kaWZpZWR1c2VyEiAKCXNldHRpbmdpZBjnu7jJASABKAlSCXNldHRpbmdpZB'
    'ImCgxzZXR0aW5ndmFsdWUY6bLY3gEgASgJUgxzZXR0aW5ndmFsdWUSGQoGc3RhdHVzGJDk+wIg'
    'ASgJUgZzdGF0dXM=');

@$core.Deprecated('Use serviceSettingNotFoundDescriptor instead')
const ServiceSettingNotFound$json = {
  '1': 'ServiceSettingNotFound',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceSettingNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceSettingNotFoundDescriptor =
    $convert.base64Decode(
        'ChZTZXJ2aWNlU2V0dGluZ05vdEZvdW5kEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use sessionDescriptor instead')
const Session$json = {
  '1': 'Session',
  '2': [
    {
      '1': 'accesstype',
      '3': 14019508,
      '4': 1,
      '5': 14,
      '6': '.ssm.AccessType',
      '10': 'accesstype'
    },
    {'1': 'details', '3': 247611974, '4': 1, '5': 9, '10': 'details'},
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'maxsessionduration',
      '3': 391414912,
      '4': 1,
      '5': 9,
      '10': 'maxsessionduration'
    },
    {
      '1': 'outputurl',
      '3': 42495998,
      '4': 1,
      '5': 11,
      '6': '.ssm.SessionManagerOutputUrl',
      '10': 'outputurl'
    },
    {'1': 'owner', '3': 455261813, '4': 1, '5': 9, '10': 'owner'},
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.SessionStatus',
      '10': 'status'
    },
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
  ],
};

/// Descriptor for `Session`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionDescriptor = $convert.base64Decode(
    'CgdTZXNzaW9uEjIKCmFjY2Vzc3R5cGUYtNfXBiABKA4yDy5zc20uQWNjZXNzVHlwZVIKYWNjZX'
    'NzdHlwZRIbCgdkZXRhaWxzGMaEiXYgASgJUgdkZXRhaWxzEiUKDGRvY3VtZW50bmFtZRjQo8c5'
    'IAEoCVIMZG9jdW1lbnRuYW1lEhsKB2VuZGRhdGUYz7P5JCABKAlSB2VuZGRhdGUSMgoSbWF4c2'
    'Vzc2lvbmR1cmF0aW9uGICJ0roBIAEoCVISbWF4c2Vzc2lvbmR1cmF0aW9uEj0KCW91dHB1dHVy'
    'bBj+36EUIAEoCzIcLnNzbS5TZXNzaW9uTWFuYWdlck91dHB1dFVybFIJb3V0cHV0dXJsEhgKBW'
    '93bmVyGPX8itkBIAEoCVIFb3duZXISGQoGcmVhc29uGLqCxQkgASgJUgZyZWFzb24SHwoJc2Vz'
    'c2lvbmlkGLuE5QkgASgJUglzZXNzaW9uaWQSIAoJc3RhcnRkYXRlGPz4oNQBIAEoCVIJc3Rhcn'
    'RkYXRlEi0KBnN0YXR1cxiQ5PsCIAEoDjISLnNzbS5TZXNzaW9uU3RhdHVzUgZzdGF0dXMSGQoG'
    'dGFyZ2V0GOnin1sgASgJUgZ0YXJnZXQ=');

@$core.Deprecated('Use sessionFilterDescriptor instead')
const SessionFilter$json = {
  '1': 'SessionFilter',
  '2': [
    {
      '1': 'key',
      '3': 135645293,
      '4': 1,
      '5': 14,
      '6': '.ssm.SessionFilterKey',
      '10': 'key'
    },
    {'1': 'value', '3': 39769035, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `SessionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionFilterDescriptor = $convert.base64Decode(
    'Cg1TZXNzaW9uRmlsdGVyEioKA2tleRjtkNdAIAEoDjIVLnNzbS5TZXNzaW9uRmlsdGVyS2V5Ug'
    'NrZXkSFwoFdmFsdWUYy6f7EiABKAlSBXZhbHVl');

@$core.Deprecated('Use sessionManagerOutputUrlDescriptor instead')
const SessionManagerOutputUrl$json = {
  '1': 'SessionManagerOutputUrl',
  '2': [
    {
      '1': 'cloudwatchoutputurl',
      '3': 60509222,
      '4': 1,
      '5': 9,
      '10': 'cloudwatchoutputurl'
    },
    {'1': 's3outputurl', '3': 490515856, '4': 1, '5': 9, '10': 's3outputurl'},
  ],
};

/// Descriptor for `SessionManagerOutputUrl`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionManagerOutputUrlDescriptor = $convert.base64Decode(
    'ChdTZXNzaW9uTWFuYWdlck91dHB1dFVybBIzChNjbG91ZHdhdGNob3V0cHV0dXJsGKaY7RwgAS'
    'gJUhNjbG91ZHdhdGNob3V0cHV0dXJsEiQKC3Mzb3V0cHV0dXJsGJDb8ukBIAEoCVILczNvdXRw'
    'dXR1cmw=');

@$core.Deprecated('Use severitySummaryDescriptor instead')
const SeveritySummary$json = {
  '1': 'SeveritySummary',
  '2': [
    {
      '1': 'criticalcount',
      '3': 113055624,
      '4': 1,
      '5': 5,
      '10': 'criticalcount'
    },
    {'1': 'highcount', '3': 243040701, '4': 1, '5': 5, '10': 'highcount'},
    {
      '1': 'informationalcount',
      '3': 58973984,
      '4': 1,
      '5': 5,
      '10': 'informationalcount'
    },
    {'1': 'lowcount', '3': 57008121, '4': 1, '5': 5, '10': 'lowcount'},
    {'1': 'mediumcount', '3': 362761286, '4': 1, '5': 5, '10': 'mediumcount'},
    {
      '1': 'unspecifiedcount',
      '3': 405758324,
      '4': 1,
      '5': 5,
      '10': 'unspecifiedcount'
    },
  ],
};

/// Descriptor for `SeveritySummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List severitySummaryDescriptor = $convert.base64Decode(
    'Cg9TZXZlcml0eVN1bW1hcnkSJwoNY3JpdGljYWxjb3VudBiIr/Q1IAEoBVINY3JpdGljYWxjb3'
    'VudBIfCgloaWdoY291bnQYvYPycyABKAVSCWhpZ2hjb3VudBIxChJpbmZvcm1hdGlvbmFsY291'
    'bnQYoL6PHCABKAVSEmluZm9ybWF0aW9uYWxjb3VudBIdCghsb3djb3VudBj5v5cbIAEoBVIIbG'
    '93Y291bnQSJAoLbWVkaXVtY291bnQYxpj9rAEgASgFUgttZWRpdW1jb3VudBIuChB1bnNwZWNp'
    'ZmllZGNvdW50GPTCvcEBIAEoBVIQdW5zcGVjaWZpZWRjb3VudA==');

@$core.Deprecated('Use startAccessRequestRequestDescriptor instead')
const StartAccessRequestRequest$json = {
  '1': 'StartAccessRequestRequest',
  '2': [
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
};

/// Descriptor for `StartAccessRequestRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAccessRequestRequestDescriptor = $convert.base64Decode(
    'ChlTdGFydEFjY2Vzc1JlcXVlc3RSZXF1ZXN0EhkKBnJlYXNvbhi6gsUJIAEoCVIGcmVhc29uEi'
    'AKBHRhZ3MYwcH2tQEgAygLMgguc3NtLlRhZ1IEdGFncxIoCgd0YXJnZXRzGIKbgn0gAygLMgsu'
    'c3NtLlRhcmdldFIHdGFyZ2V0cw==');

@$core.Deprecated('Use startAccessRequestResponseDescriptor instead')
const StartAccessRequestResponse$json = {
  '1': 'StartAccessRequestResponse',
  '2': [
    {
      '1': 'accessrequestid',
      '3': 377097210,
      '4': 1,
      '5': 9,
      '10': 'accessrequestid'
    },
  ],
};

/// Descriptor for `StartAccessRequestResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAccessRequestResponseDescriptor =
    $convert.base64Decode(
        'ChpTdGFydEFjY2Vzc1JlcXVlc3RSZXNwb25zZRIsCg9hY2Nlc3NyZXF1ZXN0aWQY+pfoswEgAS'
        'gJUg9hY2Nlc3NyZXF1ZXN0aWQ=');

@$core.Deprecated('Use startAssociationsOnceRequestDescriptor instead')
const StartAssociationsOnceRequest$json = {
  '1': 'StartAssociationsOnceRequest',
  '2': [
    {
      '1': 'associationids',
      '3': 124122183,
      '4': 3,
      '5': 9,
      '10': 'associationids'
    },
  ],
};

/// Descriptor for `StartAssociationsOnceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAssociationsOnceRequestDescriptor =
    $convert.base64Decode(
        'ChxTdGFydEFzc29jaWF0aW9uc09uY2VSZXF1ZXN0EikKDmFzc29jaWF0aW9uaWRzGMfolzsgAy'
        'gJUg5hc3NvY2lhdGlvbmlkcw==');

@$core.Deprecated('Use startAssociationsOnceResultDescriptor instead')
const StartAssociationsOnceResult$json = {
  '1': 'StartAssociationsOnceResult',
};

/// Descriptor for `StartAssociationsOnceResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAssociationsOnceResultDescriptor =
    $convert.base64Decode('ChtTdGFydEFzc29jaWF0aW9uc09uY2VSZXN1bHQ=');

@$core.Deprecated('Use startAutomationExecutionRequestDescriptor instead')
const StartAutomationExecutionRequest$json = {
  '1': 'StartAutomationExecutionRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {
      '1': 'mode',
      '3': 323909427,
      '4': 1,
      '5': 14,
      '6': '.ssm.ExecutionMode',
      '10': 'mode'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.StartAutomationExecutionRequest.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetlocationsurl',
      '3': 107583422,
      '4': 1,
      '5': 9,
      '10': 'targetlocationsurl'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.StartAutomationExecutionRequest.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targetparametername',
      '3': 351056597,
      '4': 1,
      '5': 9,
      '10': 'targetparametername'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    StartAutomationExecutionRequest_ParametersEntry$json,
    StartAutomationExecutionRequest_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use startAutomationExecutionRequestDescriptor instead')
const StartAutomationExecutionRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use startAutomationExecutionRequestDescriptor instead')
const StartAutomationExecutionRequest_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StartAutomationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAutomationExecutionRequestDescriptor = $convert.base64Decode(
    'Ch9TdGFydEF1dG9tYXRpb25FeGVjdXRpb25SZXF1ZXN0EkoKEmFsYXJtY29uZmlndXJhdGlvbh'
    'iJmbkhIAEoCzIXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlvbhIj'
    'CgtjbGllbnR0b2tlbhjM+7tBIAEoCVILY2xpZW50dG9rZW4SJQoMZG9jdW1lbnRuYW1lGNCjxz'
    'kgASgJUgxkb2N1bWVudG5hbWUSKwoPZG9jdW1lbnR2ZXJzaW9uGMnvqSggASgJUg9kb2N1bWVu'
    'dHZlcnNpb24SKQoObWF4Y29uY3VycmVuY3kY/cGODiABKAlSDm1heGNvbmN1cnJlbmN5Eh8KCW'
    '1heGVycm9ycxirwvU9IAEoCVIJbWF4ZXJyb3JzEioKBG1vZGUYs+65mgEgASgOMhIuc3NtLkV4'
    'ZWN1dGlvbk1vZGVSBG1vZGUSWAoKcGFyYW1ldGVycxj6p/7rASADKAsyNC5zc20uU3RhcnRBdX'
    'RvbWF0aW9uRXhlY3V0aW9uUmVxdWVzdC5QYXJhbWV0ZXJzRW50cnlSCnBhcmFtZXRlcnMSIAoE'
    'dGFncxjBwfa1ASADKAsyCC5zc20uVGFnUgR0YWdzEkEKD3RhcmdldGxvY2F0aW9ucxilu/GJAS'
    'ADKAsyEy5zc20uVGFyZ2V0TG9jYXRpb25SD3RhcmdldGxvY2F0aW9ucxIxChJ0YXJnZXRsb2Nh'
    'dGlvbnN1cmwYvq+mMyABKAlSEnRhcmdldGxvY2F0aW9uc3VybBJXCgp0YXJnZXRtYXBzGLi81S'
    'MgAygLMjQuc3NtLlN0YXJ0QXV0b21hdGlvbkV4ZWN1dGlvblJlcXVlc3QuVGFyZ2V0bWFwc0Vu'
    'dHJ5Ugp0YXJnZXRtYXBzEjQKE3RhcmdldHBhcmFtZXRlcm5hbWUY1eWypwEgASgJUhN0YXJnZX'
    'RwYXJhbWV0ZXJuYW1lEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRz'
    'Gj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdm'
    'FsdWU6AjgBGj0KD1RhcmdldG1hcHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgC'
    'IAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use startAutomationExecutionResultDescriptor instead')
const StartAutomationExecutionResult$json = {
  '1': 'StartAutomationExecutionResult',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
  ],
};

/// Descriptor for `StartAutomationExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startAutomationExecutionResultDescriptor =
    $convert.base64Decode(
        'Ch5TdGFydEF1dG9tYXRpb25FeGVjdXRpb25SZXN1bHQSNwoVYXV0b21hdGlvbmV4ZWN1dGlvbm'
        'lkGObv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQ=');

@$core.Deprecated('Use startChangeRequestExecutionRequestDescriptor instead')
const StartChangeRequestExecutionRequest$json = {
  '1': 'StartChangeRequestExecutionRequest',
  '2': [
    {'1': 'autoapprove', '3': 461595104, '4': 1, '5': 8, '10': 'autoapprove'},
    {
      '1': 'changedetails',
      '3': 125139630,
      '4': 1,
      '5': 9,
      '10': 'changedetails'
    },
    {
      '1': 'changerequestname',
      '3': 468779362,
      '4': 1,
      '5': 9,
      '10': 'changerequestname'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.StartChangeRequestExecutionRequest.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'runbooks',
      '3': 514418725,
      '4': 3,
      '5': 11,
      '6': '.ssm.Runbook',
      '10': 'runbooks'
    },
    {
      '1': 'scheduledendtime',
      '3': 3954269,
      '4': 1,
      '5': 9,
      '10': 'scheduledendtime'
    },
    {
      '1': 'scheduledtime',
      '3': 334708242,
      '4': 1,
      '5': 9,
      '10': 'scheduledtime'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.ssm.Tag',
      '10': 'tags'
    },
  ],
  '3': [StartChangeRequestExecutionRequest_ParametersEntry$json],
};

@$core.Deprecated('Use startChangeRequestExecutionRequestDescriptor instead')
const StartChangeRequestExecutionRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StartChangeRequestExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startChangeRequestExecutionRequestDescriptor = $convert.base64Decode(
    'CiJTdGFydENoYW5nZVJlcXVlc3RFeGVjdXRpb25SZXF1ZXN0EiQKC2F1dG9hcHByb3ZlGODDjd'
    'wBIAEoCFILYXV0b2FwcHJvdmUSJwoNY2hhbmdlZGV0YWlscxiu9dU7IAEoCVINY2hhbmdlZGV0'
    'YWlscxIwChFjaGFuZ2VyZXF1ZXN0bmFtZRjigsTfASABKAlSEWNoYW5nZXJlcXVlc3RuYW1lEi'
    'MKC2NsaWVudHRva2VuGMz7u0EgASgJUgtjbGllbnR0b2tlbhIlCgxkb2N1bWVudG5hbWUY0KPH'
    'OSABKAlSDGRvY3VtZW50bmFtZRIrCg9kb2N1bWVudHZlcnNpb24Yye+pKCABKAlSD2RvY3VtZW'
    '50dmVyc2lvbhJbCgpwYXJhbWV0ZXJzGPqn/usBIAMoCzI3LnNzbS5TdGFydENoYW5nZVJlcXVl'
    'c3RFeGVjdXRpb25SZXF1ZXN0LlBhcmFtZXRlcnNFbnRyeVIKcGFyYW1ldGVycxIsCghydW5ib2'
    '9rcxil0KX1ASADKAsyDC5zc20uUnVuYm9va1IIcnVuYm9va3MSLQoQc2NoZWR1bGVkZW5kdGlt'
    'ZRjdrPEBIAEoCVIQc2NoZWR1bGVkZW5kdGltZRIoCg1zY2hlZHVsZWR0aW1lGJL8zJ8BIAEoCV'
    'INc2NoZWR1bGVkdGltZRIgCgR0YWdzGMHB9rUBIAMoCzIILnNzbS5UYWdSBHRhZ3MaPQoPUGFy'
    'YW1ldGVyc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOA'
    'E=');

@$core.Deprecated('Use startChangeRequestExecutionResultDescriptor instead')
const StartChangeRequestExecutionResult$json = {
  '1': 'StartChangeRequestExecutionResult',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
  ],
};

/// Descriptor for `StartChangeRequestExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startChangeRequestExecutionResultDescriptor =
    $convert.base64Decode(
        'CiFTdGFydENoYW5nZVJlcXVlc3RFeGVjdXRpb25SZXN1bHQSNwoVYXV0b21hdGlvbmV4ZWN1dG'
        'lvbmlkGObv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQ=');

@$core.Deprecated('Use startExecutionPreviewRequestDescriptor instead')
const StartExecutionPreviewRequest$json = {
  '1': 'StartExecutionPreviewRequest',
  '2': [
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {
      '1': 'executioninputs',
      '3': 133189225,
      '4': 1,
      '5': 11,
      '6': '.ssm.ExecutionInputs',
      '10': 'executioninputs'
    },
  ],
};

/// Descriptor for `StartExecutionPreviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startExecutionPreviewRequestDescriptor = $convert.base64Decode(
    'ChxTdGFydEV4ZWN1dGlvblByZXZpZXdSZXF1ZXN0EiUKDGRvY3VtZW50bmFtZRjQo8c5IAEoCV'
    'IMZG9jdW1lbnRuYW1lEisKD2RvY3VtZW50dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJz'
    'aW9uEkEKD2V4ZWN1dGlvbmlucHV0cxjpnME/IAEoCzIULnNzbS5FeGVjdXRpb25JbnB1dHNSD2'
    'V4ZWN1dGlvbmlucHV0cw==');

@$core.Deprecated('Use startExecutionPreviewResponseDescriptor instead')
const StartExecutionPreviewResponse$json = {
  '1': 'StartExecutionPreviewResponse',
  '2': [
    {
      '1': 'executionpreviewid',
      '3': 35285163,
      '4': 1,
      '5': 9,
      '10': 'executionpreviewid'
    },
  ],
};

/// Descriptor for `StartExecutionPreviewResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startExecutionPreviewResponseDescriptor =
    $convert.base64Decode(
        'Ch1TdGFydEV4ZWN1dGlvblByZXZpZXdSZXNwb25zZRIxChJleGVjdXRpb25wcmV2aWV3aWQYq9'
        'HpECABKAlSEmV4ZWN1dGlvbnByZXZpZXdpZA==');

@$core.Deprecated('Use startSessionRequestDescriptor instead')
const StartSessionRequest$json = {
  '1': 'StartSessionRequest',
  '2': [
    {'1': 'documentname', '3': 120705488, '4': 1, '5': 9, '10': 'documentname'},
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.StartSessionRequest.ParametersEntry',
      '10': 'parameters'
    },
    {'1': 'reason', '3': 20005178, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'target', '3': 191361385, '4': 1, '5': 9, '10': 'target'},
  ],
  '3': [StartSessionRequest_ParametersEntry$json],
};

@$core.Deprecated('Use startSessionRequestDescriptor instead')
const StartSessionRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StartSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionRequestDescriptor = $convert.base64Decode(
    'ChNTdGFydFNlc3Npb25SZXF1ZXN0EiUKDGRvY3VtZW50bmFtZRjQo8c5IAEoCVIMZG9jdW1lbn'
    'RuYW1lEkwKCnBhcmFtZXRlcnMY+qf+6wEgAygLMiguc3NtLlN0YXJ0U2Vzc2lvblJlcXVlc3Qu'
    'UGFyYW1ldGVyc0VudHJ5UgpwYXJhbWV0ZXJzEhkKBnJlYXNvbhi6gsUJIAEoCVIGcmVhc29uEh'
    'kKBnRhcmdldBjp4p9bIAEoCVIGdGFyZ2V0Gj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASAB'
    'KAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use startSessionResponseDescriptor instead')
const StartSessionResponse$json = {
  '1': 'StartSessionResponse',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
    {'1': 'streamurl', '3': 423623039, '4': 1, '5': 9, '10': 'streamurl'},
    {'1': 'tokenvalue', '3': 97232092, '4': 1, '5': 9, '10': 'tokenvalue'},
  ],
};

/// Descriptor for `StartSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionResponseDescriptor = $convert.base64Decode(
    'ChRTdGFydFNlc3Npb25SZXNwb25zZRIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb25pZB'
    'IgCglzdHJlYW11cmwY//L/yQEgASgJUglzdHJlYW11cmwSIQoKdG9rZW52YWx1ZRjcya4uIAEo'
    'CVIKdG9rZW52YWx1ZQ==');

@$core.Deprecated('Use statusUnchangedDescriptor instead')
const StatusUnchanged$json = {
  '1': 'StatusUnchanged',
};

/// Descriptor for `StatusUnchanged`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusUnchangedDescriptor =
    $convert.base64Decode('Cg9TdGF0dXNVbmNoYW5nZWQ=');

@$core.Deprecated('Use stepExecutionDescriptor instead')
const StepExecution$json = {
  '1': 'StepExecution',
  '2': [
    {'1': 'action', '3': 175614240, '4': 1, '5': 9, '10': 'action'},
    {
      '1': 'executionendtime',
      '3': 139859196,
      '4': 1,
      '5': 9,
      '10': 'executionendtime'
    },
    {
      '1': 'executionstarttime',
      '3': 429847391,
      '4': 1,
      '5': 9,
      '10': 'executionstarttime'
    },
    {
      '1': 'failuredetails',
      '3': 409582698,
      '4': 1,
      '5': 11,
      '6': '.ssm.FailureDetails',
      '10': 'failuredetails'
    },
    {
      '1': 'failuremessage',
      '3': 353556937,
      '4': 1,
      '5': 9,
      '10': 'failuremessage'
    },
    {
      '1': 'inputs',
      '3': 499898041,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecution.InputsEntry',
      '10': 'inputs'
    },
    {'1': 'iscritical', '3': 250154393, '4': 1, '5': 8, '10': 'iscritical'},
    {'1': 'isend', '3': 276967899, '4': 1, '5': 8, '10': 'isend'},
    {'1': 'maxattempts', '3': 233732372, '4': 1, '5': 5, '10': 'maxattempts'},
    {'1': 'nextstep', '3': 346442571, '4': 1, '5': 9, '10': 'nextstep'},
    {'1': 'onfailure', '3': 424696739, '4': 1, '5': 9, '10': 'onfailure'},
    {
      '1': 'outputs',
      '3': 455868918,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecution.OutputsEntry',
      '10': 'outputs'
    },
    {
      '1': 'overriddenparameters',
      '3': 423965710,
      '4': 3,
      '5': 11,
      '6': '.ssm.StepExecution.OverriddenparametersEntry',
      '10': 'overriddenparameters'
    },
    {
      '1': 'parentstepdetails',
      '3': 428116240,
      '4': 1,
      '5': 11,
      '6': '.ssm.ParentStepDetails',
      '10': 'parentstepdetails'
    },
    {'1': 'response', '3': 363430655, '4': 1, '5': 9, '10': 'response'},
    {'1': 'responsecode', '3': 447553700, '4': 1, '5': 9, '10': 'responsecode'},
    {
      '1': 'stepexecutionid',
      '3': 47252075,
      '4': 1,
      '5': 9,
      '10': 'stepexecutionid'
    },
    {'1': 'stepname', '3': 394700557, '4': 1, '5': 9, '10': 'stepname'},
    {
      '1': 'stepstatus',
      '3': 49235738,
      '4': 1,
      '5': 14,
      '6': '.ssm.AutomationExecutionStatus',
      '10': 'stepstatus'
    },
    {
      '1': 'targetlocation',
      '3': 475808320,
      '4': 1,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocation'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'timeoutseconds',
      '3': 336148022,
      '4': 1,
      '5': 3,
      '10': 'timeoutseconds'
    },
    {
      '1': 'triggeredalarms',
      '3': 263222917,
      '4': 3,
      '5': 11,
      '6': '.ssm.AlarmStateInformation',
      '10': 'triggeredalarms'
    },
    {
      '1': 'validnextsteps',
      '3': 403306152,
      '4': 3,
      '5': 9,
      '10': 'validnextsteps'
    },
  ],
  '3': [
    StepExecution_InputsEntry$json,
    StepExecution_OutputsEntry$json,
    StepExecution_OverriddenparametersEntry$json
  ],
};

@$core.Deprecated('Use stepExecutionDescriptor instead')
const StepExecution_InputsEntry$json = {
  '1': 'InputsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use stepExecutionDescriptor instead')
const StepExecution_OutputsEntry$json = {
  '1': 'OutputsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use stepExecutionDescriptor instead')
const StepExecution_OverriddenparametersEntry$json = {
  '1': 'OverriddenparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `StepExecution`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepExecutionDescriptor = $convert.base64Decode(
    'Cg1TdGVwRXhlY3V0aW9uEhkKBmFjdGlvbhig0t5TIAEoCVIGYWN0aW9uEi0KEGV4ZWN1dGlvbm'
    'VuZHRpbWUY/KnYQiABKAlSEGV4ZWN1dGlvbmVuZHRpbWUSMgoSZXhlY3V0aW9uc3RhcnR0aW1l'
    'GN/m+8wBIAEoCVISZXhlY3V0aW9uc3RhcnR0aW1lEj8KDmZhaWx1cmVkZXRhaWxzGOr4psMBIA'
    'EoCzITLnNzbS5GYWlsdXJlRGV0YWlsc1IOZmFpbHVyZWRldGFpbHMSKgoOZmFpbHVyZW1lc3Nh'
    'Z2UYybPLqAEgASgJUg5mYWlsdXJlbWVzc2FnZRI6CgZpbnB1dHMYua2v7gEgAygLMh4uc3NtLl'
    'N0ZXBFeGVjdXRpb24uSW5wdXRzRW50cnlSBmlucHV0cxIhCgppc2NyaXRpY2FsGJmbpHcgASgI'
    'Ugppc2NyaXRpY2FsEhgKBWlzZW5kGNvjiIQBIAEoCFIFaXNlbmQSIwoLbWF4YXR0ZW1wdHMYlP'
    'K5byABKAVSC21heGF0dGVtcHRzEh4KCG5leHRzdGVwGMuWmaUBIAEoCVIIbmV4dHN0ZXASIAoJ'
    'b25mYWlsdXJlGKO3wcoBIAEoCVIJb25mYWlsdXJlEj0KB291dHB1dHMY9oOw2QEgAygLMh8uc3'
    'NtLlN0ZXBFeGVjdXRpb24uT3V0cHV0c0VudHJ5UgdvdXRwdXRzEmQKFG92ZXJyaWRkZW5wYXJh'
    'bWV0ZXJzGI7olMoBIAMoCzIsLnNzbS5TdGVwRXhlY3V0aW9uLk92ZXJyaWRkZW5wYXJhbWV0ZX'
    'JzRW50cnlSFG92ZXJyaWRkZW5wYXJhbWV0ZXJzEkgKEXBhcmVudHN0ZXBkZXRhaWxzGJCSkswB'
    'IAEoCzIWLnNzbS5QYXJlbnRTdGVwRGV0YWlsc1IRcGFyZW50c3RlcGRldGFpbHMSHgoIcmVzcG'
    '9uc2UY/4WmrQEgASgJUghyZXNwb25zZRImCgxyZXNwb25zZWNvZGUYpMG01QEgASgJUgxyZXNw'
    'b25zZWNvZGUSKwoPc3RlcGV4ZWN1dGlvbmlkGOuExBYgASgJUg9zdGVwZXhlY3V0aW9uaWQSHg'
    'oIc3RlcG5hbWUYjc6avAEgASgJUghzdGVwbmFtZRJBCgpzdGVwc3RhdHVzGJqOvRcgASgOMh4u'
    'c3NtLkF1dG9tYXRpb25FeGVjdXRpb25TdGF0dXNSCnN0ZXBzdGF0dXMSPwoOdGFyZ2V0bG9jYX'
    'Rpb24YwITx4gEgASgLMhMuc3NtLlRhcmdldExvY2F0aW9uUg50YXJnZXRsb2NhdGlvbhIoCgd0'
    'YXJnZXRzGIKbgn0gAygLMgsuc3NtLlRhcmdldFIHdGFyZ2V0cxIqCg50aW1lb3V0c2Vjb25kcx'
    'i27KSgASABKANSDnRpbWVvdXRzZWNvbmRzEkcKD3RyaWdnZXJlZGFsYXJtcxiF7cF9IAMoCzIa'
    'LnNzbS5BbGFybVN0YXRlSW5mb3JtYXRpb25SD3RyaWdnZXJlZGFsYXJtcxIqCg52YWxpZG5leH'
    'RzdGVwcxio7afAASADKAlSDnZhbGlkbmV4dHN0ZXBzGjkKC0lucHV0c0VudHJ5EhAKA2tleRgB'
    'IAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaOgoMT3V0cHV0c0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaRwoZT3ZlcnJpZGRlbnBh'
    'cmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6Aj'
    'gB');

@$core.Deprecated('Use stepExecutionFilterDescriptor instead')
const StepExecutionFilter$json = {
  '1': 'StepExecutionFilter',
  '2': [
    {
      '1': 'key',
      '3': 219859213,
      '4': 1,
      '5': 14,
      '6': '.ssm.StepExecutionFilterKey',
      '10': 'key'
    },
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `StepExecutionFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stepExecutionFilterDescriptor = $convert.base64Decode(
    'ChNTdGVwRXhlY3V0aW9uRmlsdGVyEjAKA2tleRiNkutoIAEoDjIbLnNzbS5TdGVwRXhlY3V0aW'
    '9uRmlsdGVyS2V5UgNrZXkSGQoGdmFsdWVzGNzEtGogAygJUgZ2YWx1ZXM=');

@$core.Deprecated('Use stopAutomationExecutionRequestDescriptor instead')
const StopAutomationExecutionRequest$json = {
  '1': 'StopAutomationExecutionRequest',
  '2': [
    {
      '1': 'automationexecutionid',
      '3': 12449766,
      '4': 1,
      '5': 9,
      '10': 'automationexecutionid'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.ssm.StopType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `StopAutomationExecutionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopAutomationExecutionRequestDescriptor =
    $convert.base64Decode(
        'Ch5TdG9wQXV0b21hdGlvbkV4ZWN1dGlvblJlcXVlc3QSNwoVYXV0b21hdGlvbmV4ZWN1dGlvbm'
        'lkGObv9wUgASgJUhVhdXRvbWF0aW9uZXhlY3V0aW9uaWQSJQoEdHlwZRjuoNeKASABKA4yDS5z'
        'c20uU3RvcFR5cGVSBHR5cGU=');

@$core.Deprecated('Use stopAutomationExecutionResultDescriptor instead')
const StopAutomationExecutionResult$json = {
  '1': 'StopAutomationExecutionResult',
};

/// Descriptor for `StopAutomationExecutionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopAutomationExecutionResultDescriptor =
    $convert.base64Decode('Ch1TdG9wQXV0b21hdGlvbkV4ZWN1dGlvblJlc3VsdA==');

@$core.Deprecated('Use subTypeCountLimitExceededExceptionDescriptor instead')
const SubTypeCountLimitExceededException$json = {
  '1': 'SubTypeCountLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SubTypeCountLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subTypeCountLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiJTdWJUeXBlQ291bnRMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKA'
        'lSB21lc3NhZ2U=');

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

@$core.Deprecated('Use targetDescriptor instead')
const Target$json = {
  '1': 'Target',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 9, '10': 'values'},
  ],
};

/// Descriptor for `Target`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetDescriptor = $convert.base64Decode(
    'CgZUYXJnZXQSEwoDa2V5GI2S62ggASgJUgNrZXkSGQoGdmFsdWVzGNzEtGogAygJUgZ2YWx1ZX'
    'M=');

@$core.Deprecated('Use targetInUseExceptionDescriptor instead')
const TargetInUseException$json = {
  '1': 'TargetInUseException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TargetInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChRUYXJnZXRJblVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use targetLocationDescriptor instead')
const TargetLocation$json = {
  '1': 'TargetLocation',
  '2': [
    {'1': 'accounts', '3': 292022046, '4': 3, '5': 9, '10': 'accounts'},
    {
      '1': 'excludeaccounts',
      '3': 166973682,
      '4': 3,
      '5': 9,
      '10': 'excludeaccounts'
    },
    {
      '1': 'executionrolename',
      '3': 402462979,
      '4': 1,
      '5': 9,
      '10': 'executionrolename'
    },
    {
      '1': 'includechildorganizationunits',
      '3': 323711808,
      '4': 1,
      '5': 8,
      '10': 'includechildorganizationunits'
    },
    {'1': 'regions', '3': 36200107, '4': 3, '5': 9, '10': 'regions'},
    {
      '1': 'targetlocationalarmconfiguration',
      '3': 359268383,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'targetlocationalarmconfiguration'
    },
    {
      '1': 'targetlocationmaxconcurrency',
      '3': 202765591,
      '4': 1,
      '5': 9,
      '10': 'targetlocationmaxconcurrency'
    },
    {
      '1': 'targetlocationmaxerrors',
      '3': 212923133,
      '4': 1,
      '5': 9,
      '10': 'targetlocationmaxerrors'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {
      '1': 'targetsmaxconcurrency',
      '3': 242641881,
      '4': 1,
      '5': 9,
      '10': 'targetsmaxconcurrency'
    },
    {
      '1': 'targetsmaxerrors',
      '3': 20182343,
      '4': 1,
      '5': 9,
      '10': 'targetsmaxerrors'
    },
  ],
};

/// Descriptor for `TargetLocation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetLocationDescriptor = $convert.base64Decode(
    'Cg5UYXJnZXRMb2NhdGlvbhIeCghhY2NvdW50cxiezp+LASADKAlSCGFjY291bnRzEisKD2V4Y2'
    'x1ZGVhY2NvdW50cxjyoc9PIAMoCVIPZXhjbHVkZWFjY291bnRzEjAKEWV4ZWN1dGlvbnJvbGVu'
    'YW1lGIOy9L8BIAEoCVIRZXhlY3V0aW9ucm9sZW5hbWUSSAodaW5jbHVkZWNoaWxkb3JnYW5pem'
    'F0aW9udW5pdHMYwOatmgEgASgIUh1pbmNsdWRlY2hpbGRvcmdhbml6YXRpb251bml0cxIbCgdy'
    'ZWdpb25zGKu9oREgAygJUgdyZWdpb25zEmcKIHRhcmdldGxvY2F0aW9uYWxhcm1jb25maWd1cm'
    'F0aW9uGJ+AqKsBIAEoCzIXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SIHRhcmdldGxvY2F0aW9u'
    'YWxhcm1jb25maWd1cmF0aW9uEkUKHHRhcmdldGxvY2F0aW9ubWF4Y29uY3VycmVuY3kYl+rXYC'
    'ABKAlSHHRhcmdldGxvY2F0aW9ubWF4Y29uY3VycmVuY3kSOwoXdGFyZ2V0bG9jYXRpb25tYXhl'
    'cnJvcnMY/eXDZSABKAlSF3RhcmdldGxvY2F0aW9ubWF4ZXJyb3JzEigKB3RhcmdldHMYgpuCfS'
    'ADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEjcKFXRhcmdldHNtYXhjb25jdXJyZW5jeRjZ19lz'
    'IAEoCVIVdGFyZ2V0c21heGNvbmN1cnJlbmN5Ei0KEHRhcmdldHNtYXhlcnJvcnMYx+rPCSABKA'
    'lSEHRhcmdldHNtYXhlcnJvcnM=');

@$core.Deprecated('Use targetNotConnectedDescriptor instead')
const TargetNotConnected$json = {
  '1': 'TargetNotConnected',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TargetNotConnected`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetNotConnectedDescriptor =
    $convert.base64Decode(
        'ChJUYXJnZXROb3RDb25uZWN0ZWQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use targetPreviewDescriptor instead')
const TargetPreview$json = {
  '1': 'TargetPreview',
  '2': [
    {'1': 'count', '3': 31963285, '4': 1, '5': 5, '10': 'count'},
    {'1': 'targettype', '3': 397256481, '4': 1, '5': 9, '10': 'targettype'},
  ],
};

/// Descriptor for `TargetPreview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetPreviewDescriptor = $convert.base64Decode(
    'Cg1UYXJnZXRQcmV2aWV3EhcKBWNvdW50GJXxng8gASgFUgVjb3VudBIiCgp0YXJnZXR0eXBlGK'
    'HOtr0BIAEoCVIKdGFyZ2V0dHlwZQ==');

@$core.Deprecated('Use terminateSessionRequestDescriptor instead')
const TerminateSessionRequest$json = {
  '1': 'TerminateSessionRequest',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `TerminateSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List terminateSessionRequestDescriptor =
    $convert.base64Decode(
        'ChdUZXJtaW5hdGVTZXNzaW9uUmVxdWVzdBIfCglzZXNzaW9uaWQYu4TlCSABKAlSCXNlc3Npb2'
        '5pZA==');

@$core.Deprecated('Use terminateSessionResponseDescriptor instead')
const TerminateSessionResponse$json = {
  '1': 'TerminateSessionResponse',
  '2': [
    {'1': 'sessionid', '3': 20529723, '4': 1, '5': 9, '10': 'sessionid'},
  ],
};

/// Descriptor for `TerminateSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List terminateSessionResponseDescriptor =
    $convert.base64Decode(
        'ChhUZXJtaW5hdGVTZXNzaW9uUmVzcG9uc2USHwoJc2Vzc2lvbmlkGLuE5QkgASgJUglzZXNzaW'
        '9uaWQ=');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'quotacode', '3': 306488915, '4': 1, '5': 9, '10': 'quotacode'},
    {'1': 'servicecode', '3': 474897394, '4': 1, '5': 9, '10': 'servicecode'},
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor = $convert.base64Decode(
    'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2USIAoJcX'
    'VvdGFjb2RlGNPMkpIBIAEoCVIJcXVvdGFjb2RlEiQKC3NlcnZpY2Vjb2RlGPK3ueIBIAEoCVIL'
    'c2VydmljZWNvZGU=');

@$core.Deprecated('Use tooManyTagsErrorDescriptor instead')
const TooManyTagsError$json = {
  '1': 'TooManyTagsError',
};

/// Descriptor for `TooManyTagsError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyTagsErrorDescriptor =
    $convert.base64Decode('ChBUb29NYW55VGFnc0Vycm9y');

@$core.Deprecated('Use tooManyUpdatesDescriptor instead')
const TooManyUpdates$json = {
  '1': 'TooManyUpdates',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyUpdates`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyUpdatesDescriptor = $convert.base64Decode(
    'Cg5Ub29NYW55VXBkYXRlcxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use totalSizeLimitExceededExceptionDescriptor instead')
const TotalSizeLimitExceededException$json = {
  '1': 'TotalSizeLimitExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TotalSizeLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List totalSizeLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Ub3RhbFNpemVMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use unlabelParameterVersionRequestDescriptor instead')
const UnlabelParameterVersionRequest$json = {
  '1': 'UnlabelParameterVersionRequest',
  '2': [
    {'1': 'labels', '3': 178416811, '4': 3, '5': 9, '10': 'labels'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'parameterversion',
      '3': 141162077,
      '4': 1,
      '5': 3,
      '10': 'parameterversion'
    },
  ],
};

/// Descriptor for `UnlabelParameterVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlabelParameterVersionRequestDescriptor =
    $convert.base64Decode(
        'Ch5VbmxhYmVsUGFyYW1ldGVyVmVyc2lvblJlcXVlc3QSGQoGbGFiZWxzGKvZiVUgAygJUgZsYW'
        'JlbHMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRItChBwYXJhbWV0ZXJ2ZXJzaW9uGN3sp0MgASgD'
        'UhBwYXJhbWV0ZXJ2ZXJzaW9u');

@$core.Deprecated('Use unlabelParameterVersionResultDescriptor instead')
const UnlabelParameterVersionResult$json = {
  '1': 'UnlabelParameterVersionResult',
  '2': [
    {
      '1': 'invalidlabels',
      '3': 281240648,
      '4': 3,
      '5': 9,
      '10': 'invalidlabels'
    },
    {
      '1': 'removedlabels',
      '3': 242704737,
      '4': 3,
      '5': 9,
      '10': 'removedlabels'
    },
  ],
};

/// Descriptor for `UnlabelParameterVersionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlabelParameterVersionResultDescriptor =
    $convert.base64Decode(
        'Ch1VbmxhYmVsUGFyYW1ldGVyVmVyc2lvblJlc3VsdBIoCg1pbnZhbGlkbGFiZWxzGMjIjYYBIA'
        'MoCVINaW52YWxpZGxhYmVscxInCg1yZW1vdmVkbGFiZWxzGOHC3XMgAygJUg1yZW1vdmVkbGFi'
        'ZWxz');

@$core.Deprecated('Use unsupportedCalendarExceptionDescriptor instead')
const UnsupportedCalendarException$json = {
  '1': 'UnsupportedCalendarException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedCalendarException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedCalendarExceptionDescriptor =
    $convert.base64Decode(
        'ChxVbnN1cHBvcnRlZENhbGVuZGFyRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use unsupportedFeatureRequiredExceptionDescriptor instead')
const UnsupportedFeatureRequiredException$json = {
  '1': 'UnsupportedFeatureRequiredException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedFeatureRequiredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedFeatureRequiredExceptionDescriptor =
    $convert.base64Decode(
        'CiNVbnN1cHBvcnRlZEZlYXR1cmVSZXF1aXJlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated(
    'Use unsupportedInventoryItemContextExceptionDescriptor instead')
const UnsupportedInventoryItemContextException$json = {
  '1': 'UnsupportedInventoryItemContextException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'typename', '3': 446064463, '4': 1, '5': 9, '10': 'typename'},
  ],
};

/// Descriptor for `UnsupportedInventoryItemContextException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedInventoryItemContextExceptionDescriptor =
    $convert.base64Decode(
        'CihVbnN1cHBvcnRlZEludmVudG9yeUl0ZW1Db250ZXh0RXhjZXB0aW9uEhsKB21lc3NhZ2UYhb'
        'O7cCABKAlSB21lc3NhZ2USHgoIdHlwZW5hbWUYz87Z1AEgASgJUgh0eXBlbmFtZQ==');

@$core.Deprecated(
    'Use unsupportedInventorySchemaVersionExceptionDescriptor instead')
const UnsupportedInventorySchemaVersionException$json = {
  '1': 'UnsupportedInventorySchemaVersionException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedInventorySchemaVersionException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    unsupportedInventorySchemaVersionExceptionDescriptor =
    $convert.base64Decode(
        'CipVbnN1cHBvcnRlZEludmVudG9yeVNjaGVtYVZlcnNpb25FeGNlcHRpb24SGwoHbWVzc2FnZR'
        'iFs7twIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use unsupportedOperatingSystemDescriptor instead')
const UnsupportedOperatingSystem$json = {
  '1': 'UnsupportedOperatingSystem',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperatingSystem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperatingSystemDescriptor =
    $convert.base64Decode(
        'ChpVbnN1cHBvcnRlZE9wZXJhdGluZ1N5c3RlbRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use unsupportedOperationExceptionDescriptor instead')
const UnsupportedOperationException$json = {
  '1': 'UnsupportedOperationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZE9wZXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use unsupportedParameterTypeDescriptor instead')
const UnsupportedParameterType$json = {
  '1': 'UnsupportedParameterType',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedParameterType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedParameterTypeDescriptor =
    $convert.base64Decode(
        'ChhVbnN1cHBvcnRlZFBhcmFtZXRlclR5cGUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use unsupportedPlatformTypeDescriptor instead')
const UnsupportedPlatformType$json = {
  '1': 'UnsupportedPlatformType',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedPlatformType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedPlatformTypeDescriptor =
    $convert.base64Decode(
        'ChdVbnN1cHBvcnRlZFBsYXRmb3JtVHlwZRIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use updateAssociationRequestDescriptor instead')
const UpdateAssociationRequest$json = {
  '1': 'UpdateAssociationRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'applyonlyatcroninterval',
      '3': 285867158,
      '4': 1,
      '5': 8,
      '10': 'applyonlyatcroninterval'
    },
    {
      '1': 'associationdispatchassumerole',
      '3': 281627425,
      '4': 1,
      '5': 9,
      '10': 'associationdispatchassumerole'
    },
    {
      '1': 'associationid',
      '3': 138771986,
      '4': 1,
      '5': 9,
      '10': 'associationid'
    },
    {
      '1': 'associationname',
      '3': 313608216,
      '4': 1,
      '5': 9,
      '10': 'associationname'
    },
    {
      '1': 'associationversion',
      '3': 447890705,
      '4': 1,
      '5': 9,
      '10': 'associationversion'
    },
    {
      '1': 'automationtargetparametername',
      '3': 348584826,
      '4': 1,
      '5': 9,
      '10': 'automationtargetparametername'
    },
    {
      '1': 'calendarnames',
      '3': 36075966,
      '4': 3,
      '5': 9,
      '10': 'calendarnames'
    },
    {
      '1': 'complianceseverity',
      '3': 278891158,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationComplianceSeverity',
      '10': 'complianceseverity'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'outputlocation',
      '3': 67991028,
      '4': 1,
      '5': 11,
      '6': '.ssm.InstanceAssociationOutputLocation',
      '10': 'outputlocation'
    },
    {
      '1': 'parameters',
      '3': 494900218,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateAssociationRequest.ParametersEntry',
      '10': 'parameters'
    },
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'synccompliance',
      '3': 500469318,
      '4': 1,
      '5': 14,
      '6': '.ssm.AssociationSyncCompliance',
      '10': 'synccompliance'
    },
    {
      '1': 'targetlocations',
      '3': 289168805,
      '4': 3,
      '5': 11,
      '6': '.ssm.TargetLocation',
      '10': 'targetlocations'
    },
    {
      '1': 'targetmaps',
      '3': 74800696,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateAssociationRequest.TargetmapsEntry',
      '10': 'targetmaps'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
  ],
  '3': [
    UpdateAssociationRequest_ParametersEntry$json,
    UpdateAssociationRequest_TargetmapsEntry$json
  ],
};

@$core.Deprecated('Use updateAssociationRequestDescriptor instead')
const UpdateAssociationRequest_ParametersEntry$json = {
  '1': 'ParametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use updateAssociationRequestDescriptor instead')
const UpdateAssociationRequest_TargetmapsEntry$json = {
  '1': 'TargetmapsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAssociationRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVBc3NvY2lhdGlvblJlcXVlc3QSSgoSYWxhcm1jb25maWd1cmF0aW9uGImZuSEgAS'
    'gLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9uEjwKF2FwcGx5'
    'b25seWF0Y3JvbmludGVydmFsGJb5p4gBIAEoCFIXYXBwbHlvbmx5YXRjcm9uaW50ZXJ2YWwSSA'
    'odYXNzb2NpYXRpb25kaXNwYXRjaGFzc3VtZXJvbGUYoZalhgEgASgJUh1hc3NvY2lhdGlvbmRp'
    'c3BhdGNoYXNzdW1lcm9sZRInCg1hc3NvY2lhdGlvbmlkGJL8lUIgASgJUg1hc3NvY2lhdGlvbm'
    'lkEiwKD2Fzc29jaWF0aW9ubmFtZRiYkMWVASABKAlSD2Fzc29jaWF0aW9ubmFtZRIyChJhc3Nv'
    'Y2lhdGlvbnZlcnNpb24YkYrJ1QEgASgJUhJhc3NvY2lhdGlvbnZlcnNpb24SSAodYXV0b21hdG'
    'lvbnRhcmdldHBhcmFtZXRlcm5hbWUY+vabpgEgASgJUh1hdXRvbWF0aW9udGFyZ2V0cGFyYW1l'
    'dGVybmFtZRInCg1jYWxlbmRhcm5hbWVzGL7zmREgAygJUg1jYWxlbmRhcm5hbWVzElYKEmNvbX'
    'BsaWFuY2VzZXZlcml0eRiWlf6EASABKA4yIi5zc20uQXNzb2NpYXRpb25Db21wbGlhbmNlU2V2'
    'ZXJpdHlSEmNvbXBsaWFuY2VzZXZlcml0eRIrCg9kb2N1bWVudHZlcnNpb24Yye+pKCABKAlSD2'
    'RvY3VtZW50dmVyc2lvbhIeCghkdXJhdGlvbhiukp2mASABKAVSCGR1cmF0aW9uEikKDm1heGNv'
    'bmN1cnJlbmN5GP3Bjg4gASgJUg5tYXhjb25jdXJyZW5jeRIfCgltYXhlcnJvcnMYq8L1PSABKA'
    'lSCW1heGVycm9ycxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElEKDm91dHB1dGxvY2F0aW9uGPTr'
    'tSAgASgLMiYuc3NtLkluc3RhbmNlQXNzb2NpYXRpb25PdXRwdXRMb2NhdGlvblIOb3V0cHV0bG'
    '9jYXRpb24SUQoKcGFyYW1ldGVycxj6p/7rASADKAsyLS5zc20uVXBkYXRlQXNzb2NpYXRpb25S'
    'ZXF1ZXN0LlBhcmFtZXRlcnNFbnRyeVIKcGFyYW1ldGVycxIyChJzY2hlZHVsZWV4cHJlc3Npb2'
    '4Y/5Hb1AEgASgJUhJzY2hlZHVsZWV4cHJlc3Npb24SKQoOc2NoZWR1bGVvZmZzZXQY2JHqSiAB'
    'KAVSDnNjaGVkdWxlb2Zmc2V0EkoKDnN5bmNjb21wbGlhbmNlGMac0u4BIAEoDjIeLnNzbS5Bc3'
    'NvY2lhdGlvblN5bmNDb21wbGlhbmNlUg5zeW5jY29tcGxpYW5jZRJBCg90YXJnZXRsb2NhdGlv'
    'bnMYpbvxiQEgAygLMhMuc3NtLlRhcmdldExvY2F0aW9uUg90YXJnZXRsb2NhdGlvbnMSUAoKdG'
    'FyZ2V0bWFwcxi4vNUjIAMoCzItLnNzbS5VcGRhdGVBc3NvY2lhdGlvblJlcXVlc3QuVGFyZ2V0'
    'bWFwc0VudHJ5Ugp0YXJnZXRtYXBzEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ug'
    'd0YXJnZXRzGj0KD1BhcmFtZXRlcnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgC'
    'IAEoCVIFdmFsdWU6AjgBGj0KD1RhcmdldG1hcHNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
    'V2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use updateAssociationResultDescriptor instead')
const UpdateAssociationResult$json = {
  '1': 'UpdateAssociationResult',
  '2': [
    {
      '1': 'associationdescription',
      '3': 344755863,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationDescription',
      '10': 'associationdescription'
    },
  ],
};

/// Descriptor for `UpdateAssociationResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAssociationResultDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVBc3NvY2lhdGlvblJlc3VsdBJXChZhc3NvY2lhdGlvbmRlc2NyaXB0aW9uGJedsq'
    'QBIAEoCzIbLnNzbS5Bc3NvY2lhdGlvbkRlc2NyaXB0aW9uUhZhc3NvY2lhdGlvbmRlc2NyaXB0'
    'aW9u');

@$core.Deprecated('Use updateAssociationStatusRequestDescriptor instead')
const UpdateAssociationStatusRequest$json = {
  '1': 'UpdateAssociationStatusRequest',
  '2': [
    {
      '1': 'associationstatus',
      '3': 486964487,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationStatus',
      '10': 'associationstatus'
    },
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateAssociationStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAssociationStatusRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVBc3NvY2lhdGlvblN0YXR1c1JlcXVlc3QSSAoRYXNzb2NpYXRpb25zdGF0dXMYh/'
        'qZ6AEgASgLMhYuc3NtLkFzc29jaWF0aW9uU3RhdHVzUhFhc3NvY2lhdGlvbnN0YXR1cxIhCgpp'
        'bnN0YW5jZWlkGKCt0RcgASgJUgppbnN0YW5jZWlkEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use updateAssociationStatusResultDescriptor instead')
const UpdateAssociationStatusResult$json = {
  '1': 'UpdateAssociationStatusResult',
  '2': [
    {
      '1': 'associationdescription',
      '3': 344755863,
      '4': 1,
      '5': 11,
      '6': '.ssm.AssociationDescription',
      '10': 'associationdescription'
    },
  ],
};

/// Descriptor for `UpdateAssociationStatusResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAssociationStatusResultDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVBc3NvY2lhdGlvblN0YXR1c1Jlc3VsdBJXChZhc3NvY2lhdGlvbmRlc2NyaXB0aW'
        '9uGJedsqQBIAEoCzIbLnNzbS5Bc3NvY2lhdGlvbkRlc2NyaXB0aW9uUhZhc3NvY2lhdGlvbmRl'
        'c2NyaXB0aW9u');

@$core.Deprecated('Use updateDocumentDefaultVersionRequestDescriptor instead')
const UpdateDocumentDefaultVersionRequest$json = {
  '1': 'UpdateDocumentDefaultVersionRequest',
  '2': [
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateDocumentDefaultVersionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentDefaultVersionRequestDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVEb2N1bWVudERlZmF1bHRWZXJzaW9uUmVxdWVzdBIrCg9kb2N1bWVudHZlcnNpb2'
        '4Yye+pKCABKAlSD2RvY3VtZW50dmVyc2lvbhIVCgRuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use updateDocumentDefaultVersionResultDescriptor instead')
const UpdateDocumentDefaultVersionResult$json = {
  '1': 'UpdateDocumentDefaultVersionResult',
  '2': [
    {
      '1': 'description',
      '3': 115243530,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentDefaultVersionDescription',
      '10': 'description'
    },
  ],
};

/// Descriptor for `UpdateDocumentDefaultVersionResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentDefaultVersionResultDescriptor =
    $convert.base64Decode(
        'CiJVcGRhdGVEb2N1bWVudERlZmF1bHRWZXJzaW9uUmVzdWx0EksKC2Rlc2NyaXB0aW9uGIr0+T'
        'YgASgLMiYuc3NtLkRvY3VtZW50RGVmYXVsdFZlcnNpb25EZXNjcmlwdGlvblILZGVzY3JpcHRp'
        'b24=');

@$core.Deprecated('Use updateDocumentMetadataRequestDescriptor instead')
const UpdateDocumentMetadataRequest$json = {
  '1': 'UpdateDocumentMetadataRequest',
  '2': [
    {
      '1': 'documentreviews',
      '3': 387191202,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentReviews',
      '10': 'documentreviews'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UpdateDocumentMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentMetadataRequestDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVEb2N1bWVudE1ldGFkYXRhUmVxdWVzdBJCCg9kb2N1bWVudHJldmlld3MYoqPQuA'
    'EgASgLMhQuc3NtLkRvY3VtZW50UmV2aWV3c1IPZG9jdW1lbnRyZXZpZXdzEisKD2RvY3VtZW50'
    'dmVyc2lvbhjJ76koIAEoCVIPZG9jdW1lbnR2ZXJzaW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'U=');

@$core.Deprecated('Use updateDocumentMetadataResponseDescriptor instead')
const UpdateDocumentMetadataResponse$json = {
  '1': 'UpdateDocumentMetadataResponse',
};

/// Descriptor for `UpdateDocumentMetadataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentMetadataResponseDescriptor =
    $convert.base64Decode('Ch5VcGRhdGVEb2N1bWVudE1ldGFkYXRhUmVzcG9uc2U=');

@$core.Deprecated('Use updateDocumentRequestDescriptor instead')
const UpdateDocumentRequest$json = {
  '1': 'UpdateDocumentRequest',
  '2': [
    {
      '1': 'attachments',
      '3': 498946338,
      '4': 3,
      '5': 11,
      '6': '.ssm.AttachmentsSource',
      '10': 'attachments'
    },
    {'1': 'content', '3': 23568227, '4': 1, '5': 9, '10': 'content'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {
      '1': 'documentformat',
      '3': 516934792,
      '4': 1,
      '5': 14,
      '6': '.ssm.DocumentFormat',
      '10': 'documentformat'
    },
    {
      '1': 'documentversion',
      '3': 84572105,
      '4': 1,
      '5': 9,
      '10': 'documentversion'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'targettype', '3': 397256481, '4': 1, '5': 9, '10': 'targettype'},
    {'1': 'versionname', '3': 227348949, '4': 1, '5': 9, '10': 'versionname'},
  ],
};

/// Descriptor for `UpdateDocumentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVEb2N1bWVudFJlcXVlc3QSPAoLYXR0YWNobWVudHMYoqL17QEgAygLMhYuc3NtLk'
    'F0dGFjaG1lbnRzU291cmNlUgthdHRhY2htZW50cxIbCgdjb250ZW50GOO+ngsgASgJUgdjb250'
    'ZW50EiQKC2Rpc3BsYXluYW1lGLfJsscBIAEoCVILZGlzcGxheW5hbWUSPwoOZG9jdW1lbnRmb3'
    'JtYXQYiJm/9gEgASgOMhMuc3NtLkRvY3VtZW50Rm9ybWF0Ug5kb2N1bWVudGZvcm1hdBIrCg9k'
    'b2N1bWVudHZlcnNpb24Yye+pKCABKAlSD2RvY3VtZW50dmVyc2lvbhIVCgRuYW1lGIfmgX8gAS'
    'gJUgRuYW1lEiIKCnRhcmdldHR5cGUYoc62vQEgASgJUgp0YXJnZXR0eXBlEiMKC3ZlcnNpb25u'
    'YW1lGNWjtGwgASgJUgt2ZXJzaW9ubmFtZQ==');

@$core.Deprecated('Use updateDocumentResultDescriptor instead')
const UpdateDocumentResult$json = {
  '1': 'UpdateDocumentResult',
  '2': [
    {
      '1': 'documentdescription',
      '3': 500822655,
      '4': 1,
      '5': 11,
      '6': '.ssm.DocumentDescription',
      '10': 'documentdescription'
    },
  ],
};

/// Descriptor for `UpdateDocumentResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDocumentResultDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVEb2N1bWVudFJlc3VsdBJOChNkb2N1bWVudGRlc2NyaXB0aW9uGP/k5+4BIAEoCz'
    'IYLnNzbS5Eb2N1bWVudERlc2NyaXB0aW9uUhNkb2N1bWVudGRlc2NyaXB0aW9u');

@$core.Deprecated('Use updateMaintenanceWindowRequestDescriptor instead')
const UpdateMaintenanceWindowRequest$json = {
  '1': 'UpdateMaintenanceWindowRequest',
  '2': [
    {
      '1': 'allowunassociatedtargets',
      '3': 154411300,
      '4': 1,
      '5': 8,
      '10': 'allowunassociatedtargets'
    },
    {'1': 'cutoff', '3': 498433089, '4': 1, '5': 5, '10': 'cutoff'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'replace', '3': 81088356, '4': 1, '5': 8, '10': 'replace'},
    {'1': 'schedule', '3': 66697965, '4': 1, '5': 9, '10': 'schedule'},
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'scheduletimezone',
      '3': 170037696,
      '4': 1,
      '5': 9,
      '10': 'scheduletimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `UpdateMaintenanceWindowRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowRequestDescriptor = $convert.base64Decode(
    'Ch5VcGRhdGVNYWludGVuYW5jZVdpbmRvd1JlcXVlc3QSPQoYYWxsb3d1bmFzc29jaWF0ZWR0YX'
    'JnZXRzGKTC0EkgASgIUhhhbGxvd3VuYXNzb2NpYXRlZHRhcmdldHMSGgoGY3V0b2ZmGMH41e0B'
    'IAEoBVIGY3V0b2ZmEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIeCghkdX'
    'JhdGlvbhiukp2mASABKAVSCGR1cmF0aW9uEhwKB2VuYWJsZWQYv8ib5AEgASgIUgdlbmFibGVk'
    'EhsKB2VuZGRhdGUYz7P5JCABKAlSB2VuZGRhdGUSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRIbCg'
    'dyZXBsYWNlGOSe1SYgASgIUgdyZXBsYWNlEh0KCHNjaGVkdWxlGO315h8gASgJUghzY2hlZHVs'
    'ZRIpCg5zY2hlZHVsZW9mZnNldBjYkepKIAEoBVIOc2NoZWR1bGVvZmZzZXQSLQoQc2NoZWR1bG'
    'V0aW1lem9uZRjAo4pRIAEoCVIQc2NoZWR1bGV0aW1lem9uZRIgCglzdGFydGRhdGUY/Pig1AEg'
    'ASgJUglzdGFydGRhdGUSHQoId2luZG93aWQYqeSHCSABKAlSCHdpbmRvd2lk');

@$core.Deprecated('Use updateMaintenanceWindowResultDescriptor instead')
const UpdateMaintenanceWindowResult$json = {
  '1': 'UpdateMaintenanceWindowResult',
  '2': [
    {
      '1': 'allowunassociatedtargets',
      '3': 154411300,
      '4': 1,
      '5': 8,
      '10': 'allowunassociatedtargets'
    },
    {'1': 'cutoff', '3': 498433089, '4': 1, '5': 5, '10': 'cutoff'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'duration', '3': 348604718, '4': 1, '5': 5, '10': 'duration'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'schedule', '3': 66697965, '4': 1, '5': 9, '10': 'schedule'},
    {
      '1': 'scheduleoffset',
      '3': 156928216,
      '4': 1,
      '5': 5,
      '10': 'scheduleoffset'
    },
    {
      '1': 'scheduletimezone',
      '3': 170037696,
      '4': 1,
      '5': 9,
      '10': 'scheduletimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
  ],
};

/// Descriptor for `UpdateMaintenanceWindowResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowResultDescriptor = $convert.base64Decode(
    'Ch1VcGRhdGVNYWludGVuYW5jZVdpbmRvd1Jlc3VsdBI9ChhhbGxvd3VuYXNzb2NpYXRlZHRhcm'
    'dldHMYpMLQSSABKAhSGGFsbG93dW5hc3NvY2lhdGVkdGFyZ2V0cxIaCgZjdXRvZmYYwfjV7QEg'
    'ASgFUgZjdXRvZmYSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEh4KCGR1cm'
    'F0aW9uGK6SnaYBIAEoBVIIZHVyYXRpb24SHAoHZW5hYmxlZBi/yJvkASABKAhSB2VuYWJsZWQS'
    'GwoHZW5kZGF0ZRjPs/kkIAEoCVIHZW5kZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEh0KCH'
    'NjaGVkdWxlGO315h8gASgJUghzY2hlZHVsZRIpCg5zY2hlZHVsZW9mZnNldBjYkepKIAEoBVIO'
    'c2NoZWR1bGVvZmZzZXQSLQoQc2NoZWR1bGV0aW1lem9uZRjAo4pRIAEoCVIQc2NoZWR1bGV0aW'
    '1lem9uZRIgCglzdGFydGRhdGUY/Pig1AEgASgJUglzdGFydGRhdGUSHQoId2luZG93aWQYqeSH'
    'CSABKAlSCHdpbmRvd2lk');

@$core.Deprecated('Use updateMaintenanceWindowTargetRequestDescriptor instead')
const UpdateMaintenanceWindowTargetRequest$json = {
  '1': 'UpdateMaintenanceWindowTargetRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {'1': 'replace', '3': 81088356, '4': 1, '5': 8, '10': 'replace'},
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `UpdateMaintenanceWindowTargetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowTargetRequestDescriptor = $convert.base64Decode(
    'CiRVcGRhdGVNYWludGVuYW5jZVdpbmRvd1RhcmdldFJlcXVlc3QSIwoLZGVzY3JpcHRpb24Yiv'
    'T5NiABKAlSC2Rlc2NyaXB0aW9uEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLQoQb3duZXJpbmZv'
    'cm1hdGlvbhiDmsUgIAEoCVIQb3duZXJpbmZvcm1hdGlvbhIbCgdyZXBsYWNlGOSe1SYgASgIUg'
    'dyZXBsYWNlEigKB3RhcmdldHMYgpuCfSADKAsyCy5zc20uVGFyZ2V0Ugd0YXJnZXRzEh0KCHdp'
    'bmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZBIpCg53aW5kb3d0YXJnZXRpZBjezup7IAEoCVIOd2'
    'luZG93dGFyZ2V0aWQ=');

@$core.Deprecated('Use updateMaintenanceWindowTargetResultDescriptor instead')
const UpdateMaintenanceWindowTargetResult$json = {
  '1': 'UpdateMaintenanceWindowTargetResult',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'ownerinformation',
      '3': 68242691,
      '4': 1,
      '5': 9,
      '10': 'ownerinformation'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {
      '1': 'windowtargetid',
      '3': 259696478,
      '4': 1,
      '5': 9,
      '10': 'windowtargetid'
    },
  ],
};

/// Descriptor for `UpdateMaintenanceWindowTargetResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowTargetResultDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVNYWludGVuYW5jZVdpbmRvd1RhcmdldFJlc3VsdBIjCgtkZXNjcmlwdGlvbhiK9P'
        'k2IAEoCVILZGVzY3JpcHRpb24SFQoEbmFtZRiH5oF/IAEoCVIEbmFtZRItChBvd25lcmluZm9y'
        'bWF0aW9uGIOaxSAgASgJUhBvd25lcmluZm9ybWF0aW9uEigKB3RhcmdldHMYgpuCfSADKAsyCy'
        '5zc20uVGFyZ2V0Ugd0YXJnZXRzEh0KCHdpbmRvd2lkGKnkhwkgASgJUgh3aW5kb3dpZBIpCg53'
        'aW5kb3d0YXJnZXRpZBjezup7IAEoCVIOd2luZG93dGFyZ2V0aWQ=');

@$core.Deprecated('Use updateMaintenanceWindowTaskRequestDescriptor instead')
const UpdateMaintenanceWindowTaskRequest$json = {
  '1': 'UpdateMaintenanceWindowTaskRequest',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cutoffbehavior',
      '3': 120608587,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskCutoffBehavior',
      '10': 'cutoffbehavior'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'logginginfo',
      '3': 448312415,
      '4': 1,
      '5': 11,
      '6': '.ssm.LoggingInfo',
      '10': 'logginginfo'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {'1': 'replace', '3': 81088356, '4': 1, '5': 8, '10': 'replace'},
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskinvocationparameters',
      '3': 347127635,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskInvocationParameters',
      '10': 'taskinvocationparameters'
    },
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateMaintenanceWindowTaskRequest.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
  '3': [UpdateMaintenanceWindowTaskRequest_TaskparametersEntry$json],
};

@$core.Deprecated('Use updateMaintenanceWindowTaskRequestDescriptor instead')
const UpdateMaintenanceWindowTaskRequest_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateMaintenanceWindowTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowTaskRequestDescriptor = $convert.base64Decode(
    'CiJVcGRhdGVNYWludGVuYW5jZVdpbmRvd1Rhc2tSZXF1ZXN0EkoKEmFsYXJtY29uZmlndXJhdG'
    'lvbhiJmbkhIAEoCzIXLnNzbS5BbGFybUNvbmZpZ3VyYXRpb25SEmFsYXJtY29uZmlndXJhdGlv'
    'bhJTCg5jdXRvZmZiZWhhdmlvchjLrsE5IAEoDjIoLnNzbS5NYWludGVuYW5jZVdpbmRvd1Rhc2'
    'tDdXRvZmZCZWhhdmlvclIOY3V0b2ZmYmVoYXZpb3ISIwoLZGVzY3JpcHRpb24YivT5NiABKAlS'
    'C2Rlc2NyaXB0aW9uEjYKC2xvZ2dpbmdpbmZvGN/o4tUBIAEoCzIQLnNzbS5Mb2dnaW5nSW5mb1'
    'ILbG9nZ2luZ2luZm8SKQoObWF4Y29uY3VycmVuY3kY/cGODiABKAlSDm1heGNvbmN1cnJlbmN5'
    'Eh8KCW1heGVycm9ycxirwvU9IAEoCVIJbWF4ZXJyb3JzEhUKBG5hbWUYh+aBfyABKAlSBG5hbW'
    'USHQoIcHJpb3JpdHkYqr62NCABKAVSCHByaW9yaXR5EhsKB3JlcGxhY2UY5J7VJiABKAhSB3Jl'
    'cGxhY2USKgoOc2VydmljZXJvbGVhcm4YhOPatgEgASgJUg5zZXJ2aWNlcm9sZWFybhIoCgd0YX'
    'JnZXRzGIKbgn0gAygLMgsuc3NtLlRhcmdldFIHdGFyZ2V0cxIcCgd0YXNrYXJuGOTJ+pQBIAEo'
    'CVIHdGFza2FybhJuChh0YXNraW52b2NhdGlvbnBhcmFtZXRlcnMY0/7CpQEgASgLMi4uc3NtLk'
    '1haW50ZW5hbmNlV2luZG93VGFza0ludm9jYXRpb25QYXJhbWV0ZXJzUhh0YXNraW52b2NhdGlv'
    'bnBhcmFtZXRlcnMSZwoOdGFza3BhcmFtZXRlcnMYgY/mtwEgAygLMjsuc3NtLlVwZGF0ZU1haW'
    '50ZW5hbmNlV2luZG93VGFza1JlcXVlc3QuVGFza3BhcmFtZXRlcnNFbnRyeVIOdGFza3BhcmFt'
    'ZXRlcnMSHQoId2luZG93aWQYqeSHCSABKAlSCHdpbmRvd2lkEiYKDHdpbmRvd3Rhc2tpZBjajb'
    'GaASABKAlSDHdpbmRvd3Rhc2tpZBp1ChNUYXNrcGFyYW1ldGVyc0VudHJ5EhAKA2tleRgBIAEo'
    'CVIDa2V5EkgKBXZhbHVlGAIgASgLMjIuc3NtLk1haW50ZW5hbmNlV2luZG93VGFza1BhcmFtZX'
    'RlclZhbHVlRXhwcmVzc2lvblIFdmFsdWU6AjgB');

@$core.Deprecated('Use updateMaintenanceWindowTaskResultDescriptor instead')
const UpdateMaintenanceWindowTaskResult$json = {
  '1': 'UpdateMaintenanceWindowTaskResult',
  '2': [
    {
      '1': 'alarmconfiguration',
      '3': 70143113,
      '4': 1,
      '5': 11,
      '6': '.ssm.AlarmConfiguration',
      '10': 'alarmconfiguration'
    },
    {
      '1': 'cutoffbehavior',
      '3': 120608587,
      '4': 1,
      '5': 14,
      '6': '.ssm.MaintenanceWindowTaskCutoffBehavior',
      '10': 'cutoffbehavior'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'logginginfo',
      '3': 448312415,
      '4': 1,
      '5': 11,
      '6': '.ssm.LoggingInfo',
      '10': 'logginginfo'
    },
    {
      '1': 'maxconcurrency',
      '3': 29597949,
      '4': 1,
      '5': 9,
      '10': 'maxconcurrency'
    },
    {'1': 'maxerrors', '3': 129851691, '4': 1, '5': 9, '10': 'maxerrors'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'servicerolearn',
      '3': 383168900,
      '4': 1,
      '5': 9,
      '10': 'servicerolearn'
    },
    {
      '1': 'targets',
      '3': 262180226,
      '4': 3,
      '5': 11,
      '6': '.ssm.Target',
      '10': 'targets'
    },
    {'1': 'taskarn', '3': 312386788, '4': 1, '5': 9, '10': 'taskarn'},
    {
      '1': 'taskinvocationparameters',
      '3': 347127635,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskInvocationParameters',
      '10': 'taskinvocationparameters'
    },
    {
      '1': 'taskparameters',
      '3': 385451905,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateMaintenanceWindowTaskResult.TaskparametersEntry',
      '10': 'taskparameters'
    },
    {'1': 'windowid', '3': 19001897, '4': 1, '5': 9, '10': 'windowid'},
    {'1': 'windowtaskid', '3': 323765978, '4': 1, '5': 9, '10': 'windowtaskid'},
  ],
  '3': [UpdateMaintenanceWindowTaskResult_TaskparametersEntry$json],
};

@$core.Deprecated('Use updateMaintenanceWindowTaskResultDescriptor instead')
const UpdateMaintenanceWindowTaskResult_TaskparametersEntry$json = {
  '1': 'TaskparametersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MaintenanceWindowTaskParameterValueExpression',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateMaintenanceWindowTaskResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaintenanceWindowTaskResultDescriptor = $convert.base64Decode(
    'CiFVcGRhdGVNYWludGVuYW5jZVdpbmRvd1Rhc2tSZXN1bHQSSgoSYWxhcm1jb25maWd1cmF0aW'
    '9uGImZuSEgASgLMhcuc3NtLkFsYXJtQ29uZmlndXJhdGlvblISYWxhcm1jb25maWd1cmF0aW9u'
    'ElMKDmN1dG9mZmJlaGF2aW9yGMuuwTkgASgOMiguc3NtLk1haW50ZW5hbmNlV2luZG93VGFza0'
    'N1dG9mZkJlaGF2aW9yUg5jdXRvZmZiZWhhdmlvchIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVIL'
    'ZGVzY3JpcHRpb24SNgoLbG9nZ2luZ2luZm8Y3+ji1QEgASgLMhAuc3NtLkxvZ2dpbmdJbmZvUg'
    'tsb2dnaW5naW5mbxIpCg5tYXhjb25jdXJyZW5jeRj9wY4OIAEoCVIObWF4Y29uY3VycmVuY3kS'
    'HwoJbWF4ZXJyb3JzGKvC9T0gASgJUgltYXhlcnJvcnMSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZR'
    'IdCghwcmlvcml0eRiqvrY0IAEoBVIIcHJpb3JpdHkSKgoOc2VydmljZXJvbGVhcm4YhOPatgEg'
    'ASgJUg5zZXJ2aWNlcm9sZWFybhIoCgd0YXJnZXRzGIKbgn0gAygLMgsuc3NtLlRhcmdldFIHdG'
    'FyZ2V0cxIcCgd0YXNrYXJuGOTJ+pQBIAEoCVIHdGFza2FybhJuChh0YXNraW52b2NhdGlvbnBh'
    'cmFtZXRlcnMY0/7CpQEgASgLMi4uc3NtLk1haW50ZW5hbmNlV2luZG93VGFza0ludm9jYXRpb2'
    '5QYXJhbWV0ZXJzUhh0YXNraW52b2NhdGlvbnBhcmFtZXRlcnMSZgoOdGFza3BhcmFtZXRlcnMY'
    'gY/mtwEgAygLMjouc3NtLlVwZGF0ZU1haW50ZW5hbmNlV2luZG93VGFza1Jlc3VsdC5UYXNrcG'
    'FyYW1ldGVyc0VudHJ5Ug50YXNrcGFyYW1ldGVycxIdCgh3aW5kb3dpZBip5IcJIAEoCVIId2lu'
    'ZG93aWQSJgoMd2luZG93dGFza2lkGNqNsZoBIAEoCVIMd2luZG93dGFza2lkGnUKE1Rhc2twYX'
    'JhbWV0ZXJzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSSAoFdmFsdWUYAiABKAsyMi5zc20uTWFp'
    'bnRlbmFuY2VXaW5kb3dUYXNrUGFyYW1ldGVyVmFsdWVFeHByZXNzaW9uUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use updateManagedInstanceRoleRequestDescriptor instead')
const UpdateManagedInstanceRoleRequest$json = {
  '1': 'UpdateManagedInstanceRoleRequest',
  '2': [
    {'1': 'iamrole', '3': 242424351, '4': 1, '5': 9, '10': 'iamrole'},
    {'1': 'instanceid', '3': 49567392, '4': 1, '5': 9, '10': 'instanceid'},
  ],
};

/// Descriptor for `UpdateManagedInstanceRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateManagedInstanceRoleRequestDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVNYW5hZ2VkSW5zdGFuY2VSb2xlUmVxdWVzdBIbCgdpYW1yb2xlGJ+0zHMgASgJUg'
        'dpYW1yb2xlEiEKCmluc3RhbmNlaWQYoK3RFyABKAlSCmluc3RhbmNlaWQ=');

@$core.Deprecated('Use updateManagedInstanceRoleResultDescriptor instead')
const UpdateManagedInstanceRoleResult$json = {
  '1': 'UpdateManagedInstanceRoleResult',
};

/// Descriptor for `UpdateManagedInstanceRoleResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateManagedInstanceRoleResultDescriptor =
    $convert.base64Decode('Ch9VcGRhdGVNYW5hZ2VkSW5zdGFuY2VSb2xlUmVzdWx0');

@$core.Deprecated('Use updateOpsItemRequestDescriptor instead')
const UpdateOpsItemRequest$json = {
  '1': 'UpdateOpsItemRequest',
  '2': [
    {
      '1': 'actualendtime',
      '3': 452787526,
      '4': 1,
      '5': 9,
      '10': 'actualendtime'
    },
    {
      '1': 'actualstarttime',
      '3': 532853117,
      '4': 1,
      '5': 9,
      '10': 'actualstarttime'
    },
    {'1': 'category', '3': 263447954, '4': 1, '5': 9, '10': 'category'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'notifications',
      '3': 404560144,
      '4': 3,
      '5': 11,
      '6': '.ssm.OpsItemNotification',
      '10': 'notifications'
    },
    {
      '1': 'operationaldata',
      '3': 360357940,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateOpsItemRequest.OperationaldataEntry',
      '10': 'operationaldata'
    },
    {
      '1': 'operationaldatatodelete',
      '3': 352411530,
      '4': 3,
      '5': 9,
      '10': 'operationaldatatodelete'
    },
    {'1': 'opsitemarn', '3': 222489428, '4': 1, '5': 9, '10': 'opsitemarn'},
    {'1': 'opsitemid', '3': 25520466, '4': 1, '5': 9, '10': 'opsitemid'},
    {
      '1': 'plannedendtime',
      '3': 245727820,
      '4': 1,
      '5': 9,
      '10': 'plannedendtime'
    },
    {
      '1': 'plannedstarttime',
      '3': 478079215,
      '4': 1,
      '5': 9,
      '10': 'plannedstarttime'
    },
    {'1': 'priority', '3': 109944618, '4': 1, '5': 5, '10': 'priority'},
    {
      '1': 'relatedopsitems',
      '3': 287082393,
      '4': 3,
      '5': 11,
      '6': '.ssm.RelatedOpsItem',
      '10': 'relatedopsitems'
    },
    {'1': 'severity', '3': 276886227, '4': 1, '5': 9, '10': 'severity'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.ssm.OpsItemStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 81031594, '4': 1, '5': 9, '10': 'title'},
  ],
  '3': [UpdateOpsItemRequest_OperationaldataEntry$json],
};

@$core.Deprecated('Use updateOpsItemRequestDescriptor instead')
const UpdateOpsItemRequest_OperationaldataEntry$json = {
  '1': 'OperationaldataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.OpsItemDataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateOpsItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOpsItemRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVPcHNJdGVtUmVxdWVzdBIoCg1hY3R1YWxlbmR0aW1lGMb689cBIAEoCVINYWN0dW'
    'FsZW5kdGltZRIsCg9hY3R1YWxzdGFydHRpbWUY/eKK/gEgASgJUg9hY3R1YWxzdGFydHRpbWUS'
    'HQoIY2F0ZWdvcnkYksvPfSABKAlSCGNhdGVnb3J5EiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUg'
    'tkZXNjcmlwdGlvbhJCCg1ub3RpZmljYXRpb25zGJCy9MABIAMoCzIYLnNzbS5PcHNJdGVtTm90'
    'aWZpY2F0aW9uUg1ub3RpZmljYXRpb25zElwKD29wZXJhdGlvbmFsZGF0YRi0wOqrASADKAsyLi'
    '5zc20uVXBkYXRlT3BzSXRlbVJlcXVlc3QuT3BlcmF0aW9uYWxkYXRhRW50cnlSD29wZXJhdGlv'
    'bmFsZGF0YRI8ChdvcGVyYXRpb25hbGRhdGF0b2RlbGV0ZRiKv4WoASADKAlSF29wZXJhdGlvbm'
    'FsZGF0YXRvZGVsZXRlEiEKCm9wc2l0ZW1hcm4Y1NaLaiABKAlSCm9wc2l0ZW1hcm4SHwoJb3Bz'
    'aXRlbWlkGNLSlQwgASgJUglvcHNpdGVtaWQSKQoOcGxhbm5lZGVuZHRpbWUYzISWdSABKAlSDn'
    'BsYW5uZWRlbmR0aW1lEi4KEHBsYW5uZWRzdGFydHRpbWUY79H74wEgASgJUhBwbGFubmVkc3Rh'
    'cnR0aW1lEh0KCHByaW9yaXR5GKq+tjQgASgFUghwcmlvcml0eRJBCg9yZWxhdGVkb3BzaXRlbX'
    'MYmY/yiAEgAygLMhMuc3NtLlJlbGF0ZWRPcHNJdGVtUg9yZWxhdGVkb3BzaXRlbXMSHgoIc2V2'
    'ZXJpdHkY0+WDhAEgASgJUghzZXZlcml0eRItCgZzdGF0dXMYkOT7AiABKA4yEi5zc20uT3BzSX'
    'RlbVN0YXR1c1IGc3RhdHVzEhcKBXRpdGxlGKrj0SYgASgJUgV0aXRsZRpZChRPcGVyYXRpb25h'
    'bGRhdGFFbnRyeRIQCgNrZXkYASABKAlSA2tleRIrCgV2YWx1ZRgCIAEoCzIVLnNzbS5PcHNJdG'
    'VtRGF0YVZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use updateOpsItemResponseDescriptor instead')
const UpdateOpsItemResponse$json = {
  '1': 'UpdateOpsItemResponse',
};

/// Descriptor for `UpdateOpsItemResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOpsItemResponseDescriptor =
    $convert.base64Decode('ChVVcGRhdGVPcHNJdGVtUmVzcG9uc2U=');

@$core.Deprecated('Use updateOpsMetadataRequestDescriptor instead')
const UpdateOpsMetadataRequest$json = {
  '1': 'UpdateOpsMetadataRequest',
  '2': [
    {'1': 'keystodelete', '3': 218314840, '4': 3, '5': 9, '10': 'keystodelete'},
    {
      '1': 'metadatatoupdate',
      '3': 156011881,
      '4': 3,
      '5': 11,
      '6': '.ssm.UpdateOpsMetadataRequest.MetadatatoupdateEntry',
      '10': 'metadatatoupdate'
    },
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
  ],
  '3': [UpdateOpsMetadataRequest_MetadatatoupdateEntry$json],
};

@$core.Deprecated('Use updateOpsMetadataRequestDescriptor instead')
const UpdateOpsMetadataRequest_MetadatatoupdateEntry$json = {
  '1': 'MetadatatoupdateEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.ssm.MetadataValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateOpsMetadataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOpsMetadataRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVPcHNNZXRhZGF0YVJlcXVlc3QSJQoMa2V5c3RvZGVsZXRlGNjwjGggAygJUgxrZX'
    'lzdG9kZWxldGUSYgoQbWV0YWRhdGF0b3VwZGF0ZRjpmrJKIAMoCzIzLnNzbS5VcGRhdGVPcHNN'
    'ZXRhZGF0YVJlcXVlc3QuTWV0YWRhdGF0b3VwZGF0ZUVudHJ5UhBtZXRhZGF0YXRvdXBkYXRlEi'
    'oKDm9wc21ldGFkYXRhYXJuGKK+guYBIAEoCVIOb3BzbWV0YWRhdGFhcm4aVwoVTWV0YWRhdGF0'
    'b3VwZGF0ZUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EigKBXZhbHVlGAIgASgLMhIuc3NtLk1ldG'
    'FkYXRhVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use updateOpsMetadataResultDescriptor instead')
const UpdateOpsMetadataResult$json = {
  '1': 'UpdateOpsMetadataResult',
  '2': [
    {
      '1': 'opsmetadataarn',
      '3': 482385698,
      '4': 1,
      '5': 9,
      '10': 'opsmetadataarn'
    },
  ],
};

/// Descriptor for `UpdateOpsMetadataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateOpsMetadataResultDescriptor =
    $convert.base64Decode(
        'ChdVcGRhdGVPcHNNZXRhZGF0YVJlc3VsdBIqCg5vcHNtZXRhZGF0YWFybhiivoLmASABKAlSDm'
        '9wc21ldGFkYXRhYXJu');

@$core.Deprecated('Use updatePatchBaselineRequestDescriptor instead')
const UpdatePatchBaselineRequest$json = {
  '1': 'UpdatePatchBaselineRequest',
  '2': [
    {
      '1': 'approvalrules',
      '3': 71466346,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchRuleGroup',
      '10': 'approvalrules'
    },
    {
      '1': 'approvedpatches',
      '3': 199384709,
      '4': 3,
      '5': 9,
      '10': 'approvedpatches'
    },
    {
      '1': 'approvedpatchescompliancelevel',
      '3': 63924432,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'approvedpatchescompliancelevel'
    },
    {
      '1': 'approvedpatchesenablenonsecurity',
      '3': 295555901,
      '4': 1,
      '5': 8,
      '10': 'approvedpatchesenablenonsecurity'
    },
    {
      '1': 'availablesecurityupdatescompliancestatus',
      '3': 187471858,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceStatus',
      '10': 'availablesecurityupdatescompliancestatus'
    },
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'globalfilters',
      '3': 263302754,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'globalfilters'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'rejectedpatches',
      '3': 309657116,
      '4': 3,
      '5': 9,
      '10': 'rejectedpatches'
    },
    {
      '1': 'rejectedpatchesaction',
      '3': 356538330,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchAction',
      '10': 'rejectedpatchesaction'
    },
    {'1': 'replace', '3': 81088356, '4': 1, '5': 8, '10': 'replace'},
    {
      '1': 'sources',
      '3': 46625746,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchSource',
      '10': 'sources'
    },
  ],
};

/// Descriptor for `UpdatePatchBaselineRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePatchBaselineRequestDescriptor = $convert.base64Decode(
    'ChpVcGRhdGVQYXRjaEJhc2VsaW5lUmVxdWVzdBI8Cg1hcHByb3ZhbHJ1bGVzGOr6iSIgASgLMh'
    'Muc3NtLlBhdGNoUnVsZUdyb3VwUg1hcHByb3ZhbHJ1bGVzEisKD2FwcHJvdmVkcGF0Y2hlcxiF'
    'vYlfIAMoCVIPYXBwcm92ZWRwYXRjaGVzEmQKHmFwcHJvdmVkcGF0Y2hlc2NvbXBsaWFuY2VsZX'
    'ZlbBjQ0b0eIAEoDjIZLnNzbS5QYXRjaENvbXBsaWFuY2VMZXZlbFIeYXBwcm92ZWRwYXRjaGVz'
    'Y29tcGxpYW5jZWxldmVsEk4KIGFwcHJvdmVkcGF0Y2hlc2VuYWJsZW5vbnNlY3VyaXR5GL2m94'
    'wBIAEoCFIgYXBwcm92ZWRwYXRjaGVzZW5hYmxlbm9uc2VjdXJpdHkSeQooYXZhaWxhYmxlc2Vj'
    'dXJpdHl1cGRhdGVzY29tcGxpYW5jZXN0YXR1cxjyr7JZIAEoDjIaLnNzbS5QYXRjaENvbXBsaW'
    'FuY2VTdGF0dXNSKGF2YWlsYWJsZXNlY3VyaXR5dXBkYXRlc2NvbXBsaWFuY2VzdGF0dXMSIQoK'
    'YmFzZWxpbmVpZBjQ5NsoIAEoCVIKYmFzZWxpbmVpZBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCV'
    'ILZGVzY3JpcHRpb24SPgoNZ2xvYmFsZmlsdGVycxji3MZ9IAEoCzIVLnNzbS5QYXRjaEZpbHRl'
    'ckdyb3VwUg1nbG9iYWxmaWx0ZXJzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSLAoPcmVqZWN0ZW'
    'RwYXRjaGVzGJz805MBIAMoCVIPcmVqZWN0ZWRwYXRjaGVzEkoKFXJlamVjdGVkcGF0Y2hlc2Fj'
    'dGlvbhjar4GqASABKA4yEC5zc20uUGF0Y2hBY3Rpb25SFXJlamVjdGVkcGF0Y2hlc2FjdGlvbh'
    'IbCgdyZXBsYWNlGOSe1SYgASgIUgdyZXBsYWNlEi0KB3NvdXJjZXMY0uedFiADKAsyEC5zc20u'
    'UGF0Y2hTb3VyY2VSB3NvdXJjZXM=');

@$core.Deprecated('Use updatePatchBaselineResultDescriptor instead')
const UpdatePatchBaselineResult$json = {
  '1': 'UpdatePatchBaselineResult',
  '2': [
    {
      '1': 'approvalrules',
      '3': 71466346,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchRuleGroup',
      '10': 'approvalrules'
    },
    {
      '1': 'approvedpatches',
      '3': 199384709,
      '4': 3,
      '5': 9,
      '10': 'approvedpatches'
    },
    {
      '1': 'approvedpatchescompliancelevel',
      '3': 63924432,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceLevel',
      '10': 'approvedpatchescompliancelevel'
    },
    {
      '1': 'approvedpatchesenablenonsecurity',
      '3': 295555901,
      '4': 1,
      '5': 8,
      '10': 'approvedpatchesenablenonsecurity'
    },
    {
      '1': 'availablesecurityupdatescompliancestatus',
      '3': 187471858,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchComplianceStatus',
      '10': 'availablesecurityupdatescompliancestatus'
    },
    {'1': 'baselineid', '3': 85389904, '4': 1, '5': 9, '10': 'baselineid'},
    {'1': 'createddate', '3': 416929840, '4': 1, '5': 9, '10': 'createddate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'globalfilters',
      '3': 263302754,
      '4': 1,
      '5': 11,
      '6': '.ssm.PatchFilterGroup',
      '10': 'globalfilters'
    },
    {'1': 'modifieddate', '3': 210609143, '4': 1, '5': 9, '10': 'modifieddate'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'operatingsystem',
      '3': 38829802,
      '4': 1,
      '5': 14,
      '6': '.ssm.OperatingSystem',
      '10': 'operatingsystem'
    },
    {
      '1': 'rejectedpatches',
      '3': 309657116,
      '4': 3,
      '5': 9,
      '10': 'rejectedpatches'
    },
    {
      '1': 'rejectedpatchesaction',
      '3': 356538330,
      '4': 1,
      '5': 14,
      '6': '.ssm.PatchAction',
      '10': 'rejectedpatchesaction'
    },
    {
      '1': 'sources',
      '3': 46625746,
      '4': 3,
      '5': 11,
      '6': '.ssm.PatchSource',
      '10': 'sources'
    },
  ],
};

/// Descriptor for `UpdatePatchBaselineResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePatchBaselineResultDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVQYXRjaEJhc2VsaW5lUmVzdWx0EjwKDWFwcHJvdmFscnVsZXMY6vqJIiABKAsyEy'
    '5zc20uUGF0Y2hSdWxlR3JvdXBSDWFwcHJvdmFscnVsZXMSKwoPYXBwcm92ZWRwYXRjaGVzGIW9'
    'iV8gAygJUg9hcHByb3ZlZHBhdGNoZXMSZAoeYXBwcm92ZWRwYXRjaGVzY29tcGxpYW5jZWxldm'
    'VsGNDRvR4gASgOMhkuc3NtLlBhdGNoQ29tcGxpYW5jZUxldmVsUh5hcHByb3ZlZHBhdGNoZXNj'
    'b21wbGlhbmNlbGV2ZWwSTgogYXBwcm92ZWRwYXRjaGVzZW5hYmxlbm9uc2VjdXJpdHkYvab3jA'
    'EgASgIUiBhcHByb3ZlZHBhdGNoZXNlbmFibGVub25zZWN1cml0eRJ5CihhdmFpbGFibGVzZWN1'
    'cml0eXVwZGF0ZXNjb21wbGlhbmNlc3RhdHVzGPKvslkgASgOMhouc3NtLlBhdGNoQ29tcGxpYW'
    '5jZVN0YXR1c1IoYXZhaWxhYmxlc2VjdXJpdHl1cGRhdGVzY29tcGxpYW5jZXN0YXR1cxIhCgpi'
    'YXNlbGluZWlkGNDk2yggASgJUgpiYXNlbGluZWlkEiQKC2NyZWF0ZWRkYXRlGLCw58YBIAEoCV'
    'ILY3JlYXRlZGRhdGUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEj4KDWds'
    'b2JhbGZpbHRlcnMY4tzGfSABKAsyFS5zc20uUGF0Y2hGaWx0ZXJHcm91cFINZ2xvYmFsZmlsdG'
    'VycxIlCgxtb2RpZmllZGRhdGUY98e2ZCABKAlSDG1vZGlmaWVkZGF0ZRIVCgRuYW1lGIfmgX8g'
    'ASgJUgRuYW1lEkEKD29wZXJhdGluZ3N5c3RlbRjq/cESIAEoDjIULnNzbS5PcGVyYXRpbmdTeX'
    'N0ZW1SD29wZXJhdGluZ3N5c3RlbRIsCg9yZWplY3RlZHBhdGNoZXMYnPzTkwEgAygJUg9yZWpl'
    'Y3RlZHBhdGNoZXMSSgoVcmVqZWN0ZWRwYXRjaGVzYWN0aW9uGNqvgaoBIAEoDjIQLnNzbS5QYX'
    'RjaEFjdGlvblIVcmVqZWN0ZWRwYXRjaGVzYWN0aW9uEi0KB3NvdXJjZXMY0uedFiADKAsyEC5z'
    'c20uUGF0Y2hTb3VyY2VSB3NvdXJjZXM=');

@$core.Deprecated('Use updateResourceDataSyncRequestDescriptor instead')
const UpdateResourceDataSyncRequest$json = {
  '1': 'UpdateResourceDataSyncRequest',
  '2': [
    {'1': 'syncname', '3': 369920802, '4': 1, '5': 9, '10': 'syncname'},
    {
      '1': 'syncsource',
      '3': 286486824,
      '4': 1,
      '5': 11,
      '6': '.ssm.ResourceDataSyncSource',
      '10': 'syncsource'
    },
    {'1': 'synctype', '3': 122336091, '4': 1, '5': 9, '10': 'synctype'},
  ],
};

/// Descriptor for `UpdateResourceDataSyncRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResourceDataSyncRequestDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVSZXNvdXJjZURhdGFTeW5jUmVxdWVzdBIeCghzeW5jbmFtZRiilrKwASABKAlSCH'
        'N5bmNuYW1lEj8KCnN5bmNzb3VyY2UYqOLNiAEgASgLMhsuc3NtLlJlc291cmNlRGF0YVN5bmNT'
        'b3VyY2VSCnN5bmNzb3VyY2USHQoIc3luY3R5cGUY2+aqOiABKAlSCHN5bmN0eXBl');

@$core.Deprecated('Use updateResourceDataSyncResultDescriptor instead')
const UpdateResourceDataSyncResult$json = {
  '1': 'UpdateResourceDataSyncResult',
};

/// Descriptor for `UpdateResourceDataSyncResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateResourceDataSyncResultDescriptor =
    $convert.base64Decode('ChxVcGRhdGVSZXNvdXJjZURhdGFTeW5jUmVzdWx0');

@$core.Deprecated('Use updateServiceSettingRequestDescriptor instead')
const UpdateServiceSettingRequest$json = {
  '1': 'UpdateServiceSettingRequest',
  '2': [
    {'1': 'settingid', '3': 422452711, '4': 1, '5': 9, '10': 'settingid'},
    {'1': 'settingvalue', '3': 467016041, '4': 1, '5': 9, '10': 'settingvalue'},
  ],
};

/// Descriptor for `UpdateServiceSettingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateServiceSettingRequestDescriptor =
    $convert.base64Decode(
        'ChtVcGRhdGVTZXJ2aWNlU2V0dGluZ1JlcXVlc3QSIAoJc2V0dGluZ2lkGOe7uMkBIAEoCVIJc2'
        'V0dGluZ2lkEiYKDHNldHRpbmd2YWx1ZRjpstjeASABKAlSDHNldHRpbmd2YWx1ZQ==');

@$core.Deprecated('Use updateServiceSettingResultDescriptor instead')
const UpdateServiceSettingResult$json = {
  '1': 'UpdateServiceSettingResult',
};

/// Descriptor for `UpdateServiceSettingResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateServiceSettingResultDescriptor =
    $convert.base64Decode('ChpVcGRhdGVTZXJ2aWNlU2V0dGluZ1Jlc3VsdA==');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'reasoncode', '3': 209655857, '4': 1, '5': 9, '10': 'reasoncode'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor = $convert.base64Decode(
    'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2USIQoKcm'
    'Vhc29uY29kZRixsPxjIAEoCVIKcmVhc29uY29kZQ==');
