// This is a generated file - do not edit.
//
// Generated from ssm.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AccessRequestStatus extends $pb.ProtobufEnum {
  static const AccessRequestStatus ACCESS_REQUEST_STATUS_PENDING =
      AccessRequestStatus._(
          0, _omitEnumNames ? '' : 'ACCESS_REQUEST_STATUS_PENDING');
  static const AccessRequestStatus ACCESS_REQUEST_STATUS_REVOKED =
      AccessRequestStatus._(
          1, _omitEnumNames ? '' : 'ACCESS_REQUEST_STATUS_REVOKED');
  static const AccessRequestStatus ACCESS_REQUEST_STATUS_REJECTED =
      AccessRequestStatus._(
          2, _omitEnumNames ? '' : 'ACCESS_REQUEST_STATUS_REJECTED');
  static const AccessRequestStatus ACCESS_REQUEST_STATUS_APPROVED =
      AccessRequestStatus._(
          3, _omitEnumNames ? '' : 'ACCESS_REQUEST_STATUS_APPROVED');
  static const AccessRequestStatus ACCESS_REQUEST_STATUS_EXPIRED =
      AccessRequestStatus._(
          4, _omitEnumNames ? '' : 'ACCESS_REQUEST_STATUS_EXPIRED');

  static const $core.List<AccessRequestStatus> values = <AccessRequestStatus>[
    ACCESS_REQUEST_STATUS_PENDING,
    ACCESS_REQUEST_STATUS_REVOKED,
    ACCESS_REQUEST_STATUS_REJECTED,
    ACCESS_REQUEST_STATUS_APPROVED,
    ACCESS_REQUEST_STATUS_EXPIRED,
  ];

  static final $core.List<AccessRequestStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AccessRequestStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccessRequestStatus._(super.value, super.name);
}

class AccessType extends $pb.ProtobufEnum {
  static const AccessType ACCESS_TYPE_JUSTINTIME =
      AccessType._(0, _omitEnumNames ? '' : 'ACCESS_TYPE_JUSTINTIME');
  static const AccessType ACCESS_TYPE_STANDARD =
      AccessType._(1, _omitEnumNames ? '' : 'ACCESS_TYPE_STANDARD');

  static const $core.List<AccessType> values = <AccessType>[
    ACCESS_TYPE_JUSTINTIME,
    ACCESS_TYPE_STANDARD,
  ];

  static final $core.List<AccessType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AccessType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccessType._(super.value, super.name);
}

class AssociationComplianceSeverity extends $pb.ProtobufEnum {
  static const AssociationComplianceSeverity
      ASSOCIATION_COMPLIANCE_SEVERITY_MEDIUM = AssociationComplianceSeverity._(
          0, _omitEnumNames ? '' : 'ASSOCIATION_COMPLIANCE_SEVERITY_MEDIUM');
  static const AssociationComplianceSeverity
      ASSOCIATION_COMPLIANCE_SEVERITY_UNSPECIFIED =
      AssociationComplianceSeverity._(1,
          _omitEnumNames ? '' : 'ASSOCIATION_COMPLIANCE_SEVERITY_UNSPECIFIED');
  static const AssociationComplianceSeverity
      ASSOCIATION_COMPLIANCE_SEVERITY_CRITICAL =
      AssociationComplianceSeverity._(
          2, _omitEnumNames ? '' : 'ASSOCIATION_COMPLIANCE_SEVERITY_CRITICAL');
  static const AssociationComplianceSeverity
      ASSOCIATION_COMPLIANCE_SEVERITY_LOW = AssociationComplianceSeverity._(
          3, _omitEnumNames ? '' : 'ASSOCIATION_COMPLIANCE_SEVERITY_LOW');
  static const AssociationComplianceSeverity
      ASSOCIATION_COMPLIANCE_SEVERITY_HIGH = AssociationComplianceSeverity._(
          4, _omitEnumNames ? '' : 'ASSOCIATION_COMPLIANCE_SEVERITY_HIGH');

  static const $core.List<AssociationComplianceSeverity> values =
      <AssociationComplianceSeverity>[
    ASSOCIATION_COMPLIANCE_SEVERITY_MEDIUM,
    ASSOCIATION_COMPLIANCE_SEVERITY_UNSPECIFIED,
    ASSOCIATION_COMPLIANCE_SEVERITY_CRITICAL,
    ASSOCIATION_COMPLIANCE_SEVERITY_LOW,
    ASSOCIATION_COMPLIANCE_SEVERITY_HIGH,
  ];

  static final $core.List<AssociationComplianceSeverity?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static AssociationComplianceSeverity? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationComplianceSeverity._(super.value, super.name);
}

class AssociationExecutionFilterKey extends $pb.ProtobufEnum {
  static const AssociationExecutionFilterKey
      ASSOCIATION_EXECUTION_FILTER_KEY_CREATEDTIME =
      AssociationExecutionFilterKey._(0,
          _omitEnumNames ? '' : 'ASSOCIATION_EXECUTION_FILTER_KEY_CREATEDTIME');
  static const AssociationExecutionFilterKey
      ASSOCIATION_EXECUTION_FILTER_KEY_STATUS = AssociationExecutionFilterKey._(
          1, _omitEnumNames ? '' : 'ASSOCIATION_EXECUTION_FILTER_KEY_STATUS');
  static const AssociationExecutionFilterKey
      ASSOCIATION_EXECUTION_FILTER_KEY_EXECUTIONID =
      AssociationExecutionFilterKey._(2,
          _omitEnumNames ? '' : 'ASSOCIATION_EXECUTION_FILTER_KEY_EXECUTIONID');

  static const $core.List<AssociationExecutionFilterKey> values =
      <AssociationExecutionFilterKey>[
    ASSOCIATION_EXECUTION_FILTER_KEY_CREATEDTIME,
    ASSOCIATION_EXECUTION_FILTER_KEY_STATUS,
    ASSOCIATION_EXECUTION_FILTER_KEY_EXECUTIONID,
  ];

  static final $core.List<AssociationExecutionFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AssociationExecutionFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationExecutionFilterKey._(super.value, super.name);
}

