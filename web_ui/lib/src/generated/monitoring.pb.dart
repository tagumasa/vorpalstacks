// This is a generated file - do not edit.
//
// Generated from monitoring.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'monitoring.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'monitoring.pbenum.dart';

class AlarmContributor extends $pb.GeneratedMessage {
  factory AlarmContributor({
    $core.String? contributorid,
    $core.String? statetransitionedtimestamp,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        contributorattributes,
    $core.String? statereason,
  }) {
    final result = create();
    if (contributorid != null) result.contributorid = contributorid;
    if (statetransitionedtimestamp != null)
      result.statetransitionedtimestamp = statetransitionedtimestamp;
    if (contributorattributes != null)
      result.contributorattributes.addEntries(contributorattributes);
    if (statereason != null) result.statereason = statereason;
    return result;
  }

  AlarmContributor._();

  factory AlarmContributor.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AlarmContributor.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AlarmContributor',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(33454376, _omitFieldNames ? '' : 'contributorid')
    ..aOS(46811093, _omitFieldNames ? '' : 'statetransitionedtimestamp')
    ..m<$core.String, $core.String>(
        82097856, _omitFieldNames ? '' : 'contributorattributes',
        entryClassName: 'AlarmContributor.ContributorattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('monitoring'))
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmContributor clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmContributor copyWith(void Function(AlarmContributor) updates) =>
      super.copyWith((message) => updates(message as AlarmContributor))
          as AlarmContributor;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AlarmContributor create() => AlarmContributor._();
  @$core.override
  AlarmContributor createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AlarmContributor getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AlarmContributor>(create);
  static AlarmContributor? _defaultInstance;

  @$pb.TagNumber(33454376)
  $core.String get contributorid => $_getSZ(0);
  @$pb.TagNumber(33454376)
  set contributorid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(33454376)
  $core.bool hasContributorid() => $_has(0);
  @$pb.TagNumber(33454376)
  void clearContributorid() => $_clearField(33454376);

  @$pb.TagNumber(46811093)
  $core.String get statetransitionedtimestamp => $_getSZ(1);
  @$pb.TagNumber(46811093)
  set statetransitionedtimestamp($core.String value) => $_setString(1, value);
  @$pb.TagNumber(46811093)
  $core.bool hasStatetransitionedtimestamp() => $_has(1);
  @$pb.TagNumber(46811093)
  void clearStatetransitionedtimestamp() => $_clearField(46811093);

  @$pb.TagNumber(82097856)
  $pb.PbMap<$core.String, $core.String> get contributorattributes =>
      $_getMap(2);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(3);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(3, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(3);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);
}

class AlarmHistoryItem extends $pb.GeneratedMessage {
  factory AlarmHistoryItem({
    $core.String? historydata,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>?
        alarmcontributorattributes,
    $core.String? historysummary,
    $core.String? alarmname,
    $core.String? timestamp,
    AlarmType? alarmtype,
    $core.String? alarmcontributorid,
    HistoryItemType? historyitemtype,
  }) {
    final result = create();
    if (historydata != null) result.historydata = historydata;
    if (alarmcontributorattributes != null)
      result.alarmcontributorattributes.addEntries(alarmcontributorattributes);
    if (historysummary != null) result.historysummary = historysummary;
    if (alarmname != null) result.alarmname = alarmname;
    if (timestamp != null) result.timestamp = timestamp;
    if (alarmtype != null) result.alarmtype = alarmtype;
    if (alarmcontributorid != null)
      result.alarmcontributorid = alarmcontributorid;
    if (historyitemtype != null) result.historyitemtype = historyitemtype;
    return result;
  }

  AlarmHistoryItem._();

