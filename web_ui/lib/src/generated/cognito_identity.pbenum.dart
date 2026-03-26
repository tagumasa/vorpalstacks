// This is a generated file - do not edit.
//
// Generated from cognito_identity.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AmbiguousRoleResolutionType extends $pb.ProtobufEnum {
  static const AmbiguousRoleResolutionType
      AMBIGUOUS_ROLE_RESOLUTION_TYPE_AUTHENTICATED_ROLE =
      AmbiguousRoleResolutionType._(
          0,
          _omitEnumNames
              ? ''
              : 'AMBIGUOUS_ROLE_RESOLUTION_TYPE_AUTHENTICATED_ROLE');
  static const AmbiguousRoleResolutionType AMBIGUOUS_ROLE_RESOLUTION_TYPE_DENY =
      AmbiguousRoleResolutionType._(
          1, _omitEnumNames ? '' : 'AMBIGUOUS_ROLE_RESOLUTION_TYPE_DENY');

  static const $core.List<AmbiguousRoleResolutionType> values =
      <AmbiguousRoleResolutionType>[
    AMBIGUOUS_ROLE_RESOLUTION_TYPE_AUTHENTICATED_ROLE,
    AMBIGUOUS_ROLE_RESOLUTION_TYPE_DENY,
  ];

  static final $core.List<AmbiguousRoleResolutionType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AmbiguousRoleResolutionType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AmbiguousRoleResolutionType._(super.value, super.name);
}

class ErrorCode extends $pb.ProtobufEnum {
  static const ErrorCode ERROR_CODE_ACCESS_DENIED =
      ErrorCode._(0, _omitEnumNames ? '' : 'ERROR_CODE_ACCESS_DENIED');
  static const ErrorCode ERROR_CODE_INTERNAL_SERVER_ERROR =
      ErrorCode._(1, _omitEnumNames ? '' : 'ERROR_CODE_INTERNAL_SERVER_ERROR');

  static const $core.List<ErrorCode> values = <ErrorCode>[
    ERROR_CODE_ACCESS_DENIED,
    ERROR_CODE_INTERNAL_SERVER_ERROR,
  ];

  static final $core.List<ErrorCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ErrorCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ErrorCode._(super.value, super.name);
}

class MappingRuleMatchType extends $pb.ProtobufEnum {
  static const MappingRuleMatchType MAPPING_RULE_MATCH_TYPE_CONTAINS =
      MappingRuleMatchType._(
          0, _omitEnumNames ? '' : 'MAPPING_RULE_MATCH_TYPE_CONTAINS');
  static const MappingRuleMatchType MAPPING_RULE_MATCH_TYPE_EQUALS =
      MappingRuleMatchType._(
          1, _omitEnumNames ? '' : 'MAPPING_RULE_MATCH_TYPE_EQUALS');
  static const MappingRuleMatchType MAPPING_RULE_MATCH_TYPE_STARTS_WITH =
      MappingRuleMatchType._(
          2, _omitEnumNames ? '' : 'MAPPING_RULE_MATCH_TYPE_STARTS_WITH');
  static const MappingRuleMatchType MAPPING_RULE_MATCH_TYPE_NOT_EQUAL =
      MappingRuleMatchType._(
          3, _omitEnumNames ? '' : 'MAPPING_RULE_MATCH_TYPE_NOT_EQUAL');

  static const $core.List<MappingRuleMatchType> values = <MappingRuleMatchType>[
    MAPPING_RULE_MATCH_TYPE_CONTAINS,
    MAPPING_RULE_MATCH_TYPE_EQUALS,
    MAPPING_RULE_MATCH_TYPE_STARTS_WITH,
    MAPPING_RULE_MATCH_TYPE_NOT_EQUAL,
  ];

  static final $core.List<MappingRuleMatchType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MappingRuleMatchType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MappingRuleMatchType._(super.value, super.name);
}

class RoleMappingType extends $pb.ProtobufEnum {
  static const RoleMappingType ROLE_MAPPING_TYPE_RULES =
      RoleMappingType._(0, _omitEnumNames ? '' : 'ROLE_MAPPING_TYPE_RULES');
  static const RoleMappingType ROLE_MAPPING_TYPE_TOKEN =
      RoleMappingType._(1, _omitEnumNames ? '' : 'ROLE_MAPPING_TYPE_TOKEN');

  static const $core.List<RoleMappingType> values = <RoleMappingType>[
    ROLE_MAPPING_TYPE_RULES,
    ROLE_MAPPING_TYPE_TOKEN,
  ];

  static final $core.List<RoleMappingType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RoleMappingType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RoleMappingType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
