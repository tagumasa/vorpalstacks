// This is a generated file - do not edit.
//
// Generated from iam.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AccessAdvisorUsageGranularityType extends $pb.ProtobufEnum {
  static const AccessAdvisorUsageGranularityType
      ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_SERVICE_LEVEL =
      AccessAdvisorUsageGranularityType._(
          0,
          _omitEnumNames
              ? ''
              : 'ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_SERVICE_LEVEL');
  static const AccessAdvisorUsageGranularityType
      ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_ACTION_LEVEL =
      AccessAdvisorUsageGranularityType._(
          1,
          _omitEnumNames
              ? ''
              : 'ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_ACTION_LEVEL');

  static const $core.List<AccessAdvisorUsageGranularityType> values =
      <AccessAdvisorUsageGranularityType>[
    ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_SERVICE_LEVEL,
    ACCESS_ADVISOR_USAGE_GRANULARITY_TYPE_ACTION_LEVEL,
  ];

  static final $core.List<AccessAdvisorUsageGranularityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AccessAdvisorUsageGranularityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AccessAdvisorUsageGranularityType._(super.value, super.name);
}

class ContextKeyTypeEnum extends $pb.ProtobufEnum {
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_NUMERIC_LIST =
      ContextKeyTypeEnum._(
          0, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_NUMERIC_LIST');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_BOOLEAN_LIST =
      ContextKeyTypeEnum._(
          1, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_BOOLEAN_LIST');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_BINARY =
      ContextKeyTypeEnum._(
          2, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_BINARY');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_DATE_LIST =
      ContextKeyTypeEnum._(
          3, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_DATE_LIST');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_NUMERIC =
      ContextKeyTypeEnum._(
          4, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_NUMERIC');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_STRING_LIST =
      ContextKeyTypeEnum._(
          5, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_STRING_LIST');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_STRING =
      ContextKeyTypeEnum._(
          6, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_STRING');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_IP =
      ContextKeyTypeEnum._(7, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_IP');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_BOOLEAN =
      ContextKeyTypeEnum._(
          8, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_BOOLEAN');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_IP_LIST =
      ContextKeyTypeEnum._(
          9, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_IP_LIST');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_DATE =
      ContextKeyTypeEnum._(
          10, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_DATE');
  static const ContextKeyTypeEnum CONTEXT_KEY_TYPE_ENUM_BINARY_LIST =
      ContextKeyTypeEnum._(
          11, _omitEnumNames ? '' : 'CONTEXT_KEY_TYPE_ENUM_BINARY_LIST');

  static const $core.List<ContextKeyTypeEnum> values = <ContextKeyTypeEnum>[
    CONTEXT_KEY_TYPE_ENUM_NUMERIC_LIST,
    CONTEXT_KEY_TYPE_ENUM_BOOLEAN_LIST,
    CONTEXT_KEY_TYPE_ENUM_BINARY,
    CONTEXT_KEY_TYPE_ENUM_DATE_LIST,
    CONTEXT_KEY_TYPE_ENUM_NUMERIC,
    CONTEXT_KEY_TYPE_ENUM_STRING_LIST,
    CONTEXT_KEY_TYPE_ENUM_STRING,
    CONTEXT_KEY_TYPE_ENUM_IP,
    CONTEXT_KEY_TYPE_ENUM_BOOLEAN,
    CONTEXT_KEY_TYPE_ENUM_IP_LIST,
    CONTEXT_KEY_TYPE_ENUM_DATE,
    CONTEXT_KEY_TYPE_ENUM_BINARY_LIST,
  ];

  static final $core.List<ContextKeyTypeEnum?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 11);
  static ContextKeyTypeEnum? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContextKeyTypeEnum._(super.value, super.name);
}

class DeletionTaskStatusType extends $pb.ProtobufEnum {
  static const DeletionTaskStatusType DELETION_TASK_STATUS_TYPE_IN_PROGRESS =
      DeletionTaskStatusType._(
          0, _omitEnumNames ? '' : 'DELETION_TASK_STATUS_TYPE_IN_PROGRESS');
  static const DeletionTaskStatusType DELETION_TASK_STATUS_TYPE_SUCCEEDED =
      DeletionTaskStatusType._(
          1, _omitEnumNames ? '' : 'DELETION_TASK_STATUS_TYPE_SUCCEEDED');
  static const DeletionTaskStatusType DELETION_TASK_STATUS_TYPE_FAILED =
      DeletionTaskStatusType._(
          2, _omitEnumNames ? '' : 'DELETION_TASK_STATUS_TYPE_FAILED');
  static const DeletionTaskStatusType DELETION_TASK_STATUS_TYPE_NOT_STARTED =
      DeletionTaskStatusType._(
          3, _omitEnumNames ? '' : 'DELETION_TASK_STATUS_TYPE_NOT_STARTED');

  static const $core.List<DeletionTaskStatusType> values =
      <DeletionTaskStatusType>[
    DELETION_TASK_STATUS_TYPE_IN_PROGRESS,
    DELETION_TASK_STATUS_TYPE_SUCCEEDED,
    DELETION_TASK_STATUS_TYPE_FAILED,
    DELETION_TASK_STATUS_TYPE_NOT_STARTED,
  ];

  static final $core.List<DeletionTaskStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static DeletionTaskStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeletionTaskStatusType._(super.value, super.name);
}

class EntityType extends $pb.ProtobufEnum {
  static const EntityType ENTITY_TYPE_GROUP =
      EntityType._(0, _omitEnumNames ? '' : 'ENTITY_TYPE_GROUP');
  static const EntityType ENTITY_TYPE_LOCALMANAGEDPOLICY =
      EntityType._(1, _omitEnumNames ? '' : 'ENTITY_TYPE_LOCALMANAGEDPOLICY');
  static const EntityType ENTITY_TYPE_USER =
      EntityType._(2, _omitEnumNames ? '' : 'ENTITY_TYPE_USER');
  static const EntityType ENTITY_TYPE_AWSMANAGEDPOLICY =
      EntityType._(3, _omitEnumNames ? '' : 'ENTITY_TYPE_AWSMANAGEDPOLICY');
  static const EntityType ENTITY_TYPE_ROLE =
      EntityType._(4, _omitEnumNames ? '' : 'ENTITY_TYPE_ROLE');

  static const $core.List<EntityType> values = <EntityType>[
    ENTITY_TYPE_GROUP,
    ENTITY_TYPE_LOCALMANAGEDPOLICY,
    ENTITY_TYPE_USER,
    ENTITY_TYPE_AWSMANAGEDPOLICY,
    ENTITY_TYPE_ROLE,
  ];

  static final $core.List<EntityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static EntityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EntityType._(super.value, super.name);
}

class FeatureType extends $pb.ProtobufEnum {
  static const FeatureType FEATURE_TYPE_ROOT_SESSIONS =
      FeatureType._(0, _omitEnumNames ? '' : 'FEATURE_TYPE_ROOT_SESSIONS');
  static const FeatureType FEATURE_TYPE_ROOT_CREDENTIALS_MANAGEMENT =
      FeatureType._(
          1, _omitEnumNames ? '' : 'FEATURE_TYPE_ROOT_CREDENTIALS_MANAGEMENT');

  static const $core.List<FeatureType> values = <FeatureType>[
    FEATURE_TYPE_ROOT_SESSIONS,
    FEATURE_TYPE_ROOT_CREDENTIALS_MANAGEMENT,
  ];

  static final $core.List<FeatureType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FeatureType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FeatureType._(super.value, super.name);
}

class PermissionsBoundaryAttachmentType extends $pb.ProtobufEnum {
  static const PermissionsBoundaryAttachmentType
      PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY =
      PermissionsBoundaryAttachmentType._(0,
          _omitEnumNames ? '' : 'PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY');

  static const $core.List<PermissionsBoundaryAttachmentType> values =
      <PermissionsBoundaryAttachmentType>[
    PERMISSIONS_BOUNDARY_ATTACHMENT_TYPE_POLICY,
  ];

  static final $core.List<PermissionsBoundaryAttachmentType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static PermissionsBoundaryAttachmentType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PermissionsBoundaryAttachmentType._(super.value, super.name);
}

class PolicyEvaluationDecisionType extends $pb.ProtobufEnum {
  static const PolicyEvaluationDecisionType
      POLICY_EVALUATION_DECISION_TYPE_IMPLICIT_DENY =
      PolicyEvaluationDecisionType._(
          0,
          _omitEnumNames
              ? ''
              : 'POLICY_EVALUATION_DECISION_TYPE_IMPLICIT_DENY');
  static const PolicyEvaluationDecisionType
      POLICY_EVALUATION_DECISION_TYPE_EXPLICIT_DENY =
      PolicyEvaluationDecisionType._(
          1,
          _omitEnumNames
              ? ''
              : 'POLICY_EVALUATION_DECISION_TYPE_EXPLICIT_DENY');
  static const PolicyEvaluationDecisionType
      POLICY_EVALUATION_DECISION_TYPE_ALLOWED = PolicyEvaluationDecisionType._(
          2, _omitEnumNames ? '' : 'POLICY_EVALUATION_DECISION_TYPE_ALLOWED');

  static const $core.List<PolicyEvaluationDecisionType> values =
      <PolicyEvaluationDecisionType>[
    POLICY_EVALUATION_DECISION_TYPE_IMPLICIT_DENY,
    POLICY_EVALUATION_DECISION_TYPE_EXPLICIT_DENY,
    POLICY_EVALUATION_DECISION_TYPE_ALLOWED,
  ];

  static final $core.List<PolicyEvaluationDecisionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PolicyEvaluationDecisionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicyEvaluationDecisionType._(super.value, super.name);
}

class PolicyParameterTypeEnum extends $pb.ProtobufEnum {
  static const PolicyParameterTypeEnum POLICY_PARAMETER_TYPE_ENUM_STRING_LIST =
      PolicyParameterTypeEnum._(
          0, _omitEnumNames ? '' : 'POLICY_PARAMETER_TYPE_ENUM_STRING_LIST');
  static const PolicyParameterTypeEnum POLICY_PARAMETER_TYPE_ENUM_STRING =
      PolicyParameterTypeEnum._(
          1, _omitEnumNames ? '' : 'POLICY_PARAMETER_TYPE_ENUM_STRING');

  static const $core.List<PolicyParameterTypeEnum> values =
      <PolicyParameterTypeEnum>[
    POLICY_PARAMETER_TYPE_ENUM_STRING_LIST,
    POLICY_PARAMETER_TYPE_ENUM_STRING,
  ];

  static final $core.List<PolicyParameterTypeEnum?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PolicyParameterTypeEnum? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicyParameterTypeEnum._(super.value, super.name);
}

class PolicySourceType extends $pb.ProtobufEnum {
  static const PolicySourceType POLICY_SOURCE_TYPE_GROUP =
      PolicySourceType._(0, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_GROUP');
  static const PolicySourceType POLICY_SOURCE_TYPE_USER_MANAGED =
      PolicySourceType._(
          1, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_USER_MANAGED');
  static const PolicySourceType POLICY_SOURCE_TYPE_AWS_MANAGED =
      PolicySourceType._(
          2, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_AWS_MANAGED');
  static const PolicySourceType POLICY_SOURCE_TYPE_NONE =
      PolicySourceType._(3, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_NONE');
  static const PolicySourceType POLICY_SOURCE_TYPE_RESOURCE =
      PolicySourceType._(
          4, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_RESOURCE');
  static const PolicySourceType POLICY_SOURCE_TYPE_ROLE =
      PolicySourceType._(5, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_ROLE');
  static const PolicySourceType POLICY_SOURCE_TYPE_USER =
      PolicySourceType._(6, _omitEnumNames ? '' : 'POLICY_SOURCE_TYPE_USER');

  static const $core.List<PolicySourceType> values = <PolicySourceType>[
    POLICY_SOURCE_TYPE_GROUP,
    POLICY_SOURCE_TYPE_USER_MANAGED,
    POLICY_SOURCE_TYPE_AWS_MANAGED,
    POLICY_SOURCE_TYPE_NONE,
    POLICY_SOURCE_TYPE_RESOURCE,
    POLICY_SOURCE_TYPE_ROLE,
    POLICY_SOURCE_TYPE_USER,
  ];

  static final $core.List<PolicySourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static PolicySourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicySourceType._(super.value, super.name);
}

class PolicyUsageType extends $pb.ProtobufEnum {
  static const PolicyUsageType POLICY_USAGE_TYPE_PERMISSIONSPOLICY =
      PolicyUsageType._(
          0, _omitEnumNames ? '' : 'POLICY_USAGE_TYPE_PERMISSIONSPOLICY');
  static const PolicyUsageType POLICY_USAGE_TYPE_PERMISSIONSBOUNDARY =
      PolicyUsageType._(
          1, _omitEnumNames ? '' : 'POLICY_USAGE_TYPE_PERMISSIONSBOUNDARY');

  static const $core.List<PolicyUsageType> values = <PolicyUsageType>[
    POLICY_USAGE_TYPE_PERMISSIONSPOLICY,
    POLICY_USAGE_TYPE_PERMISSIONSBOUNDARY,
  ];

  static final $core.List<PolicyUsageType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static PolicyUsageType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PolicyUsageType._(super.value, super.name);
}

class ReportFormatType extends $pb.ProtobufEnum {
  static const ReportFormatType REPORT_FORMAT_TYPE_TEXT_CSV =
      ReportFormatType._(
          0, _omitEnumNames ? '' : 'REPORT_FORMAT_TYPE_TEXT_CSV');

  static const $core.List<ReportFormatType> values = <ReportFormatType>[
    REPORT_FORMAT_TYPE_TEXT_CSV,
  ];

  static final $core.List<ReportFormatType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ReportFormatType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReportFormatType._(super.value, super.name);
}

class ReportStateType extends $pb.ProtobufEnum {
  static const ReportStateType REPORT_STATE_TYPE_STARTED =
      ReportStateType._(0, _omitEnumNames ? '' : 'REPORT_STATE_TYPE_STARTED');
  static const ReportStateType REPORT_STATE_TYPE_INPROGRESS = ReportStateType._(
      1, _omitEnumNames ? '' : 'REPORT_STATE_TYPE_INPROGRESS');
  static const ReportStateType REPORT_STATE_TYPE_COMPLETE =
      ReportStateType._(2, _omitEnumNames ? '' : 'REPORT_STATE_TYPE_COMPLETE');

  static const $core.List<ReportStateType> values = <ReportStateType>[
    REPORT_STATE_TYPE_STARTED,
    REPORT_STATE_TYPE_INPROGRESS,
    REPORT_STATE_TYPE_COMPLETE,
  ];

  static final $core.List<ReportStateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ReportStateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReportStateType._(super.value, super.name);
}

class assertionEncryptionModeType extends $pb.ProtobufEnum {
  static const assertionEncryptionModeType
      ASSERTION_ENCRYPTION_MODE_TYPE_REQUIRED = assertionEncryptionModeType._(
          0, _omitEnumNames ? '' : 'ASSERTION_ENCRYPTION_MODE_TYPE_REQUIRED');
  static const assertionEncryptionModeType
      ASSERTION_ENCRYPTION_MODE_TYPE_ALLOWED = assertionEncryptionModeType._(
          1, _omitEnumNames ? '' : 'ASSERTION_ENCRYPTION_MODE_TYPE_ALLOWED');

  static const $core.List<assertionEncryptionModeType> values =
      <assertionEncryptionModeType>[
    ASSERTION_ENCRYPTION_MODE_TYPE_REQUIRED,
    ASSERTION_ENCRYPTION_MODE_TYPE_ALLOWED,
  ];

  static final $core.List<assertionEncryptionModeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static assertionEncryptionModeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const assertionEncryptionModeType._(super.value, super.name);
}

class assignmentStatusType extends $pb.ProtobufEnum {
  static const assignmentStatusType ASSIGNMENT_STATUS_TYPE_ANY =
      assignmentStatusType._(
          0, _omitEnumNames ? '' : 'ASSIGNMENT_STATUS_TYPE_ANY');
  static const assignmentStatusType ASSIGNMENT_STATUS_TYPE_UNASSIGNED =
      assignmentStatusType._(
          1, _omitEnumNames ? '' : 'ASSIGNMENT_STATUS_TYPE_UNASSIGNED');
  static const assignmentStatusType ASSIGNMENT_STATUS_TYPE_ASSIGNED =
      assignmentStatusType._(
          2, _omitEnumNames ? '' : 'ASSIGNMENT_STATUS_TYPE_ASSIGNED');

  static const $core.List<assignmentStatusType> values = <assignmentStatusType>[
    ASSIGNMENT_STATUS_TYPE_ANY,
    ASSIGNMENT_STATUS_TYPE_UNASSIGNED,
    ASSIGNMENT_STATUS_TYPE_ASSIGNED,
  ];

  static final $core.List<assignmentStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static assignmentStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const assignmentStatusType._(super.value, super.name);
}

class encodingType extends $pb.ProtobufEnum {
  static const encodingType ENCODING_TYPE_PEM =
      encodingType._(0, _omitEnumNames ? '' : 'ENCODING_TYPE_PEM');
  static const encodingType ENCODING_TYPE_SSH =
      encodingType._(1, _omitEnumNames ? '' : 'ENCODING_TYPE_SSH');

  static const $core.List<encodingType> values = <encodingType>[
    ENCODING_TYPE_PEM,
    ENCODING_TYPE_SSH,
  ];

  static final $core.List<encodingType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static encodingType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const encodingType._(super.value, super.name);
}

class globalEndpointTokenVersion extends $pb.ProtobufEnum {
  static const globalEndpointTokenVersion
      GLOBAL_ENDPOINT_TOKEN_VERSION_V2TOKEN = globalEndpointTokenVersion._(
          0, _omitEnumNames ? '' : 'GLOBAL_ENDPOINT_TOKEN_VERSION_V2TOKEN');
  static const globalEndpointTokenVersion
      GLOBAL_ENDPOINT_TOKEN_VERSION_V1TOKEN = globalEndpointTokenVersion._(
          1, _omitEnumNames ? '' : 'GLOBAL_ENDPOINT_TOKEN_VERSION_V1TOKEN');

  static const $core.List<globalEndpointTokenVersion> values =
      <globalEndpointTokenVersion>[
    GLOBAL_ENDPOINT_TOKEN_VERSION_V2TOKEN,
    GLOBAL_ENDPOINT_TOKEN_VERSION_V1TOKEN,
  ];

  static final $core.List<globalEndpointTokenVersion?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static globalEndpointTokenVersion? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const globalEndpointTokenVersion._(super.value, super.name);
}

class jobStatusType extends $pb.ProtobufEnum {
  static const jobStatusType JOB_STATUS_TYPE_IN_PROGRESS =
      jobStatusType._(0, _omitEnumNames ? '' : 'JOB_STATUS_TYPE_IN_PROGRESS');
  static const jobStatusType JOB_STATUS_TYPE_COMPLETED =
      jobStatusType._(1, _omitEnumNames ? '' : 'JOB_STATUS_TYPE_COMPLETED');
  static const jobStatusType JOB_STATUS_TYPE_FAILED =
      jobStatusType._(2, _omitEnumNames ? '' : 'JOB_STATUS_TYPE_FAILED');

  static const $core.List<jobStatusType> values = <jobStatusType>[
    JOB_STATUS_TYPE_IN_PROGRESS,
    JOB_STATUS_TYPE_COMPLETED,
    JOB_STATUS_TYPE_FAILED,
  ];

  static final $core.List<jobStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static jobStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const jobStatusType._(super.value, super.name);
}

class permissionCheckResultType extends $pb.ProtobufEnum {
  static const permissionCheckResultType PERMISSION_CHECK_RESULT_TYPE_UNSURE =
      permissionCheckResultType._(
          0, _omitEnumNames ? '' : 'PERMISSION_CHECK_RESULT_TYPE_UNSURE');
  static const permissionCheckResultType PERMISSION_CHECK_RESULT_TYPE_ALLOWED =
      permissionCheckResultType._(
          1, _omitEnumNames ? '' : 'PERMISSION_CHECK_RESULT_TYPE_ALLOWED');
  static const permissionCheckResultType PERMISSION_CHECK_RESULT_TYPE_DENIED =
      permissionCheckResultType._(
          2, _omitEnumNames ? '' : 'PERMISSION_CHECK_RESULT_TYPE_DENIED');

  static const $core.List<permissionCheckResultType> values =
      <permissionCheckResultType>[
    PERMISSION_CHECK_RESULT_TYPE_UNSURE,
    PERMISSION_CHECK_RESULT_TYPE_ALLOWED,
    PERMISSION_CHECK_RESULT_TYPE_DENIED,
  ];

  static final $core.List<permissionCheckResultType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static permissionCheckResultType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const permissionCheckResultType._(super.value, super.name);
}

class permissionCheckStatusType extends $pb.ProtobufEnum {
  static const permissionCheckStatusType
      PERMISSION_CHECK_STATUS_TYPE_IN_PROGRESS = permissionCheckStatusType._(
          0, _omitEnumNames ? '' : 'PERMISSION_CHECK_STATUS_TYPE_IN_PROGRESS');
  static const permissionCheckStatusType PERMISSION_CHECK_STATUS_TYPE_COMPLETE =
      permissionCheckStatusType._(
          1, _omitEnumNames ? '' : 'PERMISSION_CHECK_STATUS_TYPE_COMPLETE');
  static const permissionCheckStatusType PERMISSION_CHECK_STATUS_TYPE_FAILED =
      permissionCheckStatusType._(
          2, _omitEnumNames ? '' : 'PERMISSION_CHECK_STATUS_TYPE_FAILED');

  static const $core.List<permissionCheckStatusType> values =
      <permissionCheckStatusType>[
    PERMISSION_CHECK_STATUS_TYPE_IN_PROGRESS,
    PERMISSION_CHECK_STATUS_TYPE_COMPLETE,
    PERMISSION_CHECK_STATUS_TYPE_FAILED,
  ];

  static final $core.List<permissionCheckStatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static permissionCheckStatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const permissionCheckStatusType._(super.value, super.name);
}

class policyOwnerEntityType extends $pb.ProtobufEnum {
  static const policyOwnerEntityType POLICY_OWNER_ENTITY_TYPE_GROUP =
      policyOwnerEntityType._(
          0, _omitEnumNames ? '' : 'POLICY_OWNER_ENTITY_TYPE_GROUP');
  static const policyOwnerEntityType POLICY_OWNER_ENTITY_TYPE_ROLE =
      policyOwnerEntityType._(
          1, _omitEnumNames ? '' : 'POLICY_OWNER_ENTITY_TYPE_ROLE');
  static const policyOwnerEntityType POLICY_OWNER_ENTITY_TYPE_USER =
      policyOwnerEntityType._(
          2, _omitEnumNames ? '' : 'POLICY_OWNER_ENTITY_TYPE_USER');

  static const $core.List<policyOwnerEntityType> values =
      <policyOwnerEntityType>[
    POLICY_OWNER_ENTITY_TYPE_GROUP,
    POLICY_OWNER_ENTITY_TYPE_ROLE,
    POLICY_OWNER_ENTITY_TYPE_USER,
  ];

  static final $core.List<policyOwnerEntityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static policyOwnerEntityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const policyOwnerEntityType._(super.value, super.name);
}

class policyScopeType extends $pb.ProtobufEnum {
  static const policyScopeType POLICY_SCOPE_TYPE_LOCAL =
      policyScopeType._(0, _omitEnumNames ? '' : 'POLICY_SCOPE_TYPE_LOCAL');
  static const policyScopeType POLICY_SCOPE_TYPE_ALL =
      policyScopeType._(1, _omitEnumNames ? '' : 'POLICY_SCOPE_TYPE_ALL');
  static const policyScopeType POLICY_SCOPE_TYPE_AWS =
      policyScopeType._(2, _omitEnumNames ? '' : 'POLICY_SCOPE_TYPE_AWS');

  static const $core.List<policyScopeType> values = <policyScopeType>[
    POLICY_SCOPE_TYPE_LOCAL,
    POLICY_SCOPE_TYPE_ALL,
    POLICY_SCOPE_TYPE_AWS,
  ];

  static final $core.List<policyScopeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static policyScopeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const policyScopeType._(super.value, super.name);
}

class policyType extends $pb.ProtobufEnum {
  static const policyType POLICY_TYPE_INLINE =
      policyType._(0, _omitEnumNames ? '' : 'POLICY_TYPE_INLINE');
  static const policyType POLICY_TYPE_MANAGED =
      policyType._(1, _omitEnumNames ? '' : 'POLICY_TYPE_MANAGED');

  static const $core.List<policyType> values = <policyType>[
    POLICY_TYPE_INLINE,
    POLICY_TYPE_MANAGED,
  ];

  static final $core.List<policyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static policyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const policyType._(super.value, super.name);
}

class sortKeyType extends $pb.ProtobufEnum {
  static const sortKeyType SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_DESCENDING =
      sortKeyType._(
          0,
          _omitEnumNames
              ? ''
              : 'SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_DESCENDING');
  static const sortKeyType SORT_KEY_TYPE_SERVICE_NAMESPACE_DESCENDING =
      sortKeyType._(1,
          _omitEnumNames ? '' : 'SORT_KEY_TYPE_SERVICE_NAMESPACE_DESCENDING');
  static const sortKeyType SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_ASCENDING =
      sortKeyType._(
          2,
          _omitEnumNames
              ? ''
              : 'SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_ASCENDING');
  static const sortKeyType SORT_KEY_TYPE_SERVICE_NAMESPACE_ASCENDING =
      sortKeyType._(
          3, _omitEnumNames ? '' : 'SORT_KEY_TYPE_SERVICE_NAMESPACE_ASCENDING');

  static const $core.List<sortKeyType> values = <sortKeyType>[
    SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_DESCENDING,
    SORT_KEY_TYPE_SERVICE_NAMESPACE_DESCENDING,
    SORT_KEY_TYPE_LAST_AUTHENTICATED_TIME_ASCENDING,
    SORT_KEY_TYPE_SERVICE_NAMESPACE_ASCENDING,
  ];

  static final $core.List<sortKeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static sortKeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const sortKeyType._(super.value, super.name);
}

class stateType extends $pb.ProtobufEnum {
  static const stateType STATE_TYPE_ACCEPTED =
      stateType._(0, _omitEnumNames ? '' : 'STATE_TYPE_ACCEPTED');
  static const stateType STATE_TYPE_PENDING_APPROVAL =
      stateType._(1, _omitEnumNames ? '' : 'STATE_TYPE_PENDING_APPROVAL');
  static const stateType STATE_TYPE_FINALIZED =
      stateType._(2, _omitEnumNames ? '' : 'STATE_TYPE_FINALIZED');
  static const stateType STATE_TYPE_UNASSIGNED =
      stateType._(3, _omitEnumNames ? '' : 'STATE_TYPE_UNASSIGNED');
  static const stateType STATE_TYPE_REJECTED =
      stateType._(4, _omitEnumNames ? '' : 'STATE_TYPE_REJECTED');
  static const stateType STATE_TYPE_EXPIRED =
      stateType._(5, _omitEnumNames ? '' : 'STATE_TYPE_EXPIRED');
  static const stateType STATE_TYPE_ASSIGNED =
      stateType._(6, _omitEnumNames ? '' : 'STATE_TYPE_ASSIGNED');

  static const $core.List<stateType> values = <stateType>[
    STATE_TYPE_ACCEPTED,
    STATE_TYPE_PENDING_APPROVAL,
    STATE_TYPE_FINALIZED,
    STATE_TYPE_UNASSIGNED,
    STATE_TYPE_REJECTED,
    STATE_TYPE_EXPIRED,
    STATE_TYPE_ASSIGNED,
  ];

  static final $core.List<stateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static stateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const stateType._(super.value, super.name);
}

class statusType extends $pb.ProtobufEnum {
  static const statusType STATUS_TYPE_ACTIVE =
      statusType._(0, _omitEnumNames ? '' : 'STATUS_TYPE_ACTIVE');
  static const statusType STATUS_TYPE_EXPIRED =
      statusType._(1, _omitEnumNames ? '' : 'STATUS_TYPE_EXPIRED');
  static const statusType STATUS_TYPE_INACTIVE =
      statusType._(2, _omitEnumNames ? '' : 'STATUS_TYPE_INACTIVE');

  static const $core.List<statusType> values = <statusType>[
    STATUS_TYPE_ACTIVE,
    STATUS_TYPE_EXPIRED,
    STATUS_TYPE_INACTIVE,
  ];

  static final $core.List<statusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static statusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const statusType._(super.value, super.name);
}

class summaryKeyType extends $pb.ProtobufEnum {
  static const summaryKeyType SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERGROUPQUOTA =
      summaryKeyType._(
          0,
          _omitEnumNames
              ? ''
              : 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERGROUPQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_ROLES =
      summaryKeyType._(1, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ROLES');
  static const summaryKeyType SUMMARY_KEY_TYPE_VERSIONSPERPOLICYQUOTA =
      summaryKeyType._(
          2, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_VERSIONSPERPOLICYQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_MFADEVICES =
      summaryKeyType._(3, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_MFADEVICES');
  static const summaryKeyType SUMMARY_KEY_TYPE_ACCESSKEYSPERUSERQUOTA =
      summaryKeyType._(
          4, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ACCESSKEYSPERUSERQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_POLICIES =
      summaryKeyType._(5, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_POLICIES');
  static const summaryKeyType SUMMARY_KEY_TYPE_INSTANCEPROFILES = summaryKeyType
      ._(6, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_INSTANCEPROFILES');
  static const summaryKeyType SUMMARY_KEY_TYPE_ROLEPOLICYSIZEQUOTA =
      summaryKeyType._(
          7, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ROLEPOLICYSIZEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERROLEQUOTA =
      summaryKeyType._(
          8,
          _omitEnumNames
              ? ''
              : 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERROLEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_ASSUMEROLEPOLICYSIZEQUOTA =
      summaryKeyType._(9,
          _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ASSUMEROLEPOLICYSIZEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_SIGNINGCERTIFICATESPERUSERQUOTA =
      summaryKeyType._(
          10,
          _omitEnumNames
              ? ''
              : 'SUMMARY_KEY_TYPE_SIGNINGCERTIFICATESPERUSERQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_SERVERCERTIFICATESQUOTA =
      summaryKeyType._(
          11, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_SERVERCERTIFICATESQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_INSTANCEPROFILESQUOTA =
      summaryKeyType._(
          12, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_INSTANCEPROFILESQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_ROLESQUOTA =
      summaryKeyType._(13, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ROLESQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_GROUPPOLICYSIZEQUOTA =
      summaryKeyType._(
          14, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_GROUPPOLICYSIZEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_POLICYVERSIONSINUSE =
      summaryKeyType._(
          15, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_POLICYVERSIONSINUSE');
  static const summaryKeyType SUMMARY_KEY_TYPE_POLICIESQUOTA = summaryKeyType._(
      16, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_POLICIESQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_USERSQUOTA =
      summaryKeyType._(17, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_USERSQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_USERPOLICYSIZEQUOTA =
      summaryKeyType._(
          18, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_USERPOLICYSIZEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_GLOBALENDPOINTTOKENVERSION =
      summaryKeyType._(19,
          _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_GLOBALENDPOINTTOKENVERSION');
  static const summaryKeyType SUMMARY_KEY_TYPE_POLICYSIZEQUOTA = summaryKeyType
      ._(20, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_POLICYSIZEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_SERVERCERTIFICATES =
      summaryKeyType._(
          21, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_SERVERCERTIFICATES');
  static const summaryKeyType SUMMARY_KEY_TYPE_ACCOUNTPASSWORDPRESENT =
      summaryKeyType._(
          22, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ACCOUNTPASSWORDPRESENT');
  static const summaryKeyType
      SUMMARY_KEY_TYPE_ACCOUNTSIGNINGCERTIFICATESPRESENT = summaryKeyType._(
          23,
          _omitEnumNames
              ? ''
              : 'SUMMARY_KEY_TYPE_ACCOUNTSIGNINGCERTIFICATESPRESENT');
  static const summaryKeyType SUMMARY_KEY_TYPE_USERS =
      summaryKeyType._(24, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_USERS');
  static const summaryKeyType SUMMARY_KEY_TYPE_POLICYVERSIONSINUSEQUOTA =
      summaryKeyType._(25,
          _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_POLICYVERSIONSINUSEQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_MFADEVICESINUSE = summaryKeyType
      ._(26, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_MFADEVICESINUSE');
  static const summaryKeyType SUMMARY_KEY_TYPE_GROUPS =
      summaryKeyType._(27, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_GROUPS');
  static const summaryKeyType SUMMARY_KEY_TYPE_ACCOUNTACCESSKEYSPRESENT =
      summaryKeyType._(28,
          _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ACCOUNTACCESSKEYSPRESENT');
  static const summaryKeyType SUMMARY_KEY_TYPE_GROUPSQUOTA = summaryKeyType._(
      29, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_GROUPSQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_PROVIDERS =
      summaryKeyType._(30, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_PROVIDERS');
  static const summaryKeyType SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERUSERQUOTA =
      summaryKeyType._(
          31,
          _omitEnumNames
              ? ''
              : 'SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERUSERQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_GROUPSPERUSERQUOTA =
      summaryKeyType._(
          32, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_GROUPSPERUSERQUOTA');
  static const summaryKeyType SUMMARY_KEY_TYPE_ACCOUNTMFAENABLED =
      summaryKeyType._(
          33, _omitEnumNames ? '' : 'SUMMARY_KEY_TYPE_ACCOUNTMFAENABLED');

  static const $core.List<summaryKeyType> values = <summaryKeyType>[
    SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERGROUPQUOTA,
    SUMMARY_KEY_TYPE_ROLES,
    SUMMARY_KEY_TYPE_VERSIONSPERPOLICYQUOTA,
    SUMMARY_KEY_TYPE_MFADEVICES,
    SUMMARY_KEY_TYPE_ACCESSKEYSPERUSERQUOTA,
    SUMMARY_KEY_TYPE_POLICIES,
    SUMMARY_KEY_TYPE_INSTANCEPROFILES,
    SUMMARY_KEY_TYPE_ROLEPOLICYSIZEQUOTA,
    SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERROLEQUOTA,
    SUMMARY_KEY_TYPE_ASSUMEROLEPOLICYSIZEQUOTA,
    SUMMARY_KEY_TYPE_SIGNINGCERTIFICATESPERUSERQUOTA,
    SUMMARY_KEY_TYPE_SERVERCERTIFICATESQUOTA,
    SUMMARY_KEY_TYPE_INSTANCEPROFILESQUOTA,
    SUMMARY_KEY_TYPE_ROLESQUOTA,
    SUMMARY_KEY_TYPE_GROUPPOLICYSIZEQUOTA,
    SUMMARY_KEY_TYPE_POLICYVERSIONSINUSE,
    SUMMARY_KEY_TYPE_POLICIESQUOTA,
    SUMMARY_KEY_TYPE_USERSQUOTA,
    SUMMARY_KEY_TYPE_USERPOLICYSIZEQUOTA,
    SUMMARY_KEY_TYPE_GLOBALENDPOINTTOKENVERSION,
    SUMMARY_KEY_TYPE_POLICYSIZEQUOTA,
    SUMMARY_KEY_TYPE_SERVERCERTIFICATES,
    SUMMARY_KEY_TYPE_ACCOUNTPASSWORDPRESENT,
    SUMMARY_KEY_TYPE_ACCOUNTSIGNINGCERTIFICATESPRESENT,
    SUMMARY_KEY_TYPE_USERS,
    SUMMARY_KEY_TYPE_POLICYVERSIONSINUSEQUOTA,
    SUMMARY_KEY_TYPE_MFADEVICESINUSE,
    SUMMARY_KEY_TYPE_GROUPS,
    SUMMARY_KEY_TYPE_ACCOUNTACCESSKEYSPRESENT,
    SUMMARY_KEY_TYPE_GROUPSQUOTA,
    SUMMARY_KEY_TYPE_PROVIDERS,
    SUMMARY_KEY_TYPE_ATTACHEDPOLICIESPERUSERQUOTA,
    SUMMARY_KEY_TYPE_GROUPSPERUSERQUOTA,
    SUMMARY_KEY_TYPE_ACCOUNTMFAENABLED,
  ];

  static final $core.List<summaryKeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 33);
  static summaryKeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const summaryKeyType._(super.value, super.name);
}

class summaryStateType extends $pb.ProtobufEnum {
  static const summaryStateType SUMMARY_STATE_TYPE_AVAILABLE = summaryStateType
      ._(0, _omitEnumNames ? '' : 'SUMMARY_STATE_TYPE_AVAILABLE');
  static const summaryStateType SUMMARY_STATE_TYPE_NOT_AVAILABLE =
      summaryStateType._(
          1, _omitEnumNames ? '' : 'SUMMARY_STATE_TYPE_NOT_AVAILABLE');
  static const summaryStateType SUMMARY_STATE_TYPE_NOT_SUPPORTED =
      summaryStateType._(
          2, _omitEnumNames ? '' : 'SUMMARY_STATE_TYPE_NOT_SUPPORTED');
  static const summaryStateType SUMMARY_STATE_TYPE_FAILED =
      summaryStateType._(3, _omitEnumNames ? '' : 'SUMMARY_STATE_TYPE_FAILED');

  static const $core.List<summaryStateType> values = <summaryStateType>[
    SUMMARY_STATE_TYPE_AVAILABLE,
    SUMMARY_STATE_TYPE_NOT_AVAILABLE,
    SUMMARY_STATE_TYPE_NOT_SUPPORTED,
    SUMMARY_STATE_TYPE_FAILED,
  ];

  static final $core.List<summaryStateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static summaryStateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const summaryStateType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