  factory AlarmHistoryItem.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AlarmHistoryItem.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AlarmHistoryItem',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(30118600, _omitFieldNames ? '' : 'historydata')
    ..m<$core.String, $core.String>(
        34317575, _omitFieldNames ? '' : 'alarmcontributorattributes',
        entryClassName: 'AlarmHistoryItem.AlarmcontributorattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('monitoring'))
    ..aOS(53483226, _omitFieldNames ? '' : 'historysummary')
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aE<AlarmType>(301805339, _omitFieldNames ? '' : 'alarmtype',
        enumValues: AlarmType.values)
    ..aOS(305558919, _omitFieldNames ? '' : 'alarmcontributorid')
    ..aE<HistoryItemType>(442439537, _omitFieldNames ? '' : 'historyitemtype',
        enumValues: HistoryItemType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmHistoryItem clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmHistoryItem copyWith(void Function(AlarmHistoryItem) updates) =>
      super.copyWith((message) => updates(message as AlarmHistoryItem))
          as AlarmHistoryItem;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AlarmHistoryItem create() => AlarmHistoryItem._();
  @$core.override
  AlarmHistoryItem createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AlarmHistoryItem getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AlarmHistoryItem>(create);
  static AlarmHistoryItem? _defaultInstance;

  @$pb.TagNumber(30118600)
  $core.String get historydata => $_getSZ(0);
  @$pb.TagNumber(30118600)
  set historydata($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30118600)
  $core.bool hasHistorydata() => $_has(0);
  @$pb.TagNumber(30118600)
  void clearHistorydata() => $_clearField(30118600);

  @$pb.TagNumber(34317575)
  $pb.PbMap<$core.String, $core.String> get alarmcontributorattributes =>
      $_getMap(1);

  @$pb.TagNumber(53483226)
  $core.String get historysummary => $_getSZ(2);
  @$pb.TagNumber(53483226)
  set historysummary($core.String value) => $_setString(2, value);
  @$pb.TagNumber(53483226)
  $core.bool hasHistorysummary() => $_has(2);
  @$pb.TagNumber(53483226)
  void clearHistorysummary() => $_clearField(53483226);

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(3);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(3);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(4);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(4, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(301805339)
  AlarmType get alarmtype => $_getN(5);
  @$pb.TagNumber(301805339)
  set alarmtype(AlarmType value) => $_setField(301805339, value);
  @$pb.TagNumber(301805339)
  $core.bool hasAlarmtype() => $_has(5);
  @$pb.TagNumber(301805339)
  void clearAlarmtype() => $_clearField(301805339);

  @$pb.TagNumber(305558919)
  $core.String get alarmcontributorid => $_getSZ(6);
  @$pb.TagNumber(305558919)
  set alarmcontributorid($core.String value) => $_setString(6, value);
  @$pb.TagNumber(305558919)
  $core.bool hasAlarmcontributorid() => $_has(6);
  @$pb.TagNumber(305558919)
  void clearAlarmcontributorid() => $_clearField(305558919);

  @$pb.TagNumber(442439537)
  HistoryItemType get historyitemtype => $_getN(7);
  @$pb.TagNumber(442439537)
  set historyitemtype(HistoryItemType value) => $_setField(442439537, value);
  @$pb.TagNumber(442439537)
  $core.bool hasHistoryitemtype() => $_has(7);
  @$pb.TagNumber(442439537)
  void clearHistoryitemtype() => $_clearField(442439537);
}

class AlarmMuteRuleSummary extends $pb.GeneratedMessage {
  factory AlarmMuteRuleSummary({
    AlarmMuteRuleStatus? status,
    $core.String? mutetype,
    $core.String? lastupdatedtimestamp,
    $core.String? expiredate,
    $core.String? alarmmuterulearn,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (mutetype != null) result.mutetype = mutetype;
    if (lastupdatedtimestamp != null)
      result.lastupdatedtimestamp = lastupdatedtimestamp;
    if (expiredate != null) result.expiredate = expiredate;
    if (alarmmuterulearn != null) result.alarmmuterulearn = alarmmuterulearn;
    return result;
  }

  AlarmMuteRuleSummary._();

  factory AlarmMuteRuleSummary.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AlarmMuteRuleSummary.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AlarmMuteRuleSummary',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aE<AlarmMuteRuleStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: AlarmMuteRuleStatus.values)
    ..aOS(111213245, _omitFieldNames ? '' : 'mutetype')
    ..aOS(133309845, _omitFieldNames ? '' : 'lastupdatedtimestamp')
    ..aOS(454143207, _omitFieldNames ? '' : 'expiredate')
    ..aOS(493590237, _omitFieldNames ? '' : 'alarmmuterulearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmMuteRuleSummary clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AlarmMuteRuleSummary copyWith(void Function(AlarmMuteRuleSummary) updates) =>
      super.copyWith((message) => updates(message as AlarmMuteRuleSummary))
          as AlarmMuteRuleSummary;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AlarmMuteRuleSummary create() => AlarmMuteRuleSummary._();
  @$core.override
  AlarmMuteRuleSummary createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AlarmMuteRuleSummary getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AlarmMuteRuleSummary>(create);
  static AlarmMuteRuleSummary? _defaultInstance;

  @$pb.TagNumber(6222352)
  AlarmMuteRuleStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(AlarmMuteRuleStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(111213245)
  $core.String get mutetype => $_getSZ(1);
  @$pb.TagNumber(111213245)
  set mutetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(111213245)
  $core.bool hasMutetype() => $_has(1);
  @$pb.TagNumber(111213245)
  void clearMutetype() => $_clearField(111213245);

  @$pb.TagNumber(133309845)
  $core.String get lastupdatedtimestamp => $_getSZ(2);
  @$pb.TagNumber(133309845)
  set lastupdatedtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(133309845)
  $core.bool hasLastupdatedtimestamp() => $_has(2);
  @$pb.TagNumber(133309845)
  void clearLastupdatedtimestamp() => $_clearField(133309845);

  @$pb.TagNumber(454143207)
  $core.String get expiredate => $_getSZ(3);
  @$pb.TagNumber(454143207)
  set expiredate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(454143207)
  $core.bool hasExpiredate() => $_has(3);
  @$pb.TagNumber(454143207)
  void clearExpiredate() => $_clearField(454143207);

  @$pb.TagNumber(493590237)
  $core.String get alarmmuterulearn => $_getSZ(4);
  @$pb.TagNumber(493590237)
  set alarmmuterulearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(493590237)
  $core.bool hasAlarmmuterulearn() => $_has(4);
  @$pb.TagNumber(493590237)
  void clearAlarmmuterulearn() => $_clearField(493590237);
}

class AnomalyDetector extends $pb.GeneratedMessage {
  factory AnomalyDetector({
    $core.String? metricname,
    MetricCharacteristics? metriccharacteristics,
    $core.String? stat,
    AnomalyDetectorStateValue? statevalue,
    $core.String? namespace,
    SingleMetricAnomalyDetector? singlemetricanomalydetector,
    MetricMathAnomalyDetector? metricmathanomalydetector,
    AnomalyDetectorConfiguration? configuration,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (metriccharacteristics != null)
      result.metriccharacteristics = metriccharacteristics;
    if (stat != null) result.stat = stat;
    if (statevalue != null) result.statevalue = statevalue;
    if (namespace != null) result.namespace = namespace;
    if (singlemetricanomalydetector != null)
      result.singlemetricanomalydetector = singlemetricanomalydetector;
    if (metricmathanomalydetector != null)
      result.metricmathanomalydetector = metricmathanomalydetector;
    if (configuration != null) result.configuration = configuration;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  AnomalyDetector._();

  factory AnomalyDetector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AnomalyDetector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AnomalyDetector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOM<MetricCharacteristics>(
        129404842, _omitFieldNames ? '' : 'metriccharacteristics',
        subBuilder: MetricCharacteristics.create)
    ..aOS(325549752, _omitFieldNames ? '' : 'stat')
    ..aE<AnomalyDetectorStateValue>(
        334526008, _omitFieldNames ? '' : 'statevalue',
        enumValues: AnomalyDetectorStateValue.values)
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..aOM<SingleMetricAnomalyDetector>(
        418531933, _omitFieldNames ? '' : 'singlemetricanomalydetector',
        subBuilder: SingleMetricAnomalyDetector.create)
    ..aOM<MetricMathAnomalyDetector>(
        439381475, _omitFieldNames ? '' : 'metricmathanomalydetector',
        subBuilder: MetricMathAnomalyDetector.create)
    ..aOM<AnomalyDetectorConfiguration>(
        442426458, _omitFieldNames ? '' : 'configuration',
        subBuilder: AnomalyDetectorConfiguration.create)
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AnomalyDetector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AnomalyDetector copyWith(void Function(AnomalyDetector) updates) =>
      super.copyWith((message) => updates(message as AnomalyDetector))
          as AnomalyDetector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AnomalyDetector create() => AnomalyDetector._();
  @$core.override
  AnomalyDetector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AnomalyDetector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AnomalyDetector>(create);
  static AnomalyDetector? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(129404842)
  MetricCharacteristics get metriccharacteristics => $_getN(1);
  @$pb.TagNumber(129404842)
  set metriccharacteristics(MetricCharacteristics value) =>
      $_setField(129404842, value);
  @$pb.TagNumber(129404842)
  $core.bool hasMetriccharacteristics() => $_has(1);
  @$pb.TagNumber(129404842)
  void clearMetriccharacteristics() => $_clearField(129404842);
  @$pb.TagNumber(129404842)
  MetricCharacteristics ensureMetriccharacteristics() => $_ensure(1);

  @$pb.TagNumber(325549752)
  $core.String get stat => $_getSZ(2);
  @$pb.TagNumber(325549752)
  set stat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(325549752)
  $core.bool hasStat() => $_has(2);
  @$pb.TagNumber(325549752)
  void clearStat() => $_clearField(325549752);

  @$pb.TagNumber(334526008)
  AnomalyDetectorStateValue get statevalue => $_getN(3);
  @$pb.TagNumber(334526008)
  set statevalue(AnomalyDetectorStateValue value) =>
      $_setField(334526008, value);
  @$pb.TagNumber(334526008)
  $core.bool hasStatevalue() => $_has(3);
  @$pb.TagNumber(334526008)
  void clearStatevalue() => $_clearField(334526008);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(4);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(4, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(4);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector get singlemetricanomalydetector => $_getN(5);
  @$pb.TagNumber(418531933)
  set singlemetricanomalydetector(SingleMetricAnomalyDetector value) =>
      $_setField(418531933, value);
  @$pb.TagNumber(418531933)
  $core.bool hasSinglemetricanomalydetector() => $_has(5);
  @$pb.TagNumber(418531933)
  void clearSinglemetricanomalydetector() => $_clearField(418531933);
  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector ensureSinglemetricanomalydetector() =>
      $_ensure(5);

  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector get metricmathanomalydetector => $_getN(6);
  @$pb.TagNumber(439381475)
  set metricmathanomalydetector(MetricMathAnomalyDetector value) =>
      $_setField(439381475, value);
  @$pb.TagNumber(439381475)
  $core.bool hasMetricmathanomalydetector() => $_has(6);
  @$pb.TagNumber(439381475)
  void clearMetricmathanomalydetector() => $_clearField(439381475);
  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector ensureMetricmathanomalydetector() => $_ensure(6);

  @$pb.TagNumber(442426458)
  AnomalyDetectorConfiguration get configuration => $_getN(7);
  @$pb.TagNumber(442426458)
  set configuration(AnomalyDetectorConfiguration value) =>
      $_setField(442426458, value);
  @$pb.TagNumber(442426458)
  $core.bool hasConfiguration() => $_has(7);
  @$pb.TagNumber(442426458)
  void clearConfiguration() => $_clearField(442426458);
  @$pb.TagNumber(442426458)
  AnomalyDetectorConfiguration ensureConfiguration() => $_ensure(7);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(8);
}

class AnomalyDetectorConfiguration extends $pb.GeneratedMessage {
  factory AnomalyDetectorConfiguration({
    $core.Iterable<Range>? excludedtimeranges,
    $core.String? metrictimezone,
  }) {
    final result = create();
    if (excludedtimeranges != null)
      result.excludedtimeranges.addAll(excludedtimeranges);
    if (metrictimezone != null) result.metrictimezone = metrictimezone;
    return result;
  }

  AnomalyDetectorConfiguration._();

  factory AnomalyDetectorConfiguration.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AnomalyDetectorConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AnomalyDetectorConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<Range>(31303297, _omitFieldNames ? '' : 'excludedtimeranges',
        subBuilder: Range.create)
    ..aOS(528921359, _omitFieldNames ? '' : 'metrictimezone')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AnomalyDetectorConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AnomalyDetectorConfiguration copyWith(
          void Function(AnomalyDetectorConfiguration) updates) =>
      super.copyWith(
              (message) => updates(message as AnomalyDetectorConfiguration))
          as AnomalyDetectorConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AnomalyDetectorConfiguration create() =>
      AnomalyDetectorConfiguration._();
  @$core.override
  AnomalyDetectorConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AnomalyDetectorConfiguration getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AnomalyDetectorConfiguration>(create);
  static AnomalyDetectorConfiguration? _defaultInstance;

  @$pb.TagNumber(31303297)
  $pb.PbList<Range> get excludedtimeranges => $_getList(0);

  @$pb.TagNumber(528921359)
  $core.String get metrictimezone => $_getSZ(1);
  @$pb.TagNumber(528921359)
  set metrictimezone($core.String value) => $_setString(1, value);
  @$pb.TagNumber(528921359)
  $core.bool hasMetrictimezone() => $_has(1);
  @$pb.TagNumber(528921359)
  void clearMetrictimezone() => $_clearField(528921359);
}

class CompositeAlarm extends $pb.GeneratedMessage {
  factory CompositeAlarm({
    $core.String? actionssuppressedreason,
    $core.String? alarmdescription,
    $core.String? statetransitionedtimestamp,
    $core.String? alarmconfigurationupdatedtimestamp,
    $core.String? alarmname,
    $core.int? actionssuppressorwaitperiod,
    $core.bool? actionsenabled,
    $core.int? actionssuppressorextensionperiod,
    $core.String? statereasondata,
    $core.String? alarmarn,
    StateValue? statevalue,
    $core.Iterable<$core.String>? alarmactions,
    $core.String? stateupdatedtimestamp,
    $core.String? statereason,
    $core.Iterable<$core.String>? okactions,
    ActionsSuppressedBy? actionssuppressedby,
    $core.String? actionssuppressor,
    $core.Iterable<$core.String>? insufficientdataactions,
    $core.String? alarmrule,
  }) {
    final result = create();
    if (actionssuppressedreason != null)
      result.actionssuppressedreason = actionssuppressedreason;
    if (alarmdescription != null) result.alarmdescription = alarmdescription;
    if (statetransitionedtimestamp != null)
      result.statetransitionedtimestamp = statetransitionedtimestamp;
    if (alarmconfigurationupdatedtimestamp != null)
      result.alarmconfigurationupdatedtimestamp =
          alarmconfigurationupdatedtimestamp;
    if (alarmname != null) result.alarmname = alarmname;
    if (actionssuppressorwaitperiod != null)
      result.actionssuppressorwaitperiod = actionssuppressorwaitperiod;
    if (actionsenabled != null) result.actionsenabled = actionsenabled;
    if (actionssuppressorextensionperiod != null)
      result.actionssuppressorextensionperiod =
          actionssuppressorextensionperiod;
    if (statereasondata != null) result.statereasondata = statereasondata;
    if (alarmarn != null) result.alarmarn = alarmarn;
    if (statevalue != null) result.statevalue = statevalue;
    if (alarmactions != null) result.alarmactions.addAll(alarmactions);
    if (stateupdatedtimestamp != null)
      result.stateupdatedtimestamp = stateupdatedtimestamp;
    if (statereason != null) result.statereason = statereason;
    if (okactions != null) result.okactions.addAll(okactions);
    if (actionssuppressedby != null)
      result.actionssuppressedby = actionssuppressedby;
    if (actionssuppressor != null) result.actionssuppressor = actionssuppressor;
    if (insufficientdataactions != null)
      result.insufficientdataactions.addAll(insufficientdataactions);
    if (alarmrule != null) result.alarmrule = alarmrule;
    return result;
  }

  CompositeAlarm._();

  factory CompositeAlarm.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CompositeAlarm.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CompositeAlarm',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(20815779, _omitFieldNames ? '' : 'actionssuppressedreason')
    ..aOS(30583613, _omitFieldNames ? '' : 'alarmdescription')
    ..aOS(46811093, _omitFieldNames ? '' : 'statetransitionedtimestamp')
    ..aOS(51935494, _omitFieldNames ? '' : 'alarmconfigurationupdatedtimestamp')
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aI(80227275, _omitFieldNames ? '' : 'actionssuppressorwaitperiod')
    ..aOB(126027878, _omitFieldNames ? '' : 'actionsenabled')
    ..aI(183723543, _omitFieldNames ? '' : 'actionssuppressorextensionperiod')
    ..aOS(262771075, _omitFieldNames ? '' : 'statereasondata')
    ..aOS(276019462, _omitFieldNames ? '' : 'alarmarn')
    ..aE<StateValue>(334526008, _omitFieldNames ? '' : 'statevalue',
        enumValues: StateValue.values)
    ..pPS(355779506, _omitFieldNames ? '' : 'alarmactions')
    ..aOS(375406848, _omitFieldNames ? '' : 'stateupdatedtimestamp')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..pPS(377443763, _omitFieldNames ? '' : 'okactions')
    ..aE<ActionsSuppressedBy>(
        377750064, _omitFieldNames ? '' : 'actionssuppressedby',
        enumValues: ActionsSuppressedBy.values)
    ..aOS(486417535, _omitFieldNames ? '' : 'actionssuppressor')
    ..pPS(498450778, _omitFieldNames ? '' : 'insufficientdataactions')
    ..aOS(516996973, _omitFieldNames ? '' : 'alarmrule')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompositeAlarm clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CompositeAlarm copyWith(void Function(CompositeAlarm) updates) =>
      super.copyWith((message) => updates(message as CompositeAlarm))
          as CompositeAlarm;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CompositeAlarm create() => CompositeAlarm._();
  @$core.override
  CompositeAlarm createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CompositeAlarm getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CompositeAlarm>(create);
  static CompositeAlarm? _defaultInstance;

  @$pb.TagNumber(20815779)
  $core.String get actionssuppressedreason => $_getSZ(0);
  @$pb.TagNumber(20815779)
  set actionssuppressedreason($core.String value) => $_setString(0, value);
  @$pb.TagNumber(20815779)
  $core.bool hasActionssuppressedreason() => $_has(0);
  @$pb.TagNumber(20815779)
  void clearActionssuppressedreason() => $_clearField(20815779);

  @$pb.TagNumber(30583613)
  $core.String get alarmdescription => $_getSZ(1);
  @$pb.TagNumber(30583613)
  set alarmdescription($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30583613)
  $core.bool hasAlarmdescription() => $_has(1);
  @$pb.TagNumber(30583613)
  void clearAlarmdescription() => $_clearField(30583613);

  @$pb.TagNumber(46811093)
  $core.String get statetransitionedtimestamp => $_getSZ(2);
  @$pb.TagNumber(46811093)
  set statetransitionedtimestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(46811093)
  $core.bool hasStatetransitionedtimestamp() => $_has(2);
  @$pb.TagNumber(46811093)
  void clearStatetransitionedtimestamp() => $_clearField(46811093);

  @$pb.TagNumber(51935494)
  $core.String get alarmconfigurationupdatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(51935494)
  set alarmconfigurationupdatedtimestamp($core.String value) =>
      $_setString(3, value);
  @$pb.TagNumber(51935494)
  $core.bool hasAlarmconfigurationupdatedtimestamp() => $_has(3);
  @$pb.TagNumber(51935494)
  void clearAlarmconfigurationupdatedtimestamp() => $_clearField(51935494);

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(4);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(4, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(4);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(80227275)
  $core.int get actionssuppressorwaitperiod => $_getIZ(5);
  @$pb.TagNumber(80227275)
  set actionssuppressorwaitperiod($core.int value) =>
      $_setSignedInt32(5, value);
  @$pb.TagNumber(80227275)
  $core.bool hasActionssuppressorwaitperiod() => $_has(5);
  @$pb.TagNumber(80227275)
  void clearActionssuppressorwaitperiod() => $_clearField(80227275);

  @$pb.TagNumber(126027878)
  $core.bool get actionsenabled => $_getBF(6);
  @$pb.TagNumber(126027878)
  set actionsenabled($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(126027878)
  $core.bool hasActionsenabled() => $_has(6);
  @$pb.TagNumber(126027878)
  void clearActionsenabled() => $_clearField(126027878);

  @$pb.TagNumber(183723543)
  $core.int get actionssuppressorextensionperiod => $_getIZ(7);
  @$pb.TagNumber(183723543)
  set actionssuppressorextensionperiod($core.int value) =>
      $_setSignedInt32(7, value);
  @$pb.TagNumber(183723543)
  $core.bool hasActionssuppressorextensionperiod() => $_has(7);
  @$pb.TagNumber(183723543)
  void clearActionssuppressorextensionperiod() => $_clearField(183723543);

  @$pb.TagNumber(262771075)
  $core.String get statereasondata => $_getSZ(8);
  @$pb.TagNumber(262771075)
  set statereasondata($core.String value) => $_setString(8, value);
  @$pb.TagNumber(262771075)
  $core.bool hasStatereasondata() => $_has(8);
  @$pb.TagNumber(262771075)
  void clearStatereasondata() => $_clearField(262771075);

  @$pb.TagNumber(276019462)
  $core.String get alarmarn => $_getSZ(9);
  @$pb.TagNumber(276019462)
  set alarmarn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(276019462)
  $core.bool hasAlarmarn() => $_has(9);
  @$pb.TagNumber(276019462)
  void clearAlarmarn() => $_clearField(276019462);

  @$pb.TagNumber(334526008)
  StateValue get statevalue => $_getN(10);
  @$pb.TagNumber(334526008)
  set statevalue(StateValue value) => $_setField(334526008, value);
  @$pb.TagNumber(334526008)
  $core.bool hasStatevalue() => $_has(10);
  @$pb.TagNumber(334526008)
  void clearStatevalue() => $_clearField(334526008);

  @$pb.TagNumber(355779506)
  $pb.PbList<$core.String> get alarmactions => $_getList(11);

  @$pb.TagNumber(375406848)
  $core.String get stateupdatedtimestamp => $_getSZ(12);
  @$pb.TagNumber(375406848)
  set stateupdatedtimestamp($core.String value) => $_setString(12, value);
  @$pb.TagNumber(375406848)
  $core.bool hasStateupdatedtimestamp() => $_has(12);
  @$pb.TagNumber(375406848)
  void clearStateupdatedtimestamp() => $_clearField(375406848);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(13);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(13, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(13);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(377443763)
  $pb.PbList<$core.String> get okactions => $_getList(14);

  @$pb.TagNumber(377750064)
  ActionsSuppressedBy get actionssuppressedby => $_getN(15);
  @$pb.TagNumber(377750064)
  set actionssuppressedby(ActionsSuppressedBy value) =>
      $_setField(377750064, value);
  @$pb.TagNumber(377750064)
  $core.bool hasActionssuppressedby() => $_has(15);
  @$pb.TagNumber(377750064)
  void clearActionssuppressedby() => $_clearField(377750064);

  @$pb.TagNumber(486417535)
  $core.String get actionssuppressor => $_getSZ(16);
  @$pb.TagNumber(486417535)
  set actionssuppressor($core.String value) => $_setString(16, value);
  @$pb.TagNumber(486417535)
  $core.bool hasActionssuppressor() => $_has(16);
  @$pb.TagNumber(486417535)
  void clearActionssuppressor() => $_clearField(486417535);

  @$pb.TagNumber(498450778)
  $pb.PbList<$core.String> get insufficientdataactions => $_getList(17);

  @$pb.TagNumber(516996973)
  $core.String get alarmrule => $_getSZ(18);
  @$pb.TagNumber(516996973)
  set alarmrule($core.String value) => $_setString(18, value);
  @$pb.TagNumber(516996973)
  $core.bool hasAlarmrule() => $_has(18);
  @$pb.TagNumber(516996973)
  void clearAlarmrule() => $_clearField(516996973);
}

class ConcurrentModificationException extends $pb.GeneratedMessage {
  factory ConcurrentModificationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConcurrentModificationException._();

  factory ConcurrentModificationException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConcurrentModificationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConcurrentModificationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModificationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConcurrentModificationException copyWith(
          void Function(ConcurrentModificationException) updates) =>
      super.copyWith(
              (message) => updates(message as ConcurrentModificationException))
          as ConcurrentModificationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConcurrentModificationException create() =>
      ConcurrentModificationException._();
  @$core.override
  ConcurrentModificationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConcurrentModificationException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConcurrentModificationException>(
          create);
  static ConcurrentModificationException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class ConflictException extends $pb.GeneratedMessage {
  factory ConflictException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ConflictException._();

  factory ConflictException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConflictException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConflictException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConflictException copyWith(void Function(ConflictException) updates) =>
      super.copyWith((message) => updates(message as ConflictException))
          as ConflictException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConflictException create() => ConflictException._();
  @$core.override
  ConflictException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConflictException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConflictException>(create);
  static ConflictException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class DashboardEntry extends $pb.GeneratedMessage {
  factory DashboardEntry({
    $fixnum.Int64? size,
    $core.String? dashboardarn,
    $core.String? lastmodified,
    $core.String? dashboardname,
  }) {
    final result = create();
    if (size != null) result.size = size;
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (lastmodified != null) result.lastmodified = lastmodified;
    if (dashboardname != null) result.dashboardname = dashboardname;
    return result;
  }

  DashboardEntry._();

  factory DashboardEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DashboardEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DashboardEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aInt64(105352829, _omitFieldNames ? '' : 'size')
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aOS(434048551, _omitFieldNames ? '' : 'lastmodified')
    ..aOS(506599873, _omitFieldNames ? '' : 'dashboardname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardEntry copyWith(void Function(DashboardEntry) updates) =>
      super.copyWith((message) => updates(message as DashboardEntry))
          as DashboardEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DashboardEntry create() => DashboardEntry._();
  @$core.override
  DashboardEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DashboardEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DashboardEntry>(create);
  static DashboardEntry? _defaultInstance;

  @$pb.TagNumber(105352829)
  $fixnum.Int64 get size => $_getI64(0);
  @$pb.TagNumber(105352829)
  set size($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(105352829)
  $core.bool hasSize() => $_has(0);
  @$pb.TagNumber(105352829)
  void clearSize() => $_clearField(105352829);

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(1);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(1);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(434048551)
  $core.String get lastmodified => $_getSZ(2);
  @$pb.TagNumber(434048551)
  set lastmodified($core.String value) => $_setString(2, value);
  @$pb.TagNumber(434048551)
  $core.bool hasLastmodified() => $_has(2);
  @$pb.TagNumber(434048551)
  void clearLastmodified() => $_clearField(434048551);

  @$pb.TagNumber(506599873)
  $core.String get dashboardname => $_getSZ(3);
  @$pb.TagNumber(506599873)
  set dashboardname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(506599873)
  $core.bool hasDashboardname() => $_has(3);
  @$pb.TagNumber(506599873)
  void clearDashboardname() => $_clearField(506599873);
}

class DashboardInvalidInputError extends $pb.GeneratedMessage {
  factory DashboardInvalidInputError({
    $core.Iterable<DashboardValidationMessage>? dashboardvalidationmessages,
    $core.String? message,
  }) {
    final result = create();
    if (dashboardvalidationmessages != null)
      result.dashboardvalidationmessages.addAll(dashboardvalidationmessages);
    if (message != null) result.message = message;
    return result;
  }

  DashboardInvalidInputError._();

  factory DashboardInvalidInputError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DashboardInvalidInputError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DashboardInvalidInputError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<DashboardValidationMessage>(
        36455465, _omitFieldNames ? '' : 'dashboardvalidationmessages',
        subBuilder: DashboardValidationMessage.create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardInvalidInputError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardInvalidInputError copyWith(
          void Function(DashboardInvalidInputError) updates) =>
      super.copyWith(
              (message) => updates(message as DashboardInvalidInputError))
          as DashboardInvalidInputError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DashboardInvalidInputError create() => DashboardInvalidInputError._();
  @$core.override
  DashboardInvalidInputError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DashboardInvalidInputError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DashboardInvalidInputError>(create);
  static DashboardInvalidInputError? _defaultInstance;

  @$pb.TagNumber(36455465)
  $pb.PbList<DashboardValidationMessage> get dashboardvalidationmessages =>
      $_getList(0);

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DashboardNotFoundError extends $pb.GeneratedMessage {
  factory DashboardNotFoundError({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  DashboardNotFoundError._();

  factory DashboardNotFoundError.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DashboardNotFoundError.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DashboardNotFoundError',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardNotFoundError clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardNotFoundError copyWith(
          void Function(DashboardNotFoundError) updates) =>
      super.copyWith((message) => updates(message as DashboardNotFoundError))
          as DashboardNotFoundError;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DashboardNotFoundError create() => DashboardNotFoundError._();
  @$core.override
  DashboardNotFoundError createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DashboardNotFoundError getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DashboardNotFoundError>(create);
  static DashboardNotFoundError? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class DashboardValidationMessage extends $pb.GeneratedMessage {
  factory DashboardValidationMessage({
    $core.String? message,
    $core.String? datapath,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (datapath != null) result.datapath = datapath;
    return result;
  }

  DashboardValidationMessage._();

  factory DashboardValidationMessage.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DashboardValidationMessage.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DashboardValidationMessage',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(505315415, _omitFieldNames ? '' : 'datapath')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardValidationMessage clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DashboardValidationMessage copyWith(
          void Function(DashboardValidationMessage) updates) =>
      super.copyWith(
              (message) => updates(message as DashboardValidationMessage))
          as DashboardValidationMessage;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DashboardValidationMessage create() => DashboardValidationMessage._();
  @$core.override
  DashboardValidationMessage createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DashboardValidationMessage getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DashboardValidationMessage>(create);
  static DashboardValidationMessage? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(505315415)
  $core.String get datapath => $_getSZ(1);
  @$pb.TagNumber(505315415)
  set datapath($core.String value) => $_setString(1, value);
  @$pb.TagNumber(505315415)
  $core.bool hasDatapath() => $_has(1);
  @$pb.TagNumber(505315415)
  void clearDatapath() => $_clearField(505315415);
}

class Datapoint extends $pb.GeneratedMessage {
  factory Datapoint({
    $core.double? maximum,
    $core.double? average,
    StandardUnit? unit,
    $core.String? timestamp,
    $core.Iterable<$core.MapEntry<$core.String, $core.double>>?
        extendedstatistics,
    $core.double? samplecount,
    $core.double? minimum,
    $core.double? sum,
  }) {
    final result = create();
    if (maximum != null) result.maximum = maximum;
    if (average != null) result.average = average;
    if (unit != null) result.unit = unit;
    if (timestamp != null) result.timestamp = timestamp;
    if (extendedstatistics != null)
      result.extendedstatistics.addEntries(extendedstatistics);
    if (samplecount != null) result.samplecount = samplecount;
    if (minimum != null) result.minimum = minimum;
    if (sum != null) result.sum = sum;
    return result;
  }

  Datapoint._();

  factory Datapoint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Datapoint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Datapoint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aD(43343394, _omitFieldNames ? '' : 'maximum')
    ..aD(134712777, _omitFieldNames ? '' : 'average')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..m<$core.String, $core.double>(
        365818700, _omitFieldNames ? '' : 'extendedstatistics',
        entryClassName: 'Datapoint.ExtendedstatisticsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OD,
        packageName: const $pb.PackageName('monitoring'))
    ..aD(384919839, _omitFieldNames ? '' : 'samplecount')
    ..aD(438536856, _omitFieldNames ? '' : 'minimum')
    ..aD(534303305, _omitFieldNames ? '' : 'sum')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datapoint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Datapoint copyWith(void Function(Datapoint) updates) =>
      super.copyWith((message) => updates(message as Datapoint)) as Datapoint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Datapoint create() => Datapoint._();
  @$core.override
  Datapoint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Datapoint getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Datapoint>(create);
  static Datapoint? _defaultInstance;

  @$pb.TagNumber(43343394)
  $core.double get maximum => $_getN(0);
  @$pb.TagNumber(43343394)
  set maximum($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(43343394)
  $core.bool hasMaximum() => $_has(0);
  @$pb.TagNumber(43343394)
  void clearMaximum() => $_clearField(43343394);

  @$pb.TagNumber(134712777)
  $core.double get average => $_getN(1);
  @$pb.TagNumber(134712777)
  set average($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(134712777)
  $core.bool hasAverage() => $_has(1);
  @$pb.TagNumber(134712777)
  void clearAverage() => $_clearField(134712777);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(2);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(2);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(3);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(3);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(365818700)
  $pb.PbMap<$core.String, $core.double> get extendedstatistics => $_getMap(4);

  @$pb.TagNumber(384919839)
  $core.double get samplecount => $_getN(5);
  @$pb.TagNumber(384919839)
  set samplecount($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(384919839)
  $core.bool hasSamplecount() => $_has(5);
  @$pb.TagNumber(384919839)
  void clearSamplecount() => $_clearField(384919839);

  @$pb.TagNumber(438536856)
  $core.double get minimum => $_getN(6);
  @$pb.TagNumber(438536856)
  set minimum($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(438536856)
  $core.bool hasMinimum() => $_has(6);
  @$pb.TagNumber(438536856)
  void clearMinimum() => $_clearField(438536856);

  @$pb.TagNumber(534303305)
  $core.double get sum => $_getN(7);
  @$pb.TagNumber(534303305)
  set sum($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(534303305)
  $core.bool hasSum() => $_has(7);
  @$pb.TagNumber(534303305)
  void clearSum() => $_clearField(534303305);
}

class DeleteAlarmMuteRuleInput extends $pb.GeneratedMessage {
  factory DeleteAlarmMuteRuleInput({
    $core.String? alarmmuterulename,
  }) {
    final result = create();
    if (alarmmuterulename != null) result.alarmmuterulename = alarmmuterulename;
    return result;
  }

  DeleteAlarmMuteRuleInput._();

  factory DeleteAlarmMuteRuleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAlarmMuteRuleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAlarmMuteRuleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(530877511, _omitFieldNames ? '' : 'alarmmuterulename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAlarmMuteRuleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAlarmMuteRuleInput copyWith(
          void Function(DeleteAlarmMuteRuleInput) updates) =>
      super.copyWith((message) => updates(message as DeleteAlarmMuteRuleInput))
          as DeleteAlarmMuteRuleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAlarmMuteRuleInput create() => DeleteAlarmMuteRuleInput._();
  @$core.override
  DeleteAlarmMuteRuleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAlarmMuteRuleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAlarmMuteRuleInput>(create);
  static DeleteAlarmMuteRuleInput? _defaultInstance;

  @$pb.TagNumber(530877511)
  $core.String get alarmmuterulename => $_getSZ(0);
  @$pb.TagNumber(530877511)
  set alarmmuterulename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(530877511)
  $core.bool hasAlarmmuterulename() => $_has(0);
  @$pb.TagNumber(530877511)
  void clearAlarmmuterulename() => $_clearField(530877511);
}

class DeleteAlarmsInput extends $pb.GeneratedMessage {
  factory DeleteAlarmsInput({
    $core.Iterable<$core.String>? alarmnames,
  }) {
    final result = create();
    if (alarmnames != null) result.alarmnames.addAll(alarmnames);
    return result;
  }

  DeleteAlarmsInput._();

  factory DeleteAlarmsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAlarmsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAlarmsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(90874583, _omitFieldNames ? '' : 'alarmnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAlarmsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAlarmsInput copyWith(void Function(DeleteAlarmsInput) updates) =>
      super.copyWith((message) => updates(message as DeleteAlarmsInput))
          as DeleteAlarmsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAlarmsInput create() => DeleteAlarmsInput._();
  @$core.override
  DeleteAlarmsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAlarmsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAlarmsInput>(create);
  static DeleteAlarmsInput? _defaultInstance;

  @$pb.TagNumber(90874583)
  $pb.PbList<$core.String> get alarmnames => $_getList(0);
}

class DeleteAnomalyDetectorInput extends $pb.GeneratedMessage {
  factory DeleteAnomalyDetectorInput({
    $core.String? metricname,
    $core.String? stat,
    $core.String? namespace,
    SingleMetricAnomalyDetector? singlemetricanomalydetector,
    MetricMathAnomalyDetector? metricmathanomalydetector,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (stat != null) result.stat = stat;
    if (namespace != null) result.namespace = namespace;
    if (singlemetricanomalydetector != null)
      result.singlemetricanomalydetector = singlemetricanomalydetector;
    if (metricmathanomalydetector != null)
      result.metricmathanomalydetector = metricmathanomalydetector;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  DeleteAnomalyDetectorInput._();

  factory DeleteAnomalyDetectorInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAnomalyDetectorInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAnomalyDetectorInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(325549752, _omitFieldNames ? '' : 'stat')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..aOM<SingleMetricAnomalyDetector>(
        418531933, _omitFieldNames ? '' : 'singlemetricanomalydetector',
        subBuilder: SingleMetricAnomalyDetector.create)
    ..aOM<MetricMathAnomalyDetector>(
        439381475, _omitFieldNames ? '' : 'metricmathanomalydetector',
        subBuilder: MetricMathAnomalyDetector.create)
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAnomalyDetectorInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAnomalyDetectorInput copyWith(
          void Function(DeleteAnomalyDetectorInput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteAnomalyDetectorInput))
          as DeleteAnomalyDetectorInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAnomalyDetectorInput create() => DeleteAnomalyDetectorInput._();
  @$core.override
  DeleteAnomalyDetectorInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAnomalyDetectorInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAnomalyDetectorInput>(create);
  static DeleteAnomalyDetectorInput? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(325549752)
  $core.String get stat => $_getSZ(1);
  @$pb.TagNumber(325549752)
  set stat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(325549752)
  $core.bool hasStat() => $_has(1);
  @$pb.TagNumber(325549752)
  void clearStat() => $_clearField(325549752);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(2);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(2, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(2);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector get singlemetricanomalydetector => $_getN(3);
  @$pb.TagNumber(418531933)
  set singlemetricanomalydetector(SingleMetricAnomalyDetector value) =>
      $_setField(418531933, value);
  @$pb.TagNumber(418531933)
  $core.bool hasSinglemetricanomalydetector() => $_has(3);
  @$pb.TagNumber(418531933)
  void clearSinglemetricanomalydetector() => $_clearField(418531933);
  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector ensureSinglemetricanomalydetector() =>
      $_ensure(3);

  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector get metricmathanomalydetector => $_getN(4);
  @$pb.TagNumber(439381475)
  set metricmathanomalydetector(MetricMathAnomalyDetector value) =>
      $_setField(439381475, value);
  @$pb.TagNumber(439381475)
  $core.bool hasMetricmathanomalydetector() => $_has(4);
  @$pb.TagNumber(439381475)
  void clearMetricmathanomalydetector() => $_clearField(439381475);
  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector ensureMetricmathanomalydetector() => $_ensure(4);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(5);
}

class DeleteAnomalyDetectorOutput extends $pb.GeneratedMessage {
  factory DeleteAnomalyDetectorOutput() => create();

  DeleteAnomalyDetectorOutput._();

  factory DeleteAnomalyDetectorOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteAnomalyDetectorOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteAnomalyDetectorOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAnomalyDetectorOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteAnomalyDetectorOutput copyWith(
          void Function(DeleteAnomalyDetectorOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DeleteAnomalyDetectorOutput))
          as DeleteAnomalyDetectorOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteAnomalyDetectorOutput create() =>
      DeleteAnomalyDetectorOutput._();
  @$core.override
  DeleteAnomalyDetectorOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteAnomalyDetectorOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteAnomalyDetectorOutput>(create);
  static DeleteAnomalyDetectorOutput? _defaultInstance;
}

class DeleteDashboardsInput extends $pb.GeneratedMessage {
  factory DeleteDashboardsInput({
    $core.Iterable<$core.String>? dashboardnames,
  }) {
    final result = create();
    if (dashboardnames != null) result.dashboardnames.addAll(dashboardnames);
    return result;
  }

  DeleteDashboardsInput._();

  factory DeleteDashboardsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDashboardsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDashboardsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(467563722, _omitFieldNames ? '' : 'dashboardnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardsInput copyWith(
          void Function(DeleteDashboardsInput) updates) =>
      super.copyWith((message) => updates(message as DeleteDashboardsInput))
          as DeleteDashboardsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDashboardsInput create() => DeleteDashboardsInput._();
  @$core.override
  DeleteDashboardsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDashboardsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDashboardsInput>(create);
  static DeleteDashboardsInput? _defaultInstance;

  @$pb.TagNumber(467563722)
  $pb.PbList<$core.String> get dashboardnames => $_getList(0);
}

class DeleteDashboardsOutput extends $pb.GeneratedMessage {
  factory DeleteDashboardsOutput() => create();

  DeleteDashboardsOutput._();

  factory DeleteDashboardsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteDashboardsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteDashboardsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteDashboardsOutput copyWith(
          void Function(DeleteDashboardsOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteDashboardsOutput))
          as DeleteDashboardsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteDashboardsOutput create() => DeleteDashboardsOutput._();
  @$core.override
  DeleteDashboardsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteDashboardsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteDashboardsOutput>(create);
  static DeleteDashboardsOutput? _defaultInstance;
}

class DeleteInsightRulesInput extends $pb.GeneratedMessage {
  factory DeleteInsightRulesInput({
    $core.Iterable<$core.String>? rulenames,
  }) {
    final result = create();
    if (rulenames != null) result.rulenames.addAll(rulenames);
    return result;
  }

  DeleteInsightRulesInput._();

  factory DeleteInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(267949170, _omitFieldNames ? '' : 'rulenames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteInsightRulesInput copyWith(
          void Function(DeleteInsightRulesInput) updates) =>
      super.copyWith((message) => updates(message as DeleteInsightRulesInput))
          as DeleteInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteInsightRulesInput create() => DeleteInsightRulesInput._();
  @$core.override
  DeleteInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteInsightRulesInput>(create);
  static DeleteInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(267949170)
  $pb.PbList<$core.String> get rulenames => $_getList(0);
}

class DeleteInsightRulesOutput extends $pb.GeneratedMessage {
  factory DeleteInsightRulesOutput({
    $core.Iterable<PartialFailure>? failures,
  }) {
    final result = create();
    if (failures != null) result.failures.addAll(failures);
    return result;
  }

  DeleteInsightRulesOutput._();

  factory DeleteInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<PartialFailure>(335467271, _omitFieldNames ? '' : 'failures',
        subBuilder: PartialFailure.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteInsightRulesOutput copyWith(
          void Function(DeleteInsightRulesOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteInsightRulesOutput))
          as DeleteInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteInsightRulesOutput create() => DeleteInsightRulesOutput._();
  @$core.override
  DeleteInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteInsightRulesOutput>(create);
  static DeleteInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(335467271)
  $pb.PbList<PartialFailure> get failures => $_getList(0);
}

class DeleteMetricStreamInput extends $pb.GeneratedMessage {
  factory DeleteMetricStreamInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  DeleteMetricStreamInput._();

  factory DeleteMetricStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMetricStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMetricStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMetricStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMetricStreamInput copyWith(
          void Function(DeleteMetricStreamInput) updates) =>
      super.copyWith((message) => updates(message as DeleteMetricStreamInput))
          as DeleteMetricStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMetricStreamInput create() => DeleteMetricStreamInput._();
  @$core.override
  DeleteMetricStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMetricStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMetricStreamInput>(create);
  static DeleteMetricStreamInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class DeleteMetricStreamOutput extends $pb.GeneratedMessage {
  factory DeleteMetricStreamOutput() => create();

  DeleteMetricStreamOutput._();

  factory DeleteMetricStreamOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DeleteMetricStreamOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DeleteMetricStreamOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMetricStreamOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DeleteMetricStreamOutput copyWith(
          void Function(DeleteMetricStreamOutput) updates) =>
      super.copyWith((message) => updates(message as DeleteMetricStreamOutput))
          as DeleteMetricStreamOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DeleteMetricStreamOutput create() => DeleteMetricStreamOutput._();
  @$core.override
  DeleteMetricStreamOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DeleteMetricStreamOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DeleteMetricStreamOutput>(create);
  static DeleteMetricStreamOutput? _defaultInstance;
}

class DescribeAlarmContributorsInput extends $pb.GeneratedMessage {
  factory DescribeAlarmContributorsInput({
    $core.String? alarmname,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (alarmname != null) result.alarmname = alarmname;
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  DescribeAlarmContributorsInput._();

  factory DescribeAlarmContributorsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmContributorsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmContributorsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmContributorsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmContributorsInput copyWith(
          void Function(DescribeAlarmContributorsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAlarmContributorsInput))
          as DescribeAlarmContributorsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmContributorsInput create() =>
      DescribeAlarmContributorsInput._();
  @$core.override
  DescribeAlarmContributorsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmContributorsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmContributorsInput>(create);
  static DescribeAlarmContributorsInput? _defaultInstance;

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(0);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(0);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class DescribeAlarmContributorsOutput extends $pb.GeneratedMessage {
  factory DescribeAlarmContributorsOutput({
    $core.String? nexttoken,
    $core.Iterable<AlarmContributor>? alarmcontributors,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (alarmcontributors != null)
      result.alarmcontributors.addAll(alarmcontributors);
    return result;
  }

  DescribeAlarmContributorsOutput._();

  factory DescribeAlarmContributorsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmContributorsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmContributorsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<AlarmContributor>(
        265859209, _omitFieldNames ? '' : 'alarmcontributors',
        subBuilder: AlarmContributor.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmContributorsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmContributorsOutput copyWith(
          void Function(DescribeAlarmContributorsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAlarmContributorsOutput))
          as DescribeAlarmContributorsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmContributorsOutput create() =>
      DescribeAlarmContributorsOutput._();
  @$core.override
  DescribeAlarmContributorsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmContributorsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmContributorsOutput>(
          create);
  static DescribeAlarmContributorsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(265859209)
  $pb.PbList<AlarmContributor> get alarmcontributors => $_getList(1);
}

class DescribeAlarmHistoryInput extends $pb.GeneratedMessage {
  factory DescribeAlarmHistoryInput({
    $core.String? alarmname,
    $core.String? enddate,
    $core.String? nexttoken,
    $core.int? maxrecords,
    $core.String? alarmcontributorid,
    ScanBy? scanby,
    HistoryItemType? historyitemtype,
    $core.String? startdate,
    $core.Iterable<AlarmType>? alarmtypes,
  }) {
    final result = create();
    if (alarmname != null) result.alarmname = alarmname;
    if (enddate != null) result.enddate = enddate;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxrecords != null) result.maxrecords = maxrecords;
    if (alarmcontributorid != null)
      result.alarmcontributorid = alarmcontributorid;
    if (scanby != null) result.scanby = scanby;
    if (historyitemtype != null) result.historyitemtype = historyitemtype;
    if (startdate != null) result.startdate = startdate;
    if (alarmtypes != null) result.alarmtypes.addAll(alarmtypes);
    return result;
  }

  DescribeAlarmHistoryInput._();

  factory DescribeAlarmHistoryInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmHistoryInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmHistoryInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(77486543, _omitFieldNames ? '' : 'enddate')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(220314370, _omitFieldNames ? '' : 'maxrecords')
    ..aOS(305558919, _omitFieldNames ? '' : 'alarmcontributorid')
    ..aE<ScanBy>(344142854, _omitFieldNames ? '' : 'scanby',
        enumValues: ScanBy.values)
    ..aE<HistoryItemType>(442439537, _omitFieldNames ? '' : 'historyitemtype',
        enumValues: HistoryItemType.values)
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..pc<AlarmType>(
        445751884, _omitFieldNames ? '' : 'alarmtypes', $pb.PbFieldType.KE,
        valueOf: AlarmType.valueOf,
        enumValues: AlarmType.values,
        defaultEnumValue: AlarmType.ALARM_TYPE_COMPOSITEALARM)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmHistoryInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmHistoryInput copyWith(
          void Function(DescribeAlarmHistoryInput) updates) =>
      super.copyWith((message) => updates(message as DescribeAlarmHistoryInput))
          as DescribeAlarmHistoryInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmHistoryInput create() => DescribeAlarmHistoryInput._();
  @$core.override
  DescribeAlarmHistoryInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmHistoryInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmHistoryInput>(create);
  static DescribeAlarmHistoryInput? _defaultInstance;

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(0);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(0);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(77486543)
  $core.String get enddate => $_getSZ(1);
  @$pb.TagNumber(77486543)
  set enddate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(77486543)
  $core.bool hasEnddate() => $_has(1);
  @$pb.TagNumber(77486543)
  void clearEnddate() => $_clearField(77486543);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(220314370)
  $core.int get maxrecords => $_getIZ(3);
  @$pb.TagNumber(220314370)
  set maxrecords($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(220314370)
  $core.bool hasMaxrecords() => $_has(3);
  @$pb.TagNumber(220314370)
  void clearMaxrecords() => $_clearField(220314370);

  @$pb.TagNumber(305558919)
  $core.String get alarmcontributorid => $_getSZ(4);
  @$pb.TagNumber(305558919)
  set alarmcontributorid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(305558919)
  $core.bool hasAlarmcontributorid() => $_has(4);
  @$pb.TagNumber(305558919)
  void clearAlarmcontributorid() => $_clearField(305558919);

  @$pb.TagNumber(344142854)
  ScanBy get scanby => $_getN(5);
  @$pb.TagNumber(344142854)
  set scanby(ScanBy value) => $_setField(344142854, value);
  @$pb.TagNumber(344142854)
  $core.bool hasScanby() => $_has(5);
  @$pb.TagNumber(344142854)
  void clearScanby() => $_clearField(344142854);

  @$pb.TagNumber(442439537)
  HistoryItemType get historyitemtype => $_getN(6);
  @$pb.TagNumber(442439537)
  set historyitemtype(HistoryItemType value) => $_setField(442439537, value);
  @$pb.TagNumber(442439537)
  $core.bool hasHistoryitemtype() => $_has(6);
  @$pb.TagNumber(442439537)
  void clearHistoryitemtype() => $_clearField(442439537);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(7);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(7, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(7);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(445751884)
  $pb.PbList<AlarmType> get alarmtypes => $_getList(8);
}

class DescribeAlarmHistoryOutput extends $pb.GeneratedMessage {
  factory DescribeAlarmHistoryOutput({
    $core.Iterable<AlarmHistoryItem>? alarmhistoryitems,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (alarmhistoryitems != null)
      result.alarmhistoryitems.addAll(alarmhistoryitems);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  DescribeAlarmHistoryOutput._();

  factory DescribeAlarmHistoryOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmHistoryOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmHistoryOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<AlarmHistoryItem>(
        152760951, _omitFieldNames ? '' : 'alarmhistoryitems',
        subBuilder: AlarmHistoryItem.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmHistoryOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmHistoryOutput copyWith(
          void Function(DescribeAlarmHistoryOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAlarmHistoryOutput))
          as DescribeAlarmHistoryOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmHistoryOutput create() => DescribeAlarmHistoryOutput._();
  @$core.override
  DescribeAlarmHistoryOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmHistoryOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmHistoryOutput>(create);
  static DescribeAlarmHistoryOutput? _defaultInstance;

  @$pb.TagNumber(152760951)
  $pb.PbList<AlarmHistoryItem> get alarmhistoryitems => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class DescribeAlarmsForMetricInput extends $pb.GeneratedMessage {
  factory DescribeAlarmsForMetricInput({
    Statistic? statistic,
    $core.String? metricname,
    $core.int? period,
    StandardUnit? unit,
    $core.String? extendedstatistic,
    $core.String? namespace,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (statistic != null) result.statistic = statistic;
    if (metricname != null) result.metricname = metricname;
    if (period != null) result.period = period;
    if (unit != null) result.unit = unit;
    if (extendedstatistic != null) result.extendedstatistic = extendedstatistic;
    if (namespace != null) result.namespace = namespace;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  DescribeAlarmsForMetricInput._();

  factory DescribeAlarmsForMetricInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmsForMetricInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmsForMetricInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aE<Statistic>(67293470, _omitFieldNames ? '' : 'statistic',
        enumValues: Statistic.values)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aOS(285620763, _omitFieldNames ? '' : 'extendedstatistic')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsForMetricInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsForMetricInput copyWith(
          void Function(DescribeAlarmsForMetricInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAlarmsForMetricInput))
          as DescribeAlarmsForMetricInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsForMetricInput create() =>
      DescribeAlarmsForMetricInput._();
  @$core.override
  DescribeAlarmsForMetricInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsForMetricInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmsForMetricInput>(create);
  static DescribeAlarmsForMetricInput? _defaultInstance;

  @$pb.TagNumber(67293470)
  Statistic get statistic => $_getN(0);
  @$pb.TagNumber(67293470)
  set statistic(Statistic value) => $_setField(67293470, value);
  @$pb.TagNumber(67293470)
  $core.bool hasStatistic() => $_has(0);
  @$pb.TagNumber(67293470)
  void clearStatistic() => $_clearField(67293470);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(2);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(2);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(3);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(3);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(285620763)
  $core.String get extendedstatistic => $_getSZ(4);
  @$pb.TagNumber(285620763)
  set extendedstatistic($core.String value) => $_setString(4, value);
  @$pb.TagNumber(285620763)
  $core.bool hasExtendedstatistic() => $_has(4);
  @$pb.TagNumber(285620763)
  void clearExtendedstatistic() => $_clearField(285620763);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(5);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(5, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(5);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(6);
}

class DescribeAlarmsForMetricOutput extends $pb.GeneratedMessage {
  factory DescribeAlarmsForMetricOutput({
    $core.Iterable<MetricAlarm>? metricalarms,
  }) {
    final result = create();
    if (metricalarms != null) result.metricalarms.addAll(metricalarms);
    return result;
  }

  DescribeAlarmsForMetricOutput._();

  factory DescribeAlarmsForMetricOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmsForMetricOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmsForMetricOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<MetricAlarm>(117806408, _omitFieldNames ? '' : 'metricalarms',
        subBuilder: MetricAlarm.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsForMetricOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsForMetricOutput copyWith(
          void Function(DescribeAlarmsForMetricOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAlarmsForMetricOutput))
          as DescribeAlarmsForMetricOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsForMetricOutput create() =>
      DescribeAlarmsForMetricOutput._();
  @$core.override
  DescribeAlarmsForMetricOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsForMetricOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmsForMetricOutput>(create);
  static DescribeAlarmsForMetricOutput? _defaultInstance;

  @$pb.TagNumber(117806408)
  $pb.PbList<MetricAlarm> get metricalarms => $_getList(0);
}

class DescribeAlarmsInput extends $pb.GeneratedMessage {
  factory DescribeAlarmsInput({
    $core.String? actionprefix,
    $core.Iterable<$core.String>? alarmnames,
    $core.String? alarmnameprefix,
    $core.String? childrenofalarmname,
    $core.String? nexttoken,
    $core.int? maxrecords,
    StateValue? statevalue,
    $core.Iterable<AlarmType>? alarmtypes,
    $core.String? parentsofalarmname,
  }) {
    final result = create();
    if (actionprefix != null) result.actionprefix = actionprefix;
    if (alarmnames != null) result.alarmnames.addAll(alarmnames);
    if (alarmnameprefix != null) result.alarmnameprefix = alarmnameprefix;
    if (childrenofalarmname != null)
      result.childrenofalarmname = childrenofalarmname;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxrecords != null) result.maxrecords = maxrecords;
    if (statevalue != null) result.statevalue = statevalue;
    if (alarmtypes != null) result.alarmtypes.addAll(alarmtypes);
    if (parentsofalarmname != null)
      result.parentsofalarmname = parentsofalarmname;
    return result;
  }

  DescribeAlarmsInput._();

  factory DescribeAlarmsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(67679956, _omitFieldNames ? '' : 'actionprefix')
    ..pPS(90874583, _omitFieldNames ? '' : 'alarmnames')
    ..aOS(96412406, _omitFieldNames ? '' : 'alarmnameprefix')
    ..aOS(204127700, _omitFieldNames ? '' : 'childrenofalarmname')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(220314370, _omitFieldNames ? '' : 'maxrecords')
    ..aE<StateValue>(334526008, _omitFieldNames ? '' : 'statevalue',
        enumValues: StateValue.values)
    ..pc<AlarmType>(
        445751884, _omitFieldNames ? '' : 'alarmtypes', $pb.PbFieldType.KE,
        valueOf: AlarmType.valueOf,
        enumValues: AlarmType.values,
        defaultEnumValue: AlarmType.ALARM_TYPE_COMPOSITEALARM)
    ..aOS(512433260, _omitFieldNames ? '' : 'parentsofalarmname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsInput copyWith(void Function(DescribeAlarmsInput) updates) =>
      super.copyWith((message) => updates(message as DescribeAlarmsInput))
          as DescribeAlarmsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsInput create() => DescribeAlarmsInput._();
  @$core.override
  DescribeAlarmsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmsInput>(create);
  static DescribeAlarmsInput? _defaultInstance;

  @$pb.TagNumber(67679956)
  $core.String get actionprefix => $_getSZ(0);
  @$pb.TagNumber(67679956)
  set actionprefix($core.String value) => $_setString(0, value);
  @$pb.TagNumber(67679956)
  $core.bool hasActionprefix() => $_has(0);
  @$pb.TagNumber(67679956)
  void clearActionprefix() => $_clearField(67679956);

  @$pb.TagNumber(90874583)
  $pb.PbList<$core.String> get alarmnames => $_getList(1);

  @$pb.TagNumber(96412406)
  $core.String get alarmnameprefix => $_getSZ(2);
  @$pb.TagNumber(96412406)
  set alarmnameprefix($core.String value) => $_setString(2, value);
  @$pb.TagNumber(96412406)
  $core.bool hasAlarmnameprefix() => $_has(2);
  @$pb.TagNumber(96412406)
  void clearAlarmnameprefix() => $_clearField(96412406);

  @$pb.TagNumber(204127700)
  $core.String get childrenofalarmname => $_getSZ(3);
  @$pb.TagNumber(204127700)
  set childrenofalarmname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(204127700)
  $core.bool hasChildrenofalarmname() => $_has(3);
  @$pb.TagNumber(204127700)
  void clearChildrenofalarmname() => $_clearField(204127700);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(4);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(4);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(220314370)
  $core.int get maxrecords => $_getIZ(5);
  @$pb.TagNumber(220314370)
  set maxrecords($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(220314370)
  $core.bool hasMaxrecords() => $_has(5);
  @$pb.TagNumber(220314370)
  void clearMaxrecords() => $_clearField(220314370);

  @$pb.TagNumber(334526008)
  StateValue get statevalue => $_getN(6);
  @$pb.TagNumber(334526008)
  set statevalue(StateValue value) => $_setField(334526008, value);
  @$pb.TagNumber(334526008)
  $core.bool hasStatevalue() => $_has(6);
  @$pb.TagNumber(334526008)
  void clearStatevalue() => $_clearField(334526008);

  @$pb.TagNumber(445751884)
  $pb.PbList<AlarmType> get alarmtypes => $_getList(7);

  @$pb.TagNumber(512433260)
  $core.String get parentsofalarmname => $_getSZ(8);
  @$pb.TagNumber(512433260)
  set parentsofalarmname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(512433260)
  $core.bool hasParentsofalarmname() => $_has(8);
  @$pb.TagNumber(512433260)
  void clearParentsofalarmname() => $_clearField(512433260);
}

class DescribeAlarmsOutput extends $pb.GeneratedMessage {
  factory DescribeAlarmsOutput({
    $core.Iterable<MetricAlarm>? metricalarms,
    $core.Iterable<CompositeAlarm>? compositealarms,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (metricalarms != null) result.metricalarms.addAll(metricalarms);
    if (compositealarms != null) result.compositealarms.addAll(compositealarms);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  DescribeAlarmsOutput._();

  factory DescribeAlarmsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAlarmsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAlarmsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<MetricAlarm>(117806408, _omitFieldNames ? '' : 'metricalarms',
        subBuilder: MetricAlarm.create)
    ..pPM<CompositeAlarm>(214076095, _omitFieldNames ? '' : 'compositealarms',
        subBuilder: CompositeAlarm.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAlarmsOutput copyWith(void Function(DescribeAlarmsOutput) updates) =>
      super.copyWith((message) => updates(message as DescribeAlarmsOutput))
          as DescribeAlarmsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsOutput create() => DescribeAlarmsOutput._();
  @$core.override
  DescribeAlarmsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAlarmsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAlarmsOutput>(create);
  static DescribeAlarmsOutput? _defaultInstance;

  @$pb.TagNumber(117806408)
  $pb.PbList<MetricAlarm> get metricalarms => $_getList(0);

  @$pb.TagNumber(214076095)
  $pb.PbList<CompositeAlarm> get compositealarms => $_getList(1);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class DescribeAnomalyDetectorsInput extends $pb.GeneratedMessage {
  factory DescribeAnomalyDetectorsInput({
    $core.String? metricname,
    $core.String? nexttoken,
    $core.Iterable<AnomalyDetectorType>? anomalydetectortypes,
    $core.int? maxresults,
    $core.String? namespace,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (anomalydetectortypes != null)
      result.anomalydetectortypes.addAll(anomalydetectortypes);
    if (maxresults != null) result.maxresults = maxresults;
    if (namespace != null) result.namespace = namespace;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  DescribeAnomalyDetectorsInput._();

  factory DescribeAnomalyDetectorsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAnomalyDetectorsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAnomalyDetectorsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pc<AnomalyDetectorType>(274334014,
        _omitFieldNames ? '' : 'anomalydetectortypes', $pb.PbFieldType.KE,
        valueOf: AnomalyDetectorType.valueOf,
        enumValues: AnomalyDetectorType.values,
        defaultEnumValue:
            AnomalyDetectorType.ANOMALY_DETECTOR_TYPE_SINGLE_METRIC)
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAnomalyDetectorsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAnomalyDetectorsInput copyWith(
          void Function(DescribeAnomalyDetectorsInput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAnomalyDetectorsInput))
          as DescribeAnomalyDetectorsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAnomalyDetectorsInput create() =>
      DescribeAnomalyDetectorsInput._();
  @$core.override
  DescribeAnomalyDetectorsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAnomalyDetectorsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAnomalyDetectorsInput>(create);
  static DescribeAnomalyDetectorsInput? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(274334014)
  $pb.PbList<AnomalyDetectorType> get anomalydetectortypes => $_getList(2);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(3);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(3);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(4);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(4, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(4);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(5);
}

class DescribeAnomalyDetectorsOutput extends $pb.GeneratedMessage {
  factory DescribeAnomalyDetectorsOutput({
    $core.String? nexttoken,
    $core.Iterable<AnomalyDetector>? anomalydetectors,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (anomalydetectors != null)
      result.anomalydetectors.addAll(anomalydetectors);
    return result;
  }

  DescribeAnomalyDetectorsOutput._();

  factory DescribeAnomalyDetectorsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeAnomalyDetectorsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeAnomalyDetectorsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<AnomalyDetector>(399089190, _omitFieldNames ? '' : 'anomalydetectors',
        subBuilder: AnomalyDetector.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAnomalyDetectorsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeAnomalyDetectorsOutput copyWith(
          void Function(DescribeAnomalyDetectorsOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeAnomalyDetectorsOutput))
          as DescribeAnomalyDetectorsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeAnomalyDetectorsOutput create() =>
      DescribeAnomalyDetectorsOutput._();
  @$core.override
  DescribeAnomalyDetectorsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeAnomalyDetectorsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeAnomalyDetectorsOutput>(create);
  static DescribeAnomalyDetectorsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(399089190)
  $pb.PbList<AnomalyDetector> get anomalydetectors => $_getList(1);
}

class DescribeInsightRulesInput extends $pb.GeneratedMessage {
  factory DescribeInsightRulesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  DescribeInsightRulesInput._();

  factory DescribeInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeInsightRulesInput copyWith(
          void Function(DescribeInsightRulesInput) updates) =>
      super.copyWith((message) => updates(message as DescribeInsightRulesInput))
          as DescribeInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeInsightRulesInput create() => DescribeInsightRulesInput._();
  @$core.override
  DescribeInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeInsightRulesInput>(create);
  static DescribeInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class DescribeInsightRulesOutput extends $pb.GeneratedMessage {
  factory DescribeInsightRulesOutput({
    $core.Iterable<InsightRule>? insightrules,
    $core.String? nexttoken,
  }) {
    final result = create();
    if (insightrules != null) result.insightrules.addAll(insightrules);
    if (nexttoken != null) result.nexttoken = nexttoken;
    return result;
  }

  DescribeInsightRulesOutput._();

  factory DescribeInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DescribeInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DescribeInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<InsightRule>(87263475, _omitFieldNames ? '' : 'insightrules',
        subBuilder: InsightRule.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DescribeInsightRulesOutput copyWith(
          void Function(DescribeInsightRulesOutput) updates) =>
      super.copyWith(
              (message) => updates(message as DescribeInsightRulesOutput))
          as DescribeInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DescribeInsightRulesOutput create() => DescribeInsightRulesOutput._();
  @$core.override
  DescribeInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DescribeInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DescribeInsightRulesOutput>(create);
  static DescribeInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(87263475)
  $pb.PbList<InsightRule> get insightrules => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);
}

class Dimension extends $pb.GeneratedMessage {
  factory Dimension({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  Dimension._();

  factory Dimension.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Dimension.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Dimension',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Dimension copyWith(void Function(Dimension) updates) =>
      super.copyWith((message) => updates(message as Dimension)) as Dimension;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Dimension create() => Dimension._();
  @$core.override
  Dimension createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Dimension getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Dimension>(create);
  static Dimension? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class DimensionFilter extends $pb.GeneratedMessage {
  factory DimensionFilter({
    $core.String? name,
    $core.String? value,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (value != null) result.value = value;
    return result;
  }

  DimensionFilter._();

  factory DimensionFilter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DimensionFilter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DimensionFilter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DimensionFilter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DimensionFilter copyWith(void Function(DimensionFilter) updates) =>
      super.copyWith((message) => updates(message as DimensionFilter))
          as DimensionFilter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DimensionFilter create() => DimensionFilter._();
  @$core.override
  DimensionFilter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DimensionFilter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DimensionFilter>(create);
  static DimensionFilter? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class DisableAlarmActionsInput extends $pb.GeneratedMessage {
  factory DisableAlarmActionsInput({
    $core.Iterable<$core.String>? alarmnames,
  }) {
    final result = create();
    if (alarmnames != null) result.alarmnames.addAll(alarmnames);
    return result;
  }

  DisableAlarmActionsInput._();

  factory DisableAlarmActionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableAlarmActionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableAlarmActionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(90874583, _omitFieldNames ? '' : 'alarmnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableAlarmActionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableAlarmActionsInput copyWith(
          void Function(DisableAlarmActionsInput) updates) =>
      super.copyWith((message) => updates(message as DisableAlarmActionsInput))
          as DisableAlarmActionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableAlarmActionsInput create() => DisableAlarmActionsInput._();
  @$core.override
  DisableAlarmActionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableAlarmActionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableAlarmActionsInput>(create);
  static DisableAlarmActionsInput? _defaultInstance;

  @$pb.TagNumber(90874583)
  $pb.PbList<$core.String> get alarmnames => $_getList(0);
}

class DisableInsightRulesInput extends $pb.GeneratedMessage {
  factory DisableInsightRulesInput({
    $core.Iterable<$core.String>? rulenames,
  }) {
    final result = create();
    if (rulenames != null) result.rulenames.addAll(rulenames);
    return result;
  }

  DisableInsightRulesInput._();

  factory DisableInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(267949170, _omitFieldNames ? '' : 'rulenames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableInsightRulesInput copyWith(
          void Function(DisableInsightRulesInput) updates) =>
      super.copyWith((message) => updates(message as DisableInsightRulesInput))
          as DisableInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableInsightRulesInput create() => DisableInsightRulesInput._();
  @$core.override
  DisableInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableInsightRulesInput>(create);
  static DisableInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(267949170)
  $pb.PbList<$core.String> get rulenames => $_getList(0);
}

class DisableInsightRulesOutput extends $pb.GeneratedMessage {
  factory DisableInsightRulesOutput({
    $core.Iterable<PartialFailure>? failures,
  }) {
    final result = create();
    if (failures != null) result.failures.addAll(failures);
    return result;
  }

  DisableInsightRulesOutput._();

  factory DisableInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DisableInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DisableInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<PartialFailure>(335467271, _omitFieldNames ? '' : 'failures',
        subBuilder: PartialFailure.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DisableInsightRulesOutput copyWith(
          void Function(DisableInsightRulesOutput) updates) =>
      super.copyWith((message) => updates(message as DisableInsightRulesOutput))
          as DisableInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DisableInsightRulesOutput create() => DisableInsightRulesOutput._();
  @$core.override
  DisableInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DisableInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DisableInsightRulesOutput>(create);
  static DisableInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(335467271)
  $pb.PbList<PartialFailure> get failures => $_getList(0);
}

class EnableAlarmActionsInput extends $pb.GeneratedMessage {
  factory EnableAlarmActionsInput({
    $core.Iterable<$core.String>? alarmnames,
  }) {
    final result = create();
    if (alarmnames != null) result.alarmnames.addAll(alarmnames);
    return result;
  }

  EnableAlarmActionsInput._();

  factory EnableAlarmActionsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableAlarmActionsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableAlarmActionsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(90874583, _omitFieldNames ? '' : 'alarmnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableAlarmActionsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableAlarmActionsInput copyWith(
          void Function(EnableAlarmActionsInput) updates) =>
      super.copyWith((message) => updates(message as EnableAlarmActionsInput))
          as EnableAlarmActionsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableAlarmActionsInput create() => EnableAlarmActionsInput._();
  @$core.override
  EnableAlarmActionsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableAlarmActionsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableAlarmActionsInput>(create);
  static EnableAlarmActionsInput? _defaultInstance;

  @$pb.TagNumber(90874583)
  $pb.PbList<$core.String> get alarmnames => $_getList(0);
}

class EnableInsightRulesInput extends $pb.GeneratedMessage {
  factory EnableInsightRulesInput({
    $core.Iterable<$core.String>? rulenames,
  }) {
    final result = create();
    if (rulenames != null) result.rulenames.addAll(rulenames);
    return result;
  }

  EnableInsightRulesInput._();

  factory EnableInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(267949170, _omitFieldNames ? '' : 'rulenames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableInsightRulesInput copyWith(
          void Function(EnableInsightRulesInput) updates) =>
      super.copyWith((message) => updates(message as EnableInsightRulesInput))
          as EnableInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableInsightRulesInput create() => EnableInsightRulesInput._();
  @$core.override
  EnableInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableInsightRulesInput>(create);
  static EnableInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(267949170)
  $pb.PbList<$core.String> get rulenames => $_getList(0);
}

class EnableInsightRulesOutput extends $pb.GeneratedMessage {
  factory EnableInsightRulesOutput({
    $core.Iterable<PartialFailure>? failures,
  }) {
    final result = create();
    if (failures != null) result.failures.addAll(failures);
    return result;
  }

  EnableInsightRulesOutput._();

  factory EnableInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EnableInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EnableInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<PartialFailure>(335467271, _omitFieldNames ? '' : 'failures',
        subBuilder: PartialFailure.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EnableInsightRulesOutput copyWith(
          void Function(EnableInsightRulesOutput) updates) =>
      super.copyWith((message) => updates(message as EnableInsightRulesOutput))
          as EnableInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EnableInsightRulesOutput create() => EnableInsightRulesOutput._();
  @$core.override
  EnableInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EnableInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EnableInsightRulesOutput>(create);
  static EnableInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(335467271)
  $pb.PbList<PartialFailure> get failures => $_getList(0);
}

class Entity extends $pb.GeneratedMessage {
  factory Entity({
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? attributes,
    $core.Iterable<$core.MapEntry<$core.String, $core.String>>? keyattributes,
  }) {
    final result = create();
    if (attributes != null) result.attributes.addEntries(attributes);
    if (keyattributes != null) result.keyattributes.addEntries(keyattributes);
    return result;
  }

  Entity._();

  factory Entity.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Entity.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Entity',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..m<$core.String, $core.String>(
        209638581, _omitFieldNames ? '' : 'attributes',
        entryClassName: 'Entity.AttributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('monitoring'))
    ..m<$core.String, $core.String>(
        462259374, _omitFieldNames ? '' : 'keyattributes',
        entryClassName: 'Entity.KeyattributesEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.OS,
        packageName: const $pb.PackageName('monitoring'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Entity clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Entity copyWith(void Function(Entity) updates) =>
      super.copyWith((message) => updates(message as Entity)) as Entity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Entity create() => Entity._();
  @$core.override
  Entity createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Entity getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Entity>(create);
  static Entity? _defaultInstance;

  @$pb.TagNumber(209638581)
  $pb.PbMap<$core.String, $core.String> get attributes => $_getMap(0);

  @$pb.TagNumber(462259374)
  $pb.PbMap<$core.String, $core.String> get keyattributes => $_getMap(1);
}

class EntityMetricData extends $pb.GeneratedMessage {
  factory EntityMetricData({
    Entity? entity,
    $core.Iterable<MetricDatum>? metricdata,
  }) {
    final result = create();
    if (entity != null) result.entity = entity;
    if (metricdata != null) result.metricdata.addAll(metricdata);
    return result;
  }

  EntityMetricData._();

  factory EntityMetricData.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory EntityMetricData.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'EntityMetricData',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOM<Entity>(10171131, _omitFieldNames ? '' : 'entity',
        subBuilder: Entity.create)
    ..pPM<MetricDatum>(162960562, _omitFieldNames ? '' : 'metricdata',
        subBuilder: MetricDatum.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EntityMetricData clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  EntityMetricData copyWith(void Function(EntityMetricData) updates) =>
      super.copyWith((message) => updates(message as EntityMetricData))
          as EntityMetricData;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static EntityMetricData create() => EntityMetricData._();
  @$core.override
  EntityMetricData createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static EntityMetricData getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<EntityMetricData>(create);
  static EntityMetricData? _defaultInstance;

  @$pb.TagNumber(10171131)
  Entity get entity => $_getN(0);
  @$pb.TagNumber(10171131)
  set entity(Entity value) => $_setField(10171131, value);
  @$pb.TagNumber(10171131)
  $core.bool hasEntity() => $_has(0);
  @$pb.TagNumber(10171131)
  void clearEntity() => $_clearField(10171131);
  @$pb.TagNumber(10171131)
  Entity ensureEntity() => $_ensure(0);

  @$pb.TagNumber(162960562)
  $pb.PbList<MetricDatum> get metricdata => $_getList(1);
}

class GetAlarmMuteRuleInput extends $pb.GeneratedMessage {
  factory GetAlarmMuteRuleInput({
    $core.String? alarmmuterulename,
  }) {
    final result = create();
    if (alarmmuterulename != null) result.alarmmuterulename = alarmmuterulename;
    return result;
  }

  GetAlarmMuteRuleInput._();

  factory GetAlarmMuteRuleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAlarmMuteRuleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAlarmMuteRuleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(530877511, _omitFieldNames ? '' : 'alarmmuterulename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAlarmMuteRuleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAlarmMuteRuleInput copyWith(
          void Function(GetAlarmMuteRuleInput) updates) =>
      super.copyWith((message) => updates(message as GetAlarmMuteRuleInput))
          as GetAlarmMuteRuleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAlarmMuteRuleInput create() => GetAlarmMuteRuleInput._();
  @$core.override
  GetAlarmMuteRuleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAlarmMuteRuleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAlarmMuteRuleInput>(create);
  static GetAlarmMuteRuleInput? _defaultInstance;

  @$pb.TagNumber(530877511)
  $core.String get alarmmuterulename => $_getSZ(0);
  @$pb.TagNumber(530877511)
  set alarmmuterulename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(530877511)
  $core.bool hasAlarmmuterulename() => $_has(0);
  @$pb.TagNumber(530877511)
  void clearAlarmmuterulename() => $_clearField(530877511);
}

class GetAlarmMuteRuleOutput extends $pb.GeneratedMessage {
  factory GetAlarmMuteRuleOutput({
    AlarmMuteRuleStatus? status,
    $core.String? mutetype,
    $core.String? description,
    $core.String? lastupdatedtimestamp,
    $core.String? name,
    MuteTargets? mutetargets,
    $core.String? startdate,
    $core.String? expiredate,
    Rule? rule,
    $core.String? alarmmuterulearn,
  }) {
    final result = create();
    if (status != null) result.status = status;
    if (mutetype != null) result.mutetype = mutetype;
    if (description != null) result.description = description;
    if (lastupdatedtimestamp != null)
      result.lastupdatedtimestamp = lastupdatedtimestamp;
    if (name != null) result.name = name;
    if (mutetargets != null) result.mutetargets = mutetargets;
    if (startdate != null) result.startdate = startdate;
    if (expiredate != null) result.expiredate = expiredate;
    if (rule != null) result.rule = rule;
    if (alarmmuterulearn != null) result.alarmmuterulearn = alarmmuterulearn;
    return result;
  }

  GetAlarmMuteRuleOutput._();

  factory GetAlarmMuteRuleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAlarmMuteRuleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAlarmMuteRuleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aE<AlarmMuteRuleStatus>(6222352, _omitFieldNames ? '' : 'status',
        enumValues: AlarmMuteRuleStatus.values)
    ..aOS(111213245, _omitFieldNames ? '' : 'mutetype')
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(133309845, _omitFieldNames ? '' : 'lastupdatedtimestamp')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<MuteTargets>(356717131, _omitFieldNames ? '' : 'mutetargets',
        subBuilder: MuteTargets.create)
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..aOS(454143207, _omitFieldNames ? '' : 'expiredate')
    ..aOM<Rule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: Rule.create)
    ..aOS(493590237, _omitFieldNames ? '' : 'alarmmuterulearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAlarmMuteRuleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAlarmMuteRuleOutput copyWith(
          void Function(GetAlarmMuteRuleOutput) updates) =>
      super.copyWith((message) => updates(message as GetAlarmMuteRuleOutput))
          as GetAlarmMuteRuleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAlarmMuteRuleOutput create() => GetAlarmMuteRuleOutput._();
  @$core.override
  GetAlarmMuteRuleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAlarmMuteRuleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAlarmMuteRuleOutput>(create);
  static GetAlarmMuteRuleOutput? _defaultInstance;

  @$pb.TagNumber(6222352)
  AlarmMuteRuleStatus get status => $_getN(0);
  @$pb.TagNumber(6222352)
  set status(AlarmMuteRuleStatus value) => $_setField(6222352, value);
  @$pb.TagNumber(6222352)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(6222352)
  void clearStatus() => $_clearField(6222352);

  @$pb.TagNumber(111213245)
  $core.String get mutetype => $_getSZ(1);
  @$pb.TagNumber(111213245)
  set mutetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(111213245)
  $core.bool hasMutetype() => $_has(1);
  @$pb.TagNumber(111213245)
  void clearMutetype() => $_clearField(111213245);

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(2, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(2);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(133309845)
  $core.String get lastupdatedtimestamp => $_getSZ(3);
  @$pb.TagNumber(133309845)
  set lastupdatedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(133309845)
  $core.bool hasLastupdatedtimestamp() => $_has(3);
  @$pb.TagNumber(133309845)
  void clearLastupdatedtimestamp() => $_clearField(133309845);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(356717131)
  MuteTargets get mutetargets => $_getN(5);
  @$pb.TagNumber(356717131)
  set mutetargets(MuteTargets value) => $_setField(356717131, value);
  @$pb.TagNumber(356717131)
  $core.bool hasMutetargets() => $_has(5);
  @$pb.TagNumber(356717131)
  void clearMutetargets() => $_clearField(356717131);
  @$pb.TagNumber(356717131)
  MuteTargets ensureMutetargets() => $_ensure(5);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(6);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(6, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(6);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(454143207)
  $core.String get expiredate => $_getSZ(7);
  @$pb.TagNumber(454143207)
  set expiredate($core.String value) => $_setString(7, value);
  @$pb.TagNumber(454143207)
  $core.bool hasExpiredate() => $_has(7);
  @$pb.TagNumber(454143207)
  void clearExpiredate() => $_clearField(454143207);

  @$pb.TagNumber(475696372)
  Rule get rule => $_getN(8);
  @$pb.TagNumber(475696372)
  set rule(Rule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(8);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  Rule ensureRule() => $_ensure(8);

  @$pb.TagNumber(493590237)
  $core.String get alarmmuterulearn => $_getSZ(9);
  @$pb.TagNumber(493590237)
  set alarmmuterulearn($core.String value) => $_setString(9, value);
  @$pb.TagNumber(493590237)
  $core.bool hasAlarmmuterulearn() => $_has(9);
  @$pb.TagNumber(493590237)
  void clearAlarmmuterulearn() => $_clearField(493590237);
}

class GetDashboardInput extends $pb.GeneratedMessage {
  factory GetDashboardInput({
    $core.String? dashboardname,
  }) {
    final result = create();
    if (dashboardname != null) result.dashboardname = dashboardname;
    return result;
  }

  GetDashboardInput._();

  factory GetDashboardInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDashboardInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDashboardInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(506599873, _omitFieldNames ? '' : 'dashboardname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardInput copyWith(void Function(GetDashboardInput) updates) =>
      super.copyWith((message) => updates(message as GetDashboardInput))
          as GetDashboardInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDashboardInput create() => GetDashboardInput._();
  @$core.override
  GetDashboardInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDashboardInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDashboardInput>(create);
  static GetDashboardInput? _defaultInstance;

  @$pb.TagNumber(506599873)
  $core.String get dashboardname => $_getSZ(0);
  @$pb.TagNumber(506599873)
  set dashboardname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(506599873)
  $core.bool hasDashboardname() => $_has(0);
  @$pb.TagNumber(506599873)
  void clearDashboardname() => $_clearField(506599873);
}

class GetDashboardOutput extends $pb.GeneratedMessage {
  factory GetDashboardOutput({
    $core.String? dashboardbody,
    $core.String? dashboardarn,
    $core.String? dashboardname,
  }) {
    final result = create();
    if (dashboardbody != null) result.dashboardbody = dashboardbody;
    if (dashboardarn != null) result.dashboardarn = dashboardarn;
    if (dashboardname != null) result.dashboardname = dashboardname;
    return result;
  }

  GetDashboardOutput._();

  factory GetDashboardOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetDashboardOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetDashboardOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(3403236, _omitFieldNames ? '' : 'dashboardbody')
    ..aOS(108051951, _omitFieldNames ? '' : 'dashboardarn')
    ..aOS(506599873, _omitFieldNames ? '' : 'dashboardname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetDashboardOutput copyWith(void Function(GetDashboardOutput) updates) =>
      super.copyWith((message) => updates(message as GetDashboardOutput))
          as GetDashboardOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetDashboardOutput create() => GetDashboardOutput._();
  @$core.override
  GetDashboardOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetDashboardOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetDashboardOutput>(create);
  static GetDashboardOutput? _defaultInstance;

  @$pb.TagNumber(3403236)
  $core.String get dashboardbody => $_getSZ(0);
  @$pb.TagNumber(3403236)
  set dashboardbody($core.String value) => $_setString(0, value);
  @$pb.TagNumber(3403236)
  $core.bool hasDashboardbody() => $_has(0);
  @$pb.TagNumber(3403236)
  void clearDashboardbody() => $_clearField(3403236);

  @$pb.TagNumber(108051951)
  $core.String get dashboardarn => $_getSZ(1);
  @$pb.TagNumber(108051951)
  set dashboardarn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(108051951)
  $core.bool hasDashboardarn() => $_has(1);
  @$pb.TagNumber(108051951)
  void clearDashboardarn() => $_clearField(108051951);

  @$pb.TagNumber(506599873)
  $core.String get dashboardname => $_getSZ(2);
  @$pb.TagNumber(506599873)
  set dashboardname($core.String value) => $_setString(2, value);
  @$pb.TagNumber(506599873)
  $core.bool hasDashboardname() => $_has(2);
  @$pb.TagNumber(506599873)
  void clearDashboardname() => $_clearField(506599873);
}

class GetInsightRuleReportInput extends $pb.GeneratedMessage {
  factory GetInsightRuleReportInput({
    $core.String? endtime,
    $core.int? period,
    $core.String? rulename,
    $core.int? maxcontributorcount,
    $core.String? orderby,
    $core.String? starttime,
    $core.Iterable<$core.String>? metrics,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (period != null) result.period = period;
    if (rulename != null) result.rulename = rulename;
    if (maxcontributorcount != null)
      result.maxcontributorcount = maxcontributorcount;
    if (orderby != null) result.orderby = orderby;
    if (starttime != null) result.starttime = starttime;
    if (metrics != null) result.metrics.addAll(metrics);
    return result;
  }

  GetInsightRuleReportInput._();

  factory GetInsightRuleReportInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetInsightRuleReportInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetInsightRuleReportInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aOS(214688793, _omitFieldNames ? '' : 'rulename')
    ..aI(272640614, _omitFieldNames ? '' : 'maxcontributorcount')
    ..aOS(353010019, _omitFieldNames ? '' : 'orderby')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..pPS(436365847, _omitFieldNames ? '' : 'metrics')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightRuleReportInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightRuleReportInput copyWith(
          void Function(GetInsightRuleReportInput) updates) =>
      super.copyWith((message) => updates(message as GetInsightRuleReportInput))
          as GetInsightRuleReportInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetInsightRuleReportInput create() => GetInsightRuleReportInput._();
  @$core.override
  GetInsightRuleReportInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetInsightRuleReportInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetInsightRuleReportInput>(create);
  static GetInsightRuleReportInput? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(1);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(1);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(214688793)
  $core.String get rulename => $_getSZ(2);
  @$pb.TagNumber(214688793)
  set rulename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(214688793)
  $core.bool hasRulename() => $_has(2);
  @$pb.TagNumber(214688793)
  void clearRulename() => $_clearField(214688793);

  @$pb.TagNumber(272640614)
  $core.int get maxcontributorcount => $_getIZ(3);
  @$pb.TagNumber(272640614)
  set maxcontributorcount($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(272640614)
  $core.bool hasMaxcontributorcount() => $_has(3);
  @$pb.TagNumber(272640614)
  void clearMaxcontributorcount() => $_clearField(272640614);

  @$pb.TagNumber(353010019)
  $core.String get orderby => $_getSZ(4);
  @$pb.TagNumber(353010019)
  set orderby($core.String value) => $_setString(4, value);
  @$pb.TagNumber(353010019)
  $core.bool hasOrderby() => $_has(4);
  @$pb.TagNumber(353010019)
  void clearOrderby() => $_clearField(353010019);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(5);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(5, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(5);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(436365847)
  $pb.PbList<$core.String> get metrics => $_getList(6);
}

class GetInsightRuleReportOutput extends $pb.GeneratedMessage {
  factory GetInsightRuleReportOutput({
    $core.Iterable<$core.String>? keylabels,
    $fixnum.Int64? approximateuniquecount,
    $core.Iterable<InsightRuleMetricDatapoint>? metricdatapoints,
    $core.Iterable<InsightRuleContributor>? contributors,
    $core.double? aggregatevalue,
    $core.String? aggregationstatistic,
  }) {
    final result = create();
    if (keylabels != null) result.keylabels.addAll(keylabels);
    if (approximateuniquecount != null)
      result.approximateuniquecount = approximateuniquecount;
    if (metricdatapoints != null)
      result.metricdatapoints.addAll(metricdatapoints);
    if (contributors != null) result.contributors.addAll(contributors);
    if (aggregatevalue != null) result.aggregatevalue = aggregatevalue;
    if (aggregationstatistic != null)
      result.aggregationstatistic = aggregationstatistic;
    return result;
  }

  GetInsightRuleReportOutput._();

  factory GetInsightRuleReportOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetInsightRuleReportOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetInsightRuleReportOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(16074444, _omitFieldNames ? '' : 'keylabels')
    ..aInt64(55617722, _omitFieldNames ? '' : 'approximateuniquecount')
    ..pPM<InsightRuleMetricDatapoint>(
        117991067, _omitFieldNames ? '' : 'metricdatapoints',
        subBuilder: InsightRuleMetricDatapoint.create)
    ..pPM<InsightRuleContributor>(
        168168444, _omitFieldNames ? '' : 'contributors',
        subBuilder: InsightRuleContributor.create)
    ..aD(225313242, _omitFieldNames ? '' : 'aggregatevalue')
    ..aOS(333233758, _omitFieldNames ? '' : 'aggregationstatistic')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightRuleReportOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetInsightRuleReportOutput copyWith(
          void Function(GetInsightRuleReportOutput) updates) =>
      super.copyWith(
              (message) => updates(message as GetInsightRuleReportOutput))
          as GetInsightRuleReportOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetInsightRuleReportOutput create() => GetInsightRuleReportOutput._();
  @$core.override
  GetInsightRuleReportOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetInsightRuleReportOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetInsightRuleReportOutput>(create);
  static GetInsightRuleReportOutput? _defaultInstance;

  @$pb.TagNumber(16074444)
  $pb.PbList<$core.String> get keylabels => $_getList(0);

  @$pb.TagNumber(55617722)
  $fixnum.Int64 get approximateuniquecount => $_getI64(1);
  @$pb.TagNumber(55617722)
  set approximateuniquecount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(55617722)
  $core.bool hasApproximateuniquecount() => $_has(1);
  @$pb.TagNumber(55617722)
  void clearApproximateuniquecount() => $_clearField(55617722);

  @$pb.TagNumber(117991067)
  $pb.PbList<InsightRuleMetricDatapoint> get metricdatapoints => $_getList(2);

  @$pb.TagNumber(168168444)
  $pb.PbList<InsightRuleContributor> get contributors => $_getList(3);

  @$pb.TagNumber(225313242)
  $core.double get aggregatevalue => $_getN(4);
  @$pb.TagNumber(225313242)
  set aggregatevalue($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(225313242)
  $core.bool hasAggregatevalue() => $_has(4);
  @$pb.TagNumber(225313242)
  void clearAggregatevalue() => $_clearField(225313242);

  @$pb.TagNumber(333233758)
  $core.String get aggregationstatistic => $_getSZ(5);
  @$pb.TagNumber(333233758)
  set aggregationstatistic($core.String value) => $_setString(5, value);
  @$pb.TagNumber(333233758)
  $core.bool hasAggregationstatistic() => $_has(5);
  @$pb.TagNumber(333233758)
  void clearAggregationstatistic() => $_clearField(333233758);
}

class GetMetricDataInput extends $pb.GeneratedMessage {
  factory GetMetricDataInput({
    $core.String? endtime,
    $core.Iterable<MetricDataQuery>? metricdataqueries,
    $core.int? maxdatapoints,
    LabelOptions? labeloptions,
    $core.String? nexttoken,
    ScanBy? scanby,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (metricdataqueries != null)
      result.metricdataqueries.addAll(metricdataqueries);
    if (maxdatapoints != null) result.maxdatapoints = maxdatapoints;
    if (labeloptions != null) result.labeloptions = labeloptions;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (scanby != null) result.scanby = scanby;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  GetMetricDataInput._();

  factory GetMetricDataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricDataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricDataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..pPM<MetricDataQuery>(92209504, _omitFieldNames ? '' : 'metricdataqueries',
        subBuilder: MetricDataQuery.create)
    ..aI(97504835, _omitFieldNames ? '' : 'maxdatapoints')
    ..aOM<LabelOptions>(184949902, _omitFieldNames ? '' : 'labeloptions',
        subBuilder: LabelOptions.create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aE<ScanBy>(344142854, _omitFieldNames ? '' : 'scanby',
        enumValues: ScanBy.values)
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricDataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricDataInput copyWith(void Function(GetMetricDataInput) updates) =>
      super.copyWith((message) => updates(message as GetMetricDataInput))
          as GetMetricDataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricDataInput create() => GetMetricDataInput._();
  @$core.override
  GetMetricDataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricDataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricDataInput>(create);
  static GetMetricDataInput? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(92209504)
  $pb.PbList<MetricDataQuery> get metricdataqueries => $_getList(1);

  @$pb.TagNumber(97504835)
  $core.int get maxdatapoints => $_getIZ(2);
  @$pb.TagNumber(97504835)
  set maxdatapoints($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(97504835)
  $core.bool hasMaxdatapoints() => $_has(2);
  @$pb.TagNumber(97504835)
  void clearMaxdatapoints() => $_clearField(97504835);

  @$pb.TagNumber(184949902)
  LabelOptions get labeloptions => $_getN(3);
  @$pb.TagNumber(184949902)
  set labeloptions(LabelOptions value) => $_setField(184949902, value);
  @$pb.TagNumber(184949902)
  $core.bool hasLabeloptions() => $_has(3);
  @$pb.TagNumber(184949902)
  void clearLabeloptions() => $_clearField(184949902);
  @$pb.TagNumber(184949902)
  LabelOptions ensureLabeloptions() => $_ensure(3);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(4);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(4, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(4);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(344142854)
  ScanBy get scanby => $_getN(5);
  @$pb.TagNumber(344142854)
  set scanby(ScanBy value) => $_setField(344142854, value);
  @$pb.TagNumber(344142854)
  $core.bool hasScanby() => $_has(5);
  @$pb.TagNumber(344142854)
  void clearScanby() => $_clearField(344142854);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(6);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(6);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
}

class GetMetricDataOutput extends $pb.GeneratedMessage {
  factory GetMetricDataOutput({
    $core.String? nexttoken,
    $core.Iterable<MetricDataResult>? metricdataresults,
    $core.Iterable<MessageData>? messages,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (metricdataresults != null)
      result.metricdataresults.addAll(metricdataresults);
    if (messages != null) result.messages.addAll(messages);
    return result;
  }

  GetMetricDataOutput._();

  factory GetMetricDataOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricDataOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricDataOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<MetricDataResult>(
        373784842, _omitFieldNames ? '' : 'metricdataresults',
        subBuilder: MetricDataResult.create)
    ..pPM<MessageData>(409018326, _omitFieldNames ? '' : 'messages',
        subBuilder: MessageData.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricDataOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricDataOutput copyWith(void Function(GetMetricDataOutput) updates) =>
      super.copyWith((message) => updates(message as GetMetricDataOutput))
          as GetMetricDataOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricDataOutput create() => GetMetricDataOutput._();
  @$core.override
  GetMetricDataOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricDataOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricDataOutput>(create);
  static GetMetricDataOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(373784842)
  $pb.PbList<MetricDataResult> get metricdataresults => $_getList(1);

  @$pb.TagNumber(409018326)
  $pb.PbList<MessageData> get messages => $_getList(2);
}

class GetMetricStatisticsInput extends $pb.GeneratedMessage {
  factory GetMetricStatisticsInput({
    $core.String? endtime,
    $core.String? metricname,
    $core.int? period,
    StandardUnit? unit,
    $core.String? namespace,
    $core.Iterable<$core.String>? extendedstatistics,
    $core.String? starttime,
    $core.Iterable<Dimension>? dimensions,
    $core.Iterable<Statistic>? statistics,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (metricname != null) result.metricname = metricname;
    if (period != null) result.period = period;
    if (unit != null) result.unit = unit;
    if (namespace != null) result.namespace = namespace;
    if (extendedstatistics != null)
      result.extendedstatistics.addAll(extendedstatistics);
    if (starttime != null) result.starttime = starttime;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    if (statistics != null) result.statistics.addAll(statistics);
    return result;
  }

  GetMetricStatisticsInput._();

  factory GetMetricStatisticsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricStatisticsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricStatisticsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPS(365818700, _omitFieldNames ? '' : 'extendedstatistics')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..pc<Statistic>(
        510636075, _omitFieldNames ? '' : 'statistics', $pb.PbFieldType.KE,
        valueOf: Statistic.valueOf,
        enumValues: Statistic.values,
        defaultEnumValue: Statistic.STATISTIC_SUM)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStatisticsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStatisticsInput copyWith(
          void Function(GetMetricStatisticsInput) updates) =>
      super.copyWith((message) => updates(message as GetMetricStatisticsInput))
          as GetMetricStatisticsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricStatisticsInput create() => GetMetricStatisticsInput._();
  @$core.override
  GetMetricStatisticsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricStatisticsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricStatisticsInput>(create);
  static GetMetricStatisticsInput? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(2);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(2);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(3);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(3);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(4);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(4, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(4);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(365818700)
  $pb.PbList<$core.String> get extendedstatistics => $_getList(5);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(6);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(6, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(6);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(7);

  @$pb.TagNumber(510636075)
  $pb.PbList<Statistic> get statistics => $_getList(8);
}

class GetMetricStatisticsOutput extends $pb.GeneratedMessage {
  factory GetMetricStatisticsOutput({
    $core.Iterable<Datapoint>? datapoints,
    $core.String? label,
  }) {
    final result = create();
    if (datapoints != null) result.datapoints.addAll(datapoints);
    if (label != null) result.label = label;
    return result;
  }

  GetMetricStatisticsOutput._();

  factory GetMetricStatisticsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricStatisticsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricStatisticsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<Datapoint>(37512135, _omitFieldNames ? '' : 'datapoints',
        subBuilder: Datapoint.create)
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStatisticsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStatisticsOutput copyWith(
          void Function(GetMetricStatisticsOutput) updates) =>
      super.copyWith((message) => updates(message as GetMetricStatisticsOutput))
          as GetMetricStatisticsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricStatisticsOutput create() => GetMetricStatisticsOutput._();
  @$core.override
  GetMetricStatisticsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricStatisticsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricStatisticsOutput>(create);
  static GetMetricStatisticsOutput? _defaultInstance;

  @$pb.TagNumber(37512135)
  $pb.PbList<Datapoint> get datapoints => $_getList(0);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(1);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(1, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(1);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class GetMetricStreamInput extends $pb.GeneratedMessage {
  factory GetMetricStreamInput({
    $core.String? name,
  }) {
    final result = create();
    if (name != null) result.name = name;
    return result;
  }

  GetMetricStreamInput._();

  factory GetMetricStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStreamInput copyWith(void Function(GetMetricStreamInput) updates) =>
      super.copyWith((message) => updates(message as GetMetricStreamInput))
          as GetMetricStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricStreamInput create() => GetMetricStreamInput._();
  @$core.override
  GetMetricStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricStreamInput>(create);
  static GetMetricStreamInput? _defaultInstance;

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);
}

class GetMetricStreamOutput extends $pb.GeneratedMessage {
  factory GetMetricStreamOutput({
    $core.String? firehosearn,
    $core.String? lastupdatedate,
    $core.Iterable<MetricStreamFilter>? includefilters,
    $core.Iterable<MetricStreamStatisticsConfiguration>?
        statisticsconfigurations,
    $core.String? name,
    $core.String? creationdate,
    $core.String? rolearn,
    $core.bool? includelinkedaccountsmetrics,
    $core.String? arn,
    $core.Iterable<MetricStreamFilter>? excludefilters,
    $core.String? state,
    MetricStreamOutputFormat? outputformat,
  }) {
    final result = create();
    if (firehosearn != null) result.firehosearn = firehosearn;
    if (lastupdatedate != null) result.lastupdatedate = lastupdatedate;
    if (includefilters != null) result.includefilters.addAll(includefilters);
    if (statisticsconfigurations != null)
      result.statisticsconfigurations.addAll(statisticsconfigurations);
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (rolearn != null) result.rolearn = rolearn;
    if (includelinkedaccountsmetrics != null)
      result.includelinkedaccountsmetrics = includelinkedaccountsmetrics;
    if (arn != null) result.arn = arn;
    if (excludefilters != null) result.excludefilters.addAll(excludefilters);
    if (state != null) result.state = state;
    if (outputformat != null) result.outputformat = outputformat;
    return result;
  }

  GetMetricStreamOutput._();

  factory GetMetricStreamOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricStreamOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricStreamOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(55494200, _omitFieldNames ? '' : 'firehosearn')
    ..aOS(65755725, _omitFieldNames ? '' : 'lastupdatedate')
    ..pPM<MetricStreamFilter>(90967031, _omitFieldNames ? '' : 'includefilters',
        subBuilder: MetricStreamFilter.create)
    ..pPM<MetricStreamStatisticsConfiguration>(
        175876178, _omitFieldNames ? '' : 'statisticsconfigurations',
        subBuilder: MetricStreamStatisticsConfiguration.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOB(356901952, _omitFieldNames ? '' : 'includelinkedaccountsmetrics')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..pPM<MetricStreamFilter>(
        426990841, _omitFieldNames ? '' : 'excludefilters',
        subBuilder: MetricStreamFilter.create)
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..aE<MetricStreamOutputFormat>(
        519139704, _omitFieldNames ? '' : 'outputformat',
        enumValues: MetricStreamOutputFormat.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStreamOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricStreamOutput copyWith(
          void Function(GetMetricStreamOutput) updates) =>
      super.copyWith((message) => updates(message as GetMetricStreamOutput))
          as GetMetricStreamOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricStreamOutput create() => GetMetricStreamOutput._();
  @$core.override
  GetMetricStreamOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricStreamOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricStreamOutput>(create);
  static GetMetricStreamOutput? _defaultInstance;

  @$pb.TagNumber(55494200)
  $core.String get firehosearn => $_getSZ(0);
  @$pb.TagNumber(55494200)
  set firehosearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(55494200)
  $core.bool hasFirehosearn() => $_has(0);
  @$pb.TagNumber(55494200)
  void clearFirehosearn() => $_clearField(55494200);

  @$pb.TagNumber(65755725)
  $core.String get lastupdatedate => $_getSZ(1);
  @$pb.TagNumber(65755725)
  set lastupdatedate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(65755725)
  $core.bool hasLastupdatedate() => $_has(1);
  @$pb.TagNumber(65755725)
  void clearLastupdatedate() => $_clearField(65755725);

  @$pb.TagNumber(90967031)
  $pb.PbList<MetricStreamFilter> get includefilters => $_getList(2);

  @$pb.TagNumber(175876178)
  $pb.PbList<MetricStreamStatisticsConfiguration>
      get statisticsconfigurations => $_getList(3);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(4);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(4, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(4);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(5);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(5);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(6);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(6, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(6);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(356901952)
  $core.bool get includelinkedaccountsmetrics => $_getBF(7);
  @$pb.TagNumber(356901952)
  set includelinkedaccountsmetrics($core.bool value) => $_setBool(7, value);
  @$pb.TagNumber(356901952)
  $core.bool hasIncludelinkedaccountsmetrics() => $_has(7);
  @$pb.TagNumber(356901952)
  void clearIncludelinkedaccountsmetrics() => $_clearField(356901952);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(8);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(8, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(8);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(426990841)
  $pb.PbList<MetricStreamFilter> get excludefilters => $_getList(9);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(10);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(10, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(10);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(519139704)
  MetricStreamOutputFormat get outputformat => $_getN(11);
  @$pb.TagNumber(519139704)
  set outputformat(MetricStreamOutputFormat value) =>
      $_setField(519139704, value);
  @$pb.TagNumber(519139704)
  $core.bool hasOutputformat() => $_has(11);
  @$pb.TagNumber(519139704)
  void clearOutputformat() => $_clearField(519139704);
}

class GetMetricWidgetImageInput extends $pb.GeneratedMessage {
  factory GetMetricWidgetImageInput({
    $core.String? metricwidget,
    $core.String? outputformat,
  }) {
    final result = create();
    if (metricwidget != null) result.metricwidget = metricwidget;
    if (outputformat != null) result.outputformat = outputformat;
    return result;
  }

  GetMetricWidgetImageInput._();

  factory GetMetricWidgetImageInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricWidgetImageInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricWidgetImageInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(283026730, _omitFieldNames ? '' : 'metricwidget')
    ..aOS(519139704, _omitFieldNames ? '' : 'outputformat')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricWidgetImageInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricWidgetImageInput copyWith(
          void Function(GetMetricWidgetImageInput) updates) =>
      super.copyWith((message) => updates(message as GetMetricWidgetImageInput))
          as GetMetricWidgetImageInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricWidgetImageInput create() => GetMetricWidgetImageInput._();
  @$core.override
  GetMetricWidgetImageInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricWidgetImageInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricWidgetImageInput>(create);
  static GetMetricWidgetImageInput? _defaultInstance;

  @$pb.TagNumber(283026730)
  $core.String get metricwidget => $_getSZ(0);
  @$pb.TagNumber(283026730)
  set metricwidget($core.String value) => $_setString(0, value);
  @$pb.TagNumber(283026730)
  $core.bool hasMetricwidget() => $_has(0);
  @$pb.TagNumber(283026730)
  void clearMetricwidget() => $_clearField(283026730);

  @$pb.TagNumber(519139704)
  $core.String get outputformat => $_getSZ(1);
  @$pb.TagNumber(519139704)
  set outputformat($core.String value) => $_setString(1, value);
  @$pb.TagNumber(519139704)
  $core.bool hasOutputformat() => $_has(1);
  @$pb.TagNumber(519139704)
  void clearOutputformat() => $_clearField(519139704);
}

class GetMetricWidgetImageOutput extends $pb.GeneratedMessage {
  factory GetMetricWidgetImageOutput({
    $core.List<$core.int>? metricwidgetimage,
  }) {
    final result = create();
    if (metricwidgetimage != null) result.metricwidgetimage = metricwidgetimage;
    return result;
  }

  GetMetricWidgetImageOutput._();

  factory GetMetricWidgetImageOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetMetricWidgetImageOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetMetricWidgetImageOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..a<$core.List<$core.int>>(484815311,
        _omitFieldNames ? '' : 'metricwidgetimage', $pb.PbFieldType.OY)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricWidgetImageOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetMetricWidgetImageOutput copyWith(
          void Function(GetMetricWidgetImageOutput) updates) =>
      super.copyWith(
              (message) => updates(message as GetMetricWidgetImageOutput))
          as GetMetricWidgetImageOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetMetricWidgetImageOutput create() => GetMetricWidgetImageOutput._();
  @$core.override
  GetMetricWidgetImageOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetMetricWidgetImageOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetMetricWidgetImageOutput>(create);
  static GetMetricWidgetImageOutput? _defaultInstance;

  @$pb.TagNumber(484815311)
  $core.List<$core.int> get metricwidgetimage => $_getN(0);
  @$pb.TagNumber(484815311)
  set metricwidgetimage($core.List<$core.int> value) => $_setBytes(0, value);
  @$pb.TagNumber(484815311)
  $core.bool hasMetricwidgetimage() => $_has(0);
  @$pb.TagNumber(484815311)
  void clearMetricwidgetimage() => $_clearField(484815311);
}

class InsightRule extends $pb.GeneratedMessage {
  factory InsightRule({
    $core.bool? applyontransformedlogs,
    $core.bool? managedrule,
    $core.String? definition,
    $core.String? name,
    $core.String? schema,
    $core.String? state,
  }) {
    final result = create();
    if (applyontransformedlogs != null)
      result.applyontransformedlogs = applyontransformedlogs;
    if (managedrule != null) result.managedrule = managedrule;
    if (definition != null) result.definition = definition;
    if (name != null) result.name = name;
    if (schema != null) result.schema = schema;
    if (state != null) result.state = state;
    return result;
  }

  InsightRule._();

  factory InsightRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightRule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOB(30966533, _omitFieldNames ? '' : 'applyontransformedlogs')
    ..aOB(144612369, _omitFieldNames ? '' : 'managedrule')
    ..aOS(222998209, _omitFieldNames ? '' : 'definition')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(412122455, _omitFieldNames ? '' : 'schema')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRule copyWith(void Function(InsightRule) updates) =>
      super.copyWith((message) => updates(message as InsightRule))
          as InsightRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightRule create() => InsightRule._();
  @$core.override
  InsightRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightRule>(create);
  static InsightRule? _defaultInstance;

  @$pb.TagNumber(30966533)
  $core.bool get applyontransformedlogs => $_getBF(0);
  @$pb.TagNumber(30966533)
  set applyontransformedlogs($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(30966533)
  $core.bool hasApplyontransformedlogs() => $_has(0);
  @$pb.TagNumber(30966533)
  void clearApplyontransformedlogs() => $_clearField(30966533);

  @$pb.TagNumber(144612369)
  $core.bool get managedrule => $_getBF(1);
  @$pb.TagNumber(144612369)
  set managedrule($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(144612369)
  $core.bool hasManagedrule() => $_has(1);
  @$pb.TagNumber(144612369)
  void clearManagedrule() => $_clearField(144612369);

  @$pb.TagNumber(222998209)
  $core.String get definition => $_getSZ(2);
  @$pb.TagNumber(222998209)
  set definition($core.String value) => $_setString(2, value);
  @$pb.TagNumber(222998209)
  $core.bool hasDefinition() => $_has(2);
  @$pb.TagNumber(222998209)
  void clearDefinition() => $_clearField(222998209);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(412122455)
  $core.String get schema => $_getSZ(4);
  @$pb.TagNumber(412122455)
  set schema($core.String value) => $_setString(4, value);
  @$pb.TagNumber(412122455)
  $core.bool hasSchema() => $_has(4);
  @$pb.TagNumber(412122455)
  void clearSchema() => $_clearField(412122455);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(5);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(5, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class InsightRuleContributor extends $pb.GeneratedMessage {
  factory InsightRuleContributor({
    $core.Iterable<$core.String>? keys,
    $core.Iterable<InsightRuleContributorDatapoint>? datapoints,
    $core.double? approximateaggregatevalue,
  }) {
    final result = create();
    if (keys != null) result.keys.addAll(keys);
    if (datapoints != null) result.datapoints.addAll(datapoints);
    if (approximateaggregatevalue != null)
      result.approximateaggregatevalue = approximateaggregatevalue;
    return result;
  }

  InsightRuleContributor._();

  factory InsightRuleContributor.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightRuleContributor.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightRuleContributor',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(2831086, _omitFieldNames ? '' : 'keys')
    ..pPM<InsightRuleContributorDatapoint>(
        37512135, _omitFieldNames ? '' : 'datapoints',
        subBuilder: InsightRuleContributorDatapoint.create)
    ..aD(260431894, _omitFieldNames ? '' : 'approximateaggregatevalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleContributor clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleContributor copyWith(
          void Function(InsightRuleContributor) updates) =>
      super.copyWith((message) => updates(message as InsightRuleContributor))
          as InsightRuleContributor;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightRuleContributor create() => InsightRuleContributor._();
  @$core.override
  InsightRuleContributor createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightRuleContributor getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightRuleContributor>(create);
  static InsightRuleContributor? _defaultInstance;

  @$pb.TagNumber(2831086)
  $pb.PbList<$core.String> get keys => $_getList(0);

  @$pb.TagNumber(37512135)
  $pb.PbList<InsightRuleContributorDatapoint> get datapoints => $_getList(1);

  @$pb.TagNumber(260431894)
  $core.double get approximateaggregatevalue => $_getN(2);
  @$pb.TagNumber(260431894)
  set approximateaggregatevalue($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(260431894)
  $core.bool hasApproximateaggregatevalue() => $_has(2);
  @$pb.TagNumber(260431894)
  void clearApproximateaggregatevalue() => $_clearField(260431894);
}

class InsightRuleContributorDatapoint extends $pb.GeneratedMessage {
  factory InsightRuleContributorDatapoint({
    $core.String? timestamp,
    $core.double? approximatevalue,
  }) {
    final result = create();
    if (timestamp != null) result.timestamp = timestamp;
    if (approximatevalue != null) result.approximatevalue = approximatevalue;
    return result;
  }

  InsightRuleContributorDatapoint._();

  factory InsightRuleContributorDatapoint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightRuleContributorDatapoint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightRuleContributorDatapoint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aD(428761855, _omitFieldNames ? '' : 'approximatevalue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleContributorDatapoint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleContributorDatapoint copyWith(
          void Function(InsightRuleContributorDatapoint) updates) =>
      super.copyWith(
              (message) => updates(message as InsightRuleContributorDatapoint))
          as InsightRuleContributorDatapoint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightRuleContributorDatapoint create() =>
      InsightRuleContributorDatapoint._();
  @$core.override
  InsightRuleContributorDatapoint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightRuleContributorDatapoint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightRuleContributorDatapoint>(
          create);
  static InsightRuleContributorDatapoint? _defaultInstance;

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(0);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(0);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(428761855)
  $core.double get approximatevalue => $_getN(1);
  @$pb.TagNumber(428761855)
  set approximatevalue($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(428761855)
  $core.bool hasApproximatevalue() => $_has(1);
  @$pb.TagNumber(428761855)
  void clearApproximatevalue() => $_clearField(428761855);
}

class InsightRuleMetricDatapoint extends $pb.GeneratedMessage {
  factory InsightRuleMetricDatapoint({
    $core.double? maximum,
    $core.double? average,
    $core.String? timestamp,
    $core.double? samplecount,
    $core.double? maxcontributorvalue,
    $core.double? minimum,
    $core.double? uniquecontributors,
    $core.double? sum,
  }) {
    final result = create();
    if (maximum != null) result.maximum = maximum;
    if (average != null) result.average = average;
    if (timestamp != null) result.timestamp = timestamp;
    if (samplecount != null) result.samplecount = samplecount;
    if (maxcontributorvalue != null)
      result.maxcontributorvalue = maxcontributorvalue;
    if (minimum != null) result.minimum = minimum;
    if (uniquecontributors != null)
      result.uniquecontributors = uniquecontributors;
    if (sum != null) result.sum = sum;
    return result;
  }

  InsightRuleMetricDatapoint._();

  factory InsightRuleMetricDatapoint.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InsightRuleMetricDatapoint.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InsightRuleMetricDatapoint',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aD(43343394, _omitFieldNames ? '' : 'maximum')
    ..aD(134712777, _omitFieldNames ? '' : 'average')
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aD(384919839, _omitFieldNames ? '' : 'samplecount')
    ..aD(389006680, _omitFieldNames ? '' : 'maxcontributorvalue')
    ..aD(438536856, _omitFieldNames ? '' : 'minimum')
    ..aD(504374565, _omitFieldNames ? '' : 'uniquecontributors')
    ..aD(534303305, _omitFieldNames ? '' : 'sum')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleMetricDatapoint clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InsightRuleMetricDatapoint copyWith(
          void Function(InsightRuleMetricDatapoint) updates) =>
      super.copyWith(
              (message) => updates(message as InsightRuleMetricDatapoint))
          as InsightRuleMetricDatapoint;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InsightRuleMetricDatapoint create() => InsightRuleMetricDatapoint._();
  @$core.override
  InsightRuleMetricDatapoint createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InsightRuleMetricDatapoint getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InsightRuleMetricDatapoint>(create);
  static InsightRuleMetricDatapoint? _defaultInstance;

  @$pb.TagNumber(43343394)
  $core.double get maximum => $_getN(0);
  @$pb.TagNumber(43343394)
  set maximum($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(43343394)
  $core.bool hasMaximum() => $_has(0);
  @$pb.TagNumber(43343394)
  void clearMaximum() => $_clearField(43343394);

  @$pb.TagNumber(134712777)
  $core.double get average => $_getN(1);
  @$pb.TagNumber(134712777)
  set average($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(134712777)
  $core.bool hasAverage() => $_has(1);
  @$pb.TagNumber(134712777)
  void clearAverage() => $_clearField(134712777);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(2);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(2, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(384919839)
  $core.double get samplecount => $_getN(3);
  @$pb.TagNumber(384919839)
  set samplecount($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(384919839)
  $core.bool hasSamplecount() => $_has(3);
  @$pb.TagNumber(384919839)
  void clearSamplecount() => $_clearField(384919839);

  @$pb.TagNumber(389006680)
  $core.double get maxcontributorvalue => $_getN(4);
  @$pb.TagNumber(389006680)
  set maxcontributorvalue($core.double value) => $_setDouble(4, value);
  @$pb.TagNumber(389006680)
  $core.bool hasMaxcontributorvalue() => $_has(4);
  @$pb.TagNumber(389006680)
  void clearMaxcontributorvalue() => $_clearField(389006680);

  @$pb.TagNumber(438536856)
  $core.double get minimum => $_getN(5);
  @$pb.TagNumber(438536856)
  set minimum($core.double value) => $_setDouble(5, value);
  @$pb.TagNumber(438536856)
  $core.bool hasMinimum() => $_has(5);
  @$pb.TagNumber(438536856)
  void clearMinimum() => $_clearField(438536856);

  @$pb.TagNumber(504374565)
  $core.double get uniquecontributors => $_getN(6);
  @$pb.TagNumber(504374565)
  set uniquecontributors($core.double value) => $_setDouble(6, value);
  @$pb.TagNumber(504374565)
  $core.bool hasUniquecontributors() => $_has(6);
  @$pb.TagNumber(504374565)
  void clearUniquecontributors() => $_clearField(504374565);

  @$pb.TagNumber(534303305)
  $core.double get sum => $_getN(7);
  @$pb.TagNumber(534303305)
  set sum($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(534303305)
  $core.bool hasSum() => $_has(7);
  @$pb.TagNumber(534303305)
  void clearSum() => $_clearField(534303305);
}

class InternalServiceFault extends $pb.GeneratedMessage {
  factory InternalServiceFault({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InternalServiceFault._();

  factory InternalServiceFault.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InternalServiceFault.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InternalServiceFault',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServiceFault clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InternalServiceFault copyWith(void Function(InternalServiceFault) updates) =>
      super.copyWith((message) => updates(message as InternalServiceFault))
          as InternalServiceFault;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InternalServiceFault create() => InternalServiceFault._();
  @$core.override
  InternalServiceFault createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InternalServiceFault getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InternalServiceFault>(create);
  static InternalServiceFault? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class InvalidFormatFault extends $pb.GeneratedMessage {
  factory InvalidFormatFault({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidFormatFault._();

  factory InvalidFormatFault.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidFormatFault.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidFormatFault',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidFormatFault clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidFormatFault copyWith(void Function(InvalidFormatFault) updates) =>
      super.copyWith((message) => updates(message as InvalidFormatFault))
          as InvalidFormatFault;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidFormatFault create() => InvalidFormatFault._();
  @$core.override
  InvalidFormatFault createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidFormatFault getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidFormatFault>(create);
  static InvalidFormatFault? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidNextToken extends $pb.GeneratedMessage {
  factory InvalidNextToken({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidNextToken._();

  factory InvalidNextToken.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidNextToken.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidNextToken',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidNextToken clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidNextToken copyWith(void Function(InvalidNextToken) updates) =>
      super.copyWith((message) => updates(message as InvalidNextToken))
          as InvalidNextToken;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidNextToken create() => InvalidNextToken._();
  @$core.override
  InvalidNextToken createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidNextToken getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidNextToken>(create);
  static InvalidNextToken? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidParameterCombinationException extends $pb.GeneratedMessage {
  factory InvalidParameterCombinationException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidParameterCombinationException._();

  factory InvalidParameterCombinationException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidParameterCombinationException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidParameterCombinationException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterCombinationException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterCombinationException copyWith(
          void Function(InvalidParameterCombinationException) updates) =>
      super.copyWith((message) =>
              updates(message as InvalidParameterCombinationException))
          as InvalidParameterCombinationException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidParameterCombinationException create() =>
      InvalidParameterCombinationException._();
  @$core.override
  InvalidParameterCombinationException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidParameterCombinationException getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          InvalidParameterCombinationException>(create);
  static InvalidParameterCombinationException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class InvalidParameterValueException extends $pb.GeneratedMessage {
  factory InvalidParameterValueException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  InvalidParameterValueException._();

  factory InvalidParameterValueException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory InvalidParameterValueException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'InvalidParameterValueException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterValueException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  InvalidParameterValueException copyWith(
          void Function(InvalidParameterValueException) updates) =>
      super.copyWith(
              (message) => updates(message as InvalidParameterValueException))
          as InvalidParameterValueException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static InvalidParameterValueException create() =>
      InvalidParameterValueException._();
  @$core.override
  InvalidParameterValueException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static InvalidParameterValueException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<InvalidParameterValueException>(create);
  static InvalidParameterValueException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class LabelOptions extends $pb.GeneratedMessage {
  factory LabelOptions({
    $core.String? timezone,
  }) {
    final result = create();
    if (timezone != null) result.timezone = timezone;
    return result;
  }

  LabelOptions._();

  factory LabelOptions.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LabelOptions.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LabelOptions',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(246302531, _omitFieldNames ? '' : 'timezone')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelOptions clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LabelOptions copyWith(void Function(LabelOptions) updates) =>
      super.copyWith((message) => updates(message as LabelOptions))
          as LabelOptions;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LabelOptions create() => LabelOptions._();
  @$core.override
  LabelOptions createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LabelOptions getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LabelOptions>(create);
  static LabelOptions? _defaultInstance;

  @$pb.TagNumber(246302531)
  $core.String get timezone => $_getSZ(0);
  @$pb.TagNumber(246302531)
  set timezone($core.String value) => $_setString(0, value);
  @$pb.TagNumber(246302531)
  $core.bool hasTimezone() => $_has(0);
  @$pb.TagNumber(246302531)
  void clearTimezone() => $_clearField(246302531);
}

class LimitExceededException extends $pb.GeneratedMessage {
  factory LimitExceededException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  LimitExceededException._();

  factory LimitExceededException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LimitExceededException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LimitExceededException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededException copyWith(
          void Function(LimitExceededException) updates) =>
      super.copyWith((message) => updates(message as LimitExceededException))
          as LimitExceededException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LimitExceededException create() => LimitExceededException._();
  @$core.override
  LimitExceededException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LimitExceededException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LimitExceededException>(create);
  static LimitExceededException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);
}

class LimitExceededFault extends $pb.GeneratedMessage {
  factory LimitExceededFault({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  LimitExceededFault._();

  factory LimitExceededFault.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LimitExceededFault.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LimitExceededFault',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededFault clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LimitExceededFault copyWith(void Function(LimitExceededFault) updates) =>
      super.copyWith((message) => updates(message as LimitExceededFault))
          as LimitExceededFault;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LimitExceededFault create() => LimitExceededFault._();
  @$core.override
  LimitExceededFault createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LimitExceededFault getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LimitExceededFault>(create);
  static LimitExceededFault? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ListAlarmMuteRulesInput extends $pb.GeneratedMessage {
  factory ListAlarmMuteRulesInput({
    $core.String? alarmname,
    $core.String? nexttoken,
    $core.int? maxrecords,
    $core.Iterable<AlarmMuteRuleStatus>? statuses,
  }) {
    final result = create();
    if (alarmname != null) result.alarmname = alarmname;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxrecords != null) result.maxrecords = maxrecords;
    if (statuses != null) result.statuses.addAll(statuses);
    return result;
  }

  ListAlarmMuteRulesInput._();

  factory ListAlarmMuteRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAlarmMuteRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAlarmMuteRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(220314370, _omitFieldNames ? '' : 'maxrecords')
    ..pc<AlarmMuteRuleStatus>(
        374056024, _omitFieldNames ? '' : 'statuses', $pb.PbFieldType.KE,
        valueOf: AlarmMuteRuleStatus.valueOf,
        enumValues: AlarmMuteRuleStatus.values,
        defaultEnumValue: AlarmMuteRuleStatus.ALARM_MUTE_RULE_STATUS_ACTIVE)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAlarmMuteRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAlarmMuteRulesInput copyWith(
          void Function(ListAlarmMuteRulesInput) updates) =>
      super.copyWith((message) => updates(message as ListAlarmMuteRulesInput))
          as ListAlarmMuteRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAlarmMuteRulesInput create() => ListAlarmMuteRulesInput._();
  @$core.override
  ListAlarmMuteRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAlarmMuteRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAlarmMuteRulesInput>(create);
  static ListAlarmMuteRulesInput? _defaultInstance;

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(0);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(0);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(220314370)
  $core.int get maxrecords => $_getIZ(2);
  @$pb.TagNumber(220314370)
  set maxrecords($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(220314370)
  $core.bool hasMaxrecords() => $_has(2);
  @$pb.TagNumber(220314370)
  void clearMaxrecords() => $_clearField(220314370);

  @$pb.TagNumber(374056024)
  $pb.PbList<AlarmMuteRuleStatus> get statuses => $_getList(3);
}

class ListAlarmMuteRulesOutput extends $pb.GeneratedMessage {
  factory ListAlarmMuteRulesOutput({
    $core.String? nexttoken,
    $core.Iterable<AlarmMuteRuleSummary>? alarmmuterulesummaries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (alarmmuterulesummaries != null)
      result.alarmmuterulesummaries.addAll(alarmmuterulesummaries);
    return result;
  }

  ListAlarmMuteRulesOutput._();

  factory ListAlarmMuteRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAlarmMuteRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAlarmMuteRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<AlarmMuteRuleSummary>(
        440906836, _omitFieldNames ? '' : 'alarmmuterulesummaries',
        subBuilder: AlarmMuteRuleSummary.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAlarmMuteRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAlarmMuteRulesOutput copyWith(
          void Function(ListAlarmMuteRulesOutput) updates) =>
      super.copyWith((message) => updates(message as ListAlarmMuteRulesOutput))
          as ListAlarmMuteRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAlarmMuteRulesOutput create() => ListAlarmMuteRulesOutput._();
  @$core.override
  ListAlarmMuteRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAlarmMuteRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAlarmMuteRulesOutput>(create);
  static ListAlarmMuteRulesOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(440906836)
  $pb.PbList<AlarmMuteRuleSummary> get alarmmuterulesummaries => $_getList(1);
}

class ListDashboardsInput extends $pb.GeneratedMessage {
  factory ListDashboardsInput({
    $core.String? nexttoken,
    $core.String? dashboardnameprefix,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (dashboardnameprefix != null)
      result.dashboardnameprefix = dashboardnameprefix;
    return result;
  }

  ListDashboardsInput._();

  factory ListDashboardsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDashboardsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDashboardsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOS(301015421, _omitFieldNames ? '' : 'dashboardnameprefix')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsInput copyWith(void Function(ListDashboardsInput) updates) =>
      super.copyWith((message) => updates(message as ListDashboardsInput))
          as ListDashboardsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDashboardsInput create() => ListDashboardsInput._();
  @$core.override
  ListDashboardsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDashboardsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDashboardsInput>(create);
  static ListDashboardsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(301015421)
  $core.String get dashboardnameprefix => $_getSZ(1);
  @$pb.TagNumber(301015421)
  set dashboardnameprefix($core.String value) => $_setString(1, value);
  @$pb.TagNumber(301015421)
  $core.bool hasDashboardnameprefix() => $_has(1);
  @$pb.TagNumber(301015421)
  void clearDashboardnameprefix() => $_clearField(301015421);
}

class ListDashboardsOutput extends $pb.GeneratedMessage {
  factory ListDashboardsOutput({
    $core.String? nexttoken,
    $core.Iterable<DashboardEntry>? dashboardentries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (dashboardentries != null)
      result.dashboardentries.addAll(dashboardentries);
    return result;
  }

  ListDashboardsOutput._();

  factory ListDashboardsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListDashboardsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListDashboardsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<DashboardEntry>(281104826, _omitFieldNames ? '' : 'dashboardentries',
        subBuilder: DashboardEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListDashboardsOutput copyWith(void Function(ListDashboardsOutput) updates) =>
      super.copyWith((message) => updates(message as ListDashboardsOutput))
          as ListDashboardsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListDashboardsOutput create() => ListDashboardsOutput._();
  @$core.override
  ListDashboardsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListDashboardsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListDashboardsOutput>(create);
  static ListDashboardsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(281104826)
  $pb.PbList<DashboardEntry> get dashboardentries => $_getList(1);
}

class ListManagedInsightRulesInput extends $pb.GeneratedMessage {
  factory ListManagedInsightRulesInput({
    $core.String? nexttoken,
    $core.int? maxresults,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  ListManagedInsightRulesInput._();

  factory ListManagedInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListManagedInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListManagedInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedInsightRulesInput copyWith(
          void Function(ListManagedInsightRulesInput) updates) =>
      super.copyWith(
              (message) => updates(message as ListManagedInsightRulesInput))
          as ListManagedInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListManagedInsightRulesInput create() =>
      ListManagedInsightRulesInput._();
  @$core.override
  ListManagedInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListManagedInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListManagedInsightRulesInput>(create);
  static ListManagedInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(2);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(2, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(2);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class ListManagedInsightRulesOutput extends $pb.GeneratedMessage {
  factory ListManagedInsightRulesOutput({
    $core.String? nexttoken,
    $core.Iterable<ManagedRuleDescription>? managedrules,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (managedrules != null) result.managedrules.addAll(managedrules);
    return result;
  }

  ListManagedInsightRulesOutput._();

  factory ListManagedInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListManagedInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListManagedInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<ManagedRuleDescription>(
        347090906, _omitFieldNames ? '' : 'managedrules',
        subBuilder: ManagedRuleDescription.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListManagedInsightRulesOutput copyWith(
          void Function(ListManagedInsightRulesOutput) updates) =>
      super.copyWith(
              (message) => updates(message as ListManagedInsightRulesOutput))
          as ListManagedInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListManagedInsightRulesOutput create() =>
      ListManagedInsightRulesOutput._();
  @$core.override
  ListManagedInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListManagedInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListManagedInsightRulesOutput>(create);
  static ListManagedInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(347090906)
  $pb.PbList<ManagedRuleDescription> get managedrules => $_getList(1);
}

class ListMetricStreamsInput extends $pb.GeneratedMessage {
  factory ListMetricStreamsInput({
    $core.String? nexttoken,
    $core.int? maxresults,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (maxresults != null) result.maxresults = maxresults;
    return result;
  }

  ListMetricStreamsInput._();

  factory ListMetricStreamsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMetricStreamsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMetricStreamsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aI(275174450, _omitFieldNames ? '' : 'maxresults')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricStreamsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricStreamsInput copyWith(
          void Function(ListMetricStreamsInput) updates) =>
      super.copyWith((message) => updates(message as ListMetricStreamsInput))
          as ListMetricStreamsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMetricStreamsInput create() => ListMetricStreamsInput._();
  @$core.override
  ListMetricStreamsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMetricStreamsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMetricStreamsInput>(create);
  static ListMetricStreamsInput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(275174450)
  $core.int get maxresults => $_getIZ(1);
  @$pb.TagNumber(275174450)
  set maxresults($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(275174450)
  $core.bool hasMaxresults() => $_has(1);
  @$pb.TagNumber(275174450)
  void clearMaxresults() => $_clearField(275174450);
}

class ListMetricStreamsOutput extends $pb.GeneratedMessage {
  factory ListMetricStreamsOutput({
    $core.String? nexttoken,
    $core.Iterable<MetricStreamEntry>? entries,
  }) {
    final result = create();
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (entries != null) result.entries.addAll(entries);
    return result;
  }

  ListMetricStreamsOutput._();

  factory ListMetricStreamsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMetricStreamsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMetricStreamsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<MetricStreamEntry>(481075860, _omitFieldNames ? '' : 'entries',
        subBuilder: MetricStreamEntry.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricStreamsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricStreamsOutput copyWith(
          void Function(ListMetricStreamsOutput) updates) =>
      super.copyWith((message) => updates(message as ListMetricStreamsOutput))
          as ListMetricStreamsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMetricStreamsOutput create() => ListMetricStreamsOutput._();
  @$core.override
  ListMetricStreamsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMetricStreamsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMetricStreamsOutput>(create);
  static ListMetricStreamsOutput? _defaultInstance;

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(0);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(0, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(0);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(481075860)
  $pb.PbList<MetricStreamEntry> get entries => $_getList(1);
}

class ListMetricsInput extends $pb.GeneratedMessage {
  factory ListMetricsInput({
    $core.String? metricname,
    RecentlyActive? recentlyactive,
    $core.String? nexttoken,
    $core.bool? includelinkedaccounts,
    $core.String? owningaccount,
    $core.String? namespace,
    $core.Iterable<DimensionFilter>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (recentlyactive != null) result.recentlyactive = recentlyactive;
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (includelinkedaccounts != null)
      result.includelinkedaccounts = includelinkedaccounts;
    if (owningaccount != null) result.owningaccount = owningaccount;
    if (namespace != null) result.namespace = namespace;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  ListMetricsInput._();

  factory ListMetricsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMetricsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMetricsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aE<RecentlyActive>(179617336, _omitFieldNames ? '' : 'recentlyactive',
        enumValues: RecentlyActive.values)
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..aOB(228605043, _omitFieldNames ? '' : 'includelinkedaccounts')
    ..aOS(339968073, _omitFieldNames ? '' : 'owningaccount')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPM<DimensionFilter>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: DimensionFilter.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricsInput copyWith(void Function(ListMetricsInput) updates) =>
      super.copyWith((message) => updates(message as ListMetricsInput))
          as ListMetricsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMetricsInput create() => ListMetricsInput._();
  @$core.override
  ListMetricsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMetricsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMetricsInput>(create);
  static ListMetricsInput? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(179617336)
  RecentlyActive get recentlyactive => $_getN(1);
  @$pb.TagNumber(179617336)
  set recentlyactive(RecentlyActive value) => $_setField(179617336, value);
  @$pb.TagNumber(179617336)
  $core.bool hasRecentlyactive() => $_has(1);
  @$pb.TagNumber(179617336)
  void clearRecentlyactive() => $_clearField(179617336);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(2);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(2, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(2);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(228605043)
  $core.bool get includelinkedaccounts => $_getBF(3);
  @$pb.TagNumber(228605043)
  set includelinkedaccounts($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(228605043)
  $core.bool hasIncludelinkedaccounts() => $_has(3);
  @$pb.TagNumber(228605043)
  void clearIncludelinkedaccounts() => $_clearField(228605043);

  @$pb.TagNumber(339968073)
  $core.String get owningaccount => $_getSZ(4);
  @$pb.TagNumber(339968073)
  set owningaccount($core.String value) => $_setString(4, value);
  @$pb.TagNumber(339968073)
  $core.bool hasOwningaccount() => $_has(4);
  @$pb.TagNumber(339968073)
  void clearOwningaccount() => $_clearField(339968073);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(5);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(5, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(5);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(462933457)
  $pb.PbList<DimensionFilter> get dimensions => $_getList(6);
}

class ListMetricsOutput extends $pb.GeneratedMessage {
  factory ListMetricsOutput({
    $core.Iterable<$core.String>? owningaccounts,
    $core.String? nexttoken,
    $core.Iterable<Metric>? metrics,
  }) {
    final result = create();
    if (owningaccounts != null) result.owningaccounts.addAll(owningaccounts);
    if (nexttoken != null) result.nexttoken = nexttoken;
    if (metrics != null) result.metrics.addAll(metrics);
    return result;
  }

  ListMetricsOutput._();

  factory ListMetricsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListMetricsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListMetricsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(21159138, _omitFieldNames ? '' : 'owningaccounts')
    ..aOS(216957566, _omitFieldNames ? '' : 'nexttoken')
    ..pPM<Metric>(436365847, _omitFieldNames ? '' : 'metrics',
        subBuilder: Metric.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListMetricsOutput copyWith(void Function(ListMetricsOutput) updates) =>
      super.copyWith((message) => updates(message as ListMetricsOutput))
          as ListMetricsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListMetricsOutput create() => ListMetricsOutput._();
  @$core.override
  ListMetricsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListMetricsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListMetricsOutput>(create);
  static ListMetricsOutput? _defaultInstance;

  @$pb.TagNumber(21159138)
  $pb.PbList<$core.String> get owningaccounts => $_getList(0);

  @$pb.TagNumber(216957566)
  $core.String get nexttoken => $_getSZ(1);
  @$pb.TagNumber(216957566)
  set nexttoken($core.String value) => $_setString(1, value);
  @$pb.TagNumber(216957566)
  $core.bool hasNexttoken() => $_has(1);
  @$pb.TagNumber(216957566)
  void clearNexttoken() => $_clearField(216957566);

  @$pb.TagNumber(436365847)
  $pb.PbList<Metric> get metrics => $_getList(2);
}

class ListTagsForResourceInput extends $pb.GeneratedMessage {
  factory ListTagsForResourceInput({
    $core.String? resourcearn,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  ListTagsForResourceInput._();

  factory ListTagsForResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceInput copyWith(
          void Function(ListTagsForResourceInput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForResourceInput))
          as ListTagsForResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceInput create() => ListTagsForResourceInput._();
  @$core.override
  ListTagsForResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceInput>(create);
  static ListTagsForResourceInput? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class ListTagsForResourceOutput extends $pb.GeneratedMessage {
  factory ListTagsForResourceOutput({
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  ListTagsForResourceOutput._();

  factory ListTagsForResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListTagsForResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListTagsForResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListTagsForResourceOutput copyWith(
          void Function(ListTagsForResourceOutput) updates) =>
      super.copyWith((message) => updates(message as ListTagsForResourceOutput))
          as ListTagsForResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceOutput create() => ListTagsForResourceOutput._();
  @$core.override
  ListTagsForResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListTagsForResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListTagsForResourceOutput>(create);
  static ListTagsForResourceOutput? _defaultInstance;

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(0);
}

class ManagedRule extends $pb.GeneratedMessage {
  factory ManagedRule({
    $core.String? resourcearn,
    $core.Iterable<Tag>? tags,
    $core.String? templatename,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addAll(tags);
    if (templatename != null) result.templatename = templatename;
    return result;
  }

  ManagedRule._();

  factory ManagedRule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(480529457, _omitFieldNames ? '' : 'templatename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRule copyWith(void Function(ManagedRule) updates) =>
      super.copyWith((message) => updates(message as ManagedRule))
          as ManagedRule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRule create() => ManagedRule._();
  @$core.override
  ManagedRule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRule getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRule>(create);
  static ManagedRule? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);

  @$pb.TagNumber(480529457)
  $core.String get templatename => $_getSZ(2);
  @$pb.TagNumber(480529457)
  set templatename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(480529457)
  $core.bool hasTemplatename() => $_has(2);
  @$pb.TagNumber(480529457)
  void clearTemplatename() => $_clearField(480529457);
}

class ManagedRuleDescription extends $pb.GeneratedMessage {
  factory ManagedRuleDescription({
    $core.String? resourcearn,
    ManagedRuleState? rulestate,
    $core.String? templatename,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (rulestate != null) result.rulestate = rulestate;
    if (templatename != null) result.templatename = templatename;
    return result;
  }

  ManagedRuleDescription._();

  factory ManagedRuleDescription.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleDescription.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleDescription',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..aOM<ManagedRuleState>(419921189, _omitFieldNames ? '' : 'rulestate',
        subBuilder: ManagedRuleState.create)
    ..aOS(480529457, _omitFieldNames ? '' : 'templatename')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleDescription clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleDescription copyWith(
          void Function(ManagedRuleDescription) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleDescription))
          as ManagedRuleDescription;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleDescription create() => ManagedRuleDescription._();
  @$core.override
  ManagedRuleDescription createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleDescription getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleDescription>(create);
  static ManagedRuleDescription? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(419921189)
  ManagedRuleState get rulestate => $_getN(1);
  @$pb.TagNumber(419921189)
  set rulestate(ManagedRuleState value) => $_setField(419921189, value);
  @$pb.TagNumber(419921189)
  $core.bool hasRulestate() => $_has(1);
  @$pb.TagNumber(419921189)
  void clearRulestate() => $_clearField(419921189);
  @$pb.TagNumber(419921189)
  ManagedRuleState ensureRulestate() => $_ensure(1);

  @$pb.TagNumber(480529457)
  $core.String get templatename => $_getSZ(2);
  @$pb.TagNumber(480529457)
  set templatename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(480529457)
  $core.bool hasTemplatename() => $_has(2);
  @$pb.TagNumber(480529457)
  void clearTemplatename() => $_clearField(480529457);
}

class ManagedRuleState extends $pb.GeneratedMessage {
  factory ManagedRuleState({
    $core.String? rulename,
    $core.String? state,
  }) {
    final result = create();
    if (rulename != null) result.rulename = rulename;
    if (state != null) result.state = state;
    return result;
  }

  ManagedRuleState._();

  factory ManagedRuleState.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ManagedRuleState.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ManagedRuleState',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(214688793, _omitFieldNames ? '' : 'rulename')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleState clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ManagedRuleState copyWith(void Function(ManagedRuleState) updates) =>
      super.copyWith((message) => updates(message as ManagedRuleState))
          as ManagedRuleState;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ManagedRuleState create() => ManagedRuleState._();
  @$core.override
  ManagedRuleState createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ManagedRuleState getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ManagedRuleState>(create);
  static ManagedRuleState? _defaultInstance;

  @$pb.TagNumber(214688793)
  $core.String get rulename => $_getSZ(0);
  @$pb.TagNumber(214688793)
  set rulename($core.String value) => $_setString(0, value);
  @$pb.TagNumber(214688793)
  $core.bool hasRulename() => $_has(0);
  @$pb.TagNumber(214688793)
  void clearRulename() => $_clearField(214688793);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(1);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(1, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(1);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);
}

class MessageData extends $pb.GeneratedMessage {
  factory MessageData({
    $core.String? value,
    $core.String? code,
  }) {
    final result = create();
    if (value != null) result.value = value;
    if (code != null) result.code = code;
    return result;
  }

  MessageData._();

  factory MessageData.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MessageData.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MessageData',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..aOS(425572629, _omitFieldNames ? '' : 'code')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageData clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MessageData copyWith(void Function(MessageData) updates) =>
      super.copyWith((message) => updates(message as MessageData))
          as MessageData;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MessageData create() => MessageData._();
  @$core.override
  MessageData createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MessageData getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MessageData>(create);
  static MessageData? _defaultInstance;

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(0);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(0, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(0);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(425572629)
  $core.String get code => $_getSZ(1);
  @$pb.TagNumber(425572629)
  set code($core.String value) => $_setString(1, value);
  @$pb.TagNumber(425572629)
  $core.bool hasCode() => $_has(1);
  @$pb.TagNumber(425572629)
  void clearCode() => $_clearField(425572629);
}

class Metric extends $pb.GeneratedMessage {
  factory Metric({
    $core.String? metricname,
    $core.String? namespace,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (namespace != null) result.namespace = namespace;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  Metric._();

  factory Metric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Metric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Metric',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Metric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Metric copyWith(void Function(Metric) updates) =>
      super.copyWith((message) => updates(message as Metric)) as Metric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Metric create() => Metric._();
  @$core.override
  Metric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Metric getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Metric>(create);
  static Metric? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(1);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(1);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(2);
}

class MetricAlarm extends $pb.GeneratedMessage {
  factory MetricAlarm({
    ComparisonOperator? comparisonoperator,
    $core.String? alarmdescription,
    $core.String? thresholdmetricid,
    $core.String? statetransitionedtimestamp,
    $core.String? alarmconfigurationupdatedtimestamp,
    $core.String? alarmname,
    $core.String? evaluatelowsamplecountpercentile,
    Statistic? statistic,
    $core.String? metricname,
    $core.int? period,
    $core.bool? actionsenabled,
    $core.int? datapointstoalarm,
    StandardUnit? unit,
    $core.int? evaluationperiods,
    $core.String? treatmissingdata,
    $core.String? statereasondata,
    $core.String? alarmarn,
    $core.String? extendedstatistic,
    StateValue? statevalue,
    EvaluationState? evaluationstate,
    $core.String? namespace,
    $core.Iterable<$core.String>? alarmactions,
    $core.String? stateupdatedtimestamp,
    $core.String? statereason,
    $core.Iterable<$core.String>? okactions,
    $core.double? threshold,
    $core.Iterable<MetricDataQuery>? metrics,
    $core.Iterable<Dimension>? dimensions,
    $core.Iterable<$core.String>? insufficientdataactions,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (alarmdescription != null) result.alarmdescription = alarmdescription;
    if (thresholdmetricid != null) result.thresholdmetricid = thresholdmetricid;
    if (statetransitionedtimestamp != null)
      result.statetransitionedtimestamp = statetransitionedtimestamp;
    if (alarmconfigurationupdatedtimestamp != null)
      result.alarmconfigurationupdatedtimestamp =
          alarmconfigurationupdatedtimestamp;
    if (alarmname != null) result.alarmname = alarmname;
    if (evaluatelowsamplecountpercentile != null)
      result.evaluatelowsamplecountpercentile =
          evaluatelowsamplecountpercentile;
    if (statistic != null) result.statistic = statistic;
    if (metricname != null) result.metricname = metricname;
    if (period != null) result.period = period;
    if (actionsenabled != null) result.actionsenabled = actionsenabled;
    if (datapointstoalarm != null) result.datapointstoalarm = datapointstoalarm;
    if (unit != null) result.unit = unit;
    if (evaluationperiods != null) result.evaluationperiods = evaluationperiods;
    if (treatmissingdata != null) result.treatmissingdata = treatmissingdata;
    if (statereasondata != null) result.statereasondata = statereasondata;
    if (alarmarn != null) result.alarmarn = alarmarn;
    if (extendedstatistic != null) result.extendedstatistic = extendedstatistic;
    if (statevalue != null) result.statevalue = statevalue;
    if (evaluationstate != null) result.evaluationstate = evaluationstate;
    if (namespace != null) result.namespace = namespace;
    if (alarmactions != null) result.alarmactions.addAll(alarmactions);
    if (stateupdatedtimestamp != null)
      result.stateupdatedtimestamp = stateupdatedtimestamp;
    if (statereason != null) result.statereason = statereason;
    if (okactions != null) result.okactions.addAll(okactions);
    if (threshold != null) result.threshold = threshold;
    if (metrics != null) result.metrics.addAll(metrics);
    if (dimensions != null) result.dimensions.addAll(dimensions);
    if (insufficientdataactions != null)
      result.insufficientdataactions.addAll(insufficientdataactions);
    return result;
  }

  MetricAlarm._();

  factory MetricAlarm.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricAlarm.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricAlarm',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..aOS(30583613, _omitFieldNames ? '' : 'alarmdescription')
    ..aOS(41342194, _omitFieldNames ? '' : 'thresholdmetricid')
    ..aOS(46811093, _omitFieldNames ? '' : 'statetransitionedtimestamp')
    ..aOS(51935494, _omitFieldNames ? '' : 'alarmconfigurationupdatedtimestamp')
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(64366323, _omitFieldNames ? '' : 'evaluatelowsamplecountpercentile')
    ..aE<Statistic>(67293470, _omitFieldNames ? '' : 'statistic',
        enumValues: Statistic.values)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aOB(126027878, _omitFieldNames ? '' : 'actionsenabled')
    ..aI(146378889, _omitFieldNames ? '' : 'datapointstoalarm')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aI(214794856, _omitFieldNames ? '' : 'evaluationperiods')
    ..aOS(226418150, _omitFieldNames ? '' : 'treatmissingdata')
    ..aOS(262771075, _omitFieldNames ? '' : 'statereasondata')
    ..aOS(276019462, _omitFieldNames ? '' : 'alarmarn')
    ..aOS(285620763, _omitFieldNames ? '' : 'extendedstatistic')
    ..aE<StateValue>(334526008, _omitFieldNames ? '' : 'statevalue',
        enumValues: StateValue.values)
    ..aE<EvaluationState>(338863857, _omitFieldNames ? '' : 'evaluationstate',
        enumValues: EvaluationState.values)
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPS(355779506, _omitFieldNames ? '' : 'alarmactions')
    ..aOS(375406848, _omitFieldNames ? '' : 'stateupdatedtimestamp')
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..pPS(377443763, _omitFieldNames ? '' : 'okactions')
    ..aD(394884505, _omitFieldNames ? '' : 'threshold')
    ..pPM<MetricDataQuery>(436365847, _omitFieldNames ? '' : 'metrics',
        subBuilder: MetricDataQuery.create)
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..pPS(498450778, _omitFieldNames ? '' : 'insufficientdataactions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricAlarm clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricAlarm copyWith(void Function(MetricAlarm) updates) =>
      super.copyWith((message) => updates(message as MetricAlarm))
          as MetricAlarm;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricAlarm create() => MetricAlarm._();
  @$core.override
  MetricAlarm createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricAlarm getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricAlarm>(create);
  static MetricAlarm? _defaultInstance;

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(30583613)
  $core.String get alarmdescription => $_getSZ(1);
  @$pb.TagNumber(30583613)
  set alarmdescription($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30583613)
  $core.bool hasAlarmdescription() => $_has(1);
  @$pb.TagNumber(30583613)
  void clearAlarmdescription() => $_clearField(30583613);

  @$pb.TagNumber(41342194)
  $core.String get thresholdmetricid => $_getSZ(2);
  @$pb.TagNumber(41342194)
  set thresholdmetricid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(41342194)
  $core.bool hasThresholdmetricid() => $_has(2);
  @$pb.TagNumber(41342194)
  void clearThresholdmetricid() => $_clearField(41342194);

  @$pb.TagNumber(46811093)
  $core.String get statetransitionedtimestamp => $_getSZ(3);
  @$pb.TagNumber(46811093)
  set statetransitionedtimestamp($core.String value) => $_setString(3, value);
  @$pb.TagNumber(46811093)
  $core.bool hasStatetransitionedtimestamp() => $_has(3);
  @$pb.TagNumber(46811093)
  void clearStatetransitionedtimestamp() => $_clearField(46811093);

  @$pb.TagNumber(51935494)
  $core.String get alarmconfigurationupdatedtimestamp => $_getSZ(4);
  @$pb.TagNumber(51935494)
  set alarmconfigurationupdatedtimestamp($core.String value) =>
      $_setString(4, value);
  @$pb.TagNumber(51935494)
  $core.bool hasAlarmconfigurationupdatedtimestamp() => $_has(4);
  @$pb.TagNumber(51935494)
  void clearAlarmconfigurationupdatedtimestamp() => $_clearField(51935494);

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(5);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(5, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(5);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(64366323)
  $core.String get evaluatelowsamplecountpercentile => $_getSZ(6);
  @$pb.TagNumber(64366323)
  set evaluatelowsamplecountpercentile($core.String value) =>
      $_setString(6, value);
  @$pb.TagNumber(64366323)
  $core.bool hasEvaluatelowsamplecountpercentile() => $_has(6);
  @$pb.TagNumber(64366323)
  void clearEvaluatelowsamplecountpercentile() => $_clearField(64366323);

  @$pb.TagNumber(67293470)
  Statistic get statistic => $_getN(7);
  @$pb.TagNumber(67293470)
  set statistic(Statistic value) => $_setField(67293470, value);
  @$pb.TagNumber(67293470)
  $core.bool hasStatistic() => $_has(7);
  @$pb.TagNumber(67293470)
  void clearStatistic() => $_clearField(67293470);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(8);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(8, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(8);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(9);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(9);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(126027878)
  $core.bool get actionsenabled => $_getBF(10);
  @$pb.TagNumber(126027878)
  set actionsenabled($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(126027878)
  $core.bool hasActionsenabled() => $_has(10);
  @$pb.TagNumber(126027878)
  void clearActionsenabled() => $_clearField(126027878);

  @$pb.TagNumber(146378889)
  $core.int get datapointstoalarm => $_getIZ(11);
  @$pb.TagNumber(146378889)
  set datapointstoalarm($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(146378889)
  $core.bool hasDatapointstoalarm() => $_has(11);
  @$pb.TagNumber(146378889)
  void clearDatapointstoalarm() => $_clearField(146378889);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(12);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(12);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(214794856)
  $core.int get evaluationperiods => $_getIZ(13);
  @$pb.TagNumber(214794856)
  set evaluationperiods($core.int value) => $_setSignedInt32(13, value);
  @$pb.TagNumber(214794856)
  $core.bool hasEvaluationperiods() => $_has(13);
  @$pb.TagNumber(214794856)
  void clearEvaluationperiods() => $_clearField(214794856);

  @$pb.TagNumber(226418150)
  $core.String get treatmissingdata => $_getSZ(14);
  @$pb.TagNumber(226418150)
  set treatmissingdata($core.String value) => $_setString(14, value);
  @$pb.TagNumber(226418150)
  $core.bool hasTreatmissingdata() => $_has(14);
  @$pb.TagNumber(226418150)
  void clearTreatmissingdata() => $_clearField(226418150);

  @$pb.TagNumber(262771075)
  $core.String get statereasondata => $_getSZ(15);
  @$pb.TagNumber(262771075)
  set statereasondata($core.String value) => $_setString(15, value);
  @$pb.TagNumber(262771075)
  $core.bool hasStatereasondata() => $_has(15);
  @$pb.TagNumber(262771075)
  void clearStatereasondata() => $_clearField(262771075);

  @$pb.TagNumber(276019462)
  $core.String get alarmarn => $_getSZ(16);
  @$pb.TagNumber(276019462)
  set alarmarn($core.String value) => $_setString(16, value);
  @$pb.TagNumber(276019462)
  $core.bool hasAlarmarn() => $_has(16);
  @$pb.TagNumber(276019462)
  void clearAlarmarn() => $_clearField(276019462);

  @$pb.TagNumber(285620763)
  $core.String get extendedstatistic => $_getSZ(17);
  @$pb.TagNumber(285620763)
  set extendedstatistic($core.String value) => $_setString(17, value);
  @$pb.TagNumber(285620763)
  $core.bool hasExtendedstatistic() => $_has(17);
  @$pb.TagNumber(285620763)
  void clearExtendedstatistic() => $_clearField(285620763);

  @$pb.TagNumber(334526008)
  StateValue get statevalue => $_getN(18);
  @$pb.TagNumber(334526008)
  set statevalue(StateValue value) => $_setField(334526008, value);
  @$pb.TagNumber(334526008)
  $core.bool hasStatevalue() => $_has(18);
  @$pb.TagNumber(334526008)
  void clearStatevalue() => $_clearField(334526008);

  @$pb.TagNumber(338863857)
  EvaluationState get evaluationstate => $_getN(19);
  @$pb.TagNumber(338863857)
  set evaluationstate(EvaluationState value) => $_setField(338863857, value);
  @$pb.TagNumber(338863857)
  $core.bool hasEvaluationstate() => $_has(19);
  @$pb.TagNumber(338863857)
  void clearEvaluationstate() => $_clearField(338863857);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(20);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(20, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(20);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(355779506)
  $pb.PbList<$core.String> get alarmactions => $_getList(21);

  @$pb.TagNumber(375406848)
  $core.String get stateupdatedtimestamp => $_getSZ(22);
  @$pb.TagNumber(375406848)
  set stateupdatedtimestamp($core.String value) => $_setString(22, value);
  @$pb.TagNumber(375406848)
  $core.bool hasStateupdatedtimestamp() => $_has(22);
  @$pb.TagNumber(375406848)
  void clearStateupdatedtimestamp() => $_clearField(375406848);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(23);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(23, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(23);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);

  @$pb.TagNumber(377443763)
  $pb.PbList<$core.String> get okactions => $_getList(24);

  @$pb.TagNumber(394884505)
  $core.double get threshold => $_getN(25);
  @$pb.TagNumber(394884505)
  set threshold($core.double value) => $_setDouble(25, value);
  @$pb.TagNumber(394884505)
  $core.bool hasThreshold() => $_has(25);
  @$pb.TagNumber(394884505)
  void clearThreshold() => $_clearField(394884505);

  @$pb.TagNumber(436365847)
  $pb.PbList<MetricDataQuery> get metrics => $_getList(26);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(27);

  @$pb.TagNumber(498450778)
  $pb.PbList<$core.String> get insufficientdataactions => $_getList(28);
}

class MetricCharacteristics extends $pb.GeneratedMessage {
  factory MetricCharacteristics({
    $core.bool? periodicspikes,
  }) {
    final result = create();
    if (periodicspikes != null) result.periodicspikes = periodicspikes;
    return result;
  }

  MetricCharacteristics._();

  factory MetricCharacteristics.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricCharacteristics.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricCharacteristics',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOB(454095604, _omitFieldNames ? '' : 'periodicspikes')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricCharacteristics clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricCharacteristics copyWith(
          void Function(MetricCharacteristics) updates) =>
      super.copyWith((message) => updates(message as MetricCharacteristics))
          as MetricCharacteristics;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricCharacteristics create() => MetricCharacteristics._();
  @$core.override
  MetricCharacteristics createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricCharacteristics getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricCharacteristics>(create);
  static MetricCharacteristics? _defaultInstance;

  @$pb.TagNumber(454095604)
  $core.bool get periodicspikes => $_getBF(0);
  @$pb.TagNumber(454095604)
  set periodicspikes($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(454095604)
  $core.bool hasPeriodicspikes() => $_has(0);
  @$pb.TagNumber(454095604)
  void clearPeriodicspikes() => $_clearField(454095604);
}

class MetricDataQuery extends $pb.GeneratedMessage {
  factory MetricDataQuery({
    $core.String? accountid,
    $core.int? period,
    MetricStat? metricstat,
    $core.String? expression,
    $core.String? id,
    $core.bool? returndata,
    $core.String? label,
  }) {
    final result = create();
    if (accountid != null) result.accountid = accountid;
    if (period != null) result.period = period;
    if (metricstat != null) result.metricstat = metricstat;
    if (expression != null) result.expression = expression;
    if (id != null) result.id = id;
    if (returndata != null) result.returndata = returndata;
    if (label != null) result.label = label;
    return result;
  }

  MetricDataQuery._();

  factory MetricDataQuery.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricDataQuery.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricDataQuery',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(65954002, _omitFieldNames ? '' : 'accountid')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aOM<MetricStat>(167166628, _omitFieldNames ? '' : 'metricstat',
        subBuilder: MetricStat.create)
    ..aOS(193051916, _omitFieldNames ? '' : 'expression')
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..aOB(512985704, _omitFieldNames ? '' : 'returndata')
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDataQuery clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDataQuery copyWith(void Function(MetricDataQuery) updates) =>
      super.copyWith((message) => updates(message as MetricDataQuery))
          as MetricDataQuery;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricDataQuery create() => MetricDataQuery._();
  @$core.override
  MetricDataQuery createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricDataQuery getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricDataQuery>(create);
  static MetricDataQuery? _defaultInstance;

  @$pb.TagNumber(65954002)
  $core.String get accountid => $_getSZ(0);
  @$pb.TagNumber(65954002)
  set accountid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(65954002)
  $core.bool hasAccountid() => $_has(0);
  @$pb.TagNumber(65954002)
  void clearAccountid() => $_clearField(65954002);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(1);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(1);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(167166628)
  MetricStat get metricstat => $_getN(2);
  @$pb.TagNumber(167166628)
  set metricstat(MetricStat value) => $_setField(167166628, value);
  @$pb.TagNumber(167166628)
  $core.bool hasMetricstat() => $_has(2);
  @$pb.TagNumber(167166628)
  void clearMetricstat() => $_clearField(167166628);
  @$pb.TagNumber(167166628)
  MetricStat ensureMetricstat() => $_ensure(2);

  @$pb.TagNumber(193051916)
  $core.String get expression => $_getSZ(3);
  @$pb.TagNumber(193051916)
  set expression($core.String value) => $_setString(3, value);
  @$pb.TagNumber(193051916)
  $core.bool hasExpression() => $_has(3);
  @$pb.TagNumber(193051916)
  void clearExpression() => $_clearField(193051916);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(4);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(4, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(4);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(512985704)
  $core.bool get returndata => $_getBF(5);
  @$pb.TagNumber(512985704)
  set returndata($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(512985704)
  $core.bool hasReturndata() => $_has(5);
  @$pb.TagNumber(512985704)
  void clearReturndata() => $_clearField(512985704);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(6);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(6, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(6);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class MetricDataResult extends $pb.GeneratedMessage {
  factory MetricDataResult({
    $core.Iterable<$core.String>? timestamps,
    $core.Iterable<$core.double>? values,
    StatusCode? statuscode,
    $core.String? id,
    $core.Iterable<MessageData>? messages,
    $core.String? label,
  }) {
    final result = create();
    if (timestamps != null) result.timestamps.addAll(timestamps);
    if (values != null) result.values.addAll(values);
    if (statuscode != null) result.statuscode = statuscode;
    if (id != null) result.id = id;
    if (messages != null) result.messages.addAll(messages);
    if (label != null) result.label = label;
    return result;
  }

  MetricDataResult._();

  factory MetricDataResult.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricDataResult.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricDataResult',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(213534737, _omitFieldNames ? '' : 'timestamps')
    ..p<$core.double>(
        223158876, _omitFieldNames ? '' : 'values', $pb.PbFieldType.KD)
    ..aE<StatusCode>(303830783, _omitFieldNames ? '' : 'statuscode',
        enumValues: StatusCode.values)
    ..aOS(384350465, _omitFieldNames ? '' : 'id')
    ..pPM<MessageData>(409018326, _omitFieldNames ? '' : 'messages',
        subBuilder: MessageData.create)
    ..aOS(516747934, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDataResult clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDataResult copyWith(void Function(MetricDataResult) updates) =>
      super.copyWith((message) => updates(message as MetricDataResult))
          as MetricDataResult;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricDataResult create() => MetricDataResult._();
  @$core.override
  MetricDataResult createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricDataResult getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricDataResult>(create);
  static MetricDataResult? _defaultInstance;

  @$pb.TagNumber(213534737)
  $pb.PbList<$core.String> get timestamps => $_getList(0);

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.double> get values => $_getList(1);

  @$pb.TagNumber(303830783)
  StatusCode get statuscode => $_getN(2);
  @$pb.TagNumber(303830783)
  set statuscode(StatusCode value) => $_setField(303830783, value);
  @$pb.TagNumber(303830783)
  $core.bool hasStatuscode() => $_has(2);
  @$pb.TagNumber(303830783)
  void clearStatuscode() => $_clearField(303830783);

  @$pb.TagNumber(384350465)
  $core.String get id => $_getSZ(3);
  @$pb.TagNumber(384350465)
  set id($core.String value) => $_setString(3, value);
  @$pb.TagNumber(384350465)
  $core.bool hasId() => $_has(3);
  @$pb.TagNumber(384350465)
  void clearId() => $_clearField(384350465);

  @$pb.TagNumber(409018326)
  $pb.PbList<MessageData> get messages => $_getList(4);

  @$pb.TagNumber(516747934)
  $core.String get label => $_getSZ(5);
  @$pb.TagNumber(516747934)
  set label($core.String value) => $_setString(5, value);
  @$pb.TagNumber(516747934)
  $core.bool hasLabel() => $_has(5);
  @$pb.TagNumber(516747934)
  void clearLabel() => $_clearField(516747934);
}

class MetricDatum extends $pb.GeneratedMessage {
  factory MetricDatum({
    StatisticSet? statisticvalues,
    $core.String? metricname,
    $core.Iterable<$core.double>? counts,
    StandardUnit? unit,
    $core.String? timestamp,
    $core.int? storageresolution,
    $core.Iterable<$core.double>? values,
    $core.double? value,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (statisticvalues != null) result.statisticvalues = statisticvalues;
    if (metricname != null) result.metricname = metricname;
    if (counts != null) result.counts.addAll(counts);
    if (unit != null) result.unit = unit;
    if (timestamp != null) result.timestamp = timestamp;
    if (storageresolution != null) result.storageresolution = storageresolution;
    if (values != null) result.values.addAll(values);
    if (value != null) result.value = value;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  MetricDatum._();

  factory MetricDatum.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricDatum.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricDatum',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOM<StatisticSet>(92623908, _omitFieldNames ? '' : 'statisticvalues',
        subBuilder: StatisticSet.create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..p<$core.double>(
        113775526, _omitFieldNames ? '' : 'counts', $pb.PbFieldType.KD)
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aOS(162390468, _omitFieldNames ? '' : 'timestamp')
    ..aI(194430293, _omitFieldNames ? '' : 'storageresolution')
    ..p<$core.double>(
        223158876, _omitFieldNames ? '' : 'values', $pb.PbFieldType.KD)
    ..aD(289929579, _omitFieldNames ? '' : 'value')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDatum clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricDatum copyWith(void Function(MetricDatum) updates) =>
      super.copyWith((message) => updates(message as MetricDatum))
          as MetricDatum;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricDatum create() => MetricDatum._();
  @$core.override
  MetricDatum createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricDatum getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricDatum>(create);
  static MetricDatum? _defaultInstance;

  @$pb.TagNumber(92623908)
  StatisticSet get statisticvalues => $_getN(0);
  @$pb.TagNumber(92623908)
  set statisticvalues(StatisticSet value) => $_setField(92623908, value);
  @$pb.TagNumber(92623908)
  $core.bool hasStatisticvalues() => $_has(0);
  @$pb.TagNumber(92623908)
  void clearStatisticvalues() => $_clearField(92623908);
  @$pb.TagNumber(92623908)
  StatisticSet ensureStatisticvalues() => $_ensure(0);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(113775526)
  $pb.PbList<$core.double> get counts => $_getList(2);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(3);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(3);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(162390468)
  $core.String get timestamp => $_getSZ(4);
  @$pb.TagNumber(162390468)
  set timestamp($core.String value) => $_setString(4, value);
  @$pb.TagNumber(162390468)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(162390468)
  void clearTimestamp() => $_clearField(162390468);

  @$pb.TagNumber(194430293)
  $core.int get storageresolution => $_getIZ(5);
  @$pb.TagNumber(194430293)
  set storageresolution($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(194430293)
  $core.bool hasStorageresolution() => $_has(5);
  @$pb.TagNumber(194430293)
  void clearStorageresolution() => $_clearField(194430293);

  @$pb.TagNumber(223158876)
  $pb.PbList<$core.double> get values => $_getList(6);

  @$pb.TagNumber(289929579)
  $core.double get value => $_getN(7);
  @$pb.TagNumber(289929579)
  set value($core.double value) => $_setDouble(7, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(7);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(8);
}

class MetricMathAnomalyDetector extends $pb.GeneratedMessage {
  factory MetricMathAnomalyDetector({
    $core.Iterable<MetricDataQuery>? metricdataqueries,
  }) {
    final result = create();
    if (metricdataqueries != null)
      result.metricdataqueries.addAll(metricdataqueries);
    return result;
  }

  MetricMathAnomalyDetector._();

  factory MetricMathAnomalyDetector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricMathAnomalyDetector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricMathAnomalyDetector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<MetricDataQuery>(92209504, _omitFieldNames ? '' : 'metricdataqueries',
        subBuilder: MetricDataQuery.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricMathAnomalyDetector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricMathAnomalyDetector copyWith(
          void Function(MetricMathAnomalyDetector) updates) =>
      super.copyWith((message) => updates(message as MetricMathAnomalyDetector))
          as MetricMathAnomalyDetector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricMathAnomalyDetector create() => MetricMathAnomalyDetector._();
  @$core.override
  MetricMathAnomalyDetector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricMathAnomalyDetector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricMathAnomalyDetector>(create);
  static MetricMathAnomalyDetector? _defaultInstance;

  @$pb.TagNumber(92209504)
  $pb.PbList<MetricDataQuery> get metricdataqueries => $_getList(0);
}

class MetricStat extends $pb.GeneratedMessage {
  factory MetricStat({
    $core.int? period,
    StandardUnit? unit,
    $core.String? stat,
    Metric? metric,
  }) {
    final result = create();
    if (period != null) result.period = period;
    if (unit != null) result.unit = unit;
    if (stat != null) result.stat = stat;
    if (metric != null) result.metric = metric;
    return result;
  }

  MetricStat._();

  factory MetricStat.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricStat.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricStat',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aOS(325549752, _omitFieldNames ? '' : 'stat')
    ..aOM<Metric>(386667298, _omitFieldNames ? '' : 'metric',
        subBuilder: Metric.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStat clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStat copyWith(void Function(MetricStat) updates) =>
      super.copyWith((message) => updates(message as MetricStat)) as MetricStat;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricStat create() => MetricStat._();
  @$core.override
  MetricStat createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricStat getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricStat>(create);
  static MetricStat? _defaultInstance;

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(0);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(0);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(1);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(1);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(325549752)
  $core.String get stat => $_getSZ(2);
  @$pb.TagNumber(325549752)
  set stat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(325549752)
  $core.bool hasStat() => $_has(2);
  @$pb.TagNumber(325549752)
  void clearStat() => $_clearField(325549752);

  @$pb.TagNumber(386667298)
  Metric get metric => $_getN(3);
  @$pb.TagNumber(386667298)
  set metric(Metric value) => $_setField(386667298, value);
  @$pb.TagNumber(386667298)
  $core.bool hasMetric() => $_has(3);
  @$pb.TagNumber(386667298)
  void clearMetric() => $_clearField(386667298);
  @$pb.TagNumber(386667298)
  Metric ensureMetric() => $_ensure(3);
}

class MetricStreamEntry extends $pb.GeneratedMessage {
  factory MetricStreamEntry({
    $core.String? firehosearn,
    $core.String? lastupdatedate,
    $core.String? name,
    $core.String? creationdate,
    $core.String? arn,
    $core.String? state,
    MetricStreamOutputFormat? outputformat,
  }) {
    final result = create();
    if (firehosearn != null) result.firehosearn = firehosearn;
    if (lastupdatedate != null) result.lastupdatedate = lastupdatedate;
    if (name != null) result.name = name;
    if (creationdate != null) result.creationdate = creationdate;
    if (arn != null) result.arn = arn;
    if (state != null) result.state = state;
    if (outputformat != null) result.outputformat = outputformat;
    return result;
  }

  MetricStreamEntry._();

  factory MetricStreamEntry.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricStreamEntry.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricStreamEntry',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(55494200, _omitFieldNames ? '' : 'firehosearn')
    ..aOS(65755725, _omitFieldNames ? '' : 'lastupdatedate')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(288222305, _omitFieldNames ? '' : 'creationdate')
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..aOS(502047895, _omitFieldNames ? '' : 'state')
    ..aE<MetricStreamOutputFormat>(
        519139704, _omitFieldNames ? '' : 'outputformat',
        enumValues: MetricStreamOutputFormat.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamEntry clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamEntry copyWith(void Function(MetricStreamEntry) updates) =>
      super.copyWith((message) => updates(message as MetricStreamEntry))
          as MetricStreamEntry;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricStreamEntry create() => MetricStreamEntry._();
  @$core.override
  MetricStreamEntry createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricStreamEntry getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricStreamEntry>(create);
  static MetricStreamEntry? _defaultInstance;

  @$pb.TagNumber(55494200)
  $core.String get firehosearn => $_getSZ(0);
  @$pb.TagNumber(55494200)
  set firehosearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(55494200)
  $core.bool hasFirehosearn() => $_has(0);
  @$pb.TagNumber(55494200)
  void clearFirehosearn() => $_clearField(55494200);

  @$pb.TagNumber(65755725)
  $core.String get lastupdatedate => $_getSZ(1);
  @$pb.TagNumber(65755725)
  set lastupdatedate($core.String value) => $_setString(1, value);
  @$pb.TagNumber(65755725)
  $core.bool hasLastupdatedate() => $_has(1);
  @$pb.TagNumber(65755725)
  void clearLastupdatedate() => $_clearField(65755725);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(2);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(2, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(2);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(288222305)
  $core.String get creationdate => $_getSZ(3);
  @$pb.TagNumber(288222305)
  set creationdate($core.String value) => $_setString(3, value);
  @$pb.TagNumber(288222305)
  $core.bool hasCreationdate() => $_has(3);
  @$pb.TagNumber(288222305)
  void clearCreationdate() => $_clearField(288222305);

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(4);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(4);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);

  @$pb.TagNumber(502047895)
  $core.String get state => $_getSZ(5);
  @$pb.TagNumber(502047895)
  set state($core.String value) => $_setString(5, value);
  @$pb.TagNumber(502047895)
  $core.bool hasState() => $_has(5);
  @$pb.TagNumber(502047895)
  void clearState() => $_clearField(502047895);

  @$pb.TagNumber(519139704)
  MetricStreamOutputFormat get outputformat => $_getN(6);
  @$pb.TagNumber(519139704)
  set outputformat(MetricStreamOutputFormat value) =>
      $_setField(519139704, value);
  @$pb.TagNumber(519139704)
  $core.bool hasOutputformat() => $_has(6);
  @$pb.TagNumber(519139704)
  void clearOutputformat() => $_clearField(519139704);
}

class MetricStreamFilter extends $pb.GeneratedMessage {
  factory MetricStreamFilter({
    $core.Iterable<$core.String>? metricnames,
    $core.String? namespace,
  }) {
    final result = create();
    if (metricnames != null) result.metricnames.addAll(metricnames);
    if (namespace != null) result.namespace = namespace;
    return result;
  }

  MetricStreamFilter._();

  factory MetricStreamFilter.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricStreamFilter.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricStreamFilter',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(56384300, _omitFieldNames ? '' : 'metricnames')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamFilter clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamFilter copyWith(void Function(MetricStreamFilter) updates) =>
      super.copyWith((message) => updates(message as MetricStreamFilter))
          as MetricStreamFilter;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricStreamFilter create() => MetricStreamFilter._();
  @$core.override
  MetricStreamFilter createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricStreamFilter getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricStreamFilter>(create);
  static MetricStreamFilter? _defaultInstance;

  @$pb.TagNumber(56384300)
  $pb.PbList<$core.String> get metricnames => $_getList(0);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(1);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(1);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);
}

class MetricStreamStatisticsConfiguration extends $pb.GeneratedMessage {
  factory MetricStreamStatisticsConfiguration({
    $core.Iterable<$core.String>? additionalstatistics,
    $core.Iterable<MetricStreamStatisticsMetric>? includemetrics,
  }) {
    final result = create();
    if (additionalstatistics != null)
      result.additionalstatistics.addAll(additionalstatistics);
    if (includemetrics != null) result.includemetrics.addAll(includemetrics);
    return result;
  }

  MetricStreamStatisticsConfiguration._();

  factory MetricStreamStatisticsConfiguration.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricStreamStatisticsConfiguration.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricStreamStatisticsConfiguration',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(50823856, _omitFieldNames ? '' : 'additionalstatistics')
    ..pPM<MetricStreamStatisticsMetric>(
        169287849, _omitFieldNames ? '' : 'includemetrics',
        subBuilder: MetricStreamStatisticsMetric.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamStatisticsConfiguration clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamStatisticsConfiguration copyWith(
          void Function(MetricStreamStatisticsConfiguration) updates) =>
      super.copyWith((message) =>
              updates(message as MetricStreamStatisticsConfiguration))
          as MetricStreamStatisticsConfiguration;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricStreamStatisticsConfiguration create() =>
      MetricStreamStatisticsConfiguration._();
  @$core.override
  MetricStreamStatisticsConfiguration createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricStreamStatisticsConfiguration getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          MetricStreamStatisticsConfiguration>(create);
  static MetricStreamStatisticsConfiguration? _defaultInstance;

  @$pb.TagNumber(50823856)
  $pb.PbList<$core.String> get additionalstatistics => $_getList(0);

  @$pb.TagNumber(169287849)
  $pb.PbList<MetricStreamStatisticsMetric> get includemetrics => $_getList(1);
}

class MetricStreamStatisticsMetric extends $pb.GeneratedMessage {
  factory MetricStreamStatisticsMetric({
    $core.String? metricname,
    $core.String? namespace,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (namespace != null) result.namespace = namespace;
    return result;
  }

  MetricStreamStatisticsMetric._();

  factory MetricStreamStatisticsMetric.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricStreamStatisticsMetric.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricStreamStatisticsMetric',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamStatisticsMetric clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricStreamStatisticsMetric copyWith(
          void Function(MetricStreamStatisticsMetric) updates) =>
      super.copyWith(
              (message) => updates(message as MetricStreamStatisticsMetric))
          as MetricStreamStatisticsMetric;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricStreamStatisticsMetric create() =>
      MetricStreamStatisticsMetric._();
  @$core.override
  MetricStreamStatisticsMetric createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricStreamStatisticsMetric getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricStreamStatisticsMetric>(create);
  static MetricStreamStatisticsMetric? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(1);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(1, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(1);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);
}

class MissingRequiredParameterException extends $pb.GeneratedMessage {
  factory MissingRequiredParameterException({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  MissingRequiredParameterException._();

  factory MissingRequiredParameterException.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MissingRequiredParameterException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MissingRequiredParameterException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MissingRequiredParameterException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MissingRequiredParameterException copyWith(
          void Function(MissingRequiredParameterException) updates) =>
      super.copyWith((message) =>
              updates(message as MissingRequiredParameterException))
          as MissingRequiredParameterException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MissingRequiredParameterException create() =>
      MissingRequiredParameterException._();
  @$core.override
  MissingRequiredParameterException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MissingRequiredParameterException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MissingRequiredParameterException>(
          create);
  static MissingRequiredParameterException? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class MuteTargets extends $pb.GeneratedMessage {
  factory MuteTargets({
    $core.Iterable<$core.String>? alarmnames,
  }) {
    final result = create();
    if (alarmnames != null) result.alarmnames.addAll(alarmnames);
    return result;
  }

  MuteTargets._();

  factory MuteTargets.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MuteTargets.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MuteTargets',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(90874583, _omitFieldNames ? '' : 'alarmnames')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MuteTargets clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MuteTargets copyWith(void Function(MuteTargets) updates) =>
      super.copyWith((message) => updates(message as MuteTargets))
          as MuteTargets;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MuteTargets create() => MuteTargets._();
  @$core.override
  MuteTargets createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MuteTargets getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MuteTargets>(create);
  static MuteTargets? _defaultInstance;

  @$pb.TagNumber(90874583)
  $pb.PbList<$core.String> get alarmnames => $_getList(0);
}

class PartialFailure extends $pb.GeneratedMessage {
  factory PartialFailure({
    $core.String? failurecode,
    $core.String? exceptiontype,
    $core.String? failuredescription,
    $core.String? failureresource,
  }) {
    final result = create();
    if (failurecode != null) result.failurecode = failurecode;
    if (exceptiontype != null) result.exceptiontype = exceptiontype;
    if (failuredescription != null)
      result.failuredescription = failuredescription;
    if (failureresource != null) result.failureresource = failureresource;
    return result;
  }

  PartialFailure._();

  factory PartialFailure.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PartialFailure.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PartialFailure',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(84707897, _omitFieldNames ? '' : 'failurecode')
    ..aOS(176478087, _omitFieldNames ? '' : 'exceptiontype')
    ..aOS(249093686, _omitFieldNames ? '' : 'failuredescription')
    ..aOS(315801394, _omitFieldNames ? '' : 'failureresource')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartialFailure clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PartialFailure copyWith(void Function(PartialFailure) updates) =>
      super.copyWith((message) => updates(message as PartialFailure))
          as PartialFailure;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PartialFailure create() => PartialFailure._();
  @$core.override
  PartialFailure createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PartialFailure getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PartialFailure>(create);
  static PartialFailure? _defaultInstance;

  @$pb.TagNumber(84707897)
  $core.String get failurecode => $_getSZ(0);
  @$pb.TagNumber(84707897)
  set failurecode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(84707897)
  $core.bool hasFailurecode() => $_has(0);
  @$pb.TagNumber(84707897)
  void clearFailurecode() => $_clearField(84707897);

  @$pb.TagNumber(176478087)
  $core.String get exceptiontype => $_getSZ(1);
  @$pb.TagNumber(176478087)
  set exceptiontype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(176478087)
  $core.bool hasExceptiontype() => $_has(1);
  @$pb.TagNumber(176478087)
  void clearExceptiontype() => $_clearField(176478087);

  @$pb.TagNumber(249093686)
  $core.String get failuredescription => $_getSZ(2);
  @$pb.TagNumber(249093686)
  set failuredescription($core.String value) => $_setString(2, value);
  @$pb.TagNumber(249093686)
  $core.bool hasFailuredescription() => $_has(2);
  @$pb.TagNumber(249093686)
  void clearFailuredescription() => $_clearField(249093686);

  @$pb.TagNumber(315801394)
  $core.String get failureresource => $_getSZ(3);
  @$pb.TagNumber(315801394)
  set failureresource($core.String value) => $_setString(3, value);
  @$pb.TagNumber(315801394)
  $core.bool hasFailureresource() => $_has(3);
  @$pb.TagNumber(315801394)
  void clearFailureresource() => $_clearField(315801394);
}

class PutAlarmMuteRuleInput extends $pb.GeneratedMessage {
  factory PutAlarmMuteRuleInput({
    $core.String? description,
    $core.String? name,
    MuteTargets? mutetargets,
    $core.Iterable<Tag>? tags,
    $core.String? startdate,
    $core.String? expiredate,
    Rule? rule,
  }) {
    final result = create();
    if (description != null) result.description = description;
    if (name != null) result.name = name;
    if (mutetargets != null) result.mutetargets = mutetargets;
    if (tags != null) result.tags.addAll(tags);
    if (startdate != null) result.startdate = startdate;
    if (expiredate != null) result.expiredate = expiredate;
    if (rule != null) result.rule = rule;
    return result;
  }

  PutAlarmMuteRuleInput._();

  factory PutAlarmMuteRuleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutAlarmMuteRuleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutAlarmMuteRuleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(115243530, _omitFieldNames ? '' : 'description')
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOM<MuteTargets>(356717131, _omitFieldNames ? '' : 'mutetargets',
        subBuilder: MuteTargets.create)
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(445135996, _omitFieldNames ? '' : 'startdate')
    ..aOS(454143207, _omitFieldNames ? '' : 'expiredate')
    ..aOM<Rule>(475696372, _omitFieldNames ? '' : 'rule',
        subBuilder: Rule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAlarmMuteRuleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAlarmMuteRuleInput copyWith(
          void Function(PutAlarmMuteRuleInput) updates) =>
      super.copyWith((message) => updates(message as PutAlarmMuteRuleInput))
          as PutAlarmMuteRuleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutAlarmMuteRuleInput create() => PutAlarmMuteRuleInput._();
  @$core.override
  PutAlarmMuteRuleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutAlarmMuteRuleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutAlarmMuteRuleInput>(create);
  static PutAlarmMuteRuleInput? _defaultInstance;

  @$pb.TagNumber(115243530)
  $core.String get description => $_getSZ(0);
  @$pb.TagNumber(115243530)
  set description($core.String value) => $_setString(0, value);
  @$pb.TagNumber(115243530)
  $core.bool hasDescription() => $_has(0);
  @$pb.TagNumber(115243530)
  void clearDescription() => $_clearField(115243530);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(1, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(1);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(356717131)
  MuteTargets get mutetargets => $_getN(2);
  @$pb.TagNumber(356717131)
  set mutetargets(MuteTargets value) => $_setField(356717131, value);
  @$pb.TagNumber(356717131)
  $core.bool hasMutetargets() => $_has(2);
  @$pb.TagNumber(356717131)
  void clearMutetargets() => $_clearField(356717131);
  @$pb.TagNumber(356717131)
  MuteTargets ensureMutetargets() => $_ensure(2);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);

  @$pb.TagNumber(445135996)
  $core.String get startdate => $_getSZ(4);
  @$pb.TagNumber(445135996)
  set startdate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(445135996)
  $core.bool hasStartdate() => $_has(4);
  @$pb.TagNumber(445135996)
  void clearStartdate() => $_clearField(445135996);

  @$pb.TagNumber(454143207)
  $core.String get expiredate => $_getSZ(5);
  @$pb.TagNumber(454143207)
  set expiredate($core.String value) => $_setString(5, value);
  @$pb.TagNumber(454143207)
  $core.bool hasExpiredate() => $_has(5);
  @$pb.TagNumber(454143207)
  void clearExpiredate() => $_clearField(454143207);

  @$pb.TagNumber(475696372)
  Rule get rule => $_getN(6);
  @$pb.TagNumber(475696372)
  set rule(Rule value) => $_setField(475696372, value);
  @$pb.TagNumber(475696372)
  $core.bool hasRule() => $_has(6);
  @$pb.TagNumber(475696372)
  void clearRule() => $_clearField(475696372);
  @$pb.TagNumber(475696372)
  Rule ensureRule() => $_ensure(6);
}

class PutAnomalyDetectorInput extends $pb.GeneratedMessage {
  factory PutAnomalyDetectorInput({
    $core.String? metricname,
    MetricCharacteristics? metriccharacteristics,
    $core.String? stat,
    $core.String? namespace,
    SingleMetricAnomalyDetector? singlemetricanomalydetector,
    MetricMathAnomalyDetector? metricmathanomalydetector,
    AnomalyDetectorConfiguration? configuration,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (metricname != null) result.metricname = metricname;
    if (metriccharacteristics != null)
      result.metriccharacteristics = metriccharacteristics;
    if (stat != null) result.stat = stat;
    if (namespace != null) result.namespace = namespace;
    if (singlemetricanomalydetector != null)
      result.singlemetricanomalydetector = singlemetricanomalydetector;
    if (metricmathanomalydetector != null)
      result.metricmathanomalydetector = metricmathanomalydetector;
    if (configuration != null) result.configuration = configuration;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  PutAnomalyDetectorInput._();

  factory PutAnomalyDetectorInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutAnomalyDetectorInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutAnomalyDetectorInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOM<MetricCharacteristics>(
        129404842, _omitFieldNames ? '' : 'metriccharacteristics',
        subBuilder: MetricCharacteristics.create)
    ..aOS(325549752, _omitFieldNames ? '' : 'stat')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..aOM<SingleMetricAnomalyDetector>(
        418531933, _omitFieldNames ? '' : 'singlemetricanomalydetector',
        subBuilder: SingleMetricAnomalyDetector.create)
    ..aOM<MetricMathAnomalyDetector>(
        439381475, _omitFieldNames ? '' : 'metricmathanomalydetector',
        subBuilder: MetricMathAnomalyDetector.create)
    ..aOM<AnomalyDetectorConfiguration>(
        442426458, _omitFieldNames ? '' : 'configuration',
        subBuilder: AnomalyDetectorConfiguration.create)
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAnomalyDetectorInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAnomalyDetectorInput copyWith(
          void Function(PutAnomalyDetectorInput) updates) =>
      super.copyWith((message) => updates(message as PutAnomalyDetectorInput))
          as PutAnomalyDetectorInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutAnomalyDetectorInput create() => PutAnomalyDetectorInput._();
  @$core.override
  PutAnomalyDetectorInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutAnomalyDetectorInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutAnomalyDetectorInput>(create);
  static PutAnomalyDetectorInput? _defaultInstance;

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(0);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(0);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(129404842)
  MetricCharacteristics get metriccharacteristics => $_getN(1);
  @$pb.TagNumber(129404842)
  set metriccharacteristics(MetricCharacteristics value) =>
      $_setField(129404842, value);
  @$pb.TagNumber(129404842)
  $core.bool hasMetriccharacteristics() => $_has(1);
  @$pb.TagNumber(129404842)
  void clearMetriccharacteristics() => $_clearField(129404842);
  @$pb.TagNumber(129404842)
  MetricCharacteristics ensureMetriccharacteristics() => $_ensure(1);

  @$pb.TagNumber(325549752)
  $core.String get stat => $_getSZ(2);
  @$pb.TagNumber(325549752)
  set stat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(325549752)
  $core.bool hasStat() => $_has(2);
  @$pb.TagNumber(325549752)
  void clearStat() => $_clearField(325549752);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(3);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(3, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(3);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector get singlemetricanomalydetector => $_getN(4);
  @$pb.TagNumber(418531933)
  set singlemetricanomalydetector(SingleMetricAnomalyDetector value) =>
      $_setField(418531933, value);
  @$pb.TagNumber(418531933)
  $core.bool hasSinglemetricanomalydetector() => $_has(4);
  @$pb.TagNumber(418531933)
  void clearSinglemetricanomalydetector() => $_clearField(418531933);
  @$pb.TagNumber(418531933)
  SingleMetricAnomalyDetector ensureSinglemetricanomalydetector() =>
      $_ensure(4);

  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector get metricmathanomalydetector => $_getN(5);
  @$pb.TagNumber(439381475)
  set metricmathanomalydetector(MetricMathAnomalyDetector value) =>
      $_setField(439381475, value);
  @$pb.TagNumber(439381475)
  $core.bool hasMetricmathanomalydetector() => $_has(5);
  @$pb.TagNumber(439381475)
  void clearMetricmathanomalydetector() => $_clearField(439381475);
  @$pb.TagNumber(439381475)
  MetricMathAnomalyDetector ensureMetricmathanomalydetector() => $_ensure(5);

  @$pb.TagNumber(442426458)
  AnomalyDetectorConfiguration get configuration => $_getN(6);
  @$pb.TagNumber(442426458)
  set configuration(AnomalyDetectorConfiguration value) =>
      $_setField(442426458, value);
  @$pb.TagNumber(442426458)
  $core.bool hasConfiguration() => $_has(6);
  @$pb.TagNumber(442426458)
  void clearConfiguration() => $_clearField(442426458);
  @$pb.TagNumber(442426458)
  AnomalyDetectorConfiguration ensureConfiguration() => $_ensure(6);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(7);
}

class PutAnomalyDetectorOutput extends $pb.GeneratedMessage {
  factory PutAnomalyDetectorOutput() => create();

  PutAnomalyDetectorOutput._();

  factory PutAnomalyDetectorOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutAnomalyDetectorOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutAnomalyDetectorOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAnomalyDetectorOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutAnomalyDetectorOutput copyWith(
          void Function(PutAnomalyDetectorOutput) updates) =>
      super.copyWith((message) => updates(message as PutAnomalyDetectorOutput))
          as PutAnomalyDetectorOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutAnomalyDetectorOutput create() => PutAnomalyDetectorOutput._();
  @$core.override
  PutAnomalyDetectorOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutAnomalyDetectorOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutAnomalyDetectorOutput>(create);
  static PutAnomalyDetectorOutput? _defaultInstance;
}

class PutCompositeAlarmInput extends $pb.GeneratedMessage {
  factory PutCompositeAlarmInput({
    $core.String? alarmdescription,
    $core.String? alarmname,
    $core.int? actionssuppressorwaitperiod,
    $core.bool? actionsenabled,
    $core.int? actionssuppressorextensionperiod,
    $core.Iterable<$core.String>? alarmactions,
    $core.Iterable<$core.String>? okactions,
    $core.Iterable<Tag>? tags,
    $core.String? actionssuppressor,
    $core.Iterable<$core.String>? insufficientdataactions,
    $core.String? alarmrule,
  }) {
    final result = create();
    if (alarmdescription != null) result.alarmdescription = alarmdescription;
    if (alarmname != null) result.alarmname = alarmname;
    if (actionssuppressorwaitperiod != null)
      result.actionssuppressorwaitperiod = actionssuppressorwaitperiod;
    if (actionsenabled != null) result.actionsenabled = actionsenabled;
    if (actionssuppressorextensionperiod != null)
      result.actionssuppressorextensionperiod =
          actionssuppressorextensionperiod;
    if (alarmactions != null) result.alarmactions.addAll(alarmactions);
    if (okactions != null) result.okactions.addAll(okactions);
    if (tags != null) result.tags.addAll(tags);
    if (actionssuppressor != null) result.actionssuppressor = actionssuppressor;
    if (insufficientdataactions != null)
      result.insufficientdataactions.addAll(insufficientdataactions);
    if (alarmrule != null) result.alarmrule = alarmrule;
    return result;
  }

  PutCompositeAlarmInput._();

  factory PutCompositeAlarmInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutCompositeAlarmInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutCompositeAlarmInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(30583613, _omitFieldNames ? '' : 'alarmdescription')
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aI(80227275, _omitFieldNames ? '' : 'actionssuppressorwaitperiod')
    ..aOB(126027878, _omitFieldNames ? '' : 'actionsenabled')
    ..aI(183723543, _omitFieldNames ? '' : 'actionssuppressorextensionperiod')
    ..pPS(355779506, _omitFieldNames ? '' : 'alarmactions')
    ..pPS(377443763, _omitFieldNames ? '' : 'okactions')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(486417535, _omitFieldNames ? '' : 'actionssuppressor')
    ..pPS(498450778, _omitFieldNames ? '' : 'insufficientdataactions')
    ..aOS(516996973, _omitFieldNames ? '' : 'alarmrule')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCompositeAlarmInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutCompositeAlarmInput copyWith(
          void Function(PutCompositeAlarmInput) updates) =>
      super.copyWith((message) => updates(message as PutCompositeAlarmInput))
          as PutCompositeAlarmInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutCompositeAlarmInput create() => PutCompositeAlarmInput._();
  @$core.override
  PutCompositeAlarmInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutCompositeAlarmInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutCompositeAlarmInput>(create);
  static PutCompositeAlarmInput? _defaultInstance;

  @$pb.TagNumber(30583613)
  $core.String get alarmdescription => $_getSZ(0);
  @$pb.TagNumber(30583613)
  set alarmdescription($core.String value) => $_setString(0, value);
  @$pb.TagNumber(30583613)
  $core.bool hasAlarmdescription() => $_has(0);
  @$pb.TagNumber(30583613)
  void clearAlarmdescription() => $_clearField(30583613);

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(1);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(1);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(80227275)
  $core.int get actionssuppressorwaitperiod => $_getIZ(2);
  @$pb.TagNumber(80227275)
  set actionssuppressorwaitperiod($core.int value) =>
      $_setSignedInt32(2, value);
  @$pb.TagNumber(80227275)
  $core.bool hasActionssuppressorwaitperiod() => $_has(2);
  @$pb.TagNumber(80227275)
  void clearActionssuppressorwaitperiod() => $_clearField(80227275);

  @$pb.TagNumber(126027878)
  $core.bool get actionsenabled => $_getBF(3);
  @$pb.TagNumber(126027878)
  set actionsenabled($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(126027878)
  $core.bool hasActionsenabled() => $_has(3);
  @$pb.TagNumber(126027878)
  void clearActionsenabled() => $_clearField(126027878);

  @$pb.TagNumber(183723543)
  $core.int get actionssuppressorextensionperiod => $_getIZ(4);
  @$pb.TagNumber(183723543)
  set actionssuppressorextensionperiod($core.int value) =>
      $_setSignedInt32(4, value);
  @$pb.TagNumber(183723543)
  $core.bool hasActionssuppressorextensionperiod() => $_has(4);
  @$pb.TagNumber(183723543)
  void clearActionssuppressorextensionperiod() => $_clearField(183723543);

  @$pb.TagNumber(355779506)
  $pb.PbList<$core.String> get alarmactions => $_getList(5);

  @$pb.TagNumber(377443763)
  $pb.PbList<$core.String> get okactions => $_getList(6);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(7);

  @$pb.TagNumber(486417535)
  $core.String get actionssuppressor => $_getSZ(8);
  @$pb.TagNumber(486417535)
  set actionssuppressor($core.String value) => $_setString(8, value);
  @$pb.TagNumber(486417535)
  $core.bool hasActionssuppressor() => $_has(8);
  @$pb.TagNumber(486417535)
  void clearActionssuppressor() => $_clearField(486417535);

  @$pb.TagNumber(498450778)
  $pb.PbList<$core.String> get insufficientdataactions => $_getList(9);

  @$pb.TagNumber(516996973)
  $core.String get alarmrule => $_getSZ(10);
  @$pb.TagNumber(516996973)
  set alarmrule($core.String value) => $_setString(10, value);
  @$pb.TagNumber(516996973)
  $core.bool hasAlarmrule() => $_has(10);
  @$pb.TagNumber(516996973)
  void clearAlarmrule() => $_clearField(516996973);
}

class PutDashboardInput extends $pb.GeneratedMessage {
  factory PutDashboardInput({
    $core.String? dashboardbody,
    $core.String? dashboardname,
  }) {
    final result = create();
    if (dashboardbody != null) result.dashboardbody = dashboardbody;
    if (dashboardname != null) result.dashboardname = dashboardname;
    return result;
  }

  PutDashboardInput._();

  factory PutDashboardInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutDashboardInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutDashboardInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(3403236, _omitFieldNames ? '' : 'dashboardbody')
    ..aOS(506599873, _omitFieldNames ? '' : 'dashboardname')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDashboardInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDashboardInput copyWith(void Function(PutDashboardInput) updates) =>
      super.copyWith((message) => updates(message as PutDashboardInput))
          as PutDashboardInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutDashboardInput create() => PutDashboardInput._();
  @$core.override
  PutDashboardInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutDashboardInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutDashboardInput>(create);
  static PutDashboardInput? _defaultInstance;

  @$pb.TagNumber(3403236)
  $core.String get dashboardbody => $_getSZ(0);
  @$pb.TagNumber(3403236)
  set dashboardbody($core.String value) => $_setString(0, value);
  @$pb.TagNumber(3403236)
  $core.bool hasDashboardbody() => $_has(0);
  @$pb.TagNumber(3403236)
  void clearDashboardbody() => $_clearField(3403236);

  @$pb.TagNumber(506599873)
  $core.String get dashboardname => $_getSZ(1);
  @$pb.TagNumber(506599873)
  set dashboardname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(506599873)
  $core.bool hasDashboardname() => $_has(1);
  @$pb.TagNumber(506599873)
  void clearDashboardname() => $_clearField(506599873);
}

class PutDashboardOutput extends $pb.GeneratedMessage {
  factory PutDashboardOutput({
    $core.Iterable<DashboardValidationMessage>? dashboardvalidationmessages,
  }) {
    final result = create();
    if (dashboardvalidationmessages != null)
      result.dashboardvalidationmessages.addAll(dashboardvalidationmessages);
    return result;
  }

  PutDashboardOutput._();

  factory PutDashboardOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutDashboardOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutDashboardOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<DashboardValidationMessage>(
        515270089, _omitFieldNames ? '' : 'dashboardvalidationmessages',
        subBuilder: DashboardValidationMessage.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDashboardOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutDashboardOutput copyWith(void Function(PutDashboardOutput) updates) =>
      super.copyWith((message) => updates(message as PutDashboardOutput))
          as PutDashboardOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutDashboardOutput create() => PutDashboardOutput._();
  @$core.override
  PutDashboardOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutDashboardOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutDashboardOutput>(create);
  static PutDashboardOutput? _defaultInstance;

  @$pb.TagNumber(515270089)
  $pb.PbList<DashboardValidationMessage> get dashboardvalidationmessages =>
      $_getList(0);
}

class PutInsightRuleInput extends $pb.GeneratedMessage {
  factory PutInsightRuleInput({
    $core.bool? applyontransformedlogs,
    $core.String? ruledefinition,
    $core.String? rulename,
    $core.Iterable<Tag>? tags,
    $core.String? rulestate,
  }) {
    final result = create();
    if (applyontransformedlogs != null)
      result.applyontransformedlogs = applyontransformedlogs;
    if (ruledefinition != null) result.ruledefinition = ruledefinition;
    if (rulename != null) result.rulename = rulename;
    if (tags != null) result.tags.addAll(tags);
    if (rulestate != null) result.rulestate = rulestate;
    return result;
  }

  PutInsightRuleInput._();

  factory PutInsightRuleInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutInsightRuleInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutInsightRuleInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOB(30966533, _omitFieldNames ? '' : 'applyontransformedlogs')
    ..aOS(188932559, _omitFieldNames ? '' : 'ruledefinition')
    ..aOS(214688793, _omitFieldNames ? '' : 'rulename')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aOS(419921189, _omitFieldNames ? '' : 'rulestate')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightRuleInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightRuleInput copyWith(void Function(PutInsightRuleInput) updates) =>
      super.copyWith((message) => updates(message as PutInsightRuleInput))
          as PutInsightRuleInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutInsightRuleInput create() => PutInsightRuleInput._();
  @$core.override
  PutInsightRuleInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutInsightRuleInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutInsightRuleInput>(create);
  static PutInsightRuleInput? _defaultInstance;

  @$pb.TagNumber(30966533)
  $core.bool get applyontransformedlogs => $_getBF(0);
  @$pb.TagNumber(30966533)
  set applyontransformedlogs($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(30966533)
  $core.bool hasApplyontransformedlogs() => $_has(0);
  @$pb.TagNumber(30966533)
  void clearApplyontransformedlogs() => $_clearField(30966533);

  @$pb.TagNumber(188932559)
  $core.String get ruledefinition => $_getSZ(1);
  @$pb.TagNumber(188932559)
  set ruledefinition($core.String value) => $_setString(1, value);
  @$pb.TagNumber(188932559)
  $core.bool hasRuledefinition() => $_has(1);
  @$pb.TagNumber(188932559)
  void clearRuledefinition() => $_clearField(188932559);

  @$pb.TagNumber(214688793)
  $core.String get rulename => $_getSZ(2);
  @$pb.TagNumber(214688793)
  set rulename($core.String value) => $_setString(2, value);
  @$pb.TagNumber(214688793)
  $core.bool hasRulename() => $_has(2);
  @$pb.TagNumber(214688793)
  void clearRulename() => $_clearField(214688793);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(3);

  @$pb.TagNumber(419921189)
  $core.String get rulestate => $_getSZ(4);
  @$pb.TagNumber(419921189)
  set rulestate($core.String value) => $_setString(4, value);
  @$pb.TagNumber(419921189)
  $core.bool hasRulestate() => $_has(4);
  @$pb.TagNumber(419921189)
  void clearRulestate() => $_clearField(419921189);
}

class PutInsightRuleOutput extends $pb.GeneratedMessage {
  factory PutInsightRuleOutput() => create();

  PutInsightRuleOutput._();

  factory PutInsightRuleOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutInsightRuleOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutInsightRuleOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightRuleOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutInsightRuleOutput copyWith(void Function(PutInsightRuleOutput) updates) =>
      super.copyWith((message) => updates(message as PutInsightRuleOutput))
          as PutInsightRuleOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutInsightRuleOutput create() => PutInsightRuleOutput._();
  @$core.override
  PutInsightRuleOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutInsightRuleOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutInsightRuleOutput>(create);
  static PutInsightRuleOutput? _defaultInstance;
}

class PutManagedInsightRulesInput extends $pb.GeneratedMessage {
  factory PutManagedInsightRulesInput({
    $core.Iterable<ManagedRule>? managedrules,
  }) {
    final result = create();
    if (managedrules != null) result.managedrules.addAll(managedrules);
    return result;
  }

  PutManagedInsightRulesInput._();

  factory PutManagedInsightRulesInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutManagedInsightRulesInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutManagedInsightRulesInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<ManagedRule>(347090906, _omitFieldNames ? '' : 'managedrules',
        subBuilder: ManagedRule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedInsightRulesInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedInsightRulesInput copyWith(
          void Function(PutManagedInsightRulesInput) updates) =>
      super.copyWith(
              (message) => updates(message as PutManagedInsightRulesInput))
          as PutManagedInsightRulesInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutManagedInsightRulesInput create() =>
      PutManagedInsightRulesInput._();
  @$core.override
  PutManagedInsightRulesInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutManagedInsightRulesInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutManagedInsightRulesInput>(create);
  static PutManagedInsightRulesInput? _defaultInstance;

  @$pb.TagNumber(347090906)
  $pb.PbList<ManagedRule> get managedrules => $_getList(0);
}

class PutManagedInsightRulesOutput extends $pb.GeneratedMessage {
  factory PutManagedInsightRulesOutput({
    $core.Iterable<PartialFailure>? failures,
  }) {
    final result = create();
    if (failures != null) result.failures.addAll(failures);
    return result;
  }

  PutManagedInsightRulesOutput._();

  factory PutManagedInsightRulesOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutManagedInsightRulesOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutManagedInsightRulesOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<PartialFailure>(335467271, _omitFieldNames ? '' : 'failures',
        subBuilder: PartialFailure.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedInsightRulesOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutManagedInsightRulesOutput copyWith(
          void Function(PutManagedInsightRulesOutput) updates) =>
      super.copyWith(
              (message) => updates(message as PutManagedInsightRulesOutput))
          as PutManagedInsightRulesOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutManagedInsightRulesOutput create() =>
      PutManagedInsightRulesOutput._();
  @$core.override
  PutManagedInsightRulesOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutManagedInsightRulesOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutManagedInsightRulesOutput>(create);
  static PutManagedInsightRulesOutput? _defaultInstance;

  @$pb.TagNumber(335467271)
  $pb.PbList<PartialFailure> get failures => $_getList(0);
}

class PutMetricAlarmInput extends $pb.GeneratedMessage {
  factory PutMetricAlarmInput({
    ComparisonOperator? comparisonoperator,
    $core.String? alarmdescription,
    $core.String? thresholdmetricid,
    $core.String? alarmname,
    $core.String? evaluatelowsamplecountpercentile,
    Statistic? statistic,
    $core.String? metricname,
    $core.int? period,
    $core.bool? actionsenabled,
    $core.int? datapointstoalarm,
    StandardUnit? unit,
    $core.int? evaluationperiods,
    $core.String? treatmissingdata,
    $core.String? extendedstatistic,
    $core.String? namespace,
    $core.Iterable<$core.String>? alarmactions,
    $core.Iterable<$core.String>? okactions,
    $core.Iterable<Tag>? tags,
    $core.double? threshold,
    $core.Iterable<MetricDataQuery>? metrics,
    $core.Iterable<Dimension>? dimensions,
    $core.Iterable<$core.String>? insufficientdataactions,
  }) {
    final result = create();
    if (comparisonoperator != null)
      result.comparisonoperator = comparisonoperator;
    if (alarmdescription != null) result.alarmdescription = alarmdescription;
    if (thresholdmetricid != null) result.thresholdmetricid = thresholdmetricid;
    if (alarmname != null) result.alarmname = alarmname;
    if (evaluatelowsamplecountpercentile != null)
      result.evaluatelowsamplecountpercentile =
          evaluatelowsamplecountpercentile;
    if (statistic != null) result.statistic = statistic;
    if (metricname != null) result.metricname = metricname;
    if (period != null) result.period = period;
    if (actionsenabled != null) result.actionsenabled = actionsenabled;
    if (datapointstoalarm != null) result.datapointstoalarm = datapointstoalarm;
    if (unit != null) result.unit = unit;
    if (evaluationperiods != null) result.evaluationperiods = evaluationperiods;
    if (treatmissingdata != null) result.treatmissingdata = treatmissingdata;
    if (extendedstatistic != null) result.extendedstatistic = extendedstatistic;
    if (namespace != null) result.namespace = namespace;
    if (alarmactions != null) result.alarmactions.addAll(alarmactions);
    if (okactions != null) result.okactions.addAll(okactions);
    if (tags != null) result.tags.addAll(tags);
    if (threshold != null) result.threshold = threshold;
    if (metrics != null) result.metrics.addAll(metrics);
    if (dimensions != null) result.dimensions.addAll(dimensions);
    if (insufficientdataactions != null)
      result.insufficientdataactions.addAll(insufficientdataactions);
    return result;
  }

  PutMetricAlarmInput._();

  factory PutMetricAlarmInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMetricAlarmInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMetricAlarmInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aE<ComparisonOperator>(
        7059861, _omitFieldNames ? '' : 'comparisonoperator',
        enumValues: ComparisonOperator.values)
    ..aOS(30583613, _omitFieldNames ? '' : 'alarmdescription')
    ..aOS(41342194, _omitFieldNames ? '' : 'thresholdmetricid')
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(64366323, _omitFieldNames ? '' : 'evaluatelowsamplecountpercentile')
    ..aE<Statistic>(67293470, _omitFieldNames ? '' : 'statistic',
        enumValues: Statistic.values)
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aI(119833637, _omitFieldNames ? '' : 'period')
    ..aOB(126027878, _omitFieldNames ? '' : 'actionsenabled')
    ..aI(146378889, _omitFieldNames ? '' : 'datapointstoalarm')
    ..aE<StandardUnit>(148989480, _omitFieldNames ? '' : 'unit',
        enumValues: StandardUnit.values)
    ..aI(214794856, _omitFieldNames ? '' : 'evaluationperiods')
    ..aOS(226418150, _omitFieldNames ? '' : 'treatmissingdata')
    ..aOS(285620763, _omitFieldNames ? '' : 'extendedstatistic')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPS(355779506, _omitFieldNames ? '' : 'alarmactions')
    ..pPS(377443763, _omitFieldNames ? '' : 'okactions')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..aD(394884505, _omitFieldNames ? '' : 'threshold')
    ..pPM<MetricDataQuery>(436365847, _omitFieldNames ? '' : 'metrics',
        subBuilder: MetricDataQuery.create)
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..pPS(498450778, _omitFieldNames ? '' : 'insufficientdataactions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricAlarmInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricAlarmInput copyWith(void Function(PutMetricAlarmInput) updates) =>
      super.copyWith((message) => updates(message as PutMetricAlarmInput))
          as PutMetricAlarmInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMetricAlarmInput create() => PutMetricAlarmInput._();
  @$core.override
  PutMetricAlarmInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMetricAlarmInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMetricAlarmInput>(create);
  static PutMetricAlarmInput? _defaultInstance;

  @$pb.TagNumber(7059861)
  ComparisonOperator get comparisonoperator => $_getN(0);
  @$pb.TagNumber(7059861)
  set comparisonoperator(ComparisonOperator value) =>
      $_setField(7059861, value);
  @$pb.TagNumber(7059861)
  $core.bool hasComparisonoperator() => $_has(0);
  @$pb.TagNumber(7059861)
  void clearComparisonoperator() => $_clearField(7059861);

  @$pb.TagNumber(30583613)
  $core.String get alarmdescription => $_getSZ(1);
  @$pb.TagNumber(30583613)
  set alarmdescription($core.String value) => $_setString(1, value);
  @$pb.TagNumber(30583613)
  $core.bool hasAlarmdescription() => $_has(1);
  @$pb.TagNumber(30583613)
  void clearAlarmdescription() => $_clearField(30583613);

  @$pb.TagNumber(41342194)
  $core.String get thresholdmetricid => $_getSZ(2);
  @$pb.TagNumber(41342194)
  set thresholdmetricid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(41342194)
  $core.bool hasThresholdmetricid() => $_has(2);
  @$pb.TagNumber(41342194)
  void clearThresholdmetricid() => $_clearField(41342194);

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(3);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(3, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(3);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(64366323)
  $core.String get evaluatelowsamplecountpercentile => $_getSZ(4);
  @$pb.TagNumber(64366323)
  set evaluatelowsamplecountpercentile($core.String value) =>
      $_setString(4, value);
  @$pb.TagNumber(64366323)
  $core.bool hasEvaluatelowsamplecountpercentile() => $_has(4);
  @$pb.TagNumber(64366323)
  void clearEvaluatelowsamplecountpercentile() => $_clearField(64366323);

  @$pb.TagNumber(67293470)
  Statistic get statistic => $_getN(5);
  @$pb.TagNumber(67293470)
  set statistic(Statistic value) => $_setField(67293470, value);
  @$pb.TagNumber(67293470)
  $core.bool hasStatistic() => $_has(5);
  @$pb.TagNumber(67293470)
  void clearStatistic() => $_clearField(67293470);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(6);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(6, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(6);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(119833637)
  $core.int get period => $_getIZ(7);
  @$pb.TagNumber(119833637)
  set period($core.int value) => $_setSignedInt32(7, value);
  @$pb.TagNumber(119833637)
  $core.bool hasPeriod() => $_has(7);
  @$pb.TagNumber(119833637)
  void clearPeriod() => $_clearField(119833637);

  @$pb.TagNumber(126027878)
  $core.bool get actionsenabled => $_getBF(8);
  @$pb.TagNumber(126027878)
  set actionsenabled($core.bool value) => $_setBool(8, value);
  @$pb.TagNumber(126027878)
  $core.bool hasActionsenabled() => $_has(8);
  @$pb.TagNumber(126027878)
  void clearActionsenabled() => $_clearField(126027878);

  @$pb.TagNumber(146378889)
  $core.int get datapointstoalarm => $_getIZ(9);
  @$pb.TagNumber(146378889)
  set datapointstoalarm($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(146378889)
  $core.bool hasDatapointstoalarm() => $_has(9);
  @$pb.TagNumber(146378889)
  void clearDatapointstoalarm() => $_clearField(146378889);

  @$pb.TagNumber(148989480)
  StandardUnit get unit => $_getN(10);
  @$pb.TagNumber(148989480)
  set unit(StandardUnit value) => $_setField(148989480, value);
  @$pb.TagNumber(148989480)
  $core.bool hasUnit() => $_has(10);
  @$pb.TagNumber(148989480)
  void clearUnit() => $_clearField(148989480);

  @$pb.TagNumber(214794856)
  $core.int get evaluationperiods => $_getIZ(11);
  @$pb.TagNumber(214794856)
  set evaluationperiods($core.int value) => $_setSignedInt32(11, value);
  @$pb.TagNumber(214794856)
  $core.bool hasEvaluationperiods() => $_has(11);
  @$pb.TagNumber(214794856)
  void clearEvaluationperiods() => $_clearField(214794856);

  @$pb.TagNumber(226418150)
  $core.String get treatmissingdata => $_getSZ(12);
  @$pb.TagNumber(226418150)
  set treatmissingdata($core.String value) => $_setString(12, value);
  @$pb.TagNumber(226418150)
  $core.bool hasTreatmissingdata() => $_has(12);
  @$pb.TagNumber(226418150)
  void clearTreatmissingdata() => $_clearField(226418150);

  @$pb.TagNumber(285620763)
  $core.String get extendedstatistic => $_getSZ(13);
  @$pb.TagNumber(285620763)
  set extendedstatistic($core.String value) => $_setString(13, value);
  @$pb.TagNumber(285620763)
  $core.bool hasExtendedstatistic() => $_has(13);
  @$pb.TagNumber(285620763)
  void clearExtendedstatistic() => $_clearField(285620763);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(14);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(14, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(14);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(355779506)
  $pb.PbList<$core.String> get alarmactions => $_getList(15);

  @$pb.TagNumber(377443763)
  $pb.PbList<$core.String> get okactions => $_getList(16);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(17);

  @$pb.TagNumber(394884505)
  $core.double get threshold => $_getN(18);
  @$pb.TagNumber(394884505)
  set threshold($core.double value) => $_setDouble(18, value);
  @$pb.TagNumber(394884505)
  $core.bool hasThreshold() => $_has(18);
  @$pb.TagNumber(394884505)
  void clearThreshold() => $_clearField(394884505);

  @$pb.TagNumber(436365847)
  $pb.PbList<MetricDataQuery> get metrics => $_getList(19);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(20);

  @$pb.TagNumber(498450778)
  $pb.PbList<$core.String> get insufficientdataactions => $_getList(21);
}

class PutMetricDataInput extends $pb.GeneratedMessage {
  factory PutMetricDataInput({
    $core.Iterable<MetricDatum>? metricdata,
    $core.Iterable<EntityMetricData>? entitymetricdata,
    $core.String? namespace,
    $core.bool? strictentityvalidation,
  }) {
    final result = create();
    if (metricdata != null) result.metricdata.addAll(metricdata);
    if (entitymetricdata != null)
      result.entitymetricdata.addAll(entitymetricdata);
    if (namespace != null) result.namespace = namespace;
    if (strictentityvalidation != null)
      result.strictentityvalidation = strictentityvalidation;
    return result;
  }

  PutMetricDataInput._();

  factory PutMetricDataInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMetricDataInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMetricDataInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPM<MetricDatum>(162960562, _omitFieldNames ? '' : 'metricdata',
        subBuilder: MetricDatum.create)
    ..pPM<EntityMetricData>(
        305844415, _omitFieldNames ? '' : 'entitymetricdata',
        subBuilder: EntityMetricData.create)
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..aOB(496663153, _omitFieldNames ? '' : 'strictentityvalidation')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricDataInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricDataInput copyWith(void Function(PutMetricDataInput) updates) =>
      super.copyWith((message) => updates(message as PutMetricDataInput))
          as PutMetricDataInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMetricDataInput create() => PutMetricDataInput._();
  @$core.override
  PutMetricDataInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMetricDataInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMetricDataInput>(create);
  static PutMetricDataInput? _defaultInstance;

  @$pb.TagNumber(162960562)
  $pb.PbList<MetricDatum> get metricdata => $_getList(0);

  @$pb.TagNumber(305844415)
  $pb.PbList<EntityMetricData> get entitymetricdata => $_getList(1);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(2);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(2, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(2);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(496663153)
  $core.bool get strictentityvalidation => $_getBF(3);
  @$pb.TagNumber(496663153)
  set strictentityvalidation($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(496663153)
  $core.bool hasStrictentityvalidation() => $_has(3);
  @$pb.TagNumber(496663153)
  void clearStrictentityvalidation() => $_clearField(496663153);
}

class PutMetricStreamInput extends $pb.GeneratedMessage {
  factory PutMetricStreamInput({
    $core.String? firehosearn,
    $core.Iterable<MetricStreamFilter>? includefilters,
    $core.Iterable<MetricStreamStatisticsConfiguration>?
        statisticsconfigurations,
    $core.String? name,
    $core.String? rolearn,
    $core.bool? includelinkedaccountsmetrics,
    $core.Iterable<Tag>? tags,
    $core.Iterable<MetricStreamFilter>? excludefilters,
    MetricStreamOutputFormat? outputformat,
  }) {
    final result = create();
    if (firehosearn != null) result.firehosearn = firehosearn;
    if (includefilters != null) result.includefilters.addAll(includefilters);
    if (statisticsconfigurations != null)
      result.statisticsconfigurations.addAll(statisticsconfigurations);
    if (name != null) result.name = name;
    if (rolearn != null) result.rolearn = rolearn;
    if (includelinkedaccountsmetrics != null)
      result.includelinkedaccountsmetrics = includelinkedaccountsmetrics;
    if (tags != null) result.tags.addAll(tags);
    if (excludefilters != null) result.excludefilters.addAll(excludefilters);
    if (outputformat != null) result.outputformat = outputformat;
    return result;
  }

  PutMetricStreamInput._();

  factory PutMetricStreamInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMetricStreamInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMetricStreamInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(55494200, _omitFieldNames ? '' : 'firehosearn')
    ..pPM<MetricStreamFilter>(90967031, _omitFieldNames ? '' : 'includefilters',
        subBuilder: MetricStreamFilter.create)
    ..pPM<MetricStreamStatisticsConfiguration>(
        175876178, _omitFieldNames ? '' : 'statisticsconfigurations',
        subBuilder: MetricStreamStatisticsConfiguration.create)
    ..aOS(266367751, _omitFieldNames ? '' : 'name')
    ..aOS(322567169, _omitFieldNames ? '' : 'rolearn')
    ..aOB(356901952, _omitFieldNames ? '' : 'includelinkedaccountsmetrics')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..pPM<MetricStreamFilter>(
        426990841, _omitFieldNames ? '' : 'excludefilters',
        subBuilder: MetricStreamFilter.create)
    ..aE<MetricStreamOutputFormat>(
        519139704, _omitFieldNames ? '' : 'outputformat',
        enumValues: MetricStreamOutputFormat.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricStreamInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricStreamInput copyWith(void Function(PutMetricStreamInput) updates) =>
      super.copyWith((message) => updates(message as PutMetricStreamInput))
          as PutMetricStreamInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMetricStreamInput create() => PutMetricStreamInput._();
  @$core.override
  PutMetricStreamInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMetricStreamInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMetricStreamInput>(create);
  static PutMetricStreamInput? _defaultInstance;

  @$pb.TagNumber(55494200)
  $core.String get firehosearn => $_getSZ(0);
  @$pb.TagNumber(55494200)
  set firehosearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(55494200)
  $core.bool hasFirehosearn() => $_has(0);
  @$pb.TagNumber(55494200)
  void clearFirehosearn() => $_clearField(55494200);

  @$pb.TagNumber(90967031)
  $pb.PbList<MetricStreamFilter> get includefilters => $_getList(1);

  @$pb.TagNumber(175876178)
  $pb.PbList<MetricStreamStatisticsConfiguration>
      get statisticsconfigurations => $_getList(2);

  @$pb.TagNumber(266367751)
  $core.String get name => $_getSZ(3);
  @$pb.TagNumber(266367751)
  set name($core.String value) => $_setString(3, value);
  @$pb.TagNumber(266367751)
  $core.bool hasName() => $_has(3);
  @$pb.TagNumber(266367751)
  void clearName() => $_clearField(266367751);

  @$pb.TagNumber(322567169)
  $core.String get rolearn => $_getSZ(4);
  @$pb.TagNumber(322567169)
  set rolearn($core.String value) => $_setString(4, value);
  @$pb.TagNumber(322567169)
  $core.bool hasRolearn() => $_has(4);
  @$pb.TagNumber(322567169)
  void clearRolearn() => $_clearField(322567169);

  @$pb.TagNumber(356901952)
  $core.bool get includelinkedaccountsmetrics => $_getBF(5);
  @$pb.TagNumber(356901952)
  set includelinkedaccountsmetrics($core.bool value) => $_setBool(5, value);
  @$pb.TagNumber(356901952)
  $core.bool hasIncludelinkedaccountsmetrics() => $_has(5);
  @$pb.TagNumber(356901952)
  void clearIncludelinkedaccountsmetrics() => $_clearField(356901952);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(6);

  @$pb.TagNumber(426990841)
  $pb.PbList<MetricStreamFilter> get excludefilters => $_getList(7);

  @$pb.TagNumber(519139704)
  MetricStreamOutputFormat get outputformat => $_getN(8);
  @$pb.TagNumber(519139704)
  set outputformat(MetricStreamOutputFormat value) =>
      $_setField(519139704, value);
  @$pb.TagNumber(519139704)
  $core.bool hasOutputformat() => $_has(8);
  @$pb.TagNumber(519139704)
  void clearOutputformat() => $_clearField(519139704);
}

class PutMetricStreamOutput extends $pb.GeneratedMessage {
  factory PutMetricStreamOutput({
    $core.String? arn,
  }) {
    final result = create();
    if (arn != null) result.arn = arn;
    return result;
  }

  PutMetricStreamOutput._();

  factory PutMetricStreamOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PutMetricStreamOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PutMetricStreamOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(402345373, _omitFieldNames ? '' : 'arn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricStreamOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PutMetricStreamOutput copyWith(
          void Function(PutMetricStreamOutput) updates) =>
      super.copyWith((message) => updates(message as PutMetricStreamOutput))
          as PutMetricStreamOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PutMetricStreamOutput create() => PutMetricStreamOutput._();
  @$core.override
  PutMetricStreamOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PutMetricStreamOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PutMetricStreamOutput>(create);
  static PutMetricStreamOutput? _defaultInstance;

  @$pb.TagNumber(402345373)
  $core.String get arn => $_getSZ(0);
  @$pb.TagNumber(402345373)
  set arn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(402345373)
  $core.bool hasArn() => $_has(0);
  @$pb.TagNumber(402345373)
  void clearArn() => $_clearField(402345373);
}

class Range extends $pb.GeneratedMessage {
  factory Range({
    $core.String? endtime,
    $core.String? starttime,
  }) {
    final result = create();
    if (endtime != null) result.endtime = endtime;
    if (starttime != null) result.starttime = starttime;
    return result;
  }

  Range._();

  factory Range.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Range.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Range',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(63911884, _omitFieldNames ? '' : 'endtime')
    ..aOS(370760303, _omitFieldNames ? '' : 'starttime')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Range clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Range copyWith(void Function(Range) updates) =>
      super.copyWith((message) => updates(message as Range)) as Range;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Range create() => Range._();
  @$core.override
  Range createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Range getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Range>(create);
  static Range? _defaultInstance;

  @$pb.TagNumber(63911884)
  $core.String get endtime => $_getSZ(0);
  @$pb.TagNumber(63911884)
  set endtime($core.String value) => $_setString(0, value);
  @$pb.TagNumber(63911884)
  $core.bool hasEndtime() => $_has(0);
  @$pb.TagNumber(63911884)
  void clearEndtime() => $_clearField(63911884);

  @$pb.TagNumber(370760303)
  $core.String get starttime => $_getSZ(1);
  @$pb.TagNumber(370760303)
  set starttime($core.String value) => $_setString(1, value);
  @$pb.TagNumber(370760303)
  $core.bool hasStarttime() => $_has(1);
  @$pb.TagNumber(370760303)
  void clearStarttime() => $_clearField(370760303);
}

class ResourceNotFound extends $pb.GeneratedMessage {
  factory ResourceNotFound({
    $core.String? message,
  }) {
    final result = create();
    if (message != null) result.message = message;
    return result;
  }

  ResourceNotFound._();

  factory ResourceNotFound.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceNotFound.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceNotFound',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(82970853, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFound clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFound copyWith(void Function(ResourceNotFound) updates) =>
      super.copyWith((message) => updates(message as ResourceNotFound))
          as ResourceNotFound;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceNotFound create() => ResourceNotFound._();
  @$core.override
  ResourceNotFound createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceNotFound getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceNotFound>(create);
  static ResourceNotFound? _defaultInstance;

  @$pb.TagNumber(82970853)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(82970853)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(82970853)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(82970853)
  void clearMessage() => $_clearField(82970853);
}

class ResourceNotFoundException extends $pb.GeneratedMessage {
  factory ResourceNotFoundException({
    $core.String? message,
    $core.String? resourcetype,
    $core.String? resourceid,
  }) {
    final result = create();
    if (message != null) result.message = message;
    if (resourcetype != null) result.resourcetype = resourcetype;
    if (resourceid != null) result.resourceid = resourceid;
    return result;
  }

  ResourceNotFoundException._();

  factory ResourceNotFoundException.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ResourceNotFoundException.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ResourceNotFoundException',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(235854213, _omitFieldNames ? '' : 'message')
    ..aOS(301342558, _omitFieldNames ? '' : 'resourcetype')
    ..aOS(526146833, _omitFieldNames ? '' : 'resourceid')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFoundException clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ResourceNotFoundException copyWith(
          void Function(ResourceNotFoundException) updates) =>
      super.copyWith((message) => updates(message as ResourceNotFoundException))
          as ResourceNotFoundException;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ResourceNotFoundException create() => ResourceNotFoundException._();
  @$core.override
  ResourceNotFoundException createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ResourceNotFoundException getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ResourceNotFoundException>(create);
  static ResourceNotFoundException? _defaultInstance;

  @$pb.TagNumber(235854213)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(235854213)
  set message($core.String value) => $_setString(0, value);
  @$pb.TagNumber(235854213)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(235854213)
  void clearMessage() => $_clearField(235854213);

  @$pb.TagNumber(301342558)
  $core.String get resourcetype => $_getSZ(1);
  @$pb.TagNumber(301342558)
  set resourcetype($core.String value) => $_setString(1, value);
  @$pb.TagNumber(301342558)
  $core.bool hasResourcetype() => $_has(1);
  @$pb.TagNumber(301342558)
  void clearResourcetype() => $_clearField(301342558);

  @$pb.TagNumber(526146833)
  $core.String get resourceid => $_getSZ(2);
  @$pb.TagNumber(526146833)
  set resourceid($core.String value) => $_setString(2, value);
  @$pb.TagNumber(526146833)
  $core.bool hasResourceid() => $_has(2);
  @$pb.TagNumber(526146833)
  void clearResourceid() => $_clearField(526146833);
}

class Rule extends $pb.GeneratedMessage {
  factory Rule({
    Schedule? schedule,
  }) {
    final result = create();
    if (schedule != null) result.schedule = schedule;
    return result;
  }

  Rule._();

  factory Rule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Rule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Rule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOM<Schedule>(66697965, _omitFieldNames ? '' : 'schedule',
        subBuilder: Schedule.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Rule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Rule copyWith(void Function(Rule) updates) =>
      super.copyWith((message) => updates(message as Rule)) as Rule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Rule create() => Rule._();
  @$core.override
  Rule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Rule getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Rule>(create);
  static Rule? _defaultInstance;

  @$pb.TagNumber(66697965)
  Schedule get schedule => $_getN(0);
  @$pb.TagNumber(66697965)
  set schedule(Schedule value) => $_setField(66697965, value);
  @$pb.TagNumber(66697965)
  $core.bool hasSchedule() => $_has(0);
  @$pb.TagNumber(66697965)
  void clearSchedule() => $_clearField(66697965);
  @$pb.TagNumber(66697965)
  Schedule ensureSchedule() => $_ensure(0);
}

class Schedule extends $pb.GeneratedMessage {
  factory Schedule({
    $core.String? expression,
    $core.String? timezone,
    $core.String? duration,
  }) {
    final result = create();
    if (expression != null) result.expression = expression;
    if (timezone != null) result.timezone = timezone;
    if (duration != null) result.duration = duration;
    return result;
  }

  Schedule._();

  factory Schedule.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Schedule.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Schedule',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(193051916, _omitFieldNames ? '' : 'expression')
    ..aOS(246302531, _omitFieldNames ? '' : 'timezone')
    ..aOS(348604718, _omitFieldNames ? '' : 'duration')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Schedule clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Schedule copyWith(void Function(Schedule) updates) =>
      super.copyWith((message) => updates(message as Schedule)) as Schedule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Schedule create() => Schedule._();
  @$core.override
  Schedule createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Schedule getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Schedule>(create);
  static Schedule? _defaultInstance;

  @$pb.TagNumber(193051916)
  $core.String get expression => $_getSZ(0);
  @$pb.TagNumber(193051916)
  set expression($core.String value) => $_setString(0, value);
  @$pb.TagNumber(193051916)
  $core.bool hasExpression() => $_has(0);
  @$pb.TagNumber(193051916)
  void clearExpression() => $_clearField(193051916);

  @$pb.TagNumber(246302531)
  $core.String get timezone => $_getSZ(1);
  @$pb.TagNumber(246302531)
  set timezone($core.String value) => $_setString(1, value);
  @$pb.TagNumber(246302531)
  $core.bool hasTimezone() => $_has(1);
  @$pb.TagNumber(246302531)
  void clearTimezone() => $_clearField(246302531);

  @$pb.TagNumber(348604718)
  $core.String get duration => $_getSZ(2);
  @$pb.TagNumber(348604718)
  set duration($core.String value) => $_setString(2, value);
  @$pb.TagNumber(348604718)
  $core.bool hasDuration() => $_has(2);
  @$pb.TagNumber(348604718)
  void clearDuration() => $_clearField(348604718);
}

class SetAlarmStateInput extends $pb.GeneratedMessage {
  factory SetAlarmStateInput({
    $core.String? alarmname,
    $core.String? statereasondata,
    StateValue? statevalue,
    $core.String? statereason,
  }) {
    final result = create();
    if (alarmname != null) result.alarmname = alarmname;
    if (statereasondata != null) result.statereasondata = statereasondata;
    if (statevalue != null) result.statevalue = statevalue;
    if (statereason != null) result.statereason = statereason;
    return result;
  }

  SetAlarmStateInput._();

  factory SetAlarmStateInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetAlarmStateInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetAlarmStateInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(54095842, _omitFieldNames ? '' : 'alarmname')
    ..aOS(262771075, _omitFieldNames ? '' : 'statereasondata')
    ..aE<StateValue>(334526008, _omitFieldNames ? '' : 'statevalue',
        enumValues: StateValue.values)
    ..aOS(376138483, _omitFieldNames ? '' : 'statereason')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAlarmStateInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAlarmStateInput copyWith(void Function(SetAlarmStateInput) updates) =>
      super.copyWith((message) => updates(message as SetAlarmStateInput))
          as SetAlarmStateInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetAlarmStateInput create() => SetAlarmStateInput._();
  @$core.override
  SetAlarmStateInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetAlarmStateInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetAlarmStateInput>(create);
  static SetAlarmStateInput? _defaultInstance;

  @$pb.TagNumber(54095842)
  $core.String get alarmname => $_getSZ(0);
  @$pb.TagNumber(54095842)
  set alarmname($core.String value) => $_setString(0, value);
  @$pb.TagNumber(54095842)
  $core.bool hasAlarmname() => $_has(0);
  @$pb.TagNumber(54095842)
  void clearAlarmname() => $_clearField(54095842);

  @$pb.TagNumber(262771075)
  $core.String get statereasondata => $_getSZ(1);
  @$pb.TagNumber(262771075)
  set statereasondata($core.String value) => $_setString(1, value);
  @$pb.TagNumber(262771075)
  $core.bool hasStatereasondata() => $_has(1);
  @$pb.TagNumber(262771075)
  void clearStatereasondata() => $_clearField(262771075);

  @$pb.TagNumber(334526008)
  StateValue get statevalue => $_getN(2);
  @$pb.TagNumber(334526008)
  set statevalue(StateValue value) => $_setField(334526008, value);
  @$pb.TagNumber(334526008)
  $core.bool hasStatevalue() => $_has(2);
  @$pb.TagNumber(334526008)
  void clearStatevalue() => $_clearField(334526008);

  @$pb.TagNumber(376138483)
  $core.String get statereason => $_getSZ(3);
  @$pb.TagNumber(376138483)
  set statereason($core.String value) => $_setString(3, value);
  @$pb.TagNumber(376138483)
  $core.bool hasStatereason() => $_has(3);
  @$pb.TagNumber(376138483)
  void clearStatereason() => $_clearField(376138483);
}

class SingleMetricAnomalyDetector extends $pb.GeneratedMessage {
  factory SingleMetricAnomalyDetector({
    $core.String? accountid,
    $core.String? metricname,
    $core.String? stat,
    $core.String? namespace,
    $core.Iterable<Dimension>? dimensions,
  }) {
    final result = create();
    if (accountid != null) result.accountid = accountid;
    if (metricname != null) result.metricname = metricname;
    if (stat != null) result.stat = stat;
    if (namespace != null) result.namespace = namespace;
    if (dimensions != null) result.dimensions.addAll(dimensions);
    return result;
  }

  SingleMetricAnomalyDetector._();

  factory SingleMetricAnomalyDetector.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SingleMetricAnomalyDetector.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SingleMetricAnomalyDetector',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(65954002, _omitFieldNames ? '' : 'accountid')
    ..aOS(106340219, _omitFieldNames ? '' : 'metricname')
    ..aOS(325549752, _omitFieldNames ? '' : 'stat')
    ..aOS(355353153, _omitFieldNames ? '' : 'namespace')
    ..pPM<Dimension>(462933457, _omitFieldNames ? '' : 'dimensions',
        subBuilder: Dimension.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleMetricAnomalyDetector clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SingleMetricAnomalyDetector copyWith(
          void Function(SingleMetricAnomalyDetector) updates) =>
      super.copyWith(
              (message) => updates(message as SingleMetricAnomalyDetector))
          as SingleMetricAnomalyDetector;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SingleMetricAnomalyDetector create() =>
      SingleMetricAnomalyDetector._();
  @$core.override
  SingleMetricAnomalyDetector createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SingleMetricAnomalyDetector getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SingleMetricAnomalyDetector>(create);
  static SingleMetricAnomalyDetector? _defaultInstance;

  @$pb.TagNumber(65954002)
  $core.String get accountid => $_getSZ(0);
  @$pb.TagNumber(65954002)
  set accountid($core.String value) => $_setString(0, value);
  @$pb.TagNumber(65954002)
  $core.bool hasAccountid() => $_has(0);
  @$pb.TagNumber(65954002)
  void clearAccountid() => $_clearField(65954002);

  @$pb.TagNumber(106340219)
  $core.String get metricname => $_getSZ(1);
  @$pb.TagNumber(106340219)
  set metricname($core.String value) => $_setString(1, value);
  @$pb.TagNumber(106340219)
  $core.bool hasMetricname() => $_has(1);
  @$pb.TagNumber(106340219)
  void clearMetricname() => $_clearField(106340219);

  @$pb.TagNumber(325549752)
  $core.String get stat => $_getSZ(2);
  @$pb.TagNumber(325549752)
  set stat($core.String value) => $_setString(2, value);
  @$pb.TagNumber(325549752)
  $core.bool hasStat() => $_has(2);
  @$pb.TagNumber(325549752)
  void clearStat() => $_clearField(325549752);

  @$pb.TagNumber(355353153)
  $core.String get namespace => $_getSZ(3);
  @$pb.TagNumber(355353153)
  set namespace($core.String value) => $_setString(3, value);
  @$pb.TagNumber(355353153)
  $core.bool hasNamespace() => $_has(3);
  @$pb.TagNumber(355353153)
  void clearNamespace() => $_clearField(355353153);

  @$pb.TagNumber(462933457)
  $pb.PbList<Dimension> get dimensions => $_getList(4);
}

class StartMetricStreamsInput extends $pb.GeneratedMessage {
  factory StartMetricStreamsInput({
    $core.Iterable<$core.String>? names,
  }) {
    final result = create();
    if (names != null) result.names.addAll(names);
    return result;
  }

  StartMetricStreamsInput._();

  factory StartMetricStreamsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartMetricStreamsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartMetricStreamsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(324387120, _omitFieldNames ? '' : 'names')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMetricStreamsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMetricStreamsInput copyWith(
          void Function(StartMetricStreamsInput) updates) =>
      super.copyWith((message) => updates(message as StartMetricStreamsInput))
          as StartMetricStreamsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartMetricStreamsInput create() => StartMetricStreamsInput._();
  @$core.override
  StartMetricStreamsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartMetricStreamsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartMetricStreamsInput>(create);
  static StartMetricStreamsInput? _defaultInstance;

  @$pb.TagNumber(324387120)
  $pb.PbList<$core.String> get names => $_getList(0);
}

class StartMetricStreamsOutput extends $pb.GeneratedMessage {
  factory StartMetricStreamsOutput() => create();

  StartMetricStreamsOutput._();

  factory StartMetricStreamsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StartMetricStreamsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StartMetricStreamsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMetricStreamsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StartMetricStreamsOutput copyWith(
          void Function(StartMetricStreamsOutput) updates) =>
      super.copyWith((message) => updates(message as StartMetricStreamsOutput))
          as StartMetricStreamsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartMetricStreamsOutput create() => StartMetricStreamsOutput._();
  @$core.override
  StartMetricStreamsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StartMetricStreamsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StartMetricStreamsOutput>(create);
  static StartMetricStreamsOutput? _defaultInstance;
}

class StatisticSet extends $pb.GeneratedMessage {
  factory StatisticSet({
    $core.double? maximum,
    $core.double? samplecount,
    $core.double? minimum,
    $core.double? sum,
  }) {
    final result = create();
    if (maximum != null) result.maximum = maximum;
    if (samplecount != null) result.samplecount = samplecount;
    if (minimum != null) result.minimum = minimum;
    if (sum != null) result.sum = sum;
    return result;
  }

  StatisticSet._();

  factory StatisticSet.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StatisticSet.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StatisticSet',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aD(43343394, _omitFieldNames ? '' : 'maximum')
    ..aD(384919839, _omitFieldNames ? '' : 'samplecount')
    ..aD(438536856, _omitFieldNames ? '' : 'minimum')
    ..aD(534303305, _omitFieldNames ? '' : 'sum')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StatisticSet clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StatisticSet copyWith(void Function(StatisticSet) updates) =>
      super.copyWith((message) => updates(message as StatisticSet))
          as StatisticSet;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StatisticSet create() => StatisticSet._();
  @$core.override
  StatisticSet createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StatisticSet getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StatisticSet>(create);
  static StatisticSet? _defaultInstance;

  @$pb.TagNumber(43343394)
  $core.double get maximum => $_getN(0);
  @$pb.TagNumber(43343394)
  set maximum($core.double value) => $_setDouble(0, value);
  @$pb.TagNumber(43343394)
  $core.bool hasMaximum() => $_has(0);
  @$pb.TagNumber(43343394)
  void clearMaximum() => $_clearField(43343394);

  @$pb.TagNumber(384919839)
  $core.double get samplecount => $_getN(1);
  @$pb.TagNumber(384919839)
  set samplecount($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(384919839)
  $core.bool hasSamplecount() => $_has(1);
  @$pb.TagNumber(384919839)
  void clearSamplecount() => $_clearField(384919839);

  @$pb.TagNumber(438536856)
  $core.double get minimum => $_getN(2);
  @$pb.TagNumber(438536856)
  set minimum($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(438536856)
  $core.bool hasMinimum() => $_has(2);
  @$pb.TagNumber(438536856)
  void clearMinimum() => $_clearField(438536856);

  @$pb.TagNumber(534303305)
  $core.double get sum => $_getN(3);
  @$pb.TagNumber(534303305)
  set sum($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(534303305)
  $core.bool hasSum() => $_has(3);
  @$pb.TagNumber(534303305)
  void clearSum() => $_clearField(534303305);
}

class StopMetricStreamsInput extends $pb.GeneratedMessage {
  factory StopMetricStreamsInput({
    $core.Iterable<$core.String>? names,
  }) {
    final result = create();
    if (names != null) result.names.addAll(names);
    return result;
  }

  StopMetricStreamsInput._();

  factory StopMetricStreamsInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopMetricStreamsInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopMetricStreamsInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(324387120, _omitFieldNames ? '' : 'names')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopMetricStreamsInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopMetricStreamsInput copyWith(
          void Function(StopMetricStreamsInput) updates) =>
      super.copyWith((message) => updates(message as StopMetricStreamsInput))
          as StopMetricStreamsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopMetricStreamsInput create() => StopMetricStreamsInput._();
  @$core.override
  StopMetricStreamsInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopMetricStreamsInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopMetricStreamsInput>(create);
  static StopMetricStreamsInput? _defaultInstance;

  @$pb.TagNumber(324387120)
  $pb.PbList<$core.String> get names => $_getList(0);
}

class StopMetricStreamsOutput extends $pb.GeneratedMessage {
  factory StopMetricStreamsOutput() => create();

  StopMetricStreamsOutput._();

  factory StopMetricStreamsOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory StopMetricStreamsOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'StopMetricStreamsOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopMetricStreamsOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  StopMetricStreamsOutput copyWith(
          void Function(StopMetricStreamsOutput) updates) =>
      super.copyWith((message) => updates(message as StopMetricStreamsOutput))
          as StopMetricStreamsOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StopMetricStreamsOutput create() => StopMetricStreamsOutput._();
  @$core.override
  StopMetricStreamsOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static StopMetricStreamsOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<StopMetricStreamsOutput>(create);
  static StopMetricStreamsOutput? _defaultInstance;
}

class Tag extends $pb.GeneratedMessage {
  factory Tag({
    $core.String? key,
    $core.String? value,
  }) {
    final result = create();
    if (key != null) result.key = key;
    if (value != null) result.value = value;
    return result;
  }

  Tag._();

  factory Tag.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Tag.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Tag',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(219859213, _omitFieldNames ? '' : 'key')
    ..aOS(289929579, _omitFieldNames ? '' : 'value')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tag clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Tag copyWith(void Function(Tag) updates) =>
      super.copyWith((message) => updates(message as Tag)) as Tag;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Tag create() => Tag._();
  @$core.override
  Tag createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Tag getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Tag>(create);
  static Tag? _defaultInstance;

  @$pb.TagNumber(219859213)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(219859213)
  set key($core.String value) => $_setString(0, value);
  @$pb.TagNumber(219859213)
  $core.bool hasKey() => $_has(0);
  @$pb.TagNumber(219859213)
  void clearKey() => $_clearField(219859213);

  @$pb.TagNumber(289929579)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(289929579)
  set value($core.String value) => $_setString(1, value);
  @$pb.TagNumber(289929579)
  $core.bool hasValue() => $_has(1);
  @$pb.TagNumber(289929579)
  void clearValue() => $_clearField(289929579);
}

class TagResourceInput extends $pb.GeneratedMessage {
  factory TagResourceInput({
    $core.String? resourcearn,
    $core.Iterable<Tag>? tags,
  }) {
    final result = create();
    if (resourcearn != null) result.resourcearn = resourcearn;
    if (tags != null) result.tags.addAll(tags);
    return result;
  }

  TagResourceInput._();

  factory TagResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..pPM<Tag>(381526209, _omitFieldNames ? '' : 'tags', subBuilder: Tag.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceInput copyWith(void Function(TagResourceInput) updates) =>
      super.copyWith((message) => updates(message as TagResourceInput))
          as TagResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceInput create() => TagResourceInput._();
  @$core.override
  TagResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceInput>(create);
  static TagResourceInput? _defaultInstance;

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(0);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(0, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(0);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);

  @$pb.TagNumber(381526209)
  $pb.PbList<Tag> get tags => $_getList(1);
}

class TagResourceOutput extends $pb.GeneratedMessage {
  factory TagResourceOutput() => create();

  TagResourceOutput._();

  factory TagResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TagResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TagResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TagResourceOutput copyWith(void Function(TagResourceOutput) updates) =>
      super.copyWith((message) => updates(message as TagResourceOutput))
          as TagResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TagResourceOutput create() => TagResourceOutput._();
  @$core.override
  TagResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TagResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TagResourceOutput>(create);
  static TagResourceOutput? _defaultInstance;
}

class UntagResourceInput extends $pb.GeneratedMessage {
  factory UntagResourceInput({
    $core.Iterable<$core.String>? tagkeys,
    $core.String? resourcearn,
  }) {
    final result = create();
    if (tagkeys != null) result.tagkeys.addAll(tagkeys);
    if (resourcearn != null) result.resourcearn = resourcearn;
    return result;
  }

  UntagResourceInput._();

  factory UntagResourceInput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceInput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceInput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..pPS(320659964, _omitFieldNames ? '' : 'tagkeys')
    ..aOS(369516653, _omitFieldNames ? '' : 'resourcearn')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceInput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceInput copyWith(void Function(UntagResourceInput) updates) =>
      super.copyWith((message) => updates(message as UntagResourceInput))
          as UntagResourceInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceInput create() => UntagResourceInput._();
  @$core.override
  UntagResourceInput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceInput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceInput>(create);
  static UntagResourceInput? _defaultInstance;

  @$pb.TagNumber(320659964)
  $pb.PbList<$core.String> get tagkeys => $_getList(0);

  @$pb.TagNumber(369516653)
  $core.String get resourcearn => $_getSZ(1);
  @$pb.TagNumber(369516653)
  set resourcearn($core.String value) => $_setString(1, value);
  @$pb.TagNumber(369516653)
  $core.bool hasResourcearn() => $_has(1);
  @$pb.TagNumber(369516653)
  void clearResourcearn() => $_clearField(369516653);
}

class UntagResourceOutput extends $pb.GeneratedMessage {
  factory UntagResourceOutput() => create();

  UntagResourceOutput._();

  factory UntagResourceOutput.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UntagResourceOutput.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UntagResourceOutput',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'monitoring'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceOutput clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UntagResourceOutput copyWith(void Function(UntagResourceOutput) updates) =>
      super.copyWith((message) => updates(message as UntagResourceOutput))
          as UntagResourceOutput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UntagResourceOutput create() => UntagResourceOutput._();
  @$core.override
  UntagResourceOutput createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UntagResourceOutput getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UntagResourceOutput>(create);
  static UntagResourceOutput? _defaultInstance;
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
