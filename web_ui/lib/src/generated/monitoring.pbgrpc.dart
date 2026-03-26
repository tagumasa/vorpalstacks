// This is a generated file - do not edit.
//
// Generated from monitoring.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'common.pb.dart' as $1;
import 'monitoring.pb.dart' as $0;

export 'monitoring.pb.dart';

/// CloudWatchService provides monitoring API operations.
@$pb.GrpcServiceName('monitoring.CloudWatchService')
class CloudWatchServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CloudWatchServiceClient(super.channel, {super.options, super.interceptors});

  /// Deletes a specific alarm mute rule. When you delete a mute rule, any alarms that are currently being muted by that rule are immediately unmuted. If those alarms are in an ALARM state, their configu...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteAlarmMuteRule(
    $0.DeleteAlarmMuteRuleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAlarmMuteRule, request, options: options);
  }

  /// Deletes the specified alarms. You can delete up to 100 alarms in one operation. However, this total can include no more than one composite alarm. For example, you could delete 99 metric alarms and ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteAlarms(
    $0.DeleteAlarmsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAlarms, request, options: options);
  }

  /// Deletes the specified anomaly detection model from your account. For more information about how to delete an anomaly detection model, see Deleting an anomaly detection model in the CloudWatch User ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteAnomalyDetectorOutput> deleteAnomalyDetector(
    $0.DeleteAnomalyDetectorInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAnomalyDetector, request, options: options);
  }

  /// Deletes all dashboards that you specify. You can specify up to 100 dashboards to delete. If there is an error during this call, no dashboards are deleted.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteDashboardsOutput> deleteDashboards(
    $0.DeleteDashboardsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDashboards, request, options: options);
  }

  /// Permanently deletes the specified Contributor Insights rules. If you create a rule, delete it, and then re-create it with the same name, historical data from the first time the rule was created mig...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteInsightRulesOutput> deleteInsightRules(
    $0.DeleteInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteInsightRules, request, options: options);
  }

  /// Permanently deletes the metric stream that you specify.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteMetricStreamOutput> deleteMetricStream(
    $0.DeleteMetricStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMetricStream, request, options: options);
  }

  /// Returns the information of the current alarm contributors that are in ALARM state. This operation returns details about the individual time series that contribute to the alarm's state.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAlarmContributorsOutput>
      describeAlarmContributors(
    $0.DescribeAlarmContributorsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAlarmContributors, request,
        options: options);
  }

  /// Retrieves the history for the specified alarm. You can filter the results by date range or item type. If an alarm name is not specified, the histories for either all metric alarms or all composite ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAlarmHistoryOutput> describeAlarmHistory(
    $0.DescribeAlarmHistoryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAlarmHistory, request, options: options);
  }

  /// Retrieves the specified alarms. You can filter the results by specifying a prefix for the alarm name, the alarm state, or a prefix for any action. To use this operation and return information about...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAlarmsOutput> describeAlarms(
    $0.DescribeAlarmsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAlarms, request, options: options);
  }

  /// Retrieves the alarms for the specified metric. To filter the results, specify a statistic, period, or unit. This operation retrieves only standard alarms that are based on the specified metric. It ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAlarmsForMetricOutput>
      describeAlarmsForMetric(
    $0.DescribeAlarmsForMetricInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAlarmsForMetric, request,
        options: options);
  }

  /// Lists the anomaly detection models that you have created in your account. For single metric anomaly detectors, you can list all of the models in your account or filter the results to only the model...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeAnomalyDetectorsOutput>
      describeAnomalyDetectors(
    $0.DescribeAnomalyDetectorsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAnomalyDetectors, request,
        options: options);
  }

  /// Returns a list of all the Contributor Insights rules in your account. For more information about Contributor Insights, see Using Contributor Insights to Analyze High-Cardinality Data.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeInsightRulesOutput> describeInsightRules(
    $0.DescribeInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInsightRules, request, options: options);
  }

  /// Disables the actions for the specified alarms. When an alarm's actions are disabled, the alarm actions do not execute when the alarm state changes.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> disableAlarmActions(
    $0.DisableAlarmActionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableAlarmActions, request, options: options);
  }

  /// Disables the specified Contributor Insights rules. When rules are disabled, they do not analyze log groups and do not incur costs.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DisableInsightRulesOutput> disableInsightRules(
    $0.DisableInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableInsightRules, request, options: options);
  }

  /// Enables the actions for the specified alarms.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> enableAlarmActions(
    $0.EnableAlarmActionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableAlarmActions, request, options: options);
  }

  /// Enables the specified Contributor Insights rules. When rules are enabled, they immediately begin analyzing log data.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.EnableInsightRulesOutput> enableInsightRules(
    $0.EnableInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableInsightRules, request, options: options);
  }

  /// Retrieves details for a specific alarm mute rule. This operation returns complete information about the mute rule, including its configuration, status, targeted alarms, and metadata. The returned s...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetAlarmMuteRuleOutput> getAlarmMuteRule(
    $0.GetAlarmMuteRuleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAlarmMuteRule, request, options: options);
  }

  /// Displays the details of the dashboard that you specify. To copy an existing dashboard, use GetDashboard, and then use the data returned within DashboardBody as the template for the new dashboard wh...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetDashboardOutput> getDashboard(
    $0.GetDashboardInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDashboard, request, options: options);
  }

  /// This operation returns the time series data collected by a Contributor Insights rule. The data includes the identity and number of contributors to the log group. You can also optionally return one ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetInsightRuleReportOutput> getInsightRuleReport(
    $0.GetInsightRuleReportInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInsightRuleReport, request, options: options);
  }

  /// You can use the GetMetricData API to retrieve CloudWatch metric values. The operation can also include a CloudWatch Metrics Insights query, and one or more metric math functions. A GetMetricData op...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetMetricDataOutput> getMetricData(
    $0.GetMetricDataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMetricData, request, options: options);
  }

  /// Gets statistics for the specified metric. The maximum number of data points returned from a single call is 1,440. If you request more than 1,440 data points, CloudWatch returns an error. To reduce ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetMetricStatisticsOutput> getMetricStatistics(
    $0.GetMetricStatisticsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMetricStatistics, request, options: options);
  }

  /// Returns information about the metric stream that you specify.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetMetricStreamOutput> getMetricStream(
    $0.GetMetricStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMetricStream, request, options: options);
  }

  /// You can use the GetMetricWidgetImage API to retrieve a snapshot graph of one or more Amazon CloudWatch metrics as a bitmap image. You can then embed this image into your services and products, such...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetMetricWidgetImageOutput> getMetricWidgetImage(
    $0.GetMetricWidgetImageInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMetricWidgetImage, request, options: options);
  }

  /// Lists alarm mute rules in your Amazon Web Services account and region. You can filter the results by alarm name to find all mute rules targeting a specific alarm, or by status to find rules that ar...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListAlarmMuteRulesOutput> listAlarmMuteRules(
    $0.ListAlarmMuteRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAlarmMuteRules, request, options: options);
  }

  /// Returns a list of the dashboards for your account. If you include DashboardNamePrefix, only those dashboards with names starting with the prefix are listed. Otherwise, all dashboards in your accoun...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListDashboardsOutput> listDashboards(
    $0.ListDashboardsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDashboards, request, options: options);
  }

  /// Returns a list that contains the number of managed Contributor Insights rules in your account.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListManagedInsightRulesOutput>
      listManagedInsightRules(
    $0.ListManagedInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listManagedInsightRules, request,
        options: options);
  }

  /// List the specified metrics. You can use the returned metrics with GetMetricData or GetMetricStatistics to get statistical data. Up to 500 results are returned for any one call. To retrieve addition...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListMetricsOutput> listMetrics(
    $0.ListMetricsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMetrics, request, options: options);
  }

  /// Returns a list of metric streams in this account.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListMetricStreamsOutput> listMetricStreams(
    $0.ListMetricStreamsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMetricStreams, request, options: options);
  }

  /// Displays the tags associated with a CloudWatch resource. Currently, alarms and Contributor Insights rules support tagging.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTagsForResourceOutput> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Creates or updates an alarm mute rule. Alarm mute rules automatically mute alarm actions during predefined time windows. When a mute rule is active, targeted alarms continue to evaluate metrics and...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> putAlarmMuteRule(
    $0.PutAlarmMuteRuleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAlarmMuteRule, request, options: options);
  }

  /// Creates an anomaly detection model for a CloudWatch metric. You can use the model to display a band of expected normal values when the metric is graphed. If you have enabled unified cross-account o...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutAnomalyDetectorOutput> putAnomalyDetector(
    $0.PutAnomalyDetectorInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAnomalyDetector, request, options: options);
  }

  /// Creates or updates a composite alarm. When you create a composite alarm, you specify a rule expression for the alarm that takes into account the alarm states of other alarms that you have created. ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> putCompositeAlarm(
    $0.PutCompositeAlarmInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putCompositeAlarm, request, options: options);
  }

  /// Creates a dashboard if it does not already exist, or updates an existing dashboard. If you update a dashboard, the entire contents are replaced with what you specify here. All dashboards in your ac...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutDashboardOutput> putDashboard(
    $0.PutDashboardInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDashboard, request, options: options);
  }

  /// Creates a Contributor Insights rule. Rules evaluate log events in a CloudWatch Logs log group, enabling you to find contributor data for the log events in that log group. For more information, see ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutInsightRuleOutput> putInsightRule(
    $0.PutInsightRuleInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putInsightRule, request, options: options);
  }

  /// Creates a managed Contributor Insights rule for a specified Amazon Web Services resource. When you enable a managed rule, you create a Contributor Insights rule that collects data from Amazon Web S...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutManagedInsightRulesOutput> putManagedInsightRules(
    $0.PutManagedInsightRulesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putManagedInsightRules, request,
        options: options);
  }

  /// Creates or updates an alarm and associates it with the specified metric, metric math expression, anomaly detection model, or Metrics Insights query. For more information about using a Metrics Insig...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> putMetricAlarm(
    $0.PutMetricAlarmInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMetricAlarm, request, options: options);
  }

  /// Publishes metric data to Amazon CloudWatch. CloudWatch associates the data with the specified metric. If the specified metric does not exist, CloudWatch creates the metric. When CloudWatch creates ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> putMetricData(
    $0.PutMetricDataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMetricData, request, options: options);
  }

  /// Creates or updates a metric stream. Metric streams can automatically stream CloudWatch metrics to Amazon Web Services destinations, including Amazon S3, and to many third-party solutions. For more ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutMetricStreamOutput> putMetricStream(
    $0.PutMetricStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMetricStream, request, options: options);
  }

  /// Temporarily sets the state of an alarm for testing purposes. When the updated state differs from the previous value, the action configured for the appropriate state is invoked. For example, if your...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> setAlarmState(
    $0.SetAlarmStateInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setAlarmState, request, options: options);
  }

  /// Starts the streaming of metrics for one or more of your metric streams.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StartMetricStreamsOutput> startMetricStreams(
    $0.StartMetricStreamsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startMetricStreams, request, options: options);
  }

  /// Stops the streaming of metrics for one or more of your metric streams.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StopMetricStreamsOutput> stopMetricStreams(
    $0.StopMetricStreamsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopMetricStreams, request, options: options);
  }

  /// Assigns one or more tags (key-value pairs) to the specified CloudWatch resource. Currently, the only CloudWatch resources that can be tagged are alarms and Contributor Insights rules. Tags can help...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TagResourceOutput> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes one or more tags from the specified resource.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UntagResourceOutput> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  // method descriptors

  static final _$deleteAlarmMuteRule =
      $grpc.ClientMethod<$0.DeleteAlarmMuteRuleInput, $1.Empty>(
          '/monitoring.CloudWatchService/DeleteAlarmMuteRule',
          ($0.DeleteAlarmMuteRuleInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAlarms =
      $grpc.ClientMethod<$0.DeleteAlarmsInput, $1.Empty>(
          '/monitoring.CloudWatchService/DeleteAlarms',
          ($0.DeleteAlarmsInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAnomalyDetector = $grpc.ClientMethod<
          $0.DeleteAnomalyDetectorInput, $0.DeleteAnomalyDetectorOutput>(
      '/monitoring.CloudWatchService/DeleteAnomalyDetector',
      ($0.DeleteAnomalyDetectorInput value) => value.writeToBuffer(),
      $0.DeleteAnomalyDetectorOutput.fromBuffer);
  static final _$deleteDashboards =
      $grpc.ClientMethod<$0.DeleteDashboardsInput, $0.DeleteDashboardsOutput>(
          '/monitoring.CloudWatchService/DeleteDashboards',
          ($0.DeleteDashboardsInput value) => value.writeToBuffer(),
          $0.DeleteDashboardsOutput.fromBuffer);
  static final _$deleteInsightRules = $grpc.ClientMethod<
          $0.DeleteInsightRulesInput, $0.DeleteInsightRulesOutput>(
      '/monitoring.CloudWatchService/DeleteInsightRules',
      ($0.DeleteInsightRulesInput value) => value.writeToBuffer(),
      $0.DeleteInsightRulesOutput.fromBuffer);
  static final _$deleteMetricStream = $grpc.ClientMethod<
          $0.DeleteMetricStreamInput, $0.DeleteMetricStreamOutput>(
      '/monitoring.CloudWatchService/DeleteMetricStream',
      ($0.DeleteMetricStreamInput value) => value.writeToBuffer(),
      $0.DeleteMetricStreamOutput.fromBuffer);
  static final _$describeAlarmContributors = $grpc.ClientMethod<
          $0.DescribeAlarmContributorsInput,
          $0.DescribeAlarmContributorsOutput>(
      '/monitoring.CloudWatchService/DescribeAlarmContributors',
      ($0.DescribeAlarmContributorsInput value) => value.writeToBuffer(),
      $0.DescribeAlarmContributorsOutput.fromBuffer);
  static final _$describeAlarmHistory = $grpc.ClientMethod<
          $0.DescribeAlarmHistoryInput, $0.DescribeAlarmHistoryOutput>(
      '/monitoring.CloudWatchService/DescribeAlarmHistory',
      ($0.DescribeAlarmHistoryInput value) => value.writeToBuffer(),
      $0.DescribeAlarmHistoryOutput.fromBuffer);
  static final _$describeAlarms =
      $grpc.ClientMethod<$0.DescribeAlarmsInput, $0.DescribeAlarmsOutput>(
          '/monitoring.CloudWatchService/DescribeAlarms',
          ($0.DescribeAlarmsInput value) => value.writeToBuffer(),
          $0.DescribeAlarmsOutput.fromBuffer);
  static final _$describeAlarmsForMetric = $grpc.ClientMethod<
          $0.DescribeAlarmsForMetricInput, $0.DescribeAlarmsForMetricOutput>(
      '/monitoring.CloudWatchService/DescribeAlarmsForMetric',
      ($0.DescribeAlarmsForMetricInput value) => value.writeToBuffer(),
      $0.DescribeAlarmsForMetricOutput.fromBuffer);
  static final _$describeAnomalyDetectors = $grpc.ClientMethod<
          $0.DescribeAnomalyDetectorsInput, $0.DescribeAnomalyDetectorsOutput>(
      '/monitoring.CloudWatchService/DescribeAnomalyDetectors',
      ($0.DescribeAnomalyDetectorsInput value) => value.writeToBuffer(),
      $0.DescribeAnomalyDetectorsOutput.fromBuffer);
  static final _$describeInsightRules = $grpc.ClientMethod<
          $0.DescribeInsightRulesInput, $0.DescribeInsightRulesOutput>(
      '/monitoring.CloudWatchService/DescribeInsightRules',
      ($0.DescribeInsightRulesInput value) => value.writeToBuffer(),
      $0.DescribeInsightRulesOutput.fromBuffer);
  static final _$disableAlarmActions =
      $grpc.ClientMethod<$0.DisableAlarmActionsInput, $1.Empty>(
          '/monitoring.CloudWatchService/DisableAlarmActions',
          ($0.DisableAlarmActionsInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$disableInsightRules = $grpc.ClientMethod<
          $0.DisableInsightRulesInput, $0.DisableInsightRulesOutput>(
      '/monitoring.CloudWatchService/DisableInsightRules',
      ($0.DisableInsightRulesInput value) => value.writeToBuffer(),
      $0.DisableInsightRulesOutput.fromBuffer);
  static final _$enableAlarmActions =
      $grpc.ClientMethod<$0.EnableAlarmActionsInput, $1.Empty>(
          '/monitoring.CloudWatchService/EnableAlarmActions',
          ($0.EnableAlarmActionsInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$enableInsightRules = $grpc.ClientMethod<
          $0.EnableInsightRulesInput, $0.EnableInsightRulesOutput>(
      '/monitoring.CloudWatchService/EnableInsightRules',
      ($0.EnableInsightRulesInput value) => value.writeToBuffer(),
      $0.EnableInsightRulesOutput.fromBuffer);
  static final _$getAlarmMuteRule =
      $grpc.ClientMethod<$0.GetAlarmMuteRuleInput, $0.GetAlarmMuteRuleOutput>(
          '/monitoring.CloudWatchService/GetAlarmMuteRule',
          ($0.GetAlarmMuteRuleInput value) => value.writeToBuffer(),
          $0.GetAlarmMuteRuleOutput.fromBuffer);
  static final _$getDashboard =
      $grpc.ClientMethod<$0.GetDashboardInput, $0.GetDashboardOutput>(
          '/monitoring.CloudWatchService/GetDashboard',
          ($0.GetDashboardInput value) => value.writeToBuffer(),
          $0.GetDashboardOutput.fromBuffer);
  static final _$getInsightRuleReport = $grpc.ClientMethod<
          $0.GetInsightRuleReportInput, $0.GetInsightRuleReportOutput>(
      '/monitoring.CloudWatchService/GetInsightRuleReport',
      ($0.GetInsightRuleReportInput value) => value.writeToBuffer(),
      $0.GetInsightRuleReportOutput.fromBuffer);
  static final _$getMetricData =
      $grpc.ClientMethod<$0.GetMetricDataInput, $0.GetMetricDataOutput>(
          '/monitoring.CloudWatchService/GetMetricData',
          ($0.GetMetricDataInput value) => value.writeToBuffer(),
          $0.GetMetricDataOutput.fromBuffer);
  static final _$getMetricStatistics = $grpc.ClientMethod<
          $0.GetMetricStatisticsInput, $0.GetMetricStatisticsOutput>(
      '/monitoring.CloudWatchService/GetMetricStatistics',
      ($0.GetMetricStatisticsInput value) => value.writeToBuffer(),
      $0.GetMetricStatisticsOutput.fromBuffer);
  static final _$getMetricStream =
      $grpc.ClientMethod<$0.GetMetricStreamInput, $0.GetMetricStreamOutput>(
          '/monitoring.CloudWatchService/GetMetricStream',
          ($0.GetMetricStreamInput value) => value.writeToBuffer(),
          $0.GetMetricStreamOutput.fromBuffer);
  static final _$getMetricWidgetImage = $grpc.ClientMethod<
          $0.GetMetricWidgetImageInput, $0.GetMetricWidgetImageOutput>(
      '/monitoring.CloudWatchService/GetMetricWidgetImage',
      ($0.GetMetricWidgetImageInput value) => value.writeToBuffer(),
      $0.GetMetricWidgetImageOutput.fromBuffer);
  static final _$listAlarmMuteRules = $grpc.ClientMethod<
          $0.ListAlarmMuteRulesInput, $0.ListAlarmMuteRulesOutput>(
      '/monitoring.CloudWatchService/ListAlarmMuteRules',
      ($0.ListAlarmMuteRulesInput value) => value.writeToBuffer(),
      $0.ListAlarmMuteRulesOutput.fromBuffer);
  static final _$listDashboards =
      $grpc.ClientMethod<$0.ListDashboardsInput, $0.ListDashboardsOutput>(
          '/monitoring.CloudWatchService/ListDashboards',
          ($0.ListDashboardsInput value) => value.writeToBuffer(),
          $0.ListDashboardsOutput.fromBuffer);
  static final _$listManagedInsightRules = $grpc.ClientMethod<
          $0.ListManagedInsightRulesInput, $0.ListManagedInsightRulesOutput>(
      '/monitoring.CloudWatchService/ListManagedInsightRules',
      ($0.ListManagedInsightRulesInput value) => value.writeToBuffer(),
      $0.ListManagedInsightRulesOutput.fromBuffer);
  static final _$listMetrics =
      $grpc.ClientMethod<$0.ListMetricsInput, $0.ListMetricsOutput>(
          '/monitoring.CloudWatchService/ListMetrics',
          ($0.ListMetricsInput value) => value.writeToBuffer(),
          $0.ListMetricsOutput.fromBuffer);
  static final _$listMetricStreams =
      $grpc.ClientMethod<$0.ListMetricStreamsInput, $0.ListMetricStreamsOutput>(
          '/monitoring.CloudWatchService/ListMetricStreams',
          ($0.ListMetricStreamsInput value) => value.writeToBuffer(),
          $0.ListMetricStreamsOutput.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceOutput>(
      '/monitoring.CloudWatchService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceOutput.fromBuffer);
  static final _$putAlarmMuteRule =
      $grpc.ClientMethod<$0.PutAlarmMuteRuleInput, $1.Empty>(
          '/monitoring.CloudWatchService/PutAlarmMuteRule',
          ($0.PutAlarmMuteRuleInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putAnomalyDetector = $grpc.ClientMethod<
          $0.PutAnomalyDetectorInput, $0.PutAnomalyDetectorOutput>(
      '/monitoring.CloudWatchService/PutAnomalyDetector',
      ($0.PutAnomalyDetectorInput value) => value.writeToBuffer(),
      $0.PutAnomalyDetectorOutput.fromBuffer);
  static final _$putCompositeAlarm =
      $grpc.ClientMethod<$0.PutCompositeAlarmInput, $1.Empty>(
          '/monitoring.CloudWatchService/PutCompositeAlarm',
          ($0.PutCompositeAlarmInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putDashboard =
      $grpc.ClientMethod<$0.PutDashboardInput, $0.PutDashboardOutput>(
          '/monitoring.CloudWatchService/PutDashboard',
          ($0.PutDashboardInput value) => value.writeToBuffer(),
          $0.PutDashboardOutput.fromBuffer);
  static final _$putInsightRule =
      $grpc.ClientMethod<$0.PutInsightRuleInput, $0.PutInsightRuleOutput>(
          '/monitoring.CloudWatchService/PutInsightRule',
          ($0.PutInsightRuleInput value) => value.writeToBuffer(),
          $0.PutInsightRuleOutput.fromBuffer);
  static final _$putManagedInsightRules = $grpc.ClientMethod<
          $0.PutManagedInsightRulesInput, $0.PutManagedInsightRulesOutput>(
      '/monitoring.CloudWatchService/PutManagedInsightRules',
      ($0.PutManagedInsightRulesInput value) => value.writeToBuffer(),
      $0.PutManagedInsightRulesOutput.fromBuffer);
  static final _$putMetricAlarm =
      $grpc.ClientMethod<$0.PutMetricAlarmInput, $1.Empty>(
          '/monitoring.CloudWatchService/PutMetricAlarm',
          ($0.PutMetricAlarmInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putMetricData =
      $grpc.ClientMethod<$0.PutMetricDataInput, $1.Empty>(
          '/monitoring.CloudWatchService/PutMetricData',
          ($0.PutMetricDataInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putMetricStream =
      $grpc.ClientMethod<$0.PutMetricStreamInput, $0.PutMetricStreamOutput>(
          '/monitoring.CloudWatchService/PutMetricStream',
          ($0.PutMetricStreamInput value) => value.writeToBuffer(),
          $0.PutMetricStreamOutput.fromBuffer);
  static final _$setAlarmState =
      $grpc.ClientMethod<$0.SetAlarmStateInput, $1.Empty>(
          '/monitoring.CloudWatchService/SetAlarmState',
          ($0.SetAlarmStateInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$startMetricStreams = $grpc.ClientMethod<
          $0.StartMetricStreamsInput, $0.StartMetricStreamsOutput>(
      '/monitoring.CloudWatchService/StartMetricStreams',
      ($0.StartMetricStreamsInput value) => value.writeToBuffer(),
      $0.StartMetricStreamsOutput.fromBuffer);
  static final _$stopMetricStreams =
      $grpc.ClientMethod<$0.StopMetricStreamsInput, $0.StopMetricStreamsOutput>(
          '/monitoring.CloudWatchService/StopMetricStreams',
          ($0.StopMetricStreamsInput value) => value.writeToBuffer(),
          $0.StopMetricStreamsOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $0.TagResourceOutput>(
          '/monitoring.CloudWatchService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $0.TagResourceOutput.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
          '/monitoring.CloudWatchService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $0.UntagResourceOutput.fromBuffer);
}

@$pb.GrpcServiceName('monitoring.CloudWatchService')
abstract class CloudWatchServiceBase extends $grpc.Service {
  $core.String get $name => 'monitoring.CloudWatchService';

  CloudWatchServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.DeleteAlarmMuteRuleInput, $1.Empty>(
        'DeleteAlarmMuteRule',
        deleteAlarmMuteRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAlarmMuteRuleInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAlarmsInput, $1.Empty>(
        'DeleteAlarms',
        deleteAlarms_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteAlarmsInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAnomalyDetectorInput,
            $0.DeleteAnomalyDetectorOutput>(
        'DeleteAnomalyDetector',
        deleteAnomalyDetector_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAnomalyDetectorInput.fromBuffer(value),
        ($0.DeleteAnomalyDetectorOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDashboardsInput,
            $0.DeleteDashboardsOutput>(
        'DeleteDashboards',
        deleteDashboards_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDashboardsInput.fromBuffer(value),
        ($0.DeleteDashboardsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteInsightRulesInput,
            $0.DeleteInsightRulesOutput>(
        'DeleteInsightRules',
        deleteInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteInsightRulesInput.fromBuffer(value),
        ($0.DeleteInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMetricStreamInput,
            $0.DeleteMetricStreamOutput>(
        'DeleteMetricStream',
        deleteMetricStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMetricStreamInput.fromBuffer(value),
        ($0.DeleteMetricStreamOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAlarmContributorsInput,
            $0.DescribeAlarmContributorsOutput>(
        'DescribeAlarmContributors',
        describeAlarmContributors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAlarmContributorsInput.fromBuffer(value),
        ($0.DescribeAlarmContributorsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAlarmHistoryInput,
            $0.DescribeAlarmHistoryOutput>(
        'DescribeAlarmHistory',
        describeAlarmHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAlarmHistoryInput.fromBuffer(value),
        ($0.DescribeAlarmHistoryOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeAlarmsInput, $0.DescribeAlarmsOutput>(
            'DescribeAlarms',
            describeAlarms_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeAlarmsInput.fromBuffer(value),
            ($0.DescribeAlarmsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAlarmsForMetricInput,
            $0.DescribeAlarmsForMetricOutput>(
        'DescribeAlarmsForMetric',
        describeAlarmsForMetric_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAlarmsForMetricInput.fromBuffer(value),
        ($0.DescribeAlarmsForMetricOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAnomalyDetectorsInput,
            $0.DescribeAnomalyDetectorsOutput>(
        'DescribeAnomalyDetectors',
        describeAnomalyDetectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAnomalyDetectorsInput.fromBuffer(value),
        ($0.DescribeAnomalyDetectorsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInsightRulesInput,
            $0.DescribeInsightRulesOutput>(
        'DescribeInsightRules',
        describeInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInsightRulesInput.fromBuffer(value),
        ($0.DescribeInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableAlarmActionsInput, $1.Empty>(
        'DisableAlarmActions',
        disableAlarmActions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableAlarmActionsInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableInsightRulesInput,
            $0.DisableInsightRulesOutput>(
        'DisableInsightRules',
        disableInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableInsightRulesInput.fromBuffer(value),
        ($0.DisableInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableAlarmActionsInput, $1.Empty>(
        'EnableAlarmActions',
        enableAlarmActions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableAlarmActionsInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableInsightRulesInput,
            $0.EnableInsightRulesOutput>(
        'EnableInsightRules',
        enableInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableInsightRulesInput.fromBuffer(value),
        ($0.EnableInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAlarmMuteRuleInput,
            $0.GetAlarmMuteRuleOutput>(
        'GetAlarmMuteRule',
        getAlarmMuteRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAlarmMuteRuleInput.fromBuffer(value),
        ($0.GetAlarmMuteRuleOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDashboardInput, $0.GetDashboardOutput>(
        'GetDashboard',
        getDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetDashboardInput.fromBuffer(value),
        ($0.GetDashboardOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInsightRuleReportInput,
            $0.GetInsightRuleReportOutput>(
        'GetInsightRuleReport',
        getInsightRuleReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInsightRuleReportInput.fromBuffer(value),
        ($0.GetInsightRuleReportOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetMetricDataInput, $0.GetMetricDataOutput>(
            'GetMetricData',
            getMetricData_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetMetricDataInput.fromBuffer(value),
            ($0.GetMetricDataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMetricStatisticsInput,
            $0.GetMetricStatisticsOutput>(
        'GetMetricStatistics',
        getMetricStatistics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMetricStatisticsInput.fromBuffer(value),
        ($0.GetMetricStatisticsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetMetricStreamInput, $0.GetMetricStreamOutput>(
            'GetMetricStream',
            getMetricStream_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetMetricStreamInput.fromBuffer(value),
            ($0.GetMetricStreamOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMetricWidgetImageInput,
            $0.GetMetricWidgetImageOutput>(
        'GetMetricWidgetImage',
        getMetricWidgetImage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMetricWidgetImageInput.fromBuffer(value),
        ($0.GetMetricWidgetImageOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAlarmMuteRulesInput,
            $0.ListAlarmMuteRulesOutput>(
        'ListAlarmMuteRules',
        listAlarmMuteRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAlarmMuteRulesInput.fromBuffer(value),
        ($0.ListAlarmMuteRulesOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDashboardsInput, $0.ListDashboardsOutput>(
            'ListDashboards',
            listDashboards_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDashboardsInput.fromBuffer(value),
            ($0.ListDashboardsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListManagedInsightRulesInput,
            $0.ListManagedInsightRulesOutput>(
        'ListManagedInsightRules',
        listManagedInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListManagedInsightRulesInput.fromBuffer(value),
        ($0.ListManagedInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMetricsInput, $0.ListMetricsOutput>(
        'ListMetrics',
        listMetrics_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListMetricsInput.fromBuffer(value),
        ($0.ListMetricsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMetricStreamsInput,
            $0.ListMetricStreamsOutput>(
        'ListMetricStreams',
        listMetricStreams_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMetricStreamsInput.fromBuffer(value),
        ($0.ListMetricStreamsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceOutput>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAlarmMuteRuleInput, $1.Empty>(
        'PutAlarmMuteRule',
        putAlarmMuteRule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAlarmMuteRuleInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAnomalyDetectorInput,
            $0.PutAnomalyDetectorOutput>(
        'PutAnomalyDetector',
        putAnomalyDetector_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAnomalyDetectorInput.fromBuffer(value),
        ($0.PutAnomalyDetectorOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutCompositeAlarmInput, $1.Empty>(
        'PutCompositeAlarm',
        putCompositeAlarm_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutCompositeAlarmInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDashboardInput, $0.PutDashboardOutput>(
        'PutDashboard',
        putDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutDashboardInput.fromBuffer(value),
        ($0.PutDashboardOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutInsightRuleInput, $0.PutInsightRuleOutput>(
            'PutInsightRule',
            putInsightRule_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutInsightRuleInput.fromBuffer(value),
            ($0.PutInsightRuleOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutManagedInsightRulesInput,
            $0.PutManagedInsightRulesOutput>(
        'PutManagedInsightRules',
        putManagedInsightRules_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutManagedInsightRulesInput.fromBuffer(value),
        ($0.PutManagedInsightRulesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutMetricAlarmInput, $1.Empty>(
        'PutMetricAlarm',
        putMetricAlarm_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutMetricAlarmInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutMetricDataInput, $1.Empty>(
        'PutMetricData',
        putMetricData_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutMetricDataInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutMetricStreamInput, $0.PutMetricStreamOutput>(
            'PutMetricStream',
            putMetricStream_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutMetricStreamInput.fromBuffer(value),
            ($0.PutMetricStreamOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetAlarmStateInput, $1.Empty>(
        'SetAlarmState',
        setAlarmState_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetAlarmStateInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartMetricStreamsInput,
            $0.StartMetricStreamsOutput>(
        'StartMetricStreams',
        startMetricStreams_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartMetricStreamsInput.fromBuffer(value),
        ($0.StartMetricStreamsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopMetricStreamsInput,
            $0.StopMetricStreamsOutput>(
        'StopMetricStreams',
        stopMetricStreams_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopMetricStreamsInput.fromBuffer(value),
        ($0.StopMetricStreamsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $0.TagResourceOutput>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($0.TagResourceOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceInput.fromBuffer(value),
            ($0.UntagResourceOutput value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> deleteAlarmMuteRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAlarmMuteRuleInput> $request) async {
    return deleteAlarmMuteRule($call, await $request);
  }

  $async.Future<$1.Empty> deleteAlarmMuteRule(
      $grpc.ServiceCall call, $0.DeleteAlarmMuteRuleInput request);

  $async.Future<$1.Empty> deleteAlarms_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAlarmsInput> $request) async {
    return deleteAlarms($call, await $request);
  }

  $async.Future<$1.Empty> deleteAlarms(
      $grpc.ServiceCall call, $0.DeleteAlarmsInput request);

  $async.Future<$0.DeleteAnomalyDetectorOutput> deleteAnomalyDetector_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteAnomalyDetectorInput> $request) async {
    return deleteAnomalyDetector($call, await $request);
  }

  $async.Future<$0.DeleteAnomalyDetectorOutput> deleteAnomalyDetector(
      $grpc.ServiceCall call, $0.DeleteAnomalyDetectorInput request);

  $async.Future<$0.DeleteDashboardsOutput> deleteDashboards_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDashboardsInput> $request) async {
    return deleteDashboards($call, await $request);
  }

  $async.Future<$0.DeleteDashboardsOutput> deleteDashboards(
      $grpc.ServiceCall call, $0.DeleteDashboardsInput request);

  $async.Future<$0.DeleteInsightRulesOutput> deleteInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteInsightRulesInput> $request) async {
    return deleteInsightRules($call, await $request);
  }

  $async.Future<$0.DeleteInsightRulesOutput> deleteInsightRules(
      $grpc.ServiceCall call, $0.DeleteInsightRulesInput request);

  $async.Future<$0.DeleteMetricStreamOutput> deleteMetricStream_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteMetricStreamInput> $request) async {
    return deleteMetricStream($call, await $request);
  }

  $async.Future<$0.DeleteMetricStreamOutput> deleteMetricStream(
      $grpc.ServiceCall call, $0.DeleteMetricStreamInput request);

  $async.Future<$0.DescribeAlarmContributorsOutput>
      describeAlarmContributors_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeAlarmContributorsInput> $request) async {
    return describeAlarmContributors($call, await $request);
  }

  $async.Future<$0.DescribeAlarmContributorsOutput> describeAlarmContributors(
      $grpc.ServiceCall call, $0.DescribeAlarmContributorsInput request);

  $async.Future<$0.DescribeAlarmHistoryOutput> describeAlarmHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAlarmHistoryInput> $request) async {
    return describeAlarmHistory($call, await $request);
  }

  $async.Future<$0.DescribeAlarmHistoryOutput> describeAlarmHistory(
      $grpc.ServiceCall call, $0.DescribeAlarmHistoryInput request);

  $async.Future<$0.DescribeAlarmsOutput> describeAlarms_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAlarmsInput> $request) async {
    return describeAlarms($call, await $request);
  }

  $async.Future<$0.DescribeAlarmsOutput> describeAlarms(
      $grpc.ServiceCall call, $0.DescribeAlarmsInput request);

  $async.Future<$0.DescribeAlarmsForMetricOutput> describeAlarmsForMetric_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAlarmsForMetricInput> $request) async {
    return describeAlarmsForMetric($call, await $request);
  }

  $async.Future<$0.DescribeAlarmsForMetricOutput> describeAlarmsForMetric(
      $grpc.ServiceCall call, $0.DescribeAlarmsForMetricInput request);

  $async.Future<$0.DescribeAnomalyDetectorsOutput> describeAnomalyDetectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAnomalyDetectorsInput> $request) async {
    return describeAnomalyDetectors($call, await $request);
  }

  $async.Future<$0.DescribeAnomalyDetectorsOutput> describeAnomalyDetectors(
      $grpc.ServiceCall call, $0.DescribeAnomalyDetectorsInput request);

  $async.Future<$0.DescribeInsightRulesOutput> describeInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeInsightRulesInput> $request) async {
    return describeInsightRules($call, await $request);
  }

  $async.Future<$0.DescribeInsightRulesOutput> describeInsightRules(
      $grpc.ServiceCall call, $0.DescribeInsightRulesInput request);

  $async.Future<$1.Empty> disableAlarmActions_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisableAlarmActionsInput> $request) async {
    return disableAlarmActions($call, await $request);
  }

  $async.Future<$1.Empty> disableAlarmActions(
      $grpc.ServiceCall call, $0.DisableAlarmActionsInput request);

  $async.Future<$0.DisableInsightRulesOutput> disableInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DisableInsightRulesInput> $request) async {
    return disableInsightRules($call, await $request);
  }

  $async.Future<$0.DisableInsightRulesOutput> disableInsightRules(
      $grpc.ServiceCall call, $0.DisableInsightRulesInput request);

  $async.Future<$1.Empty> enableAlarmActions_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EnableAlarmActionsInput> $request) async {
    return enableAlarmActions($call, await $request);
  }

  $async.Future<$1.Empty> enableAlarmActions(
      $grpc.ServiceCall call, $0.EnableAlarmActionsInput request);

  $async.Future<$0.EnableInsightRulesOutput> enableInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.EnableInsightRulesInput> $request) async {
    return enableInsightRules($call, await $request);
  }

  $async.Future<$0.EnableInsightRulesOutput> enableInsightRules(
      $grpc.ServiceCall call, $0.EnableInsightRulesInput request);

  $async.Future<$0.GetAlarmMuteRuleOutput> getAlarmMuteRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAlarmMuteRuleInput> $request) async {
    return getAlarmMuteRule($call, await $request);
  }

  $async.Future<$0.GetAlarmMuteRuleOutput> getAlarmMuteRule(
      $grpc.ServiceCall call, $0.GetAlarmMuteRuleInput request);

  $async.Future<$0.GetDashboardOutput> getDashboard_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDashboardInput> $request) async {
    return getDashboard($call, await $request);
  }

  $async.Future<$0.GetDashboardOutput> getDashboard(
      $grpc.ServiceCall call, $0.GetDashboardInput request);

  $async.Future<$0.GetInsightRuleReportOutput> getInsightRuleReport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInsightRuleReportInput> $request) async {
    return getInsightRuleReport($call, await $request);
  }

  $async.Future<$0.GetInsightRuleReportOutput> getInsightRuleReport(
      $grpc.ServiceCall call, $0.GetInsightRuleReportInput request);

  $async.Future<$0.GetMetricDataOutput> getMetricData_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMetricDataInput> $request) async {
    return getMetricData($call, await $request);
  }

  $async.Future<$0.GetMetricDataOutput> getMetricData(
      $grpc.ServiceCall call, $0.GetMetricDataInput request);

  $async.Future<$0.GetMetricStatisticsOutput> getMetricStatistics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMetricStatisticsInput> $request) async {
    return getMetricStatistics($call, await $request);
  }

  $async.Future<$0.GetMetricStatisticsOutput> getMetricStatistics(
      $grpc.ServiceCall call, $0.GetMetricStatisticsInput request);

  $async.Future<$0.GetMetricStreamOutput> getMetricStream_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMetricStreamInput> $request) async {
    return getMetricStream($call, await $request);
  }

  $async.Future<$0.GetMetricStreamOutput> getMetricStream(
      $grpc.ServiceCall call, $0.GetMetricStreamInput request);

  $async.Future<$0.GetMetricWidgetImageOutput> getMetricWidgetImage_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMetricWidgetImageInput> $request) async {
    return getMetricWidgetImage($call, await $request);
  }

  $async.Future<$0.GetMetricWidgetImageOutput> getMetricWidgetImage(
      $grpc.ServiceCall call, $0.GetMetricWidgetImageInput request);

  $async.Future<$0.ListAlarmMuteRulesOutput> listAlarmMuteRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAlarmMuteRulesInput> $request) async {
    return listAlarmMuteRules($call, await $request);
  }

  $async.Future<$0.ListAlarmMuteRulesOutput> listAlarmMuteRules(
      $grpc.ServiceCall call, $0.ListAlarmMuteRulesInput request);

  $async.Future<$0.ListDashboardsOutput> listDashboards_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDashboardsInput> $request) async {
    return listDashboards($call, await $request);
  }

  $async.Future<$0.ListDashboardsOutput> listDashboards(
      $grpc.ServiceCall call, $0.ListDashboardsInput request);

  $async.Future<$0.ListManagedInsightRulesOutput> listManagedInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListManagedInsightRulesInput> $request) async {
    return listManagedInsightRules($call, await $request);
  }

  $async.Future<$0.ListManagedInsightRulesOutput> listManagedInsightRules(
      $grpc.ServiceCall call, $0.ListManagedInsightRulesInput request);

  $async.Future<$0.ListMetricsOutput> listMetrics_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListMetricsInput> $request) async {
    return listMetrics($call, await $request);
  }

  $async.Future<$0.ListMetricsOutput> listMetrics(
      $grpc.ServiceCall call, $0.ListMetricsInput request);

  $async.Future<$0.ListMetricStreamsOutput> listMetricStreams_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMetricStreamsInput> $request) async {
    return listMetricStreams($call, await $request);
  }

  $async.Future<$0.ListMetricStreamsOutput> listMetricStreams(
      $grpc.ServiceCall call, $0.ListMetricStreamsInput request);

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$1.Empty> putAlarmMuteRule_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutAlarmMuteRuleInput> $request) async {
    return putAlarmMuteRule($call, await $request);
  }

  $async.Future<$1.Empty> putAlarmMuteRule(
      $grpc.ServiceCall call, $0.PutAlarmMuteRuleInput request);

  $async.Future<$0.PutAnomalyDetectorOutput> putAnomalyDetector_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutAnomalyDetectorInput> $request) async {
    return putAnomalyDetector($call, await $request);
  }

  $async.Future<$0.PutAnomalyDetectorOutput> putAnomalyDetector(
      $grpc.ServiceCall call, $0.PutAnomalyDetectorInput request);

  $async.Future<$1.Empty> putCompositeAlarm_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutCompositeAlarmInput> $request) async {
    return putCompositeAlarm($call, await $request);
  }

  $async.Future<$1.Empty> putCompositeAlarm(
      $grpc.ServiceCall call, $0.PutCompositeAlarmInput request);

  $async.Future<$0.PutDashboardOutput> putDashboard_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutDashboardInput> $request) async {
    return putDashboard($call, await $request);
  }

  $async.Future<$0.PutDashboardOutput> putDashboard(
      $grpc.ServiceCall call, $0.PutDashboardInput request);

  $async.Future<$0.PutInsightRuleOutput> putInsightRule_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutInsightRuleInput> $request) async {
    return putInsightRule($call, await $request);
  }

  $async.Future<$0.PutInsightRuleOutput> putInsightRule(
      $grpc.ServiceCall call, $0.PutInsightRuleInput request);

  $async.Future<$0.PutManagedInsightRulesOutput> putManagedInsightRules_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutManagedInsightRulesInput> $request) async {
    return putManagedInsightRules($call, await $request);
  }

  $async.Future<$0.PutManagedInsightRulesOutput> putManagedInsightRules(
      $grpc.ServiceCall call, $0.PutManagedInsightRulesInput request);

  $async.Future<$1.Empty> putMetricAlarm_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutMetricAlarmInput> $request) async {
    return putMetricAlarm($call, await $request);
  }

  $async.Future<$1.Empty> putMetricAlarm(
      $grpc.ServiceCall call, $0.PutMetricAlarmInput request);

  $async.Future<$1.Empty> putMetricData_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutMetricDataInput> $request) async {
    return putMetricData($call, await $request);
  }

  $async.Future<$1.Empty> putMetricData(
      $grpc.ServiceCall call, $0.PutMetricDataInput request);

  $async.Future<$0.PutMetricStreamOutput> putMetricStream_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutMetricStreamInput> $request) async {
    return putMetricStream($call, await $request);
  }

  $async.Future<$0.PutMetricStreamOutput> putMetricStream(
      $grpc.ServiceCall call, $0.PutMetricStreamInput request);

  $async.Future<$1.Empty> setAlarmState_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetAlarmStateInput> $request) async {
    return setAlarmState($call, await $request);
  }

  $async.Future<$1.Empty> setAlarmState(
      $grpc.ServiceCall call, $0.SetAlarmStateInput request);

  $async.Future<$0.StartMetricStreamsOutput> startMetricStreams_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartMetricStreamsInput> $request) async {
    return startMetricStreams($call, await $request);
  }

  $async.Future<$0.StartMetricStreamsOutput> startMetricStreams(
      $grpc.ServiceCall call, $0.StartMetricStreamsInput request);

  $async.Future<$0.StopMetricStreamsOutput> stopMetricStreams_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopMetricStreamsInput> $request) async {
    return stopMetricStreams($call, await $request);
  }

  $async.Future<$0.StopMetricStreamsOutput> stopMetricStreams(
      $grpc.ServiceCall call, $0.StopMetricStreamsInput request);

  $async.Future<$0.TagResourceOutput> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceOutput> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$0.UntagResourceOutput> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceOutput> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);
}