class AssociationExecutionTargetsFilterKey extends $pb.ProtobufEnum {
  static const AssociationExecutionTargetsFilterKey
      ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_STATUS =
      AssociationExecutionTargetsFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_STATUS');
  static const AssociationExecutionTargetsFilterKey
      ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCETYPE =
      AssociationExecutionTargetsFilterKey._(
          1,
          _omitEnumNames
              ? ''
              : 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCETYPE');
  static const AssociationExecutionTargetsFilterKey
      ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCEID =
      AssociationExecutionTargetsFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCEID');

  static const $core.List<AssociationExecutionTargetsFilterKey> values =
      <AssociationExecutionTargetsFilterKey>[
    ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_STATUS,
    ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCETYPE,
    ASSOCIATION_EXECUTION_TARGETS_FILTER_KEY_RESOURCEID,
  ];

  static final $core.List<AssociationExecutionTargetsFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AssociationExecutionTargetsFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationExecutionTargetsFilterKey._(super.value, super.name);
}

class AssociationFilterKey extends $pb.ProtobufEnum {
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_STATUS =
      AssociationFilterKey._(
          0, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_STATUS');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_LASTEXECUTEDBEFORE =
      AssociationFilterKey._(
          1, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_LASTEXECUTEDBEFORE');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_INSTANCEID =
      AssociationFilterKey._(
          2, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_INSTANCEID');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_LASTEXECUTEDAFTER =
      AssociationFilterKey._(
          3, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_LASTEXECUTEDAFTER');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_RESOURCEGROUPNAME =
      AssociationFilterKey._(
          4, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_RESOURCEGROUPNAME');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_NAME =
      AssociationFilterKey._(
          5, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_NAME');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_ASSOCIATIONNAME =
      AssociationFilterKey._(
          6, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_ASSOCIATIONNAME');
  static const AssociationFilterKey ASSOCIATION_FILTER_KEY_ASSOCIATIONID =
      AssociationFilterKey._(
          7, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_KEY_ASSOCIATIONID');

  static const $core.List<AssociationFilterKey> values = <AssociationFilterKey>[
    ASSOCIATION_FILTER_KEY_STATUS,
    ASSOCIATION_FILTER_KEY_LASTEXECUTEDBEFORE,
    ASSOCIATION_FILTER_KEY_INSTANCEID,
    ASSOCIATION_FILTER_KEY_LASTEXECUTEDAFTER,
    ASSOCIATION_FILTER_KEY_RESOURCEGROUPNAME,
    ASSOCIATION_FILTER_KEY_NAME,
    ASSOCIATION_FILTER_KEY_ASSOCIATIONNAME,
    ASSOCIATION_FILTER_KEY_ASSOCIATIONID,
  ];

  static final $core.List<AssociationFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static AssociationFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationFilterKey._(super.value, super.name);
}

class AssociationFilterOperatorType extends $pb.ProtobufEnum {
  static const AssociationFilterOperatorType
      ASSOCIATION_FILTER_OPERATOR_TYPE_EQUAL = AssociationFilterOperatorType._(
          0, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_OPERATOR_TYPE_EQUAL');
  static const AssociationFilterOperatorType
      ASSOCIATION_FILTER_OPERATOR_TYPE_LESSTHAN =
      AssociationFilterOperatorType._(
          1, _omitEnumNames ? '' : 'ASSOCIATION_FILTER_OPERATOR_TYPE_LESSTHAN');
  static const AssociationFilterOperatorType
      ASSOCIATION_FILTER_OPERATOR_TYPE_GREATERTHAN =
      AssociationFilterOperatorType._(2,
          _omitEnumNames ? '' : 'ASSOCIATION_FILTER_OPERATOR_TYPE_GREATERTHAN');

  static const $core.List<AssociationFilterOperatorType> values =
      <AssociationFilterOperatorType>[
    ASSOCIATION_FILTER_OPERATOR_TYPE_EQUAL,
    ASSOCIATION_FILTER_OPERATOR_TYPE_LESSTHAN,
    ASSOCIATION_FILTER_OPERATOR_TYPE_GREATERTHAN,
  ];

  static final $core.List<AssociationFilterOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AssociationFilterOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationFilterOperatorType._(super.value, super.name);
}

class AssociationStatusName extends $pb.ProtobufEnum {
  static const AssociationStatusName ASSOCIATION_STATUS_NAME_FAILED =
      AssociationStatusName._(
          0, _omitEnumNames ? '' : 'ASSOCIATION_STATUS_NAME_FAILED');
  static const AssociationStatusName ASSOCIATION_STATUS_NAME_SUCCESS =
      AssociationStatusName._(
          1, _omitEnumNames ? '' : 'ASSOCIATION_STATUS_NAME_SUCCESS');
  static const AssociationStatusName ASSOCIATION_STATUS_NAME_PENDING =
      AssociationStatusName._(
          2, _omitEnumNames ? '' : 'ASSOCIATION_STATUS_NAME_PENDING');

  static const $core.List<AssociationStatusName> values =
      <AssociationStatusName>[
    ASSOCIATION_STATUS_NAME_FAILED,
    ASSOCIATION_STATUS_NAME_SUCCESS,
    ASSOCIATION_STATUS_NAME_PENDING,
  ];

  static final $core.List<AssociationStatusName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AssociationStatusName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationStatusName._(super.value, super.name);
}

class AssociationSyncCompliance extends $pb.ProtobufEnum {
  static const AssociationSyncCompliance ASSOCIATION_SYNC_COMPLIANCE_MANUAL =
      AssociationSyncCompliance._(
          0, _omitEnumNames ? '' : 'ASSOCIATION_SYNC_COMPLIANCE_MANUAL');
  static const AssociationSyncCompliance ASSOCIATION_SYNC_COMPLIANCE_AUTO =
      AssociationSyncCompliance._(
          1, _omitEnumNames ? '' : 'ASSOCIATION_SYNC_COMPLIANCE_AUTO');

  static const $core.List<AssociationSyncCompliance> values =
      <AssociationSyncCompliance>[
    ASSOCIATION_SYNC_COMPLIANCE_MANUAL,
    ASSOCIATION_SYNC_COMPLIANCE_AUTO,
  ];

  static final $core.List<AssociationSyncCompliance?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AssociationSyncCompliance? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AssociationSyncCompliance._(super.value, super.name);
}

class AttachmentHashType extends $pb.ProtobufEnum {
  static const AttachmentHashType ATTACHMENT_HASH_TYPE_SHA256 =
      AttachmentHashType._(
          0, _omitEnumNames ? '' : 'ATTACHMENT_HASH_TYPE_SHA256');

  static const $core.List<AttachmentHashType> values = <AttachmentHashType>[
    ATTACHMENT_HASH_TYPE_SHA256,
  ];

  static final $core.List<AttachmentHashType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static AttachmentHashType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttachmentHashType._(super.value, super.name);
}

class AttachmentsSourceKey extends $pb.ProtobufEnum {
  static const AttachmentsSourceKey ATTACHMENTS_SOURCE_KEY_ATTACHMENTREFERENCE =
      AttachmentsSourceKey._(0,
          _omitEnumNames ? '' : 'ATTACHMENTS_SOURCE_KEY_ATTACHMENTREFERENCE');
  static const AttachmentsSourceKey ATTACHMENTS_SOURCE_KEY_SOURCEURL =
      AttachmentsSourceKey._(
          1, _omitEnumNames ? '' : 'ATTACHMENTS_SOURCE_KEY_SOURCEURL');
  static const AttachmentsSourceKey ATTACHMENTS_SOURCE_KEY_S3FILEURL =
      AttachmentsSourceKey._(
          2, _omitEnumNames ? '' : 'ATTACHMENTS_SOURCE_KEY_S3FILEURL');

  static const $core.List<AttachmentsSourceKey> values = <AttachmentsSourceKey>[
    ATTACHMENTS_SOURCE_KEY_ATTACHMENTREFERENCE,
    ATTACHMENTS_SOURCE_KEY_SOURCEURL,
    ATTACHMENTS_SOURCE_KEY_S3FILEURL,
  ];

  static final $core.List<AttachmentsSourceKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AttachmentsSourceKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttachmentsSourceKey._(super.value, super.name);
}

class AutomationExecutionFilterKey extends $pb.ProtobufEnum {
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_AFTER =
      AutomationExecutionFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_AFTER');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_TAG_KEY = AutomationExecutionFilterKey._(
          1, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_FILTER_KEY_TAG_KEY');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_PARENT_EXECUTION_ID =
      AutomationExecutionFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_PARENT_EXECUTION_ID');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_TARGET_RESOURCE_GROUP =
      AutomationExecutionFilterKey._(
          3,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_TARGET_RESOURCE_GROUP');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_ID =
      AutomationExecutionFilterKey._(4,
          _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_ID');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_OPS_ITEM_ID =
      AutomationExecutionFilterKey._(5,
          _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_FILTER_KEY_OPS_ITEM_ID');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_CURRENT_ACTION =
      AutomationExecutionFilterKey._(
          6,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_CURRENT_ACTION');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_BEFORE =
      AutomationExecutionFilterKey._(
          7,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_BEFORE');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_STATUS =
      AutomationExecutionFilterKey._(
          8,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_STATUS');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_DOCUMENT_NAME_PREFIX =
      AutomationExecutionFilterKey._(
          9,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_DOCUMENT_NAME_PREFIX');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_TYPE =
      AutomationExecutionFilterKey._(
          10,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_TYPE');
  static const AutomationExecutionFilterKey
      AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_SUBTYPE =
      AutomationExecutionFilterKey._(
          11,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_SUBTYPE');

  static const $core.List<AutomationExecutionFilterKey> values =
      <AutomationExecutionFilterKey>[
    AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_AFTER,
    AUTOMATION_EXECUTION_FILTER_KEY_TAG_KEY,
    AUTOMATION_EXECUTION_FILTER_KEY_PARENT_EXECUTION_ID,
    AUTOMATION_EXECUTION_FILTER_KEY_TARGET_RESOURCE_GROUP,
    AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_ID,
    AUTOMATION_EXECUTION_FILTER_KEY_OPS_ITEM_ID,
    AUTOMATION_EXECUTION_FILTER_KEY_CURRENT_ACTION,
    AUTOMATION_EXECUTION_FILTER_KEY_START_TIME_BEFORE,
    AUTOMATION_EXECUTION_FILTER_KEY_EXECUTION_STATUS,
    AUTOMATION_EXECUTION_FILTER_KEY_DOCUMENT_NAME_PREFIX,
    AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_TYPE,
    AUTOMATION_EXECUTION_FILTER_KEY_AUTOMATION_SUBTYPE,
  ];

  static final $core.List<AutomationExecutionFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 11);
  static AutomationExecutionFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AutomationExecutionFilterKey._(super.value, super.name);
}

class AutomationExecutionStatus extends $pb.ProtobufEnum {
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE =
      AutomationExecutionStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_PENDING =
      AutomationExecutionStatus._(
          1, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_PENDING');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_SUCCESS =
      AutomationExecutionStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_SUCCESS');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_TIMEDOUT =
      AutomationExecutionStatus._(
          3, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_TIMEDOUT');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_PENDING_APPROVAL =
      AutomationExecutionStatus._(4,
          _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_PENDING_APPROVAL');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED =
      AutomationExecutionStatus._(
          5,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_SUCCESS =
      AutomationExecutionStatus._(
          6, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_SUCCESS');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_INPROGRESS = AutomationExecutionStatus._(
          7, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_INPROGRESS');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_RUNBOOK_INPROGRESS =
      AutomationExecutionStatus._(
          8,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_RUNBOOK_INPROGRESS');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_EXITED =
      AutomationExecutionStatus._(
          9, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_EXITED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_CANCELLED =
      AutomationExecutionStatus._(
          10, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_CANCELLED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_REJECTED =
      AutomationExecutionStatus._(
          11, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_REJECTED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_APPROVED =
      AutomationExecutionStatus._(
          12, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_APPROVED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_SCHEDULED =
      AutomationExecutionStatus._(
          13, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_SCHEDULED');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED =
      AutomationExecutionStatus._(
          14,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_WAITING =
      AutomationExecutionStatus._(
          15, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_WAITING');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_CANCELLING = AutomationExecutionStatus._(
          16, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_CANCELLING');
  static const AutomationExecutionStatus
      AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_FAILURE =
      AutomationExecutionStatus._(
          17,
          _omitEnumNames
              ? ''
              : 'AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_FAILURE');
  static const AutomationExecutionStatus AUTOMATION_EXECUTION_STATUS_FAILED =
      AutomationExecutionStatus._(
          18, _omitEnumNames ? '' : 'AUTOMATION_EXECUTION_STATUS_FAILED');

  static const $core.List<AutomationExecutionStatus> values =
      <AutomationExecutionStatus>[
    AUTOMATION_EXECUTION_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE,
    AUTOMATION_EXECUTION_STATUS_PENDING,
    AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_SUCCESS,
    AUTOMATION_EXECUTION_STATUS_TIMEDOUT,
    AUTOMATION_EXECUTION_STATUS_PENDING_APPROVAL,
    AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED,
    AUTOMATION_EXECUTION_STATUS_SUCCESS,
    AUTOMATION_EXECUTION_STATUS_INPROGRESS,
    AUTOMATION_EXECUTION_STATUS_RUNBOOK_INPROGRESS,
    AUTOMATION_EXECUTION_STATUS_EXITED,
    AUTOMATION_EXECUTION_STATUS_CANCELLED,
    AUTOMATION_EXECUTION_STATUS_REJECTED,
    AUTOMATION_EXECUTION_STATUS_APPROVED,
    AUTOMATION_EXECUTION_STATUS_SCHEDULED,
    AUTOMATION_EXECUTION_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED,
    AUTOMATION_EXECUTION_STATUS_WAITING,
    AUTOMATION_EXECUTION_STATUS_CANCELLING,
    AUTOMATION_EXECUTION_STATUS_COMPLETED_WITH_FAILURE,
    AUTOMATION_EXECUTION_STATUS_FAILED,
  ];

  static final $core.List<AutomationExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 18);
  static AutomationExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AutomationExecutionStatus._(super.value, super.name);
}

class AutomationSubtype extends $pb.ProtobufEnum {
  static const AutomationSubtype AUTOMATION_SUBTYPE_ACCESSREQUEST =
      AutomationSubtype._(
          0, _omitEnumNames ? '' : 'AUTOMATION_SUBTYPE_ACCESSREQUEST');
  static const AutomationSubtype AUTOMATION_SUBTYPE_CHANGEREQUEST =
      AutomationSubtype._(
          1, _omitEnumNames ? '' : 'AUTOMATION_SUBTYPE_CHANGEREQUEST');

  static const $core.List<AutomationSubtype> values = <AutomationSubtype>[
    AUTOMATION_SUBTYPE_ACCESSREQUEST,
    AUTOMATION_SUBTYPE_CHANGEREQUEST,
  ];

  static final $core.List<AutomationSubtype?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AutomationSubtype? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AutomationSubtype._(super.value, super.name);
}

class AutomationType extends $pb.ProtobufEnum {
  static const AutomationType AUTOMATION_TYPE_LOCAL =
      AutomationType._(0, _omitEnumNames ? '' : 'AUTOMATION_TYPE_LOCAL');
  static const AutomationType AUTOMATION_TYPE_CROSSACCOUNT =
      AutomationType._(1, _omitEnumNames ? '' : 'AUTOMATION_TYPE_CROSSACCOUNT');

  static const $core.List<AutomationType> values = <AutomationType>[
    AUTOMATION_TYPE_LOCAL,
    AUTOMATION_TYPE_CROSSACCOUNT,
  ];

  static final $core.List<AutomationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AutomationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AutomationType._(super.value, super.name);
}

class CalendarState extends $pb.ProtobufEnum {
  static const CalendarState CALENDAR_STATE_CLOSED =
      CalendarState._(0, _omitEnumNames ? '' : 'CALENDAR_STATE_CLOSED');
  static const CalendarState CALENDAR_STATE_OPEN =
      CalendarState._(1, _omitEnumNames ? '' : 'CALENDAR_STATE_OPEN');

  static const $core.List<CalendarState> values = <CalendarState>[
    CALENDAR_STATE_CLOSED,
    CALENDAR_STATE_OPEN,
  ];

  static final $core.List<CalendarState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CalendarState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CalendarState._(super.value, super.name);
}

class CommandFilterKey extends $pb.ProtobufEnum {
  static const CommandFilterKey COMMAND_FILTER_KEY_EXECUTION_STAGE =
      CommandFilterKey._(
          0, _omitEnumNames ? '' : 'COMMAND_FILTER_KEY_EXECUTION_STAGE');
  static const CommandFilterKey COMMAND_FILTER_KEY_DOCUMENT_NAME =
      CommandFilterKey._(
          1, _omitEnumNames ? '' : 'COMMAND_FILTER_KEY_DOCUMENT_NAME');
  static const CommandFilterKey COMMAND_FILTER_KEY_STATUS =
      CommandFilterKey._(2, _omitEnumNames ? '' : 'COMMAND_FILTER_KEY_STATUS');
  static const CommandFilterKey COMMAND_FILTER_KEY_INVOKED_BEFORE =
      CommandFilterKey._(
          3, _omitEnumNames ? '' : 'COMMAND_FILTER_KEY_INVOKED_BEFORE');
  static const CommandFilterKey COMMAND_FILTER_KEY_INVOKED_AFTER =
      CommandFilterKey._(
          4, _omitEnumNames ? '' : 'COMMAND_FILTER_KEY_INVOKED_AFTER');

  static const $core.List<CommandFilterKey> values = <CommandFilterKey>[
    COMMAND_FILTER_KEY_EXECUTION_STAGE,
    COMMAND_FILTER_KEY_DOCUMENT_NAME,
    COMMAND_FILTER_KEY_STATUS,
    COMMAND_FILTER_KEY_INVOKED_BEFORE,
    COMMAND_FILTER_KEY_INVOKED_AFTER,
  ];

  static final $core.List<CommandFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static CommandFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CommandFilterKey._(super.value, super.name);
}

class CommandInvocationStatus extends $pb.ProtobufEnum {
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_PENDING =
      CommandInvocationStatus._(
          0, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_PENDING');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_TIMED_OUT =
      CommandInvocationStatus._(
          1, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_TIMED_OUT');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_SUCCESS =
      CommandInvocationStatus._(
          2, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_SUCCESS');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_IN_PROGRESS =
      CommandInvocationStatus._(
          3, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_IN_PROGRESS');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_DELAYED =
      CommandInvocationStatus._(
          4, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_DELAYED');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_CANCELLED =
      CommandInvocationStatus._(
          5, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_CANCELLED');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_CANCELLING =
      CommandInvocationStatus._(
          6, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_CANCELLING');
  static const CommandInvocationStatus COMMAND_INVOCATION_STATUS_FAILED =
      CommandInvocationStatus._(
          7, _omitEnumNames ? '' : 'COMMAND_INVOCATION_STATUS_FAILED');

  static const $core.List<CommandInvocationStatus> values =
      <CommandInvocationStatus>[
    COMMAND_INVOCATION_STATUS_PENDING,
    COMMAND_INVOCATION_STATUS_TIMED_OUT,
    COMMAND_INVOCATION_STATUS_SUCCESS,
    COMMAND_INVOCATION_STATUS_IN_PROGRESS,
    COMMAND_INVOCATION_STATUS_DELAYED,
    COMMAND_INVOCATION_STATUS_CANCELLED,
    COMMAND_INVOCATION_STATUS_CANCELLING,
    COMMAND_INVOCATION_STATUS_FAILED,
  ];

  static final $core.List<CommandInvocationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static CommandInvocationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CommandInvocationStatus._(super.value, super.name);
}

class CommandPluginStatus extends $pb.ProtobufEnum {
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_PENDING =
      CommandPluginStatus._(
          0, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_PENDING');
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_TIMED_OUT =
      CommandPluginStatus._(
          1, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_TIMED_OUT');
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_SUCCESS =
      CommandPluginStatus._(
          2, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_SUCCESS');
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_IN_PROGRESS =
      CommandPluginStatus._(
          3, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_IN_PROGRESS');
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_CANCELLED =
      CommandPluginStatus._(
          4, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_CANCELLED');
  static const CommandPluginStatus COMMAND_PLUGIN_STATUS_FAILED =
      CommandPluginStatus._(
          5, _omitEnumNames ? '' : 'COMMAND_PLUGIN_STATUS_FAILED');

  static const $core.List<CommandPluginStatus> values = <CommandPluginStatus>[
    COMMAND_PLUGIN_STATUS_PENDING,
    COMMAND_PLUGIN_STATUS_TIMED_OUT,
    COMMAND_PLUGIN_STATUS_SUCCESS,
    COMMAND_PLUGIN_STATUS_IN_PROGRESS,
    COMMAND_PLUGIN_STATUS_CANCELLED,
    COMMAND_PLUGIN_STATUS_FAILED,
  ];

  static final $core.List<CommandPluginStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static CommandPluginStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CommandPluginStatus._(super.value, super.name);
}

class CommandStatus extends $pb.ProtobufEnum {
  static const CommandStatus COMMAND_STATUS_PENDING =
      CommandStatus._(0, _omitEnumNames ? '' : 'COMMAND_STATUS_PENDING');
  static const CommandStatus COMMAND_STATUS_TIMED_OUT =
      CommandStatus._(1, _omitEnumNames ? '' : 'COMMAND_STATUS_TIMED_OUT');
  static const CommandStatus COMMAND_STATUS_SUCCESS =
      CommandStatus._(2, _omitEnumNames ? '' : 'COMMAND_STATUS_SUCCESS');
  static const CommandStatus COMMAND_STATUS_IN_PROGRESS =
      CommandStatus._(3, _omitEnumNames ? '' : 'COMMAND_STATUS_IN_PROGRESS');
  static const CommandStatus COMMAND_STATUS_CANCELLED =
      CommandStatus._(4, _omitEnumNames ? '' : 'COMMAND_STATUS_CANCELLED');
  static const CommandStatus COMMAND_STATUS_CANCELLING =
      CommandStatus._(5, _omitEnumNames ? '' : 'COMMAND_STATUS_CANCELLING');
  static const CommandStatus COMMAND_STATUS_FAILED =
      CommandStatus._(6, _omitEnumNames ? '' : 'COMMAND_STATUS_FAILED');

  static const $core.List<CommandStatus> values = <CommandStatus>[
    COMMAND_STATUS_PENDING,
    COMMAND_STATUS_TIMED_OUT,
    COMMAND_STATUS_SUCCESS,
    COMMAND_STATUS_IN_PROGRESS,
    COMMAND_STATUS_CANCELLED,
    COMMAND_STATUS_CANCELLING,
    COMMAND_STATUS_FAILED,
  ];

  static final $core.List<CommandStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static CommandStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CommandStatus._(super.value, super.name);
}

class ComplianceQueryOperatorType extends $pb.ProtobufEnum {
  static const ComplianceQueryOperatorType
      COMPLIANCE_QUERY_OPERATOR_TYPE_NOTEQUAL = ComplianceQueryOperatorType._(
          0, _omitEnumNames ? '' : 'COMPLIANCE_QUERY_OPERATOR_TYPE_NOTEQUAL');
  static const ComplianceQueryOperatorType
      COMPLIANCE_QUERY_OPERATOR_TYPE_EQUAL = ComplianceQueryOperatorType._(
          1, _omitEnumNames ? '' : 'COMPLIANCE_QUERY_OPERATOR_TYPE_EQUAL');
  static const ComplianceQueryOperatorType
      COMPLIANCE_QUERY_OPERATOR_TYPE_LESSTHAN = ComplianceQueryOperatorType._(
          2, _omitEnumNames ? '' : 'COMPLIANCE_QUERY_OPERATOR_TYPE_LESSTHAN');
  static const ComplianceQueryOperatorType
      COMPLIANCE_QUERY_OPERATOR_TYPE_GREATERTHAN =
      ComplianceQueryOperatorType._(3,
          _omitEnumNames ? '' : 'COMPLIANCE_QUERY_OPERATOR_TYPE_GREATERTHAN');
  static const ComplianceQueryOperatorType
      COMPLIANCE_QUERY_OPERATOR_TYPE_BEGINWITH = ComplianceQueryOperatorType._(
          4, _omitEnumNames ? '' : 'COMPLIANCE_QUERY_OPERATOR_TYPE_BEGINWITH');

  static const $core.List<ComplianceQueryOperatorType> values =
      <ComplianceQueryOperatorType>[
    COMPLIANCE_QUERY_OPERATOR_TYPE_NOTEQUAL,
    COMPLIANCE_QUERY_OPERATOR_TYPE_EQUAL,
    COMPLIANCE_QUERY_OPERATOR_TYPE_LESSTHAN,
    COMPLIANCE_QUERY_OPERATOR_TYPE_GREATERTHAN,
    COMPLIANCE_QUERY_OPERATOR_TYPE_BEGINWITH,
  ];

  static final $core.List<ComplianceQueryOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ComplianceQueryOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComplianceQueryOperatorType._(super.value, super.name);
}

class ComplianceSeverity extends $pb.ProtobufEnum {
  static const ComplianceSeverity COMPLIANCE_SEVERITY_INFORMATIONAL =
      ComplianceSeverity._(
          0, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_INFORMATIONAL');
  static const ComplianceSeverity COMPLIANCE_SEVERITY_MEDIUM =
      ComplianceSeverity._(
          1, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_MEDIUM');
  static const ComplianceSeverity COMPLIANCE_SEVERITY_UNSPECIFIED =
      ComplianceSeverity._(
          2, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_UNSPECIFIED');
  static const ComplianceSeverity COMPLIANCE_SEVERITY_CRITICAL =
      ComplianceSeverity._(
          3, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_CRITICAL');
  static const ComplianceSeverity COMPLIANCE_SEVERITY_LOW =
      ComplianceSeverity._(4, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_LOW');
  static const ComplianceSeverity COMPLIANCE_SEVERITY_HIGH =
      ComplianceSeverity._(5, _omitEnumNames ? '' : 'COMPLIANCE_SEVERITY_HIGH');

  static const $core.List<ComplianceSeverity> values = <ComplianceSeverity>[
    COMPLIANCE_SEVERITY_INFORMATIONAL,
    COMPLIANCE_SEVERITY_MEDIUM,
    COMPLIANCE_SEVERITY_UNSPECIFIED,
    COMPLIANCE_SEVERITY_CRITICAL,
    COMPLIANCE_SEVERITY_LOW,
    COMPLIANCE_SEVERITY_HIGH,
  ];

  static final $core.List<ComplianceSeverity?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ComplianceSeverity? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComplianceSeverity._(super.value, super.name);
}

class ComplianceStatus extends $pb.ProtobufEnum {
  static const ComplianceStatus COMPLIANCE_STATUS_COMPLIANT =
      ComplianceStatus._(
          0, _omitEnumNames ? '' : 'COMPLIANCE_STATUS_COMPLIANT');
  static const ComplianceStatus COMPLIANCE_STATUS_NONCOMPLIANT =
      ComplianceStatus._(
          1, _omitEnumNames ? '' : 'COMPLIANCE_STATUS_NONCOMPLIANT');

  static const $core.List<ComplianceStatus> values = <ComplianceStatus>[
    COMPLIANCE_STATUS_COMPLIANT,
    COMPLIANCE_STATUS_NONCOMPLIANT,
  ];

  static final $core.List<ComplianceStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ComplianceStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComplianceStatus._(super.value, super.name);
}

class ComplianceUploadType extends $pb.ProtobufEnum {
  static const ComplianceUploadType COMPLIANCE_UPLOAD_TYPE_PARTIAL =
      ComplianceUploadType._(
          0, _omitEnumNames ? '' : 'COMPLIANCE_UPLOAD_TYPE_PARTIAL');
  static const ComplianceUploadType COMPLIANCE_UPLOAD_TYPE_COMPLETE =
      ComplianceUploadType._(
          1, _omitEnumNames ? '' : 'COMPLIANCE_UPLOAD_TYPE_COMPLETE');

  static const $core.List<ComplianceUploadType> values = <ComplianceUploadType>[
    COMPLIANCE_UPLOAD_TYPE_PARTIAL,
    COMPLIANCE_UPLOAD_TYPE_COMPLETE,
  ];

  static final $core.List<ComplianceUploadType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ComplianceUploadType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ComplianceUploadType._(super.value, super.name);
}

class ConnectionStatus extends $pb.ProtobufEnum {
  static const ConnectionStatus CONNECTION_STATUS_NOT_CONNECTED =
      ConnectionStatus._(
          0, _omitEnumNames ? '' : 'CONNECTION_STATUS_NOT_CONNECTED');
  static const ConnectionStatus CONNECTION_STATUS_CONNECTED =
      ConnectionStatus._(
          1, _omitEnumNames ? '' : 'CONNECTION_STATUS_CONNECTED');

  static const $core.List<ConnectionStatus> values = <ConnectionStatus>[
    CONNECTION_STATUS_NOT_CONNECTED,
    CONNECTION_STATUS_CONNECTED,
  ];

  static final $core.List<ConnectionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ConnectionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionStatus._(super.value, super.name);
}

class DescribeActivationsFilterKeys extends $pb.ProtobufEnum {
  static const DescribeActivationsFilterKeys
      DESCRIBE_ACTIVATIONS_FILTER_KEYS_ACTIVATION_IDS =
      DescribeActivationsFilterKeys._(
          0,
          _omitEnumNames
              ? ''
              : 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_ACTIVATION_IDS');
  static const DescribeActivationsFilterKeys
      DESCRIBE_ACTIVATIONS_FILTER_KEYS_IAM_ROLE =
      DescribeActivationsFilterKeys._(
          1, _omitEnumNames ? '' : 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_IAM_ROLE');
  static const DescribeActivationsFilterKeys
      DESCRIBE_ACTIVATIONS_FILTER_KEYS_DEFAULT_INSTANCE_NAME =
      DescribeActivationsFilterKeys._(
          2,
          _omitEnumNames
              ? ''
              : 'DESCRIBE_ACTIVATIONS_FILTER_KEYS_DEFAULT_INSTANCE_NAME');

  static const $core.List<DescribeActivationsFilterKeys> values =
      <DescribeActivationsFilterKeys>[
    DESCRIBE_ACTIVATIONS_FILTER_KEYS_ACTIVATION_IDS,
    DESCRIBE_ACTIVATIONS_FILTER_KEYS_IAM_ROLE,
    DESCRIBE_ACTIVATIONS_FILTER_KEYS_DEFAULT_INSTANCE_NAME,
  ];

  static final $core.List<DescribeActivationsFilterKeys?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DescribeActivationsFilterKeys? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DescribeActivationsFilterKeys._(super.value, super.name);
}

class DocumentFilterKey extends $pb.ProtobufEnum {
  static const DocumentFilterKey DOCUMENT_FILTER_KEY_OWNER =
      DocumentFilterKey._(0, _omitEnumNames ? '' : 'DOCUMENT_FILTER_KEY_OWNER');
  static const DocumentFilterKey DOCUMENT_FILTER_KEY_PLATFORMTYPES =
      DocumentFilterKey._(
          1, _omitEnumNames ? '' : 'DOCUMENT_FILTER_KEY_PLATFORMTYPES');
  static const DocumentFilterKey DOCUMENT_FILTER_KEY_NAME =
      DocumentFilterKey._(2, _omitEnumNames ? '' : 'DOCUMENT_FILTER_KEY_NAME');
  static const DocumentFilterKey DOCUMENT_FILTER_KEY_DOCUMENTTYPE =
      DocumentFilterKey._(
          3, _omitEnumNames ? '' : 'DOCUMENT_FILTER_KEY_DOCUMENTTYPE');

  static const $core.List<DocumentFilterKey> values = <DocumentFilterKey>[
    DOCUMENT_FILTER_KEY_OWNER,
    DOCUMENT_FILTER_KEY_PLATFORMTYPES,
    DOCUMENT_FILTER_KEY_NAME,
    DOCUMENT_FILTER_KEY_DOCUMENTTYPE,
  ];

  static final $core.List<DocumentFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DocumentFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentFilterKey._(super.value, super.name);
}

class DocumentFormat extends $pb.ProtobufEnum {
  static const DocumentFormat DOCUMENT_FORMAT_JSON =
      DocumentFormat._(0, _omitEnumNames ? '' : 'DOCUMENT_FORMAT_JSON');
  static const DocumentFormat DOCUMENT_FORMAT_YAML =
      DocumentFormat._(1, _omitEnumNames ? '' : 'DOCUMENT_FORMAT_YAML');
  static const DocumentFormat DOCUMENT_FORMAT_TEXT =
      DocumentFormat._(2, _omitEnumNames ? '' : 'DOCUMENT_FORMAT_TEXT');

  static const $core.List<DocumentFormat> values = <DocumentFormat>[
    DOCUMENT_FORMAT_JSON,
    DOCUMENT_FORMAT_YAML,
    DOCUMENT_FORMAT_TEXT,
  ];

  static final $core.List<DocumentFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DocumentFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentFormat._(super.value, super.name);
}

class DocumentHashType extends $pb.ProtobufEnum {
  static const DocumentHashType DOCUMENT_HASH_TYPE_SHA256 =
      DocumentHashType._(0, _omitEnumNames ? '' : 'DOCUMENT_HASH_TYPE_SHA256');
  static const DocumentHashType DOCUMENT_HASH_TYPE_SHA1 =
      DocumentHashType._(1, _omitEnumNames ? '' : 'DOCUMENT_HASH_TYPE_SHA1');

  static const $core.List<DocumentHashType> values = <DocumentHashType>[
    DOCUMENT_HASH_TYPE_SHA256,
    DOCUMENT_HASH_TYPE_SHA1,
  ];

  static final $core.List<DocumentHashType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DocumentHashType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentHashType._(super.value, super.name);
}

class DocumentMetadataEnum extends $pb.ProtobufEnum {
  static const DocumentMetadataEnum DOCUMENT_METADATA_ENUM_DOCUMENTREVIEWS =
      DocumentMetadataEnum._(
          0, _omitEnumNames ? '' : 'DOCUMENT_METADATA_ENUM_DOCUMENTREVIEWS');

  static const $core.List<DocumentMetadataEnum> values = <DocumentMetadataEnum>[
    DOCUMENT_METADATA_ENUM_DOCUMENTREVIEWS,
  ];

  static final $core.List<DocumentMetadataEnum?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static DocumentMetadataEnum? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentMetadataEnum._(super.value, super.name);
}

class DocumentParameterType extends $pb.ProtobufEnum {
  static const DocumentParameterType DOCUMENT_PARAMETER_TYPE_STRINGLIST =
      DocumentParameterType._(
          0, _omitEnumNames ? '' : 'DOCUMENT_PARAMETER_TYPE_STRINGLIST');
  static const DocumentParameterType DOCUMENT_PARAMETER_TYPE_STRING =
      DocumentParameterType._(
          1, _omitEnumNames ? '' : 'DOCUMENT_PARAMETER_TYPE_STRING');

  static const $core.List<DocumentParameterType> values =
      <DocumentParameterType>[
    DOCUMENT_PARAMETER_TYPE_STRINGLIST,
    DOCUMENT_PARAMETER_TYPE_STRING,
  ];

  static final $core.List<DocumentParameterType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DocumentParameterType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentParameterType._(super.value, super.name);
}

class DocumentPermissionType extends $pb.ProtobufEnum {
  static const DocumentPermissionType DOCUMENT_PERMISSION_TYPE_SHARE =
      DocumentPermissionType._(
          0, _omitEnumNames ? '' : 'DOCUMENT_PERMISSION_TYPE_SHARE');

  static const $core.List<DocumentPermissionType> values =
      <DocumentPermissionType>[
    DOCUMENT_PERMISSION_TYPE_SHARE,
  ];

  static final $core.List<DocumentPermissionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static DocumentPermissionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentPermissionType._(super.value, super.name);
}

class DocumentReviewAction extends $pb.ProtobufEnum {
  static const DocumentReviewAction DOCUMENT_REVIEW_ACTION_SENDFORREVIEW =
      DocumentReviewAction._(
          0, _omitEnumNames ? '' : 'DOCUMENT_REVIEW_ACTION_SENDFORREVIEW');
  static const DocumentReviewAction DOCUMENT_REVIEW_ACTION_REJECT =
      DocumentReviewAction._(
          1, _omitEnumNames ? '' : 'DOCUMENT_REVIEW_ACTION_REJECT');
  static const DocumentReviewAction DOCUMENT_REVIEW_ACTION_APPROVE =
      DocumentReviewAction._(
          2, _omitEnumNames ? '' : 'DOCUMENT_REVIEW_ACTION_APPROVE');
  static const DocumentReviewAction DOCUMENT_REVIEW_ACTION_UPDATEREVIEW =
      DocumentReviewAction._(
          3, _omitEnumNames ? '' : 'DOCUMENT_REVIEW_ACTION_UPDATEREVIEW');

  static const $core.List<DocumentReviewAction> values = <DocumentReviewAction>[
    DOCUMENT_REVIEW_ACTION_SENDFORREVIEW,
    DOCUMENT_REVIEW_ACTION_REJECT,
    DOCUMENT_REVIEW_ACTION_APPROVE,
    DOCUMENT_REVIEW_ACTION_UPDATEREVIEW,
  ];

  static final $core.List<DocumentReviewAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DocumentReviewAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentReviewAction._(super.value, super.name);
}

class DocumentReviewCommentType extends $pb.ProtobufEnum {
  static const DocumentReviewCommentType DOCUMENT_REVIEW_COMMENT_TYPE_COMMENT =
      DocumentReviewCommentType._(
          0, _omitEnumNames ? '' : 'DOCUMENT_REVIEW_COMMENT_TYPE_COMMENT');

  static const $core.List<DocumentReviewCommentType> values =
      <DocumentReviewCommentType>[
    DOCUMENT_REVIEW_COMMENT_TYPE_COMMENT,
  ];

  static final $core.List<DocumentReviewCommentType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static DocumentReviewCommentType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentReviewCommentType._(super.value, super.name);
}

class DocumentStatus extends $pb.ProtobufEnum {
  static const DocumentStatus DOCUMENT_STATUS_ACTIVE =
      DocumentStatus._(0, _omitEnumNames ? '' : 'DOCUMENT_STATUS_ACTIVE');
  static const DocumentStatus DOCUMENT_STATUS_UPDATING =
      DocumentStatus._(1, _omitEnumNames ? '' : 'DOCUMENT_STATUS_UPDATING');
  static const DocumentStatus DOCUMENT_STATUS_FAILED =
      DocumentStatus._(2, _omitEnumNames ? '' : 'DOCUMENT_STATUS_FAILED');
  static const DocumentStatus DOCUMENT_STATUS_CREATING =
      DocumentStatus._(3, _omitEnumNames ? '' : 'DOCUMENT_STATUS_CREATING');
  static const DocumentStatus DOCUMENT_STATUS_DELETING =
      DocumentStatus._(4, _omitEnumNames ? '' : 'DOCUMENT_STATUS_DELETING');

  static const $core.List<DocumentStatus> values = <DocumentStatus>[
    DOCUMENT_STATUS_ACTIVE,
    DOCUMENT_STATUS_UPDATING,
    DOCUMENT_STATUS_FAILED,
    DOCUMENT_STATUS_CREATING,
    DOCUMENT_STATUS_DELETING,
  ];

  static final $core.List<DocumentStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static DocumentStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentStatus._(super.value, super.name);
}

class DocumentType extends $pb.ProtobufEnum {
  static const DocumentType DOCUMENT_TYPE_PROBLEMANALYSISTEMPLATE =
      DocumentType._(
          0, _omitEnumNames ? '' : 'DOCUMENT_TYPE_PROBLEMANALYSISTEMPLATE');
  static const DocumentType DOCUMENT_TYPE_AUTOAPPROVALPOLICY = DocumentType._(
      1, _omitEnumNames ? '' : 'DOCUMENT_TYPE_AUTOAPPROVALPOLICY');
  static const DocumentType DOCUMENT_TYPE_APPLICATIONCONFIGURATION =
      DocumentType._(
          2, _omitEnumNames ? '' : 'DOCUMENT_TYPE_APPLICATIONCONFIGURATION');
  static const DocumentType DOCUMENT_TYPE_MANUALAPPROVALPOLICY = DocumentType._(
      3, _omitEnumNames ? '' : 'DOCUMENT_TYPE_MANUALAPPROVALPOLICY');
  static const DocumentType DOCUMENT_TYPE_DEPLOYMENTSTRATEGY = DocumentType._(
      4, _omitEnumNames ? '' : 'DOCUMENT_TYPE_DEPLOYMENTSTRATEGY');
  static const DocumentType DOCUMENT_TYPE_PROBLEMANALYSIS =
      DocumentType._(5, _omitEnumNames ? '' : 'DOCUMENT_TYPE_PROBLEMANALYSIS');
  static const DocumentType DOCUMENT_TYPE_AUTOMATION =
      DocumentType._(6, _omitEnumNames ? '' : 'DOCUMENT_TYPE_AUTOMATION');
  static const DocumentType DOCUMENT_TYPE_CHANGECALENDAR =
      DocumentType._(7, _omitEnumNames ? '' : 'DOCUMENT_TYPE_CHANGECALENDAR');
  static const DocumentType DOCUMENT_TYPE_APPLICATIONCONFIGURATIONSCHEMA =
      DocumentType._(8,
          _omitEnumNames ? '' : 'DOCUMENT_TYPE_APPLICATIONCONFIGURATIONSCHEMA');
  static const DocumentType DOCUMENT_TYPE_QUICKSETUP =
      DocumentType._(9, _omitEnumNames ? '' : 'DOCUMENT_TYPE_QUICKSETUP');
  static const DocumentType DOCUMENT_TYPE_CONFORMANCEPACKTEMPLATE =
      DocumentType._(
          10, _omitEnumNames ? '' : 'DOCUMENT_TYPE_CONFORMANCEPACKTEMPLATE');
  static const DocumentType DOCUMENT_TYPE_POLICY =
      DocumentType._(11, _omitEnumNames ? '' : 'DOCUMENT_TYPE_POLICY');
  static const DocumentType DOCUMENT_TYPE_PACKAGE =
      DocumentType._(12, _omitEnumNames ? '' : 'DOCUMENT_TYPE_PACKAGE');
  static const DocumentType DOCUMENT_TYPE_SESSION =
      DocumentType._(13, _omitEnumNames ? '' : 'DOCUMENT_TYPE_SESSION');
  static const DocumentType DOCUMENT_TYPE_COMMAND =
      DocumentType._(14, _omitEnumNames ? '' : 'DOCUMENT_TYPE_COMMAND');
  static const DocumentType DOCUMENT_TYPE_CLOUDFORMATION =
      DocumentType._(15, _omitEnumNames ? '' : 'DOCUMENT_TYPE_CLOUDFORMATION');
  static const DocumentType DOCUMENT_TYPE_CHANGETEMPLATE =
      DocumentType._(16, _omitEnumNames ? '' : 'DOCUMENT_TYPE_CHANGETEMPLATE');

  static const $core.List<DocumentType> values = <DocumentType>[
    DOCUMENT_TYPE_PROBLEMANALYSISTEMPLATE,
    DOCUMENT_TYPE_AUTOAPPROVALPOLICY,
    DOCUMENT_TYPE_APPLICATIONCONFIGURATION,
    DOCUMENT_TYPE_MANUALAPPROVALPOLICY,
    DOCUMENT_TYPE_DEPLOYMENTSTRATEGY,
    DOCUMENT_TYPE_PROBLEMANALYSIS,
    DOCUMENT_TYPE_AUTOMATION,
    DOCUMENT_TYPE_CHANGECALENDAR,
    DOCUMENT_TYPE_APPLICATIONCONFIGURATIONSCHEMA,
    DOCUMENT_TYPE_QUICKSETUP,
    DOCUMENT_TYPE_CONFORMANCEPACKTEMPLATE,
    DOCUMENT_TYPE_POLICY,
    DOCUMENT_TYPE_PACKAGE,
    DOCUMENT_TYPE_SESSION,
    DOCUMENT_TYPE_COMMAND,
    DOCUMENT_TYPE_CLOUDFORMATION,
    DOCUMENT_TYPE_CHANGETEMPLATE,
  ];

  static final $core.List<DocumentType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 16);
  static DocumentType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DocumentType._(super.value, super.name);
}

class ExecutionMode extends $pb.ProtobufEnum {
  static const ExecutionMode EXECUTION_MODE_AUTO =
      ExecutionMode._(0, _omitEnumNames ? '' : 'EXECUTION_MODE_AUTO');
  static const ExecutionMode EXECUTION_MODE_INTERACTIVE =
      ExecutionMode._(1, _omitEnumNames ? '' : 'EXECUTION_MODE_INTERACTIVE');

  static const $core.List<ExecutionMode> values = <ExecutionMode>[
    EXECUTION_MODE_AUTO,
    EXECUTION_MODE_INTERACTIVE,
  ];

  static final $core.List<ExecutionMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExecutionMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionMode._(super.value, super.name);
}

class ExecutionPreviewStatus extends $pb.ProtobufEnum {
  static const ExecutionPreviewStatus EXECUTION_PREVIEW_STATUS_PENDING =
      ExecutionPreviewStatus._(
          0, _omitEnumNames ? '' : 'EXECUTION_PREVIEW_STATUS_PENDING');
  static const ExecutionPreviewStatus EXECUTION_PREVIEW_STATUS_SUCCESS =
      ExecutionPreviewStatus._(
          1, _omitEnumNames ? '' : 'EXECUTION_PREVIEW_STATUS_SUCCESS');
  static const ExecutionPreviewStatus EXECUTION_PREVIEW_STATUS_IN_PROGRESS =
      ExecutionPreviewStatus._(
          2, _omitEnumNames ? '' : 'EXECUTION_PREVIEW_STATUS_IN_PROGRESS');
  static const ExecutionPreviewStatus EXECUTION_PREVIEW_STATUS_FAILED =
      ExecutionPreviewStatus._(
          3, _omitEnumNames ? '' : 'EXECUTION_PREVIEW_STATUS_FAILED');

  static const $core.List<ExecutionPreviewStatus> values =
      <ExecutionPreviewStatus>[
    EXECUTION_PREVIEW_STATUS_PENDING,
    EXECUTION_PREVIEW_STATUS_SUCCESS,
    EXECUTION_PREVIEW_STATUS_IN_PROGRESS,
    EXECUTION_PREVIEW_STATUS_FAILED,
  ];

  static final $core.List<ExecutionPreviewStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ExecutionPreviewStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExecutionPreviewStatus._(super.value, super.name);
}

class ExternalAlarmState extends $pb.ProtobufEnum {
  static const ExternalAlarmState EXTERNAL_ALARM_STATE_ALARM =
      ExternalAlarmState._(
          0, _omitEnumNames ? '' : 'EXTERNAL_ALARM_STATE_ALARM');
  static const ExternalAlarmState EXTERNAL_ALARM_STATE_UNKNOWN =
      ExternalAlarmState._(
          1, _omitEnumNames ? '' : 'EXTERNAL_ALARM_STATE_UNKNOWN');

  static const $core.List<ExternalAlarmState> values = <ExternalAlarmState>[
    EXTERNAL_ALARM_STATE_ALARM,
    EXTERNAL_ALARM_STATE_UNKNOWN,
  ];

  static final $core.List<ExternalAlarmState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExternalAlarmState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExternalAlarmState._(super.value, super.name);
}

class Fault extends $pb.ProtobufEnum {
  static const Fault FAULT_CLIENT =
      Fault._(0, _omitEnumNames ? '' : 'FAULT_CLIENT');
  static const Fault FAULT_UNKNOWN =
      Fault._(1, _omitEnumNames ? '' : 'FAULT_UNKNOWN');
  static const Fault FAULT_SERVER =
      Fault._(2, _omitEnumNames ? '' : 'FAULT_SERVER');

  static const $core.List<Fault> values = <Fault>[
    FAULT_CLIENT,
    FAULT_UNKNOWN,
    FAULT_SERVER,
  ];

  static final $core.List<Fault?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Fault? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Fault._(super.value, super.name);
}

class ImpactType extends $pb.ProtobufEnum {
  static const ImpactType IMPACT_TYPE_UNDETERMINED =
      ImpactType._(0, _omitEnumNames ? '' : 'IMPACT_TYPE_UNDETERMINED');
  static const ImpactType IMPACT_TYPE_NON_MUTATING =
      ImpactType._(1, _omitEnumNames ? '' : 'IMPACT_TYPE_NON_MUTATING');
  static const ImpactType IMPACT_TYPE_MUTATING =
      ImpactType._(2, _omitEnumNames ? '' : 'IMPACT_TYPE_MUTATING');

  static const $core.List<ImpactType> values = <ImpactType>[
    IMPACT_TYPE_UNDETERMINED,
    IMPACT_TYPE_NON_MUTATING,
    IMPACT_TYPE_MUTATING,
  ];

  static final $core.List<ImpactType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ImpactType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImpactType._(super.value, super.name);
}

class InstanceInformationFilterKey extends $pb.ProtobufEnum {
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_ASSOCIATION_STATUS =
      InstanceInformationFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'INSTANCE_INFORMATION_FILTER_KEY_ASSOCIATION_STATUS');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_INSTANCE_IDS =
      InstanceInformationFilterKey._(1,
          _omitEnumNames ? '' : 'INSTANCE_INFORMATION_FILTER_KEY_INSTANCE_IDS');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_ACTIVATION_IDS =
      InstanceInformationFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'INSTANCE_INFORMATION_FILTER_KEY_ACTIVATION_IDS');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_IAM_ROLE = InstanceInformationFilterKey._(
          3, _omitEnumNames ? '' : 'INSTANCE_INFORMATION_FILTER_KEY_IAM_ROLE');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_AGENT_VERSION =
      InstanceInformationFilterKey._(
          4,
          _omitEnumNames
              ? ''
              : 'INSTANCE_INFORMATION_FILTER_KEY_AGENT_VERSION');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_PING_STATUS =
      InstanceInformationFilterKey._(5,
          _omitEnumNames ? '' : 'INSTANCE_INFORMATION_FILTER_KEY_PING_STATUS');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_PLATFORM_TYPES =
      InstanceInformationFilterKey._(
          6,
          _omitEnumNames
              ? ''
              : 'INSTANCE_INFORMATION_FILTER_KEY_PLATFORM_TYPES');
  static const InstanceInformationFilterKey
      INSTANCE_INFORMATION_FILTER_KEY_RESOURCE_TYPE =
      InstanceInformationFilterKey._(
          7,
          _omitEnumNames
              ? ''
              : 'INSTANCE_INFORMATION_FILTER_KEY_RESOURCE_TYPE');

  static const $core.List<InstanceInformationFilterKey> values =
      <InstanceInformationFilterKey>[
    INSTANCE_INFORMATION_FILTER_KEY_ASSOCIATION_STATUS,
    INSTANCE_INFORMATION_FILTER_KEY_INSTANCE_IDS,
    INSTANCE_INFORMATION_FILTER_KEY_ACTIVATION_IDS,
    INSTANCE_INFORMATION_FILTER_KEY_IAM_ROLE,
    INSTANCE_INFORMATION_FILTER_KEY_AGENT_VERSION,
    INSTANCE_INFORMATION_FILTER_KEY_PING_STATUS,
    INSTANCE_INFORMATION_FILTER_KEY_PLATFORM_TYPES,
    INSTANCE_INFORMATION_FILTER_KEY_RESOURCE_TYPE,
  ];

  static final $core.List<InstanceInformationFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static InstanceInformationFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InstanceInformationFilterKey._(super.value, super.name);
}

class InstancePatchStateOperatorType extends $pb.ProtobufEnum {
  static const InstancePatchStateOperatorType
      INSTANCE_PATCH_STATE_OPERATOR_TYPE_LESS_THAN =
      InstancePatchStateOperatorType._(0,
          _omitEnumNames ? '' : 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_LESS_THAN');
  static const InstancePatchStateOperatorType
      INSTANCE_PATCH_STATE_OPERATOR_TYPE_GREATER_THAN =
      InstancePatchStateOperatorType._(
          1,
          _omitEnumNames
              ? ''
              : 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_GREATER_THAN');
  static const InstancePatchStateOperatorType
      INSTANCE_PATCH_STATE_OPERATOR_TYPE_NOT_EQUAL =
      InstancePatchStateOperatorType._(2,
          _omitEnumNames ? '' : 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_NOT_EQUAL');
  static const InstancePatchStateOperatorType
      INSTANCE_PATCH_STATE_OPERATOR_TYPE_EQUAL =
      InstancePatchStateOperatorType._(
          3, _omitEnumNames ? '' : 'INSTANCE_PATCH_STATE_OPERATOR_TYPE_EQUAL');

  static const $core.List<InstancePatchStateOperatorType> values =
      <InstancePatchStateOperatorType>[
    INSTANCE_PATCH_STATE_OPERATOR_TYPE_LESS_THAN,
    INSTANCE_PATCH_STATE_OPERATOR_TYPE_GREATER_THAN,
    INSTANCE_PATCH_STATE_OPERATOR_TYPE_NOT_EQUAL,
    INSTANCE_PATCH_STATE_OPERATOR_TYPE_EQUAL,
  ];

  static final $core.List<InstancePatchStateOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static InstancePatchStateOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InstancePatchStateOperatorType._(super.value, super.name);
}

class InstancePropertyFilterKey extends $pb.ProtobufEnum {
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_ASSOCIATION_STATUS =
      InstancePropertyFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'INSTANCE_PROPERTY_FILTER_KEY_ASSOCIATION_STATUS');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_INSTANCE_IDS = InstancePropertyFilterKey._(
          1, _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_INSTANCE_IDS');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_DOCUMENT_NAME = InstancePropertyFilterKey._(
          2,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_DOCUMENT_NAME');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_ACTIVATION_IDS = InstancePropertyFilterKey._(
          3,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_ACTIVATION_IDS');
  static const InstancePropertyFilterKey INSTANCE_PROPERTY_FILTER_KEY_IAM_ROLE =
      InstancePropertyFilterKey._(
          4, _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_IAM_ROLE');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_AGENT_VERSION = InstancePropertyFilterKey._(
          5,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_AGENT_VERSION');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_PING_STATUS = InstancePropertyFilterKey._(
          6, _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_PING_STATUS');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_PLATFORM_TYPES = InstancePropertyFilterKey._(
          7,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_PLATFORM_TYPES');
  static const InstancePropertyFilterKey
      INSTANCE_PROPERTY_FILTER_KEY_RESOURCE_TYPE = InstancePropertyFilterKey._(
          8,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_KEY_RESOURCE_TYPE');

  static const $core.List<InstancePropertyFilterKey> values =
      <InstancePropertyFilterKey>[
    INSTANCE_PROPERTY_FILTER_KEY_ASSOCIATION_STATUS,
    INSTANCE_PROPERTY_FILTER_KEY_INSTANCE_IDS,
    INSTANCE_PROPERTY_FILTER_KEY_DOCUMENT_NAME,
    INSTANCE_PROPERTY_FILTER_KEY_ACTIVATION_IDS,
    INSTANCE_PROPERTY_FILTER_KEY_IAM_ROLE,
    INSTANCE_PROPERTY_FILTER_KEY_AGENT_VERSION,
    INSTANCE_PROPERTY_FILTER_KEY_PING_STATUS,
    INSTANCE_PROPERTY_FILTER_KEY_PLATFORM_TYPES,
    INSTANCE_PROPERTY_FILTER_KEY_RESOURCE_TYPE,
  ];

  static final $core.List<InstancePropertyFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static InstancePropertyFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InstancePropertyFilterKey._(super.value, super.name);
}

class InstancePropertyFilterOperator extends $pb.ProtobufEnum {
  static const InstancePropertyFilterOperator
      INSTANCE_PROPERTY_FILTER_OPERATOR_LESS_THAN =
      InstancePropertyFilterOperator._(0,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_OPERATOR_LESS_THAN');
  static const InstancePropertyFilterOperator
      INSTANCE_PROPERTY_FILTER_OPERATOR_GREATER_THAN =
      InstancePropertyFilterOperator._(
          1,
          _omitEnumNames
              ? ''
              : 'INSTANCE_PROPERTY_FILTER_OPERATOR_GREATER_THAN');
  static const InstancePropertyFilterOperator
      INSTANCE_PROPERTY_FILTER_OPERATOR_BEGIN_WITH =
      InstancePropertyFilterOperator._(2,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_OPERATOR_BEGIN_WITH');
  static const InstancePropertyFilterOperator
      INSTANCE_PROPERTY_FILTER_OPERATOR_NOT_EQUAL =
      InstancePropertyFilterOperator._(3,
          _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_OPERATOR_NOT_EQUAL');
  static const InstancePropertyFilterOperator
      INSTANCE_PROPERTY_FILTER_OPERATOR_EQUAL =
      InstancePropertyFilterOperator._(
          4, _omitEnumNames ? '' : 'INSTANCE_PROPERTY_FILTER_OPERATOR_EQUAL');

  static const $core.List<InstancePropertyFilterOperator> values =
      <InstancePropertyFilterOperator>[
    INSTANCE_PROPERTY_FILTER_OPERATOR_LESS_THAN,
    INSTANCE_PROPERTY_FILTER_OPERATOR_GREATER_THAN,
    INSTANCE_PROPERTY_FILTER_OPERATOR_BEGIN_WITH,
    INSTANCE_PROPERTY_FILTER_OPERATOR_NOT_EQUAL,
    INSTANCE_PROPERTY_FILTER_OPERATOR_EQUAL,
  ];

  static final $core.List<InstancePropertyFilterOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static InstancePropertyFilterOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InstancePropertyFilterOperator._(super.value, super.name);
}

class InventoryAttributeDataType extends $pb.ProtobufEnum {
  static const InventoryAttributeDataType INVENTORY_ATTRIBUTE_DATA_TYPE_STRING =
      InventoryAttributeDataType._(
          0, _omitEnumNames ? '' : 'INVENTORY_ATTRIBUTE_DATA_TYPE_STRING');
  static const InventoryAttributeDataType INVENTORY_ATTRIBUTE_DATA_TYPE_NUMBER =
      InventoryAttributeDataType._(
          1, _omitEnumNames ? '' : 'INVENTORY_ATTRIBUTE_DATA_TYPE_NUMBER');

  static const $core.List<InventoryAttributeDataType> values =
      <InventoryAttributeDataType>[
    INVENTORY_ATTRIBUTE_DATA_TYPE_STRING,
    INVENTORY_ATTRIBUTE_DATA_TYPE_NUMBER,
  ];

  static final $core.List<InventoryAttributeDataType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventoryAttributeDataType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryAttributeDataType._(super.value, super.name);
}

class InventoryDeletionStatus extends $pb.ProtobufEnum {
  static const InventoryDeletionStatus INVENTORY_DELETION_STATUS_IN_PROGRESS =
      InventoryDeletionStatus._(
          0, _omitEnumNames ? '' : 'INVENTORY_DELETION_STATUS_IN_PROGRESS');
  static const InventoryDeletionStatus INVENTORY_DELETION_STATUS_COMPLETE =
      InventoryDeletionStatus._(
          1, _omitEnumNames ? '' : 'INVENTORY_DELETION_STATUS_COMPLETE');

  static const $core.List<InventoryDeletionStatus> values =
      <InventoryDeletionStatus>[
    INVENTORY_DELETION_STATUS_IN_PROGRESS,
    INVENTORY_DELETION_STATUS_COMPLETE,
  ];

  static final $core.List<InventoryDeletionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventoryDeletionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryDeletionStatus._(super.value, super.name);
}

class InventoryQueryOperatorType extends $pb.ProtobufEnum {
  static const InventoryQueryOperatorType
      INVENTORY_QUERY_OPERATOR_TYPE_LESS_THAN = InventoryQueryOperatorType._(
          0, _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_LESS_THAN');
  static const InventoryQueryOperatorType INVENTORY_QUERY_OPERATOR_TYPE_EXISTS =
      InventoryQueryOperatorType._(
          1, _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_EXISTS');
  static const InventoryQueryOperatorType
      INVENTORY_QUERY_OPERATOR_TYPE_GREATER_THAN = InventoryQueryOperatorType._(
          2,
          _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_GREATER_THAN');
  static const InventoryQueryOperatorType
      INVENTORY_QUERY_OPERATOR_TYPE_BEGIN_WITH = InventoryQueryOperatorType._(
          3, _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_BEGIN_WITH');
  static const InventoryQueryOperatorType
      INVENTORY_QUERY_OPERATOR_TYPE_NOT_EQUAL = InventoryQueryOperatorType._(
          4, _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_NOT_EQUAL');
  static const InventoryQueryOperatorType INVENTORY_QUERY_OPERATOR_TYPE_EQUAL =
      InventoryQueryOperatorType._(
          5, _omitEnumNames ? '' : 'INVENTORY_QUERY_OPERATOR_TYPE_EQUAL');

  static const $core.List<InventoryQueryOperatorType> values =
      <InventoryQueryOperatorType>[
    INVENTORY_QUERY_OPERATOR_TYPE_LESS_THAN,
    INVENTORY_QUERY_OPERATOR_TYPE_EXISTS,
    INVENTORY_QUERY_OPERATOR_TYPE_GREATER_THAN,
    INVENTORY_QUERY_OPERATOR_TYPE_BEGIN_WITH,
    INVENTORY_QUERY_OPERATOR_TYPE_NOT_EQUAL,
    INVENTORY_QUERY_OPERATOR_TYPE_EQUAL,
  ];

  static final $core.List<InventoryQueryOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static InventoryQueryOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventoryQueryOperatorType._(super.value, super.name);
}

class InventorySchemaDeleteOption extends $pb.ProtobufEnum {
  static const InventorySchemaDeleteOption
      INVENTORY_SCHEMA_DELETE_OPTION_DISABLE_SCHEMA =
      InventorySchemaDeleteOption._(
          0,
          _omitEnumNames
              ? ''
              : 'INVENTORY_SCHEMA_DELETE_OPTION_DISABLE_SCHEMA');
  static const InventorySchemaDeleteOption
      INVENTORY_SCHEMA_DELETE_OPTION_DELETE_SCHEMA =
      InventorySchemaDeleteOption._(1,
          _omitEnumNames ? '' : 'INVENTORY_SCHEMA_DELETE_OPTION_DELETE_SCHEMA');

  static const $core.List<InventorySchemaDeleteOption> values =
      <InventorySchemaDeleteOption>[
    INVENTORY_SCHEMA_DELETE_OPTION_DISABLE_SCHEMA,
    INVENTORY_SCHEMA_DELETE_OPTION_DELETE_SCHEMA,
  ];

  static final $core.List<InventorySchemaDeleteOption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static InventorySchemaDeleteOption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InventorySchemaDeleteOption._(super.value, super.name);
}

class LastResourceDataSyncStatus extends $pb.ProtobufEnum {
  static const LastResourceDataSyncStatus
      LAST_RESOURCE_DATA_SYNC_STATUS_SUCCESSFUL = LastResourceDataSyncStatus._(
          0, _omitEnumNames ? '' : 'LAST_RESOURCE_DATA_SYNC_STATUS_SUCCESSFUL');
  static const LastResourceDataSyncStatus
      LAST_RESOURCE_DATA_SYNC_STATUS_INPROGRESS = LastResourceDataSyncStatus._(
          1, _omitEnumNames ? '' : 'LAST_RESOURCE_DATA_SYNC_STATUS_INPROGRESS');
  static const LastResourceDataSyncStatus
      LAST_RESOURCE_DATA_SYNC_STATUS_FAILED = LastResourceDataSyncStatus._(
          2, _omitEnumNames ? '' : 'LAST_RESOURCE_DATA_SYNC_STATUS_FAILED');

  static const $core.List<LastResourceDataSyncStatus> values =
      <LastResourceDataSyncStatus>[
    LAST_RESOURCE_DATA_SYNC_STATUS_SUCCESSFUL,
    LAST_RESOURCE_DATA_SYNC_STATUS_INPROGRESS,
    LAST_RESOURCE_DATA_SYNC_STATUS_FAILED,
  ];

  static final $core.List<LastResourceDataSyncStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static LastResourceDataSyncStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const LastResourceDataSyncStatus._(super.value, super.name);
}

class MaintenanceWindowExecutionStatus extends $pb.ProtobufEnum {
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_TIMEDOUT =
      MaintenanceWindowExecutionStatus._(0,
          _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_TIMEDOUT');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_SKIPPEDOVERLAPPING =
      MaintenanceWindowExecutionStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_SKIPPEDOVERLAPPING');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_FAILED =
      MaintenanceWindowExecutionStatus._(2,
          _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_FAILED');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLED =
      MaintenanceWindowExecutionStatus._(
          3,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLED');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLING =
      MaintenanceWindowExecutionStatus._(
          4,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLING');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_SUCCESS =
      MaintenanceWindowExecutionStatus._(5,
          _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_SUCCESS');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_INPROGRESS =
      MaintenanceWindowExecutionStatus._(
          6,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_INPROGRESS');
  static const MaintenanceWindowExecutionStatus
      MAINTENANCE_WINDOW_EXECUTION_STATUS_PENDING =
      MaintenanceWindowExecutionStatus._(7,
          _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_EXECUTION_STATUS_PENDING');

  static const $core.List<MaintenanceWindowExecutionStatus> values =
      <MaintenanceWindowExecutionStatus>[
    MAINTENANCE_WINDOW_EXECUTION_STATUS_TIMEDOUT,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_SKIPPEDOVERLAPPING,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_FAILED,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLED,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_CANCELLING,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_SUCCESS,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_INPROGRESS,
    MAINTENANCE_WINDOW_EXECUTION_STATUS_PENDING,
  ];

  static final $core.List<MaintenanceWindowExecutionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static MaintenanceWindowExecutionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MaintenanceWindowExecutionStatus._(super.value, super.name);
}

class MaintenanceWindowResourceType extends $pb.ProtobufEnum {
  static const MaintenanceWindowResourceType
      MAINTENANCE_WINDOW_RESOURCE_TYPE_RESOURCEGROUP =
      MaintenanceWindowResourceType._(
          0,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_RESOURCE_TYPE_RESOURCEGROUP');
  static const MaintenanceWindowResourceType
      MAINTENANCE_WINDOW_RESOURCE_TYPE_INSTANCE =
      MaintenanceWindowResourceType._(
          1, _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_RESOURCE_TYPE_INSTANCE');

  static const $core.List<MaintenanceWindowResourceType> values =
      <MaintenanceWindowResourceType>[
    MAINTENANCE_WINDOW_RESOURCE_TYPE_RESOURCEGROUP,
    MAINTENANCE_WINDOW_RESOURCE_TYPE_INSTANCE,
  ];

  static final $core.List<MaintenanceWindowResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MaintenanceWindowResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MaintenanceWindowResourceType._(super.value, super.name);
}

class MaintenanceWindowTaskCutoffBehavior extends $pb.ProtobufEnum {
  static const MaintenanceWindowTaskCutoffBehavior
      MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CANCELTASK =
      MaintenanceWindowTaskCutoffBehavior._(
          0,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CANCELTASK');
  static const MaintenanceWindowTaskCutoffBehavior
      MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CONTINUETASK =
      MaintenanceWindowTaskCutoffBehavior._(
          1,
          _omitEnumNames
              ? ''
              : 'MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CONTINUETASK');

  static const $core.List<MaintenanceWindowTaskCutoffBehavior> values =
      <MaintenanceWindowTaskCutoffBehavior>[
    MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CANCELTASK,
    MAINTENANCE_WINDOW_TASK_CUTOFF_BEHAVIOR_CONTINUETASK,
  ];

  static final $core.List<MaintenanceWindowTaskCutoffBehavior?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MaintenanceWindowTaskCutoffBehavior? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MaintenanceWindowTaskCutoffBehavior._(super.value, super.name);
}

class MaintenanceWindowTaskType extends $pb.ProtobufEnum {
  static const MaintenanceWindowTaskType
      MAINTENANCE_WINDOW_TASK_TYPE_STEPFUNCTIONS = MaintenanceWindowTaskType._(
          0,
          _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_TASK_TYPE_STEPFUNCTIONS');
  static const MaintenanceWindowTaskType
      MAINTENANCE_WINDOW_TASK_TYPE_RUNCOMMAND = MaintenanceWindowTaskType._(
          1, _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_TASK_TYPE_RUNCOMMAND');
  static const MaintenanceWindowTaskType
      MAINTENANCE_WINDOW_TASK_TYPE_AUTOMATION = MaintenanceWindowTaskType._(
          2, _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_TASK_TYPE_AUTOMATION');
  static const MaintenanceWindowTaskType MAINTENANCE_WINDOW_TASK_TYPE_LAMBDA =
      MaintenanceWindowTaskType._(
          3, _omitEnumNames ? '' : 'MAINTENANCE_WINDOW_TASK_TYPE_LAMBDA');

  static const $core.List<MaintenanceWindowTaskType> values =
      <MaintenanceWindowTaskType>[
    MAINTENANCE_WINDOW_TASK_TYPE_STEPFUNCTIONS,
    MAINTENANCE_WINDOW_TASK_TYPE_RUNCOMMAND,
    MAINTENANCE_WINDOW_TASK_TYPE_AUTOMATION,
    MAINTENANCE_WINDOW_TASK_TYPE_LAMBDA,
  ];

  static final $core.List<MaintenanceWindowTaskType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MaintenanceWindowTaskType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MaintenanceWindowTaskType._(super.value, super.name);
}

class ManagedStatus extends $pb.ProtobufEnum {
  static const ManagedStatus MANAGED_STATUS_UNMANAGED =
      ManagedStatus._(0, _omitEnumNames ? '' : 'MANAGED_STATUS_UNMANAGED');
  static const ManagedStatus MANAGED_STATUS_MANAGED =
      ManagedStatus._(1, _omitEnumNames ? '' : 'MANAGED_STATUS_MANAGED');
  static const ManagedStatus MANAGED_STATUS_ALL =
      ManagedStatus._(2, _omitEnumNames ? '' : 'MANAGED_STATUS_ALL');

  static const $core.List<ManagedStatus> values = <ManagedStatus>[
    MANAGED_STATUS_UNMANAGED,
    MANAGED_STATUS_MANAGED,
    MANAGED_STATUS_ALL,
  ];

  static final $core.List<ManagedStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ManagedStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ManagedStatus._(super.value, super.name);
}

class NodeAggregatorType extends $pb.ProtobufEnum {
  static const NodeAggregatorType NODE_AGGREGATOR_TYPE_COUNT =
      NodeAggregatorType._(
          0, _omitEnumNames ? '' : 'NODE_AGGREGATOR_TYPE_COUNT');

  static const $core.List<NodeAggregatorType> values = <NodeAggregatorType>[
    NODE_AGGREGATOR_TYPE_COUNT,
  ];

  static final $core.List<NodeAggregatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static NodeAggregatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NodeAggregatorType._(super.value, super.name);
}

class NodeAttributeName extends $pb.ProtobufEnum {
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_PLATFORM_VERSION =
      NodeAttributeName._(
          0, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_PLATFORM_VERSION');
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_REGION =
      NodeAttributeName._(
          1, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_REGION');
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_PLATFORM_TYPE =
      NodeAttributeName._(
          2, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_PLATFORM_TYPE');
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_AGENT_VERSION =
      NodeAttributeName._(
          3, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_AGENT_VERSION');
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_PLATFORM_NAME =
      NodeAttributeName._(
          4, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_PLATFORM_NAME');
  static const NodeAttributeName NODE_ATTRIBUTE_NAME_RESOURCE_TYPE =
      NodeAttributeName._(
          5, _omitEnumNames ? '' : 'NODE_ATTRIBUTE_NAME_RESOURCE_TYPE');

  static const $core.List<NodeAttributeName> values = <NodeAttributeName>[
    NODE_ATTRIBUTE_NAME_PLATFORM_VERSION,
    NODE_ATTRIBUTE_NAME_REGION,
    NODE_ATTRIBUTE_NAME_PLATFORM_TYPE,
    NODE_ATTRIBUTE_NAME_AGENT_VERSION,
    NODE_ATTRIBUTE_NAME_PLATFORM_NAME,
    NODE_ATTRIBUTE_NAME_RESOURCE_TYPE,
  ];

  static final $core.List<NodeAttributeName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static NodeAttributeName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NodeAttributeName._(super.value, super.name);
}

class NodeFilterKey extends $pb.ProtobufEnum {
  static const NodeFilterKey NODE_FILTER_KEY_PLATFORM_VERSION = NodeFilterKey._(
      0, _omitEnumNames ? '' : 'NODE_FILTER_KEY_PLATFORM_VERSION');
  static const NodeFilterKey NODE_FILTER_KEY_REGION =
      NodeFilterKey._(1, _omitEnumNames ? '' : 'NODE_FILTER_KEY_REGION');
  static const NodeFilterKey NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_ID =
      NodeFilterKey._(
          2, _omitEnumNames ? '' : 'NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_ID');
  static const NodeFilterKey NODE_FILTER_KEY_IP_ADDRESS =
      NodeFilterKey._(3, _omitEnumNames ? '' : 'NODE_FILTER_KEY_IP_ADDRESS');
  static const NodeFilterKey NODE_FILTER_KEY_PLATFORM_TYPE =
      NodeFilterKey._(4, _omitEnumNames ? '' : 'NODE_FILTER_KEY_PLATFORM_TYPE');
  static const NodeFilterKey NODE_FILTER_KEY_ACCOUNT_ID =
      NodeFilterKey._(5, _omitEnumNames ? '' : 'NODE_FILTER_KEY_ACCOUNT_ID');
  static const NodeFilterKey NODE_FILTER_KEY_MANAGED_STATUS = NodeFilterKey._(
      6, _omitEnumNames ? '' : 'NODE_FILTER_KEY_MANAGED_STATUS');
  static const NodeFilterKey NODE_FILTER_KEY_INSTANCE_STATUS = NodeFilterKey._(
      7, _omitEnumNames ? '' : 'NODE_FILTER_KEY_INSTANCE_STATUS');
  static const NodeFilterKey NODE_FILTER_KEY_AGENT_VERSION =
      NodeFilterKey._(8, _omitEnumNames ? '' : 'NODE_FILTER_KEY_AGENT_VERSION');
  static const NodeFilterKey NODE_FILTER_KEY_COMPUTER_NAME =
      NodeFilterKey._(9, _omitEnumNames ? '' : 'NODE_FILTER_KEY_COMPUTER_NAME');
  static const NodeFilterKey NODE_FILTER_KEY_AGENT_TYPE =
      NodeFilterKey._(10, _omitEnumNames ? '' : 'NODE_FILTER_KEY_AGENT_TYPE');
  static const NodeFilterKey NODE_FILTER_KEY_PLATFORM_NAME = NodeFilterKey._(
      11, _omitEnumNames ? '' : 'NODE_FILTER_KEY_PLATFORM_NAME');
  static const NodeFilterKey NODE_FILTER_KEY_RESOURCE_TYPE = NodeFilterKey._(
      12, _omitEnumNames ? '' : 'NODE_FILTER_KEY_RESOURCE_TYPE');
  static const NodeFilterKey NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_PATH =
      NodeFilterKey._(
          13, _omitEnumNames ? '' : 'NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_PATH');
  static const NodeFilterKey NODE_FILTER_KEY_INSTANCE_ID =
      NodeFilterKey._(14, _omitEnumNames ? '' : 'NODE_FILTER_KEY_INSTANCE_ID');

  static const $core.List<NodeFilterKey> values = <NodeFilterKey>[
    NODE_FILTER_KEY_PLATFORM_VERSION,
    NODE_FILTER_KEY_REGION,
    NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_ID,
    NODE_FILTER_KEY_IP_ADDRESS,
    NODE_FILTER_KEY_PLATFORM_TYPE,
    NODE_FILTER_KEY_ACCOUNT_ID,
    NODE_FILTER_KEY_MANAGED_STATUS,
    NODE_FILTER_KEY_INSTANCE_STATUS,
    NODE_FILTER_KEY_AGENT_VERSION,
    NODE_FILTER_KEY_COMPUTER_NAME,
    NODE_FILTER_KEY_AGENT_TYPE,
    NODE_FILTER_KEY_PLATFORM_NAME,
    NODE_FILTER_KEY_RESOURCE_TYPE,
    NODE_FILTER_KEY_ORGANIZATIONAL_UNIT_PATH,
    NODE_FILTER_KEY_INSTANCE_ID,
  ];

  static final $core.List<NodeFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 14);
  static NodeFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NodeFilterKey._(super.value, super.name);
}

class NodeFilterOperatorType extends $pb.ProtobufEnum {
  static const NodeFilterOperatorType NODE_FILTER_OPERATOR_TYPE_BEGIN_WITH =
      NodeFilterOperatorType._(
          0, _omitEnumNames ? '' : 'NODE_FILTER_OPERATOR_TYPE_BEGIN_WITH');
  static const NodeFilterOperatorType NODE_FILTER_OPERATOR_TYPE_NOT_EQUAL =
      NodeFilterOperatorType._(
          1, _omitEnumNames ? '' : 'NODE_FILTER_OPERATOR_TYPE_NOT_EQUAL');
  static const NodeFilterOperatorType NODE_FILTER_OPERATOR_TYPE_EQUAL =
      NodeFilterOperatorType._(
          2, _omitEnumNames ? '' : 'NODE_FILTER_OPERATOR_TYPE_EQUAL');

  static const $core.List<NodeFilterOperatorType> values =
      <NodeFilterOperatorType>[
    NODE_FILTER_OPERATOR_TYPE_BEGIN_WITH,
    NODE_FILTER_OPERATOR_TYPE_NOT_EQUAL,
    NODE_FILTER_OPERATOR_TYPE_EQUAL,
  ];

  static final $core.List<NodeFilterOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static NodeFilterOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NodeFilterOperatorType._(super.value, super.name);
}

class NodeTypeName extends $pb.ProtobufEnum {
  static const NodeTypeName NODE_TYPE_NAME_INSTANCE =
      NodeTypeName._(0, _omitEnumNames ? '' : 'NODE_TYPE_NAME_INSTANCE');

  static const $core.List<NodeTypeName> values = <NodeTypeName>[
    NODE_TYPE_NAME_INSTANCE,
  ];

  static final $core.List<NodeTypeName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static NodeTypeName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NodeTypeName._(super.value, super.name);
}

class NotificationEvent extends $pb.ProtobufEnum {
  static const NotificationEvent NOTIFICATION_EVENT_TIMED_OUT =
      NotificationEvent._(
          0, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_TIMED_OUT');
  static const NotificationEvent NOTIFICATION_EVENT_SUCCESS =
      NotificationEvent._(
          1, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_SUCCESS');
  static const NotificationEvent NOTIFICATION_EVENT_IN_PROGRESS =
      NotificationEvent._(
          2, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_IN_PROGRESS');
  static const NotificationEvent NOTIFICATION_EVENT_CANCELLED =
      NotificationEvent._(
          3, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_CANCELLED');
  static const NotificationEvent NOTIFICATION_EVENT_ALL =
      NotificationEvent._(4, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_ALL');
  static const NotificationEvent NOTIFICATION_EVENT_FAILED =
      NotificationEvent._(5, _omitEnumNames ? '' : 'NOTIFICATION_EVENT_FAILED');

  static const $core.List<NotificationEvent> values = <NotificationEvent>[
    NOTIFICATION_EVENT_TIMED_OUT,
    NOTIFICATION_EVENT_SUCCESS,
    NOTIFICATION_EVENT_IN_PROGRESS,
    NOTIFICATION_EVENT_CANCELLED,
    NOTIFICATION_EVENT_ALL,
    NOTIFICATION_EVENT_FAILED,
  ];

  static final $core.List<NotificationEvent?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static NotificationEvent? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotificationEvent._(super.value, super.name);
}

class NotificationType extends $pb.ProtobufEnum {
  static const NotificationType NOTIFICATION_TYPE_INVOCATION =
      NotificationType._(
          0, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_INVOCATION');
  static const NotificationType NOTIFICATION_TYPE_COMMAND =
      NotificationType._(1, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_COMMAND');

  static const $core.List<NotificationType> values = <NotificationType>[
    NOTIFICATION_TYPE_INVOCATION,
    NOTIFICATION_TYPE_COMMAND,
  ];

  static final $core.List<NotificationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static NotificationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotificationType._(super.value, super.name);
}

class OperatingSystem extends $pb.ProtobufEnum {
  static const OperatingSystem OPERATING_SYSTEM_MACOS =
      OperatingSystem._(0, _omitEnumNames ? '' : 'OPERATING_SYSTEM_MACOS');
  static const OperatingSystem OPERATING_SYSTEM_REDHATENTERPRISELINUX =
      OperatingSystem._(
          1, _omitEnumNames ? '' : 'OPERATING_SYSTEM_REDHATENTERPRISELINUX');
  static const OperatingSystem OPERATING_SYSTEM_ALMALINUX =
      OperatingSystem._(2, _omitEnumNames ? '' : 'OPERATING_SYSTEM_ALMALINUX');
  static const OperatingSystem OPERATING_SYSTEM_DEBIAN =
      OperatingSystem._(3, _omitEnumNames ? '' : 'OPERATING_SYSTEM_DEBIAN');
  static const OperatingSystem OPERATING_SYSTEM_WINDOWS =
      OperatingSystem._(4, _omitEnumNames ? '' : 'OPERATING_SYSTEM_WINDOWS');
  static const OperatingSystem OPERATING_SYSTEM_AMAZONLINUX = OperatingSystem._(
      5, _omitEnumNames ? '' : 'OPERATING_SYSTEM_AMAZONLINUX');
  static const OperatingSystem OPERATING_SYSTEM_AMAZONLINUX2022 =
      OperatingSystem._(
          6, _omitEnumNames ? '' : 'OPERATING_SYSTEM_AMAZONLINUX2022');
  static const OperatingSystem OPERATING_SYSTEM_CENTOS =
      OperatingSystem._(7, _omitEnumNames ? '' : 'OPERATING_SYSTEM_CENTOS');
  static const OperatingSystem OPERATING_SYSTEM_AMAZONLINUX2 =
      OperatingSystem._(
          8, _omitEnumNames ? '' : 'OPERATING_SYSTEM_AMAZONLINUX2');
  static const OperatingSystem OPERATING_SYSTEM_ORACLELINUX = OperatingSystem._(
      9, _omitEnumNames ? '' : 'OPERATING_SYSTEM_ORACLELINUX');
  static const OperatingSystem OPERATING_SYSTEM_SUSE =
      OperatingSystem._(10, _omitEnumNames ? '' : 'OPERATING_SYSTEM_SUSE');
  static const OperatingSystem OPERATING_SYSTEM_RASPBIAN =
      OperatingSystem._(11, _omitEnumNames ? '' : 'OPERATING_SYSTEM_RASPBIAN');
  static const OperatingSystem OPERATING_SYSTEM_UBUNTU =
      OperatingSystem._(12, _omitEnumNames ? '' : 'OPERATING_SYSTEM_UBUNTU');
  static const OperatingSystem OPERATING_SYSTEM_ROCKY_LINUX = OperatingSystem._(
      13, _omitEnumNames ? '' : 'OPERATING_SYSTEM_ROCKY_LINUX');
  static const OperatingSystem OPERATING_SYSTEM_AMAZONLINUX2023 =
      OperatingSystem._(
          14, _omitEnumNames ? '' : 'OPERATING_SYSTEM_AMAZONLINUX2023');

  static const $core.List<OperatingSystem> values = <OperatingSystem>[
    OPERATING_SYSTEM_MACOS,
    OPERATING_SYSTEM_REDHATENTERPRISELINUX,
    OPERATING_SYSTEM_ALMALINUX,
    OPERATING_SYSTEM_DEBIAN,
    OPERATING_SYSTEM_WINDOWS,
    OPERATING_SYSTEM_AMAZONLINUX,
    OPERATING_SYSTEM_AMAZONLINUX2022,
    OPERATING_SYSTEM_CENTOS,
    OPERATING_SYSTEM_AMAZONLINUX2,
    OPERATING_SYSTEM_ORACLELINUX,
    OPERATING_SYSTEM_SUSE,
    OPERATING_SYSTEM_RASPBIAN,
    OPERATING_SYSTEM_UBUNTU,
    OPERATING_SYSTEM_ROCKY_LINUX,
    OPERATING_SYSTEM_AMAZONLINUX2023,
  ];

  static final $core.List<OperatingSystem?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 14);
  static OperatingSystem? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OperatingSystem._(super.value, super.name);
}

class OpsFilterOperatorType extends $pb.ProtobufEnum {
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_LESS_THAN =
      OpsFilterOperatorType._(
          0, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_LESS_THAN');
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_EXISTS =
      OpsFilterOperatorType._(
          1, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_EXISTS');
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_GREATER_THAN =
      OpsFilterOperatorType._(
          2, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_GREATER_THAN');
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_BEGIN_WITH =
      OpsFilterOperatorType._(
          3, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_BEGIN_WITH');
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_NOT_EQUAL =
      OpsFilterOperatorType._(
          4, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_NOT_EQUAL');
  static const OpsFilterOperatorType OPS_FILTER_OPERATOR_TYPE_EQUAL =
      OpsFilterOperatorType._(
          5, _omitEnumNames ? '' : 'OPS_FILTER_OPERATOR_TYPE_EQUAL');

  static const $core.List<OpsFilterOperatorType> values =
      <OpsFilterOperatorType>[
    OPS_FILTER_OPERATOR_TYPE_LESS_THAN,
    OPS_FILTER_OPERATOR_TYPE_EXISTS,
    OPS_FILTER_OPERATOR_TYPE_GREATER_THAN,
    OPS_FILTER_OPERATOR_TYPE_BEGIN_WITH,
    OPS_FILTER_OPERATOR_TYPE_NOT_EQUAL,
    OPS_FILTER_OPERATOR_TYPE_EQUAL,
  ];

  static final $core.List<OpsFilterOperatorType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static OpsFilterOperatorType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsFilterOperatorType._(super.value, super.name);
}

class OpsItemDataType extends $pb.ProtobufEnum {
  static const OpsItemDataType OPS_ITEM_DATA_TYPE_SEARCHABLE_STRING =
      OpsItemDataType._(
          0, _omitEnumNames ? '' : 'OPS_ITEM_DATA_TYPE_SEARCHABLE_STRING');
  static const OpsItemDataType OPS_ITEM_DATA_TYPE_STRING =
      OpsItemDataType._(1, _omitEnumNames ? '' : 'OPS_ITEM_DATA_TYPE_STRING');

  static const $core.List<OpsItemDataType> values = <OpsItemDataType>[
    OPS_ITEM_DATA_TYPE_SEARCHABLE_STRING,
    OPS_ITEM_DATA_TYPE_STRING,
  ];

  static final $core.List<OpsItemDataType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static OpsItemDataType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemDataType._(super.value, super.name);
}

class OpsItemEventFilterKey extends $pb.ProtobufEnum {
  static const OpsItemEventFilterKey OPS_ITEM_EVENT_FILTER_KEY_OPSITEM_ID =
      OpsItemEventFilterKey._(
          0, _omitEnumNames ? '' : 'OPS_ITEM_EVENT_FILTER_KEY_OPSITEM_ID');

  static const $core.List<OpsItemEventFilterKey> values =
      <OpsItemEventFilterKey>[
    OPS_ITEM_EVENT_FILTER_KEY_OPSITEM_ID,
  ];

  static final $core.List<OpsItemEventFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OpsItemEventFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemEventFilterKey._(super.value, super.name);
}

class OpsItemEventFilterOperator extends $pb.ProtobufEnum {
  static const OpsItemEventFilterOperator OPS_ITEM_EVENT_FILTER_OPERATOR_EQUAL =
      OpsItemEventFilterOperator._(
          0, _omitEnumNames ? '' : 'OPS_ITEM_EVENT_FILTER_OPERATOR_EQUAL');

  static const $core.List<OpsItemEventFilterOperator> values =
      <OpsItemEventFilterOperator>[
    OPS_ITEM_EVENT_FILTER_OPERATOR_EQUAL,
  ];

  static final $core.List<OpsItemEventFilterOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OpsItemEventFilterOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemEventFilterOperator._(super.value, super.name);
}

class OpsItemFilterKey extends $pb.ProtobufEnum {
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_ACTUAL_START_TIME =
      OpsItemFilterKey._(
          0, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_ACTUAL_START_TIME');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_OPS_ITEM_ID =
      OpsItemFilterKey._(
          1,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_OPS_ITEM_ID');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TARGETS_RESOURCE_GROUP =
      OpsItemFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TARGETS_RESOURCE_GROUP');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TEMPLATE =
      OpsItemFilterKey._(3,
          _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TEMPLATE');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_VALUE =
      OpsItemFilterKey._(4,
          _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_VALUE');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_NAME = OpsItemFilterKey._(
          5,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_NAME');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_ARN = OpsItemFilterKey._(
          6,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_ARN');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_PRIORITY =
      OpsItemFilterKey._(
          7, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_PRIORITY');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_ARN = OpsItemFilterKey._(
          8,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_ARN');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_REGION = OpsItemFilterKey._(
          9,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_REGION');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_CREATED_TIME =
      OpsItemFilterKey._(
          10, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_CREATED_TIME');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA =
      OpsItemFilterKey._(
          11, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_ACCOUNT_ID =
      OpsItemFilterKey._(
          12, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_ACCOUNT_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_RESOURCE_ID =
      OpsItemFilterKey._(
          13, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_RESOURCE_ID');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_ACCOUNT_ID = OpsItemFilterKey._(
          14,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_ACCOUNT_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_STATUS = OpsItemFilterKey._(
      15, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_STATUS');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ID = OpsItemFilterKey._(
          16,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ID =
      OpsItemFilterKey._(
          17,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_KEY =
      OpsItemFilterKey._(
          18, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_KEY');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_SOURCE = OpsItemFilterKey._(
      19, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_SOURCE');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_AUTOMATION_ID =
      OpsItemFilterKey._(
          20, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_AUTOMATION_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_CREATED_BY =
      OpsItemFilterKey._(
          21, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_CREATED_BY');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_PLANNED_START_TIME =
      OpsItemFilterKey._(
          22, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_PLANNED_START_TIME');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ARN = OpsItemFilterKey._(
          23,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ARN');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_TITLE =
      OpsItemFilterKey._(24, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_TITLE');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_PLANNED_END_TIME =
      OpsItemFilterKey._(
          25, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_PLANNED_END_TIME');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_NAME = OpsItemFilterKey._(
          26,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_NAME');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_IS_REPLICA =
      OpsItemFilterKey._(
          27,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_IS_REPLICA');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ARN = OpsItemFilterKey._(
          28,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ARN');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_INSIGHT_TYPE =
      OpsItemFilterKey._(
          29, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_INSIGHT_TYPE');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_LAST_MODIFIED_TIME =
      OpsItemFilterKey._(
          30, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_LAST_MODIFIED_TIME');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_SEVERITY =
      OpsItemFilterKey._(
          31, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_SEVERITY');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_ACTUAL_END_TIME =
      OpsItemFilterKey._(
          32, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_ACTUAL_END_TIME');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_CATEGORY =
      OpsItemFilterKey._(
          33, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_CATEGORY');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_OPSITEM_TYPE =
      OpsItemFilterKey._(
          34, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_OPSITEM_TYPE');
  static const OpsItemFilterKey
      OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_TARGET_RESOURCE_ID =
      OpsItemFilterKey._(
          35,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_TARGET_RESOURCE_ID');
  static const OpsItemFilterKey OPS_ITEM_FILTER_KEY_OPSITEM_ID =
      OpsItemFilterKey._(
          36, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_KEY_OPSITEM_ID');

  static const $core.List<OpsItemFilterKey> values = <OpsItemFilterKey>[
    OPS_ITEM_FILTER_KEY_ACTUAL_START_TIME,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_OPS_ITEM_ID,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TARGETS_RESOURCE_GROUP,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_TEMPLATE,
    OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_VALUE,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_NAME,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_ARN,
    OPS_ITEM_FILTER_KEY_PRIORITY,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_REQUESTER_ARN,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_REGION,
    OPS_ITEM_FILTER_KEY_CREATED_TIME,
    OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA,
    OPS_ITEM_FILTER_KEY_ACCOUNT_ID,
    OPS_ITEM_FILTER_KEY_RESOURCE_ID,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_SOURCE_ACCOUNT_ID,
    OPS_ITEM_FILTER_KEY_STATUS,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ID,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ID,
    OPS_ITEM_FILTER_KEY_OPERATIONAL_DATA_KEY,
    OPS_ITEM_FILTER_KEY_SOURCE,
    OPS_ITEM_FILTER_KEY_AUTOMATION_ID,
    OPS_ITEM_FILTER_KEY_CREATED_BY,
    OPS_ITEM_FILTER_KEY_PLANNED_START_TIME,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_APPROVER_ARN,
    OPS_ITEM_FILTER_KEY_TITLE,
    OPS_ITEM_FILTER_KEY_PLANNED_END_TIME,
    OPS_ITEM_FILTER_KEY_CHANGE_REQUEST_APPROVER_NAME,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_IS_REPLICA,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_REQUESTER_ARN,
    OPS_ITEM_FILTER_KEY_INSIGHT_TYPE,
    OPS_ITEM_FILTER_KEY_LAST_MODIFIED_TIME,
    OPS_ITEM_FILTER_KEY_SEVERITY,
    OPS_ITEM_FILTER_KEY_ACTUAL_END_TIME,
    OPS_ITEM_FILTER_KEY_CATEGORY,
    OPS_ITEM_FILTER_KEY_OPSITEM_TYPE,
    OPS_ITEM_FILTER_KEY_ACCESS_REQUEST_TARGET_RESOURCE_ID,
    OPS_ITEM_FILTER_KEY_OPSITEM_ID,
  ];

  static final $core.List<OpsItemFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 36);
  static OpsItemFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemFilterKey._(super.value, super.name);
}

class OpsItemFilterOperator extends $pb.ProtobufEnum {
  static const OpsItemFilterOperator OPS_ITEM_FILTER_OPERATOR_LESS_THAN =
      OpsItemFilterOperator._(
          0, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_OPERATOR_LESS_THAN');
  static const OpsItemFilterOperator OPS_ITEM_FILTER_OPERATOR_GREATER_THAN =
      OpsItemFilterOperator._(
          1, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_OPERATOR_GREATER_THAN');
  static const OpsItemFilterOperator OPS_ITEM_FILTER_OPERATOR_CONTAINS =
      OpsItemFilterOperator._(
          2, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_OPERATOR_CONTAINS');
  static const OpsItemFilterOperator OPS_ITEM_FILTER_OPERATOR_EQUAL =
      OpsItemFilterOperator._(
          3, _omitEnumNames ? '' : 'OPS_ITEM_FILTER_OPERATOR_EQUAL');

  static const $core.List<OpsItemFilterOperator> values =
      <OpsItemFilterOperator>[
    OPS_ITEM_FILTER_OPERATOR_LESS_THAN,
    OPS_ITEM_FILTER_OPERATOR_GREATER_THAN,
    OPS_ITEM_FILTER_OPERATOR_CONTAINS,
    OPS_ITEM_FILTER_OPERATOR_EQUAL,
  ];

  static final $core.List<OpsItemFilterOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static OpsItemFilterOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemFilterOperator._(super.value, super.name);
}

class OpsItemRelatedItemsFilterKey extends $pb.ProtobufEnum {
  static const OpsItemRelatedItemsFilterKey
      OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_URI =
      OpsItemRelatedItemsFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_URI');
  static const OpsItemRelatedItemsFilterKey
      OPS_ITEM_RELATED_ITEMS_FILTER_KEY_ASSOCIATION_ID =
      OpsItemRelatedItemsFilterKey._(
          1,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_ASSOCIATION_ID');
  static const OpsItemRelatedItemsFilterKey
      OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_TYPE =
      OpsItemRelatedItemsFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_TYPE');

  static const $core.List<OpsItemRelatedItemsFilterKey> values =
      <OpsItemRelatedItemsFilterKey>[
    OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_URI,
    OPS_ITEM_RELATED_ITEMS_FILTER_KEY_ASSOCIATION_ID,
    OPS_ITEM_RELATED_ITEMS_FILTER_KEY_RESOURCE_TYPE,
  ];

  static final $core.List<OpsItemRelatedItemsFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static OpsItemRelatedItemsFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemRelatedItemsFilterKey._(super.value, super.name);
}

class OpsItemRelatedItemsFilterOperator extends $pb.ProtobufEnum {
  static const OpsItemRelatedItemsFilterOperator
      OPS_ITEM_RELATED_ITEMS_FILTER_OPERATOR_EQUAL =
      OpsItemRelatedItemsFilterOperator._(0,
          _omitEnumNames ? '' : 'OPS_ITEM_RELATED_ITEMS_FILTER_OPERATOR_EQUAL');

  static const $core.List<OpsItemRelatedItemsFilterOperator> values =
      <OpsItemRelatedItemsFilterOperator>[
    OPS_ITEM_RELATED_ITEMS_FILTER_OPERATOR_EQUAL,
  ];

  static final $core.List<OpsItemRelatedItemsFilterOperator?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static OpsItemRelatedItemsFilterOperator? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemRelatedItemsFilterOperator._(super.value, super.name);
}

class OpsItemStatus extends $pb.ProtobufEnum {
  static const OpsItemStatus OPS_ITEM_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE =
      OpsItemStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE');
  static const OpsItemStatus OPS_ITEM_STATUS_PENDING =
      OpsItemStatus._(1, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_PENDING');
  static const OpsItemStatus OPS_ITEM_STATUS_COMPLETED_WITH_SUCCESS =
      OpsItemStatus._(
          2, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_COMPLETED_WITH_SUCCESS');
  static const OpsItemStatus OPS_ITEM_STATUS_RUNBOOK_IN_PROGRESS =
      OpsItemStatus._(
          3, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_RUNBOOK_IN_PROGRESS');
  static const OpsItemStatus OPS_ITEM_STATUS_REVOKED =
      OpsItemStatus._(4, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_REVOKED');
  static const OpsItemStatus OPS_ITEM_STATUS_TIMED_OUT =
      OpsItemStatus._(5, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_TIMED_OUT');
  static const OpsItemStatus OPS_ITEM_STATUS_CLOSED =
      OpsItemStatus._(6, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_CLOSED');
  static const OpsItemStatus OPS_ITEM_STATUS_PENDING_APPROVAL = OpsItemStatus._(
      7, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_PENDING_APPROVAL');
  static const OpsItemStatus OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED =
      OpsItemStatus._(
          8,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED');
  static const OpsItemStatus OPS_ITEM_STATUS_OPEN =
      OpsItemStatus._(9, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_OPEN');
  static const OpsItemStatus OPS_ITEM_STATUS_IN_PROGRESS =
      OpsItemStatus._(10, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_IN_PROGRESS');
  static const OpsItemStatus OPS_ITEM_STATUS_RESOLVED =
      OpsItemStatus._(11, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_RESOLVED');
  static const OpsItemStatus OPS_ITEM_STATUS_CANCELLED =
      OpsItemStatus._(12, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_CANCELLED');
  static const OpsItemStatus OPS_ITEM_STATUS_REJECTED =
      OpsItemStatus._(13, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_REJECTED');
  static const OpsItemStatus OPS_ITEM_STATUS_APPROVED =
      OpsItemStatus._(14, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_APPROVED');
  static const OpsItemStatus OPS_ITEM_STATUS_SCHEDULED =
      OpsItemStatus._(15, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_SCHEDULED');
  static const OpsItemStatus OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED =
      OpsItemStatus._(
          16,
          _omitEnumNames
              ? ''
              : 'OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED');
  static const OpsItemStatus OPS_ITEM_STATUS_CANCELLING =
      OpsItemStatus._(17, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_CANCELLING');
  static const OpsItemStatus OPS_ITEM_STATUS_COMPLETED_WITH_FAILURE =
      OpsItemStatus._(
          18, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_COMPLETED_WITH_FAILURE');
  static const OpsItemStatus OPS_ITEM_STATUS_FAILED =
      OpsItemStatus._(19, _omitEnumNames ? '' : 'OPS_ITEM_STATUS_FAILED');

  static const $core.List<OpsItemStatus> values = <OpsItemStatus>[
    OPS_ITEM_STATUS_PENDING_CHANGE_CALENDAR_OVERRIDE,
    OPS_ITEM_STATUS_PENDING,
    OPS_ITEM_STATUS_COMPLETED_WITH_SUCCESS,
    OPS_ITEM_STATUS_RUNBOOK_IN_PROGRESS,
    OPS_ITEM_STATUS_REVOKED,
    OPS_ITEM_STATUS_TIMED_OUT,
    OPS_ITEM_STATUS_CLOSED,
    OPS_ITEM_STATUS_PENDING_APPROVAL,
    OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_APPROVED,
    OPS_ITEM_STATUS_OPEN,
    OPS_ITEM_STATUS_IN_PROGRESS,
    OPS_ITEM_STATUS_RESOLVED,
    OPS_ITEM_STATUS_CANCELLED,
    OPS_ITEM_STATUS_REJECTED,
    OPS_ITEM_STATUS_APPROVED,
    OPS_ITEM_STATUS_SCHEDULED,
    OPS_ITEM_STATUS_CHANGE_CALENDAR_OVERRIDE_REJECTED,
    OPS_ITEM_STATUS_CANCELLING,
    OPS_ITEM_STATUS_COMPLETED_WITH_FAILURE,
    OPS_ITEM_STATUS_FAILED,
  ];

  static final $core.List<OpsItemStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 19);
  static OpsItemStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OpsItemStatus._(super.value, super.name);
}

class ParameterTier extends $pb.ProtobufEnum {
  static const ParameterTier PARAMETER_TIER_STANDARD =
      ParameterTier._(0, _omitEnumNames ? '' : 'PARAMETER_TIER_STANDARD');
  static const ParameterTier PARAMETER_TIER_ADVANCED =
      ParameterTier._(1, _omitEnumNames ? '' : 'PARAMETER_TIER_ADVANCED');
  static const ParameterTier PARAMETER_TIER_INTELLIGENT_TIERING =
      ParameterTier._(
          2, _omitEnumNames ? '' : 'PARAMETER_TIER_INTELLIGENT_TIERING');

  static const $core.List<ParameterTier> values = <ParameterTier>[
    PARAMETER_TIER_STANDARD,
    PARAMETER_TIER_ADVANCED,
    PARAMETER_TIER_INTELLIGENT_TIERING,
  ];

  static final $core.List<ParameterTier?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ParameterTier? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParameterTier._(super.value, super.name);
}

class ParameterType extends $pb.ProtobufEnum {
  static const ParameterType PARAMETER_TYPE_SECURE_STRING =
      ParameterType._(0, _omitEnumNames ? '' : 'PARAMETER_TYPE_SECURE_STRING');
  static const ParameterType PARAMETER_TYPE_STRING_LIST =
      ParameterType._(1, _omitEnumNames ? '' : 'PARAMETER_TYPE_STRING_LIST');
  static const ParameterType PARAMETER_TYPE_STRING =
      ParameterType._(2, _omitEnumNames ? '' : 'PARAMETER_TYPE_STRING');

  static const $core.List<ParameterType> values = <ParameterType>[
    PARAMETER_TYPE_SECURE_STRING,
    PARAMETER_TYPE_STRING_LIST,
    PARAMETER_TYPE_STRING,
  ];

  static final $core.List<ParameterType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ParameterType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParameterType._(super.value, super.name);
}

class ParametersFilterKey extends $pb.ProtobufEnum {
  static const ParametersFilterKey PARAMETERS_FILTER_KEY_KEY_ID =
      ParametersFilterKey._(
          0, _omitEnumNames ? '' : 'PARAMETERS_FILTER_KEY_KEY_ID');
  static const ParametersFilterKey PARAMETERS_FILTER_KEY_NAME =
      ParametersFilterKey._(
          1, _omitEnumNames ? '' : 'PARAMETERS_FILTER_KEY_NAME');
  static const ParametersFilterKey PARAMETERS_FILTER_KEY_TYPE =
      ParametersFilterKey._(
          2, _omitEnumNames ? '' : 'PARAMETERS_FILTER_KEY_TYPE');

  static const $core.List<ParametersFilterKey> values = <ParametersFilterKey>[
    PARAMETERS_FILTER_KEY_KEY_ID,
    PARAMETERS_FILTER_KEY_NAME,
    PARAMETERS_FILTER_KEY_TYPE,
  ];

  static final $core.List<ParametersFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ParametersFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ParametersFilterKey._(super.value, super.name);
}

class PatchAction extends $pb.ProtobufEnum {
  static const PatchAction PATCH_ACTION_ALLOWASDEPENDENCY =
      PatchAction._(0, _omitEnumNames ? '' : 'PATCH_ACTION_ALLOWASDEPENDENCY');
  static const PatchAction PATCH_ACTION_BLOCK =
      PatchAction._(1, _omitEnumNames ? '' : 'PATCH_ACTION_BLOCK');

  static const $core.List<PatchAction> values = <PatchAction>[
    PATCH_ACTION_ALLOWASDEPENDENCY,
    PATCH_ACTION_BLOCK,
  ];

  static final $core.List<PatchAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PatchAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchAction._(super.value, super.name);
}

class PatchComplianceDataState extends $pb.ProtobufEnum {
  static const PatchComplianceDataState
      PATCH_COMPLIANCE_DATA_STATE_AVAILABLESECURITYUPDATE =
      PatchComplianceDataState._(
          0,
          _omitEnumNames
              ? ''
              : 'PATCH_COMPLIANCE_DATA_STATE_AVAILABLESECURITYUPDATE');
  static const PatchComplianceDataState PATCH_COMPLIANCE_DATA_STATE_INSTALLED =
      PatchComplianceDataState._(
          1, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_DATA_STATE_INSTALLED');
  static const PatchComplianceDataState PATCH_COMPLIANCE_DATA_STATE_FAILED =
      PatchComplianceDataState._(
          2, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_DATA_STATE_FAILED');
  static const PatchComplianceDataState PATCH_COMPLIANCE_DATA_STATE_MISSING =
      PatchComplianceDataState._(
          3, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_DATA_STATE_MISSING');
  static const PatchComplianceDataState
      PATCH_COMPLIANCE_DATA_STATE_INSTALLEDPENDINGREBOOT =
      PatchComplianceDataState._(
          4,
          _omitEnumNames
              ? ''
              : 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDPENDINGREBOOT');
  static const PatchComplianceDataState
      PATCH_COMPLIANCE_DATA_STATE_INSTALLEDREJECTED =
      PatchComplianceDataState._(
          5,
          _omitEnumNames
              ? ''
              : 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDREJECTED');
  static const PatchComplianceDataState
      PATCH_COMPLIANCE_DATA_STATE_NOTAPPLICABLE = PatchComplianceDataState._(
          6, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_DATA_STATE_NOTAPPLICABLE');
  static const PatchComplianceDataState
      PATCH_COMPLIANCE_DATA_STATE_INSTALLEDOTHER = PatchComplianceDataState._(7,
          _omitEnumNames ? '' : 'PATCH_COMPLIANCE_DATA_STATE_INSTALLEDOTHER');

  static const $core.List<PatchComplianceDataState> values =
      <PatchComplianceDataState>[
    PATCH_COMPLIANCE_DATA_STATE_AVAILABLESECURITYUPDATE,
    PATCH_COMPLIANCE_DATA_STATE_INSTALLED,
    PATCH_COMPLIANCE_DATA_STATE_FAILED,
    PATCH_COMPLIANCE_DATA_STATE_MISSING,
    PATCH_COMPLIANCE_DATA_STATE_INSTALLEDPENDINGREBOOT,
    PATCH_COMPLIANCE_DATA_STATE_INSTALLEDREJECTED,
    PATCH_COMPLIANCE_DATA_STATE_NOTAPPLICABLE,
    PATCH_COMPLIANCE_DATA_STATE_INSTALLEDOTHER,
  ];

  static final $core.List<PatchComplianceDataState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static PatchComplianceDataState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchComplianceDataState._(super.value, super.name);
}

class PatchComplianceLevel extends $pb.ProtobufEnum {
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_INFORMATIONAL =
      PatchComplianceLevel._(
          0, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_INFORMATIONAL');
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_MEDIUM =
      PatchComplianceLevel._(
          1, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_MEDIUM');
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_UNSPECIFIED =
      PatchComplianceLevel._(
          2, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_UNSPECIFIED');
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_CRITICAL =
      PatchComplianceLevel._(
          3, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_CRITICAL');
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_LOW =
      PatchComplianceLevel._(
          4, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_LOW');
  static const PatchComplianceLevel PATCH_COMPLIANCE_LEVEL_HIGH =
      PatchComplianceLevel._(
          5, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_LEVEL_HIGH');

  static const $core.List<PatchComplianceLevel> values = <PatchComplianceLevel>[
    PATCH_COMPLIANCE_LEVEL_INFORMATIONAL,
    PATCH_COMPLIANCE_LEVEL_MEDIUM,
    PATCH_COMPLIANCE_LEVEL_UNSPECIFIED,
    PATCH_COMPLIANCE_LEVEL_CRITICAL,
    PATCH_COMPLIANCE_LEVEL_LOW,
    PATCH_COMPLIANCE_LEVEL_HIGH,
  ];

  static final $core.List<PatchComplianceLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static PatchComplianceLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchComplianceLevel._(super.value, super.name);
}

class PatchComplianceStatus extends $pb.ProtobufEnum {
  static const PatchComplianceStatus PATCH_COMPLIANCE_STATUS_COMPLIANT =
      PatchComplianceStatus._(
          0, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_STATUS_COMPLIANT');
  static const PatchComplianceStatus PATCH_COMPLIANCE_STATUS_NONCOMPLIANT =
      PatchComplianceStatus._(
          1, _omitEnumNames ? '' : 'PATCH_COMPLIANCE_STATUS_NONCOMPLIANT');

  static const $core.List<PatchComplianceStatus> values =
      <PatchComplianceStatus>[
    PATCH_COMPLIANCE_STATUS_COMPLIANT,
    PATCH_COMPLIANCE_STATUS_NONCOMPLIANT,
  ];

  static final $core.List<PatchComplianceStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PatchComplianceStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchComplianceStatus._(super.value, super.name);
}

class PatchDeploymentStatus extends $pb.ProtobufEnum {
  static const PatchDeploymentStatus PATCH_DEPLOYMENT_STATUS_APPROVED =
      PatchDeploymentStatus._(
          0, _omitEnumNames ? '' : 'PATCH_DEPLOYMENT_STATUS_APPROVED');
  static const PatchDeploymentStatus PATCH_DEPLOYMENT_STATUS_PENDINGAPPROVAL =
      PatchDeploymentStatus._(
          1, _omitEnumNames ? '' : 'PATCH_DEPLOYMENT_STATUS_PENDINGAPPROVAL');
  static const PatchDeploymentStatus PATCH_DEPLOYMENT_STATUS_EXPLICITREJECTED =
      PatchDeploymentStatus._(
          2, _omitEnumNames ? '' : 'PATCH_DEPLOYMENT_STATUS_EXPLICITREJECTED');
  static const PatchDeploymentStatus PATCH_DEPLOYMENT_STATUS_EXPLICITAPPROVED =
      PatchDeploymentStatus._(
          3, _omitEnumNames ? '' : 'PATCH_DEPLOYMENT_STATUS_EXPLICITAPPROVED');

  static const $core.List<PatchDeploymentStatus> values =
      <PatchDeploymentStatus>[
    PATCH_DEPLOYMENT_STATUS_APPROVED,
    PATCH_DEPLOYMENT_STATUS_PENDINGAPPROVAL,
    PATCH_DEPLOYMENT_STATUS_EXPLICITREJECTED,
    PATCH_DEPLOYMENT_STATUS_EXPLICITAPPROVED,
  ];

  static final $core.List<PatchDeploymentStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static PatchDeploymentStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchDeploymentStatus._(super.value, super.name);
}

class PatchFilterKey extends $pb.ProtobufEnum {
  static const PatchFilterKey PATCH_FILTER_KEY_PRIORITY =
      PatchFilterKey._(0, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_PRIORITY');
  static const PatchFilterKey PATCH_FILTER_KEY_PRODUCT =
      PatchFilterKey._(1, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_PRODUCT');
  static const PatchFilterKey PATCH_FILTER_KEY_MSRCSEVERITY = PatchFilterKey._(
      2, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_MSRCSEVERITY');
  static const PatchFilterKey PATCH_FILTER_KEY_ADVISORYID =
      PatchFilterKey._(3, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_ADVISORYID');
  static const PatchFilterKey PATCH_FILTER_KEY_PATCHID =
      PatchFilterKey._(4, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_PATCHID');
  static const PatchFilterKey PATCH_FILTER_KEY_PATCHSET =
      PatchFilterKey._(5, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_PATCHSET');
  static const PatchFilterKey PATCH_FILTER_KEY_CVEID =
      PatchFilterKey._(6, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_CVEID');
  static const PatchFilterKey PATCH_FILTER_KEY_RELEASE =
      PatchFilterKey._(7, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_RELEASE');
  static const PatchFilterKey PATCH_FILTER_KEY_REPOSITORY =
      PatchFilterKey._(8, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_REPOSITORY');
  static const PatchFilterKey PATCH_FILTER_KEY_ARCH =
      PatchFilterKey._(9, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_ARCH');
  static const PatchFilterKey PATCH_FILTER_KEY_CLASSIFICATION =
      PatchFilterKey._(
          10, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_CLASSIFICATION');
  static const PatchFilterKey PATCH_FILTER_KEY_VERSION =
      PatchFilterKey._(11, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_VERSION');
  static const PatchFilterKey PATCH_FILTER_KEY_SECTION =
      PatchFilterKey._(12, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_SECTION');
  static const PatchFilterKey PATCH_FILTER_KEY_EPOCH =
      PatchFilterKey._(13, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_EPOCH');
  static const PatchFilterKey PATCH_FILTER_KEY_PRODUCTFAMILY = PatchFilterKey._(
      14, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_PRODUCTFAMILY');
  static const PatchFilterKey PATCH_FILTER_KEY_SECURITY =
      PatchFilterKey._(15, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_SECURITY');
  static const PatchFilterKey PATCH_FILTER_KEY_NAME =
      PatchFilterKey._(16, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_NAME');
  static const PatchFilterKey PATCH_FILTER_KEY_BUGZILLAID =
      PatchFilterKey._(17, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_BUGZILLAID');
  static const PatchFilterKey PATCH_FILTER_KEY_SEVERITY =
      PatchFilterKey._(18, _omitEnumNames ? '' : 'PATCH_FILTER_KEY_SEVERITY');

  static const $core.List<PatchFilterKey> values = <PatchFilterKey>[
    PATCH_FILTER_KEY_PRIORITY,
    PATCH_FILTER_KEY_PRODUCT,
    PATCH_FILTER_KEY_MSRCSEVERITY,
    PATCH_FILTER_KEY_ADVISORYID,
    PATCH_FILTER_KEY_PATCHID,
    PATCH_FILTER_KEY_PATCHSET,
    PATCH_FILTER_KEY_CVEID,
    PATCH_FILTER_KEY_RELEASE,
    PATCH_FILTER_KEY_REPOSITORY,
    PATCH_FILTER_KEY_ARCH,
    PATCH_FILTER_KEY_CLASSIFICATION,
    PATCH_FILTER_KEY_VERSION,
    PATCH_FILTER_KEY_SECTION,
    PATCH_FILTER_KEY_EPOCH,
    PATCH_FILTER_KEY_PRODUCTFAMILY,
    PATCH_FILTER_KEY_SECURITY,
    PATCH_FILTER_KEY_NAME,
    PATCH_FILTER_KEY_BUGZILLAID,
    PATCH_FILTER_KEY_SEVERITY,
  ];

  static final $core.List<PatchFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 18);
  static PatchFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchFilterKey._(super.value, super.name);
}

class PatchOperationType extends $pb.ProtobufEnum {
  static const PatchOperationType PATCH_OPERATION_TYPE_SCAN =
      PatchOperationType._(
          0, _omitEnumNames ? '' : 'PATCH_OPERATION_TYPE_SCAN');
  static const PatchOperationType PATCH_OPERATION_TYPE_INSTALL =
      PatchOperationType._(
          1, _omitEnumNames ? '' : 'PATCH_OPERATION_TYPE_INSTALL');

  static const $core.List<PatchOperationType> values = <PatchOperationType>[
    PATCH_OPERATION_TYPE_SCAN,
    PATCH_OPERATION_TYPE_INSTALL,
  ];

  static final $core.List<PatchOperationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PatchOperationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchOperationType._(super.value, super.name);
}

class PatchProperty extends $pb.ProtobufEnum {
  static const PatchProperty PATCH_PROPERTY_PRODUCT =
      PatchProperty._(0, _omitEnumNames ? '' : 'PATCH_PROPERTY_PRODUCT');
  static const PatchProperty PATCH_PROPERTY_PATCHPRODUCTFAMILY =
      PatchProperty._(
          1, _omitEnumNames ? '' : 'PATCH_PROPERTY_PATCHPRODUCTFAMILY');
  static const PatchProperty PATCH_PROPERTY_PATCHCLASSIFICATION =
      PatchProperty._(
          2, _omitEnumNames ? '' : 'PATCH_PROPERTY_PATCHCLASSIFICATION');
  static const PatchProperty PATCH_PROPERTY_PATCHMSRCSEVERITY = PatchProperty._(
      3, _omitEnumNames ? '' : 'PATCH_PROPERTY_PATCHMSRCSEVERITY');
  static const PatchProperty PATCH_PROPERTY_PATCHPRIORITY =
      PatchProperty._(4, _omitEnumNames ? '' : 'PATCH_PROPERTY_PATCHPRIORITY');
  static const PatchProperty PATCH_PROPERTY_PATCHSEVERITY =
      PatchProperty._(5, _omitEnumNames ? '' : 'PATCH_PROPERTY_PATCHSEVERITY');

  static const $core.List<PatchProperty> values = <PatchProperty>[
    PATCH_PROPERTY_PRODUCT,
    PATCH_PROPERTY_PATCHPRODUCTFAMILY,
    PATCH_PROPERTY_PATCHCLASSIFICATION,
    PATCH_PROPERTY_PATCHMSRCSEVERITY,
    PATCH_PROPERTY_PATCHPRIORITY,
    PATCH_PROPERTY_PATCHSEVERITY,
  ];

  static final $core.List<PatchProperty?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static PatchProperty? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchProperty._(super.value, super.name);
}

class PatchSet extends $pb.ProtobufEnum {
  static const PatchSet PATCH_SET_OS =
      PatchSet._(0, _omitEnumNames ? '' : 'PATCH_SET_OS');
  static const PatchSet PATCH_SET_APPLICATION =
      PatchSet._(1, _omitEnumNames ? '' : 'PATCH_SET_APPLICATION');

  static const $core.List<PatchSet> values = <PatchSet>[
    PATCH_SET_OS,
    PATCH_SET_APPLICATION,
  ];

  static final $core.List<PatchSet?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PatchSet? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PatchSet._(super.value, super.name);
}

class PingStatus extends $pb.ProtobufEnum {
  static const PingStatus PING_STATUS_ONLINE =
      PingStatus._(0, _omitEnumNames ? '' : 'PING_STATUS_ONLINE');
  static const PingStatus PING_STATUS_CONNECTION_LOST =
      PingStatus._(1, _omitEnumNames ? '' : 'PING_STATUS_CONNECTION_LOST');
  static const PingStatus PING_STATUS_INACTIVE =
      PingStatus._(2, _omitEnumNames ? '' : 'PING_STATUS_INACTIVE');

  static const $core.List<PingStatus> values = <PingStatus>[
    PING_STATUS_ONLINE,
    PING_STATUS_CONNECTION_LOST,
    PING_STATUS_INACTIVE,
  ];

  static final $core.List<PingStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PingStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PingStatus._(super.value, super.name);
}

class PlatformType extends $pb.ProtobufEnum {
  static const PlatformType PLATFORM_TYPE_LINUX =
      PlatformType._(0, _omitEnumNames ? '' : 'PLATFORM_TYPE_LINUX');
  static const PlatformType PLATFORM_TYPE_MACOS =
      PlatformType._(1, _omitEnumNames ? '' : 'PLATFORM_TYPE_MACOS');
  static const PlatformType PLATFORM_TYPE_WINDOWS =
      PlatformType._(2, _omitEnumNames ? '' : 'PLATFORM_TYPE_WINDOWS');

  static const $core.List<PlatformType> values = <PlatformType>[
    PLATFORM_TYPE_LINUX,
    PLATFORM_TYPE_MACOS,
    PLATFORM_TYPE_WINDOWS,
  ];

  static final $core.List<PlatformType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PlatformType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PlatformType._(super.value, super.name);
}

class RebootOption extends $pb.ProtobufEnum {
  static const RebootOption REBOOT_OPTION_REBOOT_IF_NEEDED =
      RebootOption._(0, _omitEnumNames ? '' : 'REBOOT_OPTION_REBOOT_IF_NEEDED');
  static const RebootOption REBOOT_OPTION_NO_REBOOT =
      RebootOption._(1, _omitEnumNames ? '' : 'REBOOT_OPTION_NO_REBOOT');

  static const $core.List<RebootOption> values = <RebootOption>[
    REBOOT_OPTION_REBOOT_IF_NEEDED,
    REBOOT_OPTION_NO_REBOOT,
  ];

  static final $core.List<RebootOption?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RebootOption? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RebootOption._(super.value, super.name);
}

class ResourceDataSyncS3Format extends $pb.ProtobufEnum {
  static const ResourceDataSyncS3Format
      RESOURCE_DATA_SYNC_S3_FORMAT_JSON_SERDE = ResourceDataSyncS3Format._(
          0, _omitEnumNames ? '' : 'RESOURCE_DATA_SYNC_S3_FORMAT_JSON_SERDE');

  static const $core.List<ResourceDataSyncS3Format> values =
      <ResourceDataSyncS3Format>[
    RESOURCE_DATA_SYNC_S3_FORMAT_JSON_SERDE,
  ];

  static final $core.List<ResourceDataSyncS3Format?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ResourceDataSyncS3Format? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceDataSyncS3Format._(super.value, super.name);
}

class ResourceType extends $pb.ProtobufEnum {
  static const ResourceType RESOURCE_TYPE_EC2_INSTANCE =
      ResourceType._(0, _omitEnumNames ? '' : 'RESOURCE_TYPE_EC2_INSTANCE');
  static const ResourceType RESOURCE_TYPE_MANAGED_INSTANCE =
      ResourceType._(1, _omitEnumNames ? '' : 'RESOURCE_TYPE_MANAGED_INSTANCE');

  static const $core.List<ResourceType> values = <ResourceType>[
    RESOURCE_TYPE_EC2_INSTANCE,
    RESOURCE_TYPE_MANAGED_INSTANCE,
  ];

  static final $core.List<ResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceType._(super.value, super.name);
}

class ResourceTypeForTagging extends $pb.ProtobufEnum {
  static const ResourceTypeForTagging
      RESOURCE_TYPE_FOR_TAGGING_MAINTENANCE_WINDOW = ResourceTypeForTagging._(0,
          _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_MAINTENANCE_WINDOW');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_OPS_ITEM =
      ResourceTypeForTagging._(
          1, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_OPS_ITEM');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_DOCUMENT =
      ResourceTypeForTagging._(
          2, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_DOCUMENT');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_AUTOMATION =
      ResourceTypeForTagging._(
          3, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_AUTOMATION');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_PATCH_BASELINE =
      ResourceTypeForTagging._(
          4, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_PATCH_BASELINE');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_OPSMETADATA =
      ResourceTypeForTagging._(
          5, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_OPSMETADATA');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_PARAMETER =
      ResourceTypeForTagging._(
          6, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_PARAMETER');
  static const ResourceTypeForTagging
      RESOURCE_TYPE_FOR_TAGGING_MANAGED_INSTANCE = ResourceTypeForTagging._(7,
          _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_MANAGED_INSTANCE');
  static const ResourceTypeForTagging RESOURCE_TYPE_FOR_TAGGING_ASSOCIATION =
      ResourceTypeForTagging._(
          8, _omitEnumNames ? '' : 'RESOURCE_TYPE_FOR_TAGGING_ASSOCIATION');

  static const $core.List<ResourceTypeForTagging> values =
      <ResourceTypeForTagging>[
    RESOURCE_TYPE_FOR_TAGGING_MAINTENANCE_WINDOW,
    RESOURCE_TYPE_FOR_TAGGING_OPS_ITEM,
    RESOURCE_TYPE_FOR_TAGGING_DOCUMENT,
    RESOURCE_TYPE_FOR_TAGGING_AUTOMATION,
    RESOURCE_TYPE_FOR_TAGGING_PATCH_BASELINE,
    RESOURCE_TYPE_FOR_TAGGING_OPSMETADATA,
    RESOURCE_TYPE_FOR_TAGGING_PARAMETER,
    RESOURCE_TYPE_FOR_TAGGING_MANAGED_INSTANCE,
    RESOURCE_TYPE_FOR_TAGGING_ASSOCIATION,
  ];

  static final $core.List<ResourceTypeForTagging?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static ResourceTypeForTagging? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceTypeForTagging._(super.value, super.name);
}

class ReviewStatus extends $pb.ProtobufEnum {
  static const ReviewStatus REVIEW_STATUS_PENDING =
      ReviewStatus._(0, _omitEnumNames ? '' : 'REVIEW_STATUS_PENDING');
  static const ReviewStatus REVIEW_STATUS_REJECTED =
      ReviewStatus._(1, _omitEnumNames ? '' : 'REVIEW_STATUS_REJECTED');
  static const ReviewStatus REVIEW_STATUS_NOT_REVIEWED =
      ReviewStatus._(2, _omitEnumNames ? '' : 'REVIEW_STATUS_NOT_REVIEWED');
  static const ReviewStatus REVIEW_STATUS_APPROVED =
      ReviewStatus._(3, _omitEnumNames ? '' : 'REVIEW_STATUS_APPROVED');

  static const $core.List<ReviewStatus> values = <ReviewStatus>[
    REVIEW_STATUS_PENDING,
    REVIEW_STATUS_REJECTED,
    REVIEW_STATUS_NOT_REVIEWED,
    REVIEW_STATUS_APPROVED,
  ];

  static final $core.List<ReviewStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ReviewStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReviewStatus._(super.value, super.name);
}

class SessionFilterKey extends $pb.ProtobufEnum {
  static const SessionFilterKey SESSION_FILTER_KEY_SESSION_ID =
      SessionFilterKey._(
          0, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_SESSION_ID');
  static const SessionFilterKey SESSION_FILTER_KEY_STATUS =
      SessionFilterKey._(1, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_STATUS');
  static const SessionFilterKey SESSION_FILTER_KEY_ACCESS_TYPE =
      SessionFilterKey._(
          2, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_ACCESS_TYPE');
  static const SessionFilterKey SESSION_FILTER_KEY_INVOKED_BEFORE =
      SessionFilterKey._(
          3, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_INVOKED_BEFORE');
  static const SessionFilterKey SESSION_FILTER_KEY_OWNER =
      SessionFilterKey._(4, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_OWNER');
  static const SessionFilterKey SESSION_FILTER_KEY_TARGET_ID =
      SessionFilterKey._(
          5, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_TARGET_ID');
  static const SessionFilterKey SESSION_FILTER_KEY_INVOKED_AFTER =
      SessionFilterKey._(
          6, _omitEnumNames ? '' : 'SESSION_FILTER_KEY_INVOKED_AFTER');

  static const $core.List<SessionFilterKey> values = <SessionFilterKey>[
    SESSION_FILTER_KEY_SESSION_ID,
    SESSION_FILTER_KEY_STATUS,
    SESSION_FILTER_KEY_ACCESS_TYPE,
    SESSION_FILTER_KEY_INVOKED_BEFORE,
    SESSION_FILTER_KEY_OWNER,
    SESSION_FILTER_KEY_TARGET_ID,
    SESSION_FILTER_KEY_INVOKED_AFTER,
  ];

  static final $core.List<SessionFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static SessionFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SessionFilterKey._(super.value, super.name);
}

class SessionState extends $pb.ProtobufEnum {
  static const SessionState SESSION_STATE_HISTORY =
      SessionState._(0, _omitEnumNames ? '' : 'SESSION_STATE_HISTORY');
  static const SessionState SESSION_STATE_ACTIVE =
      SessionState._(1, _omitEnumNames ? '' : 'SESSION_STATE_ACTIVE');

  static const $core.List<SessionState> values = <SessionState>[
    SESSION_STATE_HISTORY,
    SESSION_STATE_ACTIVE,
  ];

  static final $core.List<SessionState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SessionState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SessionState._(super.value, super.name);
}

class SessionStatus extends $pb.ProtobufEnum {
  static const SessionStatus SESSION_STATUS_TERMINATING =
      SessionStatus._(0, _omitEnumNames ? '' : 'SESSION_STATUS_TERMINATING');
  static const SessionStatus SESSION_STATUS_CONNECTING =
      SessionStatus._(1, _omitEnumNames ? '' : 'SESSION_STATUS_CONNECTING');
  static const SessionStatus SESSION_STATUS_TERMINATED =
      SessionStatus._(2, _omitEnumNames ? '' : 'SESSION_STATUS_TERMINATED');
  static const SessionStatus SESSION_STATUS_DISCONNECTED =
      SessionStatus._(3, _omitEnumNames ? '' : 'SESSION_STATUS_DISCONNECTED');
  static const SessionStatus SESSION_STATUS_CONNECTED =
      SessionStatus._(4, _omitEnumNames ? '' : 'SESSION_STATUS_CONNECTED');
  static const SessionStatus SESSION_STATUS_FAILED =
      SessionStatus._(5, _omitEnumNames ? '' : 'SESSION_STATUS_FAILED');

  static const $core.List<SessionStatus> values = <SessionStatus>[
    SESSION_STATUS_TERMINATING,
    SESSION_STATUS_CONNECTING,
    SESSION_STATUS_TERMINATED,
    SESSION_STATUS_DISCONNECTED,
    SESSION_STATUS_CONNECTED,
    SESSION_STATUS_FAILED,
  ];

  static final $core.List<SessionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static SessionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SessionStatus._(super.value, super.name);
}

class SignalType extends $pb.ProtobufEnum {
  static const SignalType SIGNAL_TYPE_REVOKE =
      SignalType._(0, _omitEnumNames ? '' : 'SIGNAL_TYPE_REVOKE');
  static const SignalType SIGNAL_TYPE_STOP_STEP =
      SignalType._(1, _omitEnumNames ? '' : 'SIGNAL_TYPE_STOP_STEP');
  static const SignalType SIGNAL_TYPE_START_STEP =
      SignalType._(2, _omitEnumNames ? '' : 'SIGNAL_TYPE_START_STEP');
  static const SignalType SIGNAL_TYPE_REJECT =
      SignalType._(3, _omitEnumNames ? '' : 'SIGNAL_TYPE_REJECT');
  static const SignalType SIGNAL_TYPE_APPROVE =
      SignalType._(4, _omitEnumNames ? '' : 'SIGNAL_TYPE_APPROVE');
  static const SignalType SIGNAL_TYPE_RESUME =
      SignalType._(5, _omitEnumNames ? '' : 'SIGNAL_TYPE_RESUME');

  static const $core.List<SignalType> values = <SignalType>[
    SIGNAL_TYPE_REVOKE,
    SIGNAL_TYPE_STOP_STEP,
    SIGNAL_TYPE_START_STEP,
    SIGNAL_TYPE_REJECT,
    SIGNAL_TYPE_APPROVE,
    SIGNAL_TYPE_RESUME,
  ];

  static final $core.List<SignalType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static SignalType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SignalType._(super.value, super.name);
}

class SourceType extends $pb.ProtobufEnum {
  static const SourceType SOURCE_TYPE_AWS_SSM_MANAGEDINSTANCE = SourceType._(
      0, _omitEnumNames ? '' : 'SOURCE_TYPE_AWS_SSM_MANAGEDINSTANCE');
  static const SourceType SOURCE_TYPE_AWS_EC2_INSTANCE =
      SourceType._(1, _omitEnumNames ? '' : 'SOURCE_TYPE_AWS_EC2_INSTANCE');
  static const SourceType SOURCE_TYPE_AWS_IOT_THING =
      SourceType._(2, _omitEnumNames ? '' : 'SOURCE_TYPE_AWS_IOT_THING');

  static const $core.List<SourceType> values = <SourceType>[
    SOURCE_TYPE_AWS_SSM_MANAGEDINSTANCE,
    SOURCE_TYPE_AWS_EC2_INSTANCE,
    SOURCE_TYPE_AWS_IOT_THING,
  ];

  static final $core.List<SourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SourceType._(super.value, super.name);
}

class StepExecutionFilterKey extends $pb.ProtobufEnum {
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_START_TIME_AFTER = StepExecutionFilterKey._(0,
          _omitEnumNames ? '' : 'STEP_EXECUTION_FILTER_KEY_START_TIME_AFTER');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_STATUS =
      StepExecutionFilterKey._(
          1,
          _omitEnumNames
              ? ''
              : 'STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_STATUS');
  static const StepExecutionFilterKey STEP_EXECUTION_FILTER_KEY_STEP_NAME =
      StepExecutionFilterKey._(
          2, _omitEnumNames ? '' : 'STEP_EXECUTION_FILTER_KEY_STEP_NAME');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_ID = StepExecutionFilterKey._(3,
          _omitEnumNames ? '' : 'STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_ID');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATION =
      StepExecutionFilterKey._(
          4,
          _omitEnumNames
              ? ''
              : 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATION');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_START_TIME_BEFORE = StepExecutionFilterKey._(5,
          _omitEnumNames ? '' : 'STEP_EXECUTION_FILTER_KEY_START_TIME_BEFORE');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATOR_VALUE =
      StepExecutionFilterKey._(
          6,
          _omitEnumNames
              ? ''
              : 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATOR_VALUE');
  static const StepExecutionFilterKey
      STEP_EXECUTION_FILTER_KEY_PARENT_STEP_EXECUTION_ID =
      StepExecutionFilterKey._(
          7,
          _omitEnumNames
              ? ''
              : 'STEP_EXECUTION_FILTER_KEY_PARENT_STEP_EXECUTION_ID');
  static const StepExecutionFilterKey STEP_EXECUTION_FILTER_KEY_ACTION =
      StepExecutionFilterKey._(
          8, _omitEnumNames ? '' : 'STEP_EXECUTION_FILTER_KEY_ACTION');

  static const $core.List<StepExecutionFilterKey> values =
      <StepExecutionFilterKey>[
    STEP_EXECUTION_FILTER_KEY_START_TIME_AFTER,
    STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_STATUS,
    STEP_EXECUTION_FILTER_KEY_STEP_NAME,
    STEP_EXECUTION_FILTER_KEY_STEP_EXECUTION_ID,
    STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATION,
    STEP_EXECUTION_FILTER_KEY_START_TIME_BEFORE,
    STEP_EXECUTION_FILTER_KEY_PARENT_STEP_ITERATOR_VALUE,
    STEP_EXECUTION_FILTER_KEY_PARENT_STEP_EXECUTION_ID,
    STEP_EXECUTION_FILTER_KEY_ACTION,
  ];

  static final $core.List<StepExecutionFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static StepExecutionFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StepExecutionFilterKey._(super.value, super.name);
}

class StopType extends $pb.ProtobufEnum {
  static const StopType STOP_TYPE_COMPLETE =
      StopType._(0, _omitEnumNames ? '' : 'STOP_TYPE_COMPLETE');
  static const StopType STOP_TYPE_CANCEL =
      StopType._(1, _omitEnumNames ? '' : 'STOP_TYPE_CANCEL');

  static const $core.List<StopType> values = <StopType>[
    STOP_TYPE_COMPLETE,
    STOP_TYPE_CANCEL,
  ];

  static final $core.List<StopType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static StopType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StopType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
