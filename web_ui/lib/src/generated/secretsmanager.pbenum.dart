// This is a generated file - do not edit.
//
// Generated from secretsmanager.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class FilterNameStringType extends $pb.ProtobufEnum {
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_TAG_KEY =
      FilterNameStringType._(
          0, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_TAG_KEY');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_ALL =
      FilterNameStringType._(
          1, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_ALL');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_PRIMARY_REGION =
      FilterNameStringType._(
          2, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_PRIMARY_REGION');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_OWNING_SERVICE =
      FilterNameStringType._(
          3, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_OWNING_SERVICE');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_NAME =
      FilterNameStringType._(
          4, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_NAME');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_TAG_VALUE =
      FilterNameStringType._(
          5, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_TAG_VALUE');
  static const FilterNameStringType FILTER_NAME_STRING_TYPE_DESCRIPTION =
      FilterNameStringType._(
          6, _omitEnumNames ? '' : 'FILTER_NAME_STRING_TYPE_DESCRIPTION');

  static const $core.List<FilterNameStringType> values = <FilterNameStringType>[
    FILTER_NAME_STRING_TYPE_TAG_KEY,
    FILTER_NAME_STRING_TYPE_ALL,
    FILTER_NAME_STRING_TYPE_PRIMARY_REGION,
    FILTER_NAME_STRING_TYPE_OWNING_SERVICE,
    FILTER_NAME_STRING_TYPE_NAME,
    FILTER_NAME_STRING_TYPE_TAG_VALUE,
    FILTER_NAME_STRING_TYPE_DESCRIPTION,
  ];

  static final $core.List<FilterNameStringType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static FilterNameStringType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FilterNameStringType._(super.value, super.name);
}

class SortByType extends $pb.ProtobufEnum {
  static const SortByType SORT_BY_TYPE_LAST_ACCESSED_DATE =
      SortByType._(0, _omitEnumNames ? '' : 'SORT_BY_TYPE_LAST_ACCESSED_DATE');
  static const SortByType SORT_BY_TYPE_CREATED_DATE =
      SortByType._(1, _omitEnumNames ? '' : 'SORT_BY_TYPE_CREATED_DATE');
  static const SortByType SORT_BY_TYPE_LAST_CHANGED_DATE =
      SortByType._(2, _omitEnumNames ? '' : 'SORT_BY_TYPE_LAST_CHANGED_DATE');
  static const SortByType SORT_BY_TYPE_NAME =
      SortByType._(3, _omitEnumNames ? '' : 'SORT_BY_TYPE_NAME');

  static const $core.List<SortByType> values = <SortByType>[
    SORT_BY_TYPE_LAST_ACCESSED_DATE,
    SORT_BY_TYPE_CREATED_DATE,
    SORT_BY_TYPE_LAST_CHANGED_DATE,
    SORT_BY_TYPE_NAME,
  ];

  static final $core.List<SortByType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static SortByType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SortByType._(super.value, super.name);
}

class SortOrderType extends $pb.ProtobufEnum {
  static const SortOrderType SORT_ORDER_TYPE_ASC =
      SortOrderType._(0, _omitEnumNames ? '' : 'SORT_ORDER_TYPE_ASC');
  static const SortOrderType SORT_ORDER_TYPE_DESC =
      SortOrderType._(1, _omitEnumNames ? '' : 'SORT_ORDER_TYPE_DESC');

  static const $core.List<SortOrderType> values = <SortOrderType>[
    SORT_ORDER_TYPE_ASC,
    SORT_ORDER_TYPE_DESC,
  ];

  static final $core.List<SortOrderType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SortOrderType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SortOrderType._(super.value, super.name);
}

class StatusType extends $pb.ProtobufEnum {
  static const StatusType STATUS_TYPE_INSYNC =
      StatusType._(0, _omitEnumNames ? '' : 'STATUS_TYPE_INSYNC');
  static const StatusType STATUS_TYPE_FAILED =
      StatusType._(1, _omitEnumNames ? '' : 'STATUS_TYPE_FAILED');
  static const StatusType STATUS_TYPE_INPROGRESS =
      StatusType._(2, _omitEnumNames ? '' : 'STATUS_TYPE_INPROGRESS');

  static const $core.List<StatusType> values = <StatusType>[
    STATUS_TYPE_INSYNC,
    STATUS_TYPE_FAILED,
    STATUS_TYPE_INPROGRESS,
  ];

  static final $core.List<StatusType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static StatusType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const StatusType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
