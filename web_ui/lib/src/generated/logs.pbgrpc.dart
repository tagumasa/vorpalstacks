// This is a generated file - do not edit.
//
// Generated from logs.proto.

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
import 'logs.pb.dart' as $0;

export 'logs.pb.dart';

/// CloudWatchLogsService provides logs API operations.
@$pb.GrpcServiceName('logs.CloudWatchLogsService')
class CloudWatchLogsServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CloudWatchLogsServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Associates the specified KMS key with either one log group in the account, or with all stored CloudWatch Logs query insights results in the account. When you use AssociateKmsKey, you specify either...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> associateKmsKey(
    $0.AssociateKmsKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateKmsKey, request, options: options);
  }

  /// Associates a data source with an S3 Table Integration for query access in the 'logs' namespace. This enables querying log data using analytics engines that support Iceberg such as Amazon Athena, Am...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AssociateSourceToS3TableIntegrationResponse>
      associateSourceToS3TableIntegration(
    $0.AssociateSourceToS3TableIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateSourceToS3TableIntegration, request,
        options: options);
  }

  /// Cancels the specified export task. The task must be in the PENDING or RUNNING state.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> cancelExportTask(
    $0.CancelExportTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelExportTask, request, options: options);
  }

  /// Cancels an active import task and stops importing data from the CloudTrail Lake Event Data Store.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelImportTaskResponse> cancelImportTask(
    $0.CancelImportTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelImportTask, request, options: options);
  }

  /// Creates a delivery. A delivery is a connection between a logical delivery source and a logical delivery destination that you have already created. Only some Amazon Web Services services support bei...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateDeliveryResponse> createDelivery(
    $0.CreateDeliveryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDelivery, request, options: options);
  }

  /// Creates an export task so that you can efficiently export data from a log group to an Amazon S3 bucket. When you perform a CreateExportTask operation, you must use credentials that have permission ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateExportTaskResponse> createExportTask(
    $0.CreateExportTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createExportTask, request, options: options);
  }

  /// Starts an import from a data source to CloudWatch Log and creates a managed log group as the destination for the imported data. Currently, CloudTrail Event Data Store is the only supported data sou...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateImportTaskResponse> createImportTask(
    $0.CreateImportTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createImportTask, request, options: options);
  }

  /// Creates an anomaly detector that regularly scans one or more log groups and look for patterns and anomalies in the logs. An anomaly detector can help surface issues by automatically discovering ano...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateLogAnomalyDetectorResponse>
      createLogAnomalyDetector(
    $0.CreateLogAnomalyDetectorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createLogAnomalyDetector, request,
        options: options);
  }

  /// Creates a log group with the specified name. You can create up to 1,000,000 log groups per Region per account. You must use the following guidelines when naming a log group: Log group names must be...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> createLogGroup(
    $0.CreateLogGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createLogGroup, request, options: options);
  }

  /// Creates a log stream for the specified log group. A log stream is a sequence of log events that originate from a single source, such as an application instance or a resource that is being monitored...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> createLogStream(
    $0.CreateLogStreamRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createLogStream, request, options: options);
  }

  /// Creates a scheduled query that runs CloudWatch Logs Insights queries at regular intervals. Scheduled queries enable proactive monitoring by automatically executing queries to detect patterns and an...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateScheduledQueryResponse> createScheduledQuery(
    $0.CreateScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createScheduledQuery, request, options: options);
  }

  /// Deletes a CloudWatch Logs account policy. This stops the account-wide policy from applying to log groups or data sources in the account. If you delete a data protection policy or subscription filte...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteAccountPolicy(
    $0.DeleteAccountPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAccountPolicy, request, options: options);
  }

  /// Deletes the data protection policy from the specified log group. For more information about data protection policies, see PutDataProtectionPolicy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDataProtectionPolicy(
    $0.DeleteDataProtectionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDataProtectionPolicy, request,
        options: options);
  }

  /// Deletes a delivery. A delivery is a connection between a logical delivery source and a logical delivery destination. Deleting a delivery only deletes the connection between the delivery source and ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDelivery(
    $0.DeleteDeliveryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDelivery, request, options: options);
  }

  /// Deletes a delivery destination. A delivery is a connection between a logical delivery source and a logical delivery destination. You can't delete a delivery destination if any current deliveries ar...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDeliveryDestination(
    $0.DeleteDeliveryDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDeliveryDestination, request,
        options: options);
  }

  /// Deletes a delivery destination policy. For more information about these policies, see PutDeliveryDestinationPolicy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDeliveryDestinationPolicy(
    $0.DeleteDeliveryDestinationPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDeliveryDestinationPolicy, request,
        options: options);
  }

  /// Deletes a delivery source. A delivery is a connection between a logical delivery source and a logical delivery destination. You can't delete a delivery source if any current deliveries are associat...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDeliverySource(
    $0.DeleteDeliverySourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDeliverySource, request, options: options);
  }

  /// Deletes the specified destination, and eventually disables all the subscription filters that publish to it. This operation does not delete the physical resource encapsulated by the destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteDestination(
    $0.DeleteDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDestination, request, options: options);
  }

  /// Deletes a log-group level field index policy that was applied to a single log group. The indexing of the log events that happened before you delete the policy will still be used for as many as 30 d...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteIndexPolicyResponse> deleteIndexPolicy(
    $0.DeleteIndexPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIndexPolicy, request, options: options);
  }

  /// Deletes the integration between CloudWatch Logs and OpenSearch Service. If your integration has active vended logs dashboards, you must specify true for the force parameter, otherwise the operation...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteIntegrationResponse> deleteIntegration(
    $0.DeleteIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIntegration, request, options: options);
  }

  /// Deletes the specified CloudWatch Logs anomaly detector.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteLogAnomalyDetector(
    $0.DeleteLogAnomalyDetectorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteLogAnomalyDetector, request,
        options: options);
  }

  /// Deletes the specified log group and permanently deletes all the archived log events associated with the log group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteLogGroup(
    $0.DeleteLogGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteLogGroup, request, options: options);
  }

  /// Deletes the specified log stream and permanently deletes all the archived log events associated with the log stream.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteLogStream(
    $0.DeleteLogStreamRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteLogStream, request, options: options);
  }

  /// Deletes the specified metric filter.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteMetricFilter(
    $0.DeleteMetricFilterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMetricFilter, request, options: options);
  }

  /// Deletes a saved CloudWatch Logs Insights query definition. A query definition contains details about a saved CloudWatch Logs Insights query. Each DeleteQueryDefinition operation can delete one quer...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteQueryDefinitionResponse> deleteQueryDefinition(
    $0.DeleteQueryDefinitionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteQueryDefinition, request, options: options);
  }

  /// Deletes a resource policy from this account. This revokes the access of the identities in that policy to put log events to this account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteResourcePolicy(
    $0.DeleteResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Deletes the specified retention policy. Log events do not expire if they belong to log groups without a retention policy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteRetentionPolicy(
    $0.DeleteRetentionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRetentionPolicy, request, options: options);
  }

  /// Deletes a scheduled query and stops all future executions. This operation also removes any configured actions and associated resources.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteScheduledQueryResponse> deleteScheduledQuery(
    $0.DeleteScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteScheduledQuery, request, options: options);
  }

  /// Deletes the specified subscription filter.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteSubscriptionFilter(
    $0.DeleteSubscriptionFilterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSubscriptionFilter, request,
        options: options);
  }

  /// Deletes the log transformer for the specified log group. As soon as you do this, the transformation of incoming log events according to that transformer stops. If this account has an account-level ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteTransformer(
    $0.DeleteTransformerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTransformer, request, options: options);
  }

  /// Returns a list of all CloudWatch Logs account policies in the account. To use this operation, you must be signed on with the correct permissions depending on the type of policy that you are retriev...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAccountPoliciesResponse>
      describeAccountPolicies(
    $0.DescribeAccountPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAccountPolicies, request,
        options: options);
  }

  /// Use this operation to return the valid and default values that are used when creating delivery sources, delivery destinations, and deliveries. For more information about deliveries, see CreateDeliv...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeConfigurationTemplatesResponse>
      describeConfigurationTemplates(
    $0.DescribeConfigurationTemplatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeConfigurationTemplates, request,
        options: options);
  }

  /// Retrieves a list of the deliveries that have been created in the account. A delivery is a connection between a delivery source and a delivery destination . A delivery source represents an Amazon We...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDeliveriesResponse> describeDeliveries(
    $0.DescribeDeliveriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDeliveries, request, options: options);
  }

  /// Retrieves a list of the delivery destinations that have been created in the account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDeliveryDestinationsResponse>
      describeDeliveryDestinations(
    $0.DescribeDeliveryDestinationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDeliveryDestinations, request,
        options: options);
  }

  /// Retrieves a list of the delivery sources that have been created in the account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDeliverySourcesResponse>
      describeDeliverySources(
    $0.DescribeDeliverySourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDeliverySources, request,
        options: options);
  }

  /// Lists all your destinations. The results are ASCII-sorted by destination name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDestinationsResponse> describeDestinations(
    $0.DescribeDestinationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDestinations, request, options: options);
  }

  /// Lists the specified export tasks. You can list all your export tasks or filter the results based on task ID or task status.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeExportTasksResponse> describeExportTasks(
    $0.DescribeExportTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeExportTasks, request, options: options);
  }

  /// Returns a list of custom and default field indexes which are discovered in log data. For more information about field index policies, see PutIndexPolicy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeFieldIndexesResponse> describeFieldIndexes(
    $0.DescribeFieldIndexesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeFieldIndexes, request, options: options);
  }

  /// Gets detailed information about the individual batches within an import task, including their status and any error messages. For CloudTrail Event Data Store sources, a batch refers to a subset of s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeImportTaskBatchesResponse>
      describeImportTaskBatches(
    $0.DescribeImportTaskBatchesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeImportTaskBatches, request,
        options: options);
  }

  /// Lists and describes import tasks, with optional filtering by import status and source ARN.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeImportTasksResponse> describeImportTasks(
    $0.DescribeImportTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeImportTasks, request, options: options);
  }

  /// Returns the field index policies of the specified log group. For more information about field index policies, see PutIndexPolicy. If a specified log group has a log-group level index policy, that p...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeIndexPoliciesResponse> describeIndexPolicies(
    $0.DescribeIndexPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeIndexPolicies, request, options: options);
  }

  /// Returns information about log groups, including data sources that ingest into each log group. You can return all your log groups or filter the results by prefix. The results are ASCII-sorted by log...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeLogGroupsResponse> describeLogGroups(
    $0.DescribeLogGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeLogGroups, request, options: options);
  }

  /// Lists the log streams for the specified log group. You can list all the log streams or filter the results by prefix. You can also control how the results are ordered. You can specify the log group ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeLogStreamsResponse> describeLogStreams(
    $0.DescribeLogStreamsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeLogStreams, request, options: options);
  }

  /// Lists the specified metric filters. You can list all of the metric filters or filter the results by log name, prefix, metric name, or metric namespace. The results are ASCII-sorted by filter name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMetricFiltersResponse> describeMetricFilters(
    $0.DescribeMetricFiltersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMetricFilters, request, options: options);
  }

  /// Returns a list of CloudWatch Logs Insights queries that are scheduled, running, or have been run recently in this account. You can request all queries or limit it to queries of a specific log group...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeQueriesResponse> describeQueries(
    $0.DescribeQueriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeQueries, request, options: options);
  }

  /// This operation returns a paginated list of your saved CloudWatch Logs Insights query definitions. You can retrieve query definitions from the current account or from a source account that is linked...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeQueryDefinitionsResponse>
      describeQueryDefinitions(
    $0.DescribeQueryDefinitionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeQueryDefinitions, request,
        options: options);
  }

  /// Lists the resource policies in this account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeResourcePoliciesResponse>
      describeResourcePolicies(
    $0.DescribeResourcePoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeResourcePolicies, request,
        options: options);
  }

  /// Lists the subscription filters for the specified log group. You can list all the subscription filters or filter the results by prefix. The results are ASCII-sorted by filter name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeSubscriptionFiltersResponse>
      describeSubscriptionFilters(
    $0.DescribeSubscriptionFiltersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeSubscriptionFilters, request,
        options: options);
  }

  /// Disassociates the specified KMS key from the specified log group or from all CloudWatch Logs Insights query results in the account. When you use DisassociateKmsKey, you specify either the logGroupN...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> disassociateKmsKey(
    $0.DisassociateKmsKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateKmsKey, request, options: options);
  }

  /// Disassociates a data source from an S3 Table Integration, removing query access and deleting all associated data from the integration.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DisassociateSourceFromS3TableIntegrationResponse>
      disassociateSourceFromS3TableIntegration(
    $0.DisassociateSourceFromS3TableIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateSourceFromS3TableIntegration, request,
        options: options);
  }

  /// Lists log events from the specified log group. You can list all the log events or filter the results using one or more of the following: A filter pattern A time range The log stream name, or a log ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.FilterLogEventsResponse> filterLogEvents(
    $0.FilterLogEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$filterLogEvents, request, options: options);
  }

  /// Returns information about a log group data protection policy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDataProtectionPolicyResponse>
      getDataProtectionPolicy(
    $0.GetDataProtectionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDataProtectionPolicy, request,
        options: options);
  }

  /// Returns complete information about one logical delivery. A delivery is a connection between a delivery source and a delivery destination . A delivery source represents an Amazon Web Services resour...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeliveryResponse> getDelivery(
    $0.GetDeliveryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDelivery, request, options: options);
  }

  /// Retrieves complete information about one delivery destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeliveryDestinationResponse>
      getDeliveryDestination(
    $0.GetDeliveryDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeliveryDestination, request,
        options: options);
  }

  /// Retrieves the delivery destination policy assigned to the delivery destination that you specify. For more information about delivery destinations and their policies, see PutDeliveryDestinationPolicy.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeliveryDestinationPolicyResponse>
      getDeliveryDestinationPolicy(
    $0.GetDeliveryDestinationPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeliveryDestinationPolicy, request,
        options: options);
  }

  /// Retrieves complete information about one delivery source.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeliverySourceResponse> getDeliverySource(
    $0.GetDeliverySourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeliverySource, request, options: options);
  }

  /// Returns information about one integration between CloudWatch Logs and OpenSearch Service.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIntegrationResponse> getIntegration(
    $0.GetIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIntegration, request, options: options);
  }

  /// Retrieves information about the log anomaly detector that you specify. The KMS key ARN detected is valid.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogAnomalyDetectorResponse> getLogAnomalyDetector(
    $0.GetLogAnomalyDetectorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogAnomalyDetector, request, options: options);
  }

  /// Lists log events from the specified log stream. You can list all of the log events or filter using a time range. GetLogEvents is a paginated operation. Each page returned can contain up to 1 MB of ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogEventsResponse> getLogEvents(
    $0.GetLogEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogEvents, request, options: options);
  }

  /// Discovers available fields for a specific data source and type. The response includes any field modifications introduced through pipelines, such as new fields or changed field types.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogFieldsResponse> getLogFields(
    $0.GetLogFieldsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogFields, request, options: options);
  }

  /// Returns a list of the fields that are included in log events in the specified log group. Includes the percentage of log events that contain each field. The search is limited to a time period that y...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogGroupFieldsResponse> getLogGroupFields(
    $0.GetLogGroupFieldsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogGroupFields, request, options: options);
  }

  /// Retrieves a large logging object (LLO) and streams it back. This API is used to fetch the content of large portions of log events that have been ingested through the PutOpenTelemetryLogs API. When ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogObjectResponse> getLogObject(
    $0.GetLogObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogObject, request, options: options);
  }

  /// Retrieves all of the fields and values of a single log event. All fields are retrieved, even if the original query that produced the logRecordPointer retrieved only a subset of fields. Fields are r...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogRecordResponse> getLogRecord(
    $0.GetLogRecordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogRecord, request, options: options);
  }

  /// Returns the results from the specified query. Only the fields requested in the query are returned, along with a @ptr field, which is the identifier for the log record. You can use the value of @ptr...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetQueryResultsResponse> getQueryResults(
    $0.GetQueryResultsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryResults, request, options: options);
  }

  /// Retrieves details about a specific scheduled query, including its configuration, execution status, and metadata.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetScheduledQueryResponse> getScheduledQuery(
    $0.GetScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getScheduledQuery, request, options: options);
  }

  /// Retrieves the execution history of a scheduled query within a specified time range, including query results and destination processing status.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetScheduledQueryHistoryResponse>
      getScheduledQueryHistory(
    $0.GetScheduledQueryHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getScheduledQueryHistory, request,
        options: options);
  }

  /// Returns the information about the log transformer associated with this log group. This operation returns data only for transformers created at the log group level. To get information for an account...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTransformerResponse> getTransformer(
    $0.GetTransformerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTransformer, request, options: options);
  }

  /// Returns an aggregate summary of all log groups in the Region grouped by specified data source characteristics. Supports optional filtering by log group class, name patterns, and data sources. If yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAggregateLogGroupSummariesResponse>
      listAggregateLogGroupSummaries(
    $0.ListAggregateLogGroupSummariesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAggregateLogGroupSummaries, request,
        options: options);
  }

  /// Returns a list of anomalies that log anomaly detectors have found. For details about the structure format of each anomaly object that is returned, see the example in this section.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAnomaliesResponse> listAnomalies(
    $0.ListAnomaliesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAnomalies, request, options: options);
  }

  /// Returns a list of integrations between CloudWatch Logs and other services in this account. Currently, only one integration can be created in an account, and this integration must be with OpenSearch...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIntegrationsResponse> listIntegrations(
    $0.ListIntegrationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIntegrations, request, options: options);
  }

  /// Retrieves a list of the log anomaly detectors in the account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListLogAnomalyDetectorsResponse>
      listLogAnomalyDetectors(
    $0.ListLogAnomalyDetectorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLogAnomalyDetectors, request,
        options: options);
  }

  /// Returns a list of log groups in the Region in your account. If you are performing this action in a monitoring account, you can choose to also return log groups from source accounts that are linked ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListLogGroupsResponse> listLogGroups(
    $0.ListLogGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLogGroups, request, options: options);
  }

  /// Returns a list of the log groups that were analyzed during a single CloudWatch Logs Insights query. This can be useful for queries that use log group name prefixes or the filterIndex command, becau...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListLogGroupsForQueryResponse> listLogGroupsForQuery(
    $0.ListLogGroupsForQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLogGroupsForQuery, request, options: options);
  }

  /// Lists all scheduled queries in your account and region. You can filter results by state to show only enabled or disabled queries.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListScheduledQueriesResponse> listScheduledQueries(
    $0.ListScheduledQueriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listScheduledQueries, request, options: options);
  }

  /// Returns a list of data source associations for a specified S3 Table Integration, showing which data sources are currently associated for query access.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSourcesForS3TableIntegrationResponse>
      listSourcesForS3TableIntegration(
    $0.ListSourcesForS3TableIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSourcesForS3TableIntegration, request,
        options: options);
  }

  /// Displays the tags associated with a CloudWatch Logs resource. Currently, log groups and destinations support tagging.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// The ListTagsLogGroup operation is on the path to deprecation. We recommend that you use ListTagsForResource instead. Lists the tags for the specified log group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsLogGroupResponse> listTagsLogGroup(
    $0.ListTagsLogGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsLogGroup, request, options: options);
  }

  /// Creates an account-level data protection policy, subscription filter policy, field index policy, transformer policy, or metric extraction policy that applies to all log groups, a subset of log grou...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutAccountPolicyResponse> putAccountPolicy(
    $0.PutAccountPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountPolicy, request, options: options);
  }

  /// Enables or disables bearer token authentication for the specified log group. When enabled on a log group, bearer token authentication is enabled on operations until it is explicitly disabled. For i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putBearerTokenAuthentication(
    $0.PutBearerTokenAuthenticationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBearerTokenAuthentication, request,
        options: options);
  }

  /// Creates a data protection policy for the specified log group. A data protection policy can help safeguard sensitive data that's ingested by the log group by auditing and masking the sensitive log d...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutDataProtectionPolicyResponse>
      putDataProtectionPolicy(
    $0.PutDataProtectionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDataProtectionPolicy, request,
        options: options);
  }

  /// Creates or updates a logical delivery destination. A delivery destination is an Amazon Web Services resource that represents an Amazon Web Services service that logs can be sent to. CloudWatch Logs...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutDeliveryDestinationResponse>
      putDeliveryDestination(
    $0.PutDeliveryDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDeliveryDestination, request,
        options: options);
  }

  /// Creates and assigns an IAM policy that grants permissions to CloudWatch Logs to deliver logs cross-account to a specified destination in this account. To configure the delivery of logs from an Amaz...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutDeliveryDestinationPolicyResponse>
      putDeliveryDestinationPolicy(
    $0.PutDeliveryDestinationPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDeliveryDestinationPolicy, request,
        options: options);
  }

  /// Creates or updates a logical delivery source. A delivery source represents an Amazon Web Services resource that sends logs to an logs delivery destination. The destination can be CloudWatch Logs, A...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutDeliverySourceResponse> putDeliverySource(
    $0.PutDeliverySourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDeliverySource, request, options: options);
  }

  /// Creates or updates a destination. This operation is used only to create destinations for cross-account subscriptions. A destination encapsulates a physical resource (such as an Amazon Kinesis strea...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutDestinationResponse> putDestination(
    $0.PutDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDestination, request, options: options);
  }

  /// Creates or updates an access policy associated with an existing destination. An access policy is an IAM policy document that is used to authorize claims to register a subscription filter against a ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putDestinationPolicy(
    $0.PutDestinationPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDestinationPolicy, request, options: options);
  }

  /// Creates or updates a field index policy for the specified log group. Only log groups in the Standard log class support field index policies. For more information about log classes, see Log classes....
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutIndexPolicyResponse> putIndexPolicy(
    $0.PutIndexPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putIndexPolicy, request, options: options);
  }

  /// Creates an integration between CloudWatch Logs and another service in this account. Currently, only integrations with OpenSearch Service are supported, and currently you can have only one integrati...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutIntegrationResponse> putIntegration(
    $0.PutIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putIntegration, request, options: options);
  }

  /// Uploads a batch of log events to the specified log stream. The sequence token is now ignored in PutLogEvents actions. PutLogEvents actions are always accepted and never return InvalidSequenceTokenE...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutLogEventsResponse> putLogEvents(
    $0.PutLogEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putLogEvents, request, options: options);
  }

  /// Enables or disables deletion protection for the specified log group. When enabled on a log group, deletion protection blocks all deletion operations until it is explicitly disabled. For information...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putLogGroupDeletionProtection(
    $0.PutLogGroupDeletionProtectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putLogGroupDeletionProtection, request,
        options: options);
  }

  /// Creates or updates a metric filter and associates it with the specified log group. With metric filters, you can configure rules to extract metric data from log events ingested through PutLogEvents....
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putMetricFilter(
    $0.PutMetricFilterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMetricFilter, request, options: options);
  }

  /// Creates or updates a query definition for CloudWatch Logs Insights. For more information, see Analyzing Log Data with CloudWatch Logs Insights. To update a query definition, specify its queryDefini...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutQueryDefinitionResponse> putQueryDefinition(
    $0.PutQueryDefinitionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putQueryDefinition, request, options: options);
  }

  /// Creates or updates a resource policy allowing other Amazon Web Services services to put log events to this account, such as Amazon Route 53. This API has the following restrictions: Supported actio...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutResourcePolicyResponse> putResourcePolicy(
    $0.PutResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Sets the retention of the specified log group. With a retention policy, you can configure the number of days for which to retain log events in the specified log group. CloudWatch Logs doesn't immed...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putRetentionPolicy(
    $0.PutRetentionPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRetentionPolicy, request, options: options);
  }

  /// Creates or updates a subscription filter and associates it with the specified log group. With subscription filters, you can subscribe to a real-time stream of log events ingested through PutLogEven...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putSubscriptionFilter(
    $0.PutSubscriptionFilterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putSubscriptionFilter, request, options: options);
  }

  /// Creates or updates a log transformer for a single log group. You use log transformers to transform log events into a different format, making them easier for you to process and analyze. You can als...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putTransformer(
    $0.PutTransformerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putTransformer, request, options: options);
  }

  /// Starts a Live Tail streaming session for one or more log groups. A Live Tail session returns a stream of log events that have been recently ingested in the log groups. For more information, see Use...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartLiveTailResponse> startLiveTail(
    $0.StartLiveTailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startLiveTail, request, options: options);
  }

  /// Starts a query of one or more log groups or data sources using CloudWatch Logs Insights. You specify the log groups or data sources and time range to query and the query string to use. You can quer...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartQueryResponse> startQuery(
    $0.StartQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startQuery, request, options: options);
  }

  /// Stops a CloudWatch Logs Insights query that is in progress. If the query has already ended, the operation returns an error indicating that the specified query is not running. This operation can be ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopQueryResponse> stopQuery(
    $0.StopQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopQuery, request, options: options);
  }

  /// The TagLogGroup operation is on the path to deprecation. We recommend that you use TagResource instead. Adds or updates the specified tags for the specified log group. To list the tags for a log gr...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> tagLogGroup(
    $0.TagLogGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagLogGroup, request, options: options);
  }

  /// Assigns one or more tags (key-value pairs) to the specified CloudWatch Logs resource. Currently, the only CloudWatch Logs resources that can be tagged are log groups and destinations. Tags can help...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Tests the filter pattern of a metric filter against a sample of log event messages. You can use this operation to validate the correctness of a metric filter pattern.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TestMetricFilterResponse> testMetricFilter(
    $0.TestMetricFilterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testMetricFilter, request, options: options);
  }

  /// Use this operation to test a log transformer. You enter the transformer configuration and a set of log events to test with. The operation responds with an array that includes the original log event...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TestTransformerResponse> testTransformer(
    $0.TestTransformerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testTransformer, request, options: options);
  }

  /// The UntagLogGroup operation is on the path to deprecation. We recommend that you use UntagResource instead. Removes the specified tags from the specified log group. To list the tags for a log group...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> untagLogGroup(
    $0.UntagLogGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagLogGroup, request, options: options);
  }

  /// Removes one or more tags from the specified resource.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Use this operation to suppress anomaly detection for a specified anomaly or pattern. If you suppress an anomaly, CloudWatch Logs won't report new occurrences of that anomaly and won't update that a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateAnomaly(
    $0.UpdateAnomalyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAnomaly, request, options: options);
  }

  /// Use this operation to update the configuration of a delivery to change either the S3 path pattern or the format of the delivered logs. You can't use this operation to change the source or destinati...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDeliveryConfigurationResponse>
      updateDeliveryConfiguration(
    $0.UpdateDeliveryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDeliveryConfiguration, request,
        options: options);
  }

  /// Updates an existing log anomaly detector.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateLogAnomalyDetector(
    $0.UpdateLogAnomalyDetectorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateLogAnomalyDetector, request,
        options: options);
  }

  /// Updates an existing scheduled query with new configuration. This operation uses PUT semantics, allowing modification of query parameters, schedule, and destinations.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateScheduledQueryResponse> updateScheduledQuery(
    $0.UpdateScheduledQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateScheduledQuery, request, options: options);
  }

  // method descriptors

  static final _$associateKmsKey =
      $grpc.ClientMethod<$0.AssociateKmsKeyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/AssociateKmsKey',
          ($0.AssociateKmsKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$associateSourceToS3TableIntegration = $grpc.ClientMethod<
          $0.AssociateSourceToS3TableIntegrationRequest,
          $0.AssociateSourceToS3TableIntegrationResponse>(
      '/logs.CloudWatchLogsService/AssociateSourceToS3TableIntegration',
      ($0.AssociateSourceToS3TableIntegrationRequest value) =>
          value.writeToBuffer(),
      $0.AssociateSourceToS3TableIntegrationResponse.fromBuffer);
  static final _$cancelExportTask =
      $grpc.ClientMethod<$0.CancelExportTaskRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/CancelExportTask',
          ($0.CancelExportTaskRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$cancelImportTask = $grpc.ClientMethod<
          $0.CancelImportTaskRequest, $0.CancelImportTaskResponse>(
      '/logs.CloudWatchLogsService/CancelImportTask',
      ($0.CancelImportTaskRequest value) => value.writeToBuffer(),
      $0.CancelImportTaskResponse.fromBuffer);
  static final _$createDelivery =
      $grpc.ClientMethod<$0.CreateDeliveryRequest, $0.CreateDeliveryResponse>(
          '/logs.CloudWatchLogsService/CreateDelivery',
          ($0.CreateDeliveryRequest value) => value.writeToBuffer(),
          $0.CreateDeliveryResponse.fromBuffer);
  static final _$createExportTask = $grpc.ClientMethod<
          $0.CreateExportTaskRequest, $0.CreateExportTaskResponse>(
      '/logs.CloudWatchLogsService/CreateExportTask',
      ($0.CreateExportTaskRequest value) => value.writeToBuffer(),
      $0.CreateExportTaskResponse.fromBuffer);
  static final _$createImportTask = $grpc.ClientMethod<
          $0.CreateImportTaskRequest, $0.CreateImportTaskResponse>(
      '/logs.CloudWatchLogsService/CreateImportTask',
      ($0.CreateImportTaskRequest value) => value.writeToBuffer(),
      $0.CreateImportTaskResponse.fromBuffer);
  static final _$createLogAnomalyDetector = $grpc.ClientMethod<
          $0.CreateLogAnomalyDetectorRequest,
          $0.CreateLogAnomalyDetectorResponse>(
      '/logs.CloudWatchLogsService/CreateLogAnomalyDetector',
      ($0.CreateLogAnomalyDetectorRequest value) => value.writeToBuffer(),
      $0.CreateLogAnomalyDetectorResponse.fromBuffer);
  static final _$createLogGroup =
      $grpc.ClientMethod<$0.CreateLogGroupRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/CreateLogGroup',
          ($0.CreateLogGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createLogStream =
      $grpc.ClientMethod<$0.CreateLogStreamRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/CreateLogStream',
          ($0.CreateLogStreamRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createScheduledQuery = $grpc.ClientMethod<
          $0.CreateScheduledQueryRequest, $0.CreateScheduledQueryResponse>(
      '/logs.CloudWatchLogsService/CreateScheduledQuery',
      ($0.CreateScheduledQueryRequest value) => value.writeToBuffer(),
      $0.CreateScheduledQueryResponse.fromBuffer);
  static final _$deleteAccountPolicy =
      $grpc.ClientMethod<$0.DeleteAccountPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteAccountPolicy',
          ($0.DeleteAccountPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDataProtectionPolicy =
      $grpc.ClientMethod<$0.DeleteDataProtectionPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDataProtectionPolicy',
          ($0.DeleteDataProtectionPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDelivery =
      $grpc.ClientMethod<$0.DeleteDeliveryRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDelivery',
          ($0.DeleteDeliveryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDeliveryDestination =
      $grpc.ClientMethod<$0.DeleteDeliveryDestinationRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDeliveryDestination',
          ($0.DeleteDeliveryDestinationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDeliveryDestinationPolicy =
      $grpc.ClientMethod<$0.DeleteDeliveryDestinationPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDeliveryDestinationPolicy',
          ($0.DeleteDeliveryDestinationPolicyRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDeliverySource =
      $grpc.ClientMethod<$0.DeleteDeliverySourceRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDeliverySource',
          ($0.DeleteDeliverySourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDestination =
      $grpc.ClientMethod<$0.DeleteDestinationRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteDestination',
          ($0.DeleteDestinationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteIndexPolicy = $grpc.ClientMethod<
          $0.DeleteIndexPolicyRequest, $0.DeleteIndexPolicyResponse>(
      '/logs.CloudWatchLogsService/DeleteIndexPolicy',
      ($0.DeleteIndexPolicyRequest value) => value.writeToBuffer(),
      $0.DeleteIndexPolicyResponse.fromBuffer);
  static final _$deleteIntegration = $grpc.ClientMethod<
          $0.DeleteIntegrationRequest, $0.DeleteIntegrationResponse>(
      '/logs.CloudWatchLogsService/DeleteIntegration',
      ($0.DeleteIntegrationRequest value) => value.writeToBuffer(),
      $0.DeleteIntegrationResponse.fromBuffer);
  static final _$deleteLogAnomalyDetector =
      $grpc.ClientMethod<$0.DeleteLogAnomalyDetectorRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteLogAnomalyDetector',
          ($0.DeleteLogAnomalyDetectorRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteLogGroup =
      $grpc.ClientMethod<$0.DeleteLogGroupRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteLogGroup',
          ($0.DeleteLogGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteLogStream =
      $grpc.ClientMethod<$0.DeleteLogStreamRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteLogStream',
          ($0.DeleteLogStreamRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteMetricFilter =
      $grpc.ClientMethod<$0.DeleteMetricFilterRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteMetricFilter',
          ($0.DeleteMetricFilterRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteQueryDefinition = $grpc.ClientMethod<
          $0.DeleteQueryDefinitionRequest, $0.DeleteQueryDefinitionResponse>(
      '/logs.CloudWatchLogsService/DeleteQueryDefinition',
      ($0.DeleteQueryDefinitionRequest value) => value.writeToBuffer(),
      $0.DeleteQueryDefinitionResponse.fromBuffer);
  static final _$deleteResourcePolicy =
      $grpc.ClientMethod<$0.DeleteResourcePolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteResourcePolicy',
          ($0.DeleteResourcePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRetentionPolicy =
      $grpc.ClientMethod<$0.DeleteRetentionPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteRetentionPolicy',
          ($0.DeleteRetentionPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteScheduledQuery = $grpc.ClientMethod<
          $0.DeleteScheduledQueryRequest, $0.DeleteScheduledQueryResponse>(
      '/logs.CloudWatchLogsService/DeleteScheduledQuery',
      ($0.DeleteScheduledQueryRequest value) => value.writeToBuffer(),
      $0.DeleteScheduledQueryResponse.fromBuffer);
  static final _$deleteSubscriptionFilter =
      $grpc.ClientMethod<$0.DeleteSubscriptionFilterRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteSubscriptionFilter',
          ($0.DeleteSubscriptionFilterRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteTransformer =
      $grpc.ClientMethod<$0.DeleteTransformerRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DeleteTransformer',
          ($0.DeleteTransformerRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeAccountPolicies = $grpc.ClientMethod<
          $0.DescribeAccountPoliciesRequest,
          $0.DescribeAccountPoliciesResponse>(
      '/logs.CloudWatchLogsService/DescribeAccountPolicies',
      ($0.DescribeAccountPoliciesRequest value) => value.writeToBuffer(),
      $0.DescribeAccountPoliciesResponse.fromBuffer);
  static final _$describeConfigurationTemplates = $grpc.ClientMethod<
          $0.DescribeConfigurationTemplatesRequest,
          $0.DescribeConfigurationTemplatesResponse>(
      '/logs.CloudWatchLogsService/DescribeConfigurationTemplates',
      ($0.DescribeConfigurationTemplatesRequest value) => value.writeToBuffer(),
      $0.DescribeConfigurationTemplatesResponse.fromBuffer);
  static final _$describeDeliveries = $grpc.ClientMethod<
          $0.DescribeDeliveriesRequest, $0.DescribeDeliveriesResponse>(
      '/logs.CloudWatchLogsService/DescribeDeliveries',
      ($0.DescribeDeliveriesRequest value) => value.writeToBuffer(),
      $0.DescribeDeliveriesResponse.fromBuffer);
  static final _$describeDeliveryDestinations = $grpc.ClientMethod<
          $0.DescribeDeliveryDestinationsRequest,
          $0.DescribeDeliveryDestinationsResponse>(
      '/logs.CloudWatchLogsService/DescribeDeliveryDestinations',
      ($0.DescribeDeliveryDestinationsRequest value) => value.writeToBuffer(),
      $0.DescribeDeliveryDestinationsResponse.fromBuffer);
  static final _$describeDeliverySources = $grpc.ClientMethod<
          $0.DescribeDeliverySourcesRequest,
          $0.DescribeDeliverySourcesResponse>(
      '/logs.CloudWatchLogsService/DescribeDeliverySources',
      ($0.DescribeDeliverySourcesRequest value) => value.writeToBuffer(),
      $0.DescribeDeliverySourcesResponse.fromBuffer);
  static final _$describeDestinations = $grpc.ClientMethod<
          $0.DescribeDestinationsRequest, $0.DescribeDestinationsResponse>(
      '/logs.CloudWatchLogsService/DescribeDestinations',
      ($0.DescribeDestinationsRequest value) => value.writeToBuffer(),
      $0.DescribeDestinationsResponse.fromBuffer);
  static final _$describeExportTasks = $grpc.ClientMethod<
          $0.DescribeExportTasksRequest, $0.DescribeExportTasksResponse>(
      '/logs.CloudWatchLogsService/DescribeExportTasks',
      ($0.DescribeExportTasksRequest value) => value.writeToBuffer(),
      $0.DescribeExportTasksResponse.fromBuffer);
  static final _$describeFieldIndexes = $grpc.ClientMethod<
          $0.DescribeFieldIndexesRequest, $0.DescribeFieldIndexesResponse>(
      '/logs.CloudWatchLogsService/DescribeFieldIndexes',
      ($0.DescribeFieldIndexesRequest value) => value.writeToBuffer(),
      $0.DescribeFieldIndexesResponse.fromBuffer);
  static final _$describeImportTaskBatches = $grpc.ClientMethod<
          $0.DescribeImportTaskBatchesRequest,
          $0.DescribeImportTaskBatchesResponse>(
      '/logs.CloudWatchLogsService/DescribeImportTaskBatches',
      ($0.DescribeImportTaskBatchesRequest value) => value.writeToBuffer(),
      $0.DescribeImportTaskBatchesResponse.fromBuffer);
  static final _$describeImportTasks = $grpc.ClientMethod<
          $0.DescribeImportTasksRequest, $0.DescribeImportTasksResponse>(
      '/logs.CloudWatchLogsService/DescribeImportTasks',
      ($0.DescribeImportTasksRequest value) => value.writeToBuffer(),
      $0.DescribeImportTasksResponse.fromBuffer);
  static final _$describeIndexPolicies = $grpc.ClientMethod<
          $0.DescribeIndexPoliciesRequest, $0.DescribeIndexPoliciesResponse>(
      '/logs.CloudWatchLogsService/DescribeIndexPolicies',
      ($0.DescribeIndexPoliciesRequest value) => value.writeToBuffer(),
      $0.DescribeIndexPoliciesResponse.fromBuffer);
  static final _$describeLogGroups = $grpc.ClientMethod<
          $0.DescribeLogGroupsRequest, $0.DescribeLogGroupsResponse>(
      '/logs.CloudWatchLogsService/DescribeLogGroups',
      ($0.DescribeLogGroupsRequest value) => value.writeToBuffer(),
      $0.DescribeLogGroupsResponse.fromBuffer);
  static final _$describeLogStreams = $grpc.ClientMethod<
          $0.DescribeLogStreamsRequest, $0.DescribeLogStreamsResponse>(
      '/logs.CloudWatchLogsService/DescribeLogStreams',
      ($0.DescribeLogStreamsRequest value) => value.writeToBuffer(),
      $0.DescribeLogStreamsResponse.fromBuffer);
  static final _$describeMetricFilters = $grpc.ClientMethod<
          $0.DescribeMetricFiltersRequest, $0.DescribeMetricFiltersResponse>(
      '/logs.CloudWatchLogsService/DescribeMetricFilters',
      ($0.DescribeMetricFiltersRequest value) => value.writeToBuffer(),
      $0.DescribeMetricFiltersResponse.fromBuffer);
  static final _$describeQueries =
      $grpc.ClientMethod<$0.DescribeQueriesRequest, $0.DescribeQueriesResponse>(
          '/logs.CloudWatchLogsService/DescribeQueries',
          ($0.DescribeQueriesRequest value) => value.writeToBuffer(),
          $0.DescribeQueriesResponse.fromBuffer);
  static final _$describeQueryDefinitions = $grpc.ClientMethod<
          $0.DescribeQueryDefinitionsRequest,
          $0.DescribeQueryDefinitionsResponse>(
      '/logs.CloudWatchLogsService/DescribeQueryDefinitions',
      ($0.DescribeQueryDefinitionsRequest value) => value.writeToBuffer(),
      $0.DescribeQueryDefinitionsResponse.fromBuffer);
  static final _$describeResourcePolicies = $grpc.ClientMethod<
          $0.DescribeResourcePoliciesRequest,
          $0.DescribeResourcePoliciesResponse>(
      '/logs.CloudWatchLogsService/DescribeResourcePolicies',
      ($0.DescribeResourcePoliciesRequest value) => value.writeToBuffer(),
      $0.DescribeResourcePoliciesResponse.fromBuffer);
  static final _$describeSubscriptionFilters = $grpc.ClientMethod<
          $0.DescribeSubscriptionFiltersRequest,
          $0.DescribeSubscriptionFiltersResponse>(
      '/logs.CloudWatchLogsService/DescribeSubscriptionFilters',
      ($0.DescribeSubscriptionFiltersRequest value) => value.writeToBuffer(),
      $0.DescribeSubscriptionFiltersResponse.fromBuffer);
  static final _$disassociateKmsKey =
      $grpc.ClientMethod<$0.DisassociateKmsKeyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/DisassociateKmsKey',
          ($0.DisassociateKmsKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$disassociateSourceFromS3TableIntegration = $grpc.ClientMethod<
          $0.DisassociateSourceFromS3TableIntegrationRequest,
          $0.DisassociateSourceFromS3TableIntegrationResponse>(
      '/logs.CloudWatchLogsService/DisassociateSourceFromS3TableIntegration',
      ($0.DisassociateSourceFromS3TableIntegrationRequest value) =>
          value.writeToBuffer(),
      $0.DisassociateSourceFromS3TableIntegrationResponse.fromBuffer);
  static final _$filterLogEvents =
      $grpc.ClientMethod<$0.FilterLogEventsRequest, $0.FilterLogEventsResponse>(
          '/logs.CloudWatchLogsService/FilterLogEvents',
          ($0.FilterLogEventsRequest value) => value.writeToBuffer(),
          $0.FilterLogEventsResponse.fromBuffer);
  static final _$getDataProtectionPolicy = $grpc.ClientMethod<
          $0.GetDataProtectionPolicyRequest,
          $0.GetDataProtectionPolicyResponse>(
      '/logs.CloudWatchLogsService/GetDataProtectionPolicy',
      ($0.GetDataProtectionPolicyRequest value) => value.writeToBuffer(),
      $0.GetDataProtectionPolicyResponse.fromBuffer);
  static final _$getDelivery =
      $grpc.ClientMethod<$0.GetDeliveryRequest, $0.GetDeliveryResponse>(
          '/logs.CloudWatchLogsService/GetDelivery',
          ($0.GetDeliveryRequest value) => value.writeToBuffer(),
          $0.GetDeliveryResponse.fromBuffer);
  static final _$getDeliveryDestination = $grpc.ClientMethod<
          $0.GetDeliveryDestinationRequest, $0.GetDeliveryDestinationResponse>(
      '/logs.CloudWatchLogsService/GetDeliveryDestination',
      ($0.GetDeliveryDestinationRequest value) => value.writeToBuffer(),
      $0.GetDeliveryDestinationResponse.fromBuffer);
  static final _$getDeliveryDestinationPolicy = $grpc.ClientMethod<
          $0.GetDeliveryDestinationPolicyRequest,
          $0.GetDeliveryDestinationPolicyResponse>(
      '/logs.CloudWatchLogsService/GetDeliveryDestinationPolicy',
      ($0.GetDeliveryDestinationPolicyRequest value) => value.writeToBuffer(),
      $0.GetDeliveryDestinationPolicyResponse.fromBuffer);
  static final _$getDeliverySource = $grpc.ClientMethod<
          $0.GetDeliverySourceRequest, $0.GetDeliverySourceResponse>(
      '/logs.CloudWatchLogsService/GetDeliverySource',
      ($0.GetDeliverySourceRequest value) => value.writeToBuffer(),
      $0.GetDeliverySourceResponse.fromBuffer);
  static final _$getIntegration =
      $grpc.ClientMethod<$0.GetIntegrationRequest, $0.GetIntegrationResponse>(
          '/logs.CloudWatchLogsService/GetIntegration',
          ($0.GetIntegrationRequest value) => value.writeToBuffer(),
          $0.GetIntegrationResponse.fromBuffer);
  static final _$getLogAnomalyDetector = $grpc.ClientMethod<
          $0.GetLogAnomalyDetectorRequest, $0.GetLogAnomalyDetectorResponse>(
      '/logs.CloudWatchLogsService/GetLogAnomalyDetector',
      ($0.GetLogAnomalyDetectorRequest value) => value.writeToBuffer(),
      $0.GetLogAnomalyDetectorResponse.fromBuffer);
  static final _$getLogEvents =
      $grpc.ClientMethod<$0.GetLogEventsRequest, $0.GetLogEventsResponse>(
          '/logs.CloudWatchLogsService/GetLogEvents',
          ($0.GetLogEventsRequest value) => value.writeToBuffer(),
          $0.GetLogEventsResponse.fromBuffer);
  static final _$getLogFields =
      $grpc.ClientMethod<$0.GetLogFieldsRequest, $0.GetLogFieldsResponse>(
          '/logs.CloudWatchLogsService/GetLogFields',
          ($0.GetLogFieldsRequest value) => value.writeToBuffer(),
          $0.GetLogFieldsResponse.fromBuffer);
  static final _$getLogGroupFields = $grpc.ClientMethod<
          $0.GetLogGroupFieldsRequest, $0.GetLogGroupFieldsResponse>(
      '/logs.CloudWatchLogsService/GetLogGroupFields',
      ($0.GetLogGroupFieldsRequest value) => value.writeToBuffer(),
      $0.GetLogGroupFieldsResponse.fromBuffer);
  static final _$getLogObject =
      $grpc.ClientMethod<$0.GetLogObjectRequest, $0.GetLogObjectResponse>(
          '/logs.CloudWatchLogsService/GetLogObject',
          ($0.GetLogObjectRequest value) => value.writeToBuffer(),
          $0.GetLogObjectResponse.fromBuffer);
  static final _$getLogRecord =
      $grpc.ClientMethod<$0.GetLogRecordRequest, $0.GetLogRecordResponse>(
          '/logs.CloudWatchLogsService/GetLogRecord',
          ($0.GetLogRecordRequest value) => value.writeToBuffer(),
          $0.GetLogRecordResponse.fromBuffer);
  static final _$getQueryResults =
      $grpc.ClientMethod<$0.GetQueryResultsRequest, $0.GetQueryResultsResponse>(
          '/logs.CloudWatchLogsService/GetQueryResults',
          ($0.GetQueryResultsRequest value) => value.writeToBuffer(),
          $0.GetQueryResultsResponse.fromBuffer);
  static final _$getScheduledQuery = $grpc.ClientMethod<
          $0.GetScheduledQueryRequest, $0.GetScheduledQueryResponse>(
      '/logs.CloudWatchLogsService/GetScheduledQuery',
      ($0.GetScheduledQueryRequest value) => value.writeToBuffer(),
      $0.GetScheduledQueryResponse.fromBuffer);
  static final _$getScheduledQueryHistory = $grpc.ClientMethod<
          $0.GetScheduledQueryHistoryRequest,
          $0.GetScheduledQueryHistoryResponse>(
      '/logs.CloudWatchLogsService/GetScheduledQueryHistory',
      ($0.GetScheduledQueryHistoryRequest value) => value.writeToBuffer(),
      $0.GetScheduledQueryHistoryResponse.fromBuffer);
  static final _$getTransformer =
      $grpc.ClientMethod<$0.GetTransformerRequest, $0.GetTransformerResponse>(
          '/logs.CloudWatchLogsService/GetTransformer',
          ($0.GetTransformerRequest value) => value.writeToBuffer(),
          $0.GetTransformerResponse.fromBuffer);
  static final _$listAggregateLogGroupSummaries = $grpc.ClientMethod<
          $0.ListAggregateLogGroupSummariesRequest,
          $0.ListAggregateLogGroupSummariesResponse>(
      '/logs.CloudWatchLogsService/ListAggregateLogGroupSummaries',
      ($0.ListAggregateLogGroupSummariesRequest value) => value.writeToBuffer(),
      $0.ListAggregateLogGroupSummariesResponse.fromBuffer);
  static final _$listAnomalies =
      $grpc.ClientMethod<$0.ListAnomaliesRequest, $0.ListAnomaliesResponse>(
          '/logs.CloudWatchLogsService/ListAnomalies',
          ($0.ListAnomaliesRequest value) => value.writeToBuffer(),
          $0.ListAnomaliesResponse.fromBuffer);
  static final _$listIntegrations = $grpc.ClientMethod<
          $0.ListIntegrationsRequest, $0.ListIntegrationsResponse>(
      '/logs.CloudWatchLogsService/ListIntegrations',
      ($0.ListIntegrationsRequest value) => value.writeToBuffer(),
      $0.ListIntegrationsResponse.fromBuffer);
  static final _$listLogAnomalyDetectors = $grpc.ClientMethod<
          $0.ListLogAnomalyDetectorsRequest,
          $0.ListLogAnomalyDetectorsResponse>(
      '/logs.CloudWatchLogsService/ListLogAnomalyDetectors',
      ($0.ListLogAnomalyDetectorsRequest value) => value.writeToBuffer(),
      $0.ListLogAnomalyDetectorsResponse.fromBuffer);
  static final _$listLogGroups =
      $grpc.ClientMethod<$0.ListLogGroupsRequest, $0.ListLogGroupsResponse>(
          '/logs.CloudWatchLogsService/ListLogGroups',
          ($0.ListLogGroupsRequest value) => value.writeToBuffer(),
          $0.ListLogGroupsResponse.fromBuffer);
  static final _$listLogGroupsForQuery = $grpc.ClientMethod<
          $0.ListLogGroupsForQueryRequest, $0.ListLogGroupsForQueryResponse>(
      '/logs.CloudWatchLogsService/ListLogGroupsForQuery',
      ($0.ListLogGroupsForQueryRequest value) => value.writeToBuffer(),
      $0.ListLogGroupsForQueryResponse.fromBuffer);
  static final _$listScheduledQueries = $grpc.ClientMethod<
          $0.ListScheduledQueriesRequest, $0.ListScheduledQueriesResponse>(
      '/logs.CloudWatchLogsService/ListScheduledQueries',
      ($0.ListScheduledQueriesRequest value) => value.writeToBuffer(),
      $0.ListScheduledQueriesResponse.fromBuffer);
  static final _$listSourcesForS3TableIntegration = $grpc.ClientMethod<
          $0.ListSourcesForS3TableIntegrationRequest,
          $0.ListSourcesForS3TableIntegrationResponse>(
      '/logs.CloudWatchLogsService/ListSourcesForS3TableIntegration',
      ($0.ListSourcesForS3TableIntegrationRequest value) =>
          value.writeToBuffer(),
      $0.ListSourcesForS3TableIntegrationResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/logs.CloudWatchLogsService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTagsLogGroup = $grpc.ClientMethod<
          $0.ListTagsLogGroupRequest, $0.ListTagsLogGroupResponse>(
      '/logs.CloudWatchLogsService/ListTagsLogGroup',
      ($0.ListTagsLogGroupRequest value) => value.writeToBuffer(),
      $0.ListTagsLogGroupResponse.fromBuffer);
  static final _$putAccountPolicy = $grpc.ClientMethod<
          $0.PutAccountPolicyRequest, $0.PutAccountPolicyResponse>(
      '/logs.CloudWatchLogsService/PutAccountPolicy',
      ($0.PutAccountPolicyRequest value) => value.writeToBuffer(),
      $0.PutAccountPolicyResponse.fromBuffer);
  static final _$putBearerTokenAuthentication =
      $grpc.ClientMethod<$0.PutBearerTokenAuthenticationRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutBearerTokenAuthentication',
          ($0.PutBearerTokenAuthenticationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putDataProtectionPolicy = $grpc.ClientMethod<
          $0.PutDataProtectionPolicyRequest,
          $0.PutDataProtectionPolicyResponse>(
      '/logs.CloudWatchLogsService/PutDataProtectionPolicy',
      ($0.PutDataProtectionPolicyRequest value) => value.writeToBuffer(),
      $0.PutDataProtectionPolicyResponse.fromBuffer);
  static final _$putDeliveryDestination = $grpc.ClientMethod<
          $0.PutDeliveryDestinationRequest, $0.PutDeliveryDestinationResponse>(
      '/logs.CloudWatchLogsService/PutDeliveryDestination',
      ($0.PutDeliveryDestinationRequest value) => value.writeToBuffer(),
      $0.PutDeliveryDestinationResponse.fromBuffer);
  static final _$putDeliveryDestinationPolicy = $grpc.ClientMethod<
          $0.PutDeliveryDestinationPolicyRequest,
          $0.PutDeliveryDestinationPolicyResponse>(
      '/logs.CloudWatchLogsService/PutDeliveryDestinationPolicy',
      ($0.PutDeliveryDestinationPolicyRequest value) => value.writeToBuffer(),
      $0.PutDeliveryDestinationPolicyResponse.fromBuffer);
  static final _$putDeliverySource = $grpc.ClientMethod<
          $0.PutDeliverySourceRequest, $0.PutDeliverySourceResponse>(
      '/logs.CloudWatchLogsService/PutDeliverySource',
      ($0.PutDeliverySourceRequest value) => value.writeToBuffer(),
      $0.PutDeliverySourceResponse.fromBuffer);
  static final _$putDestination =
      $grpc.ClientMethod<$0.PutDestinationRequest, $0.PutDestinationResponse>(
          '/logs.CloudWatchLogsService/PutDestination',
          ($0.PutDestinationRequest value) => value.writeToBuffer(),
          $0.PutDestinationResponse.fromBuffer);
  static final _$putDestinationPolicy =
      $grpc.ClientMethod<$0.PutDestinationPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutDestinationPolicy',
          ($0.PutDestinationPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putIndexPolicy =
      $grpc.ClientMethod<$0.PutIndexPolicyRequest, $0.PutIndexPolicyResponse>(
          '/logs.CloudWatchLogsService/PutIndexPolicy',
          ($0.PutIndexPolicyRequest value) => value.writeToBuffer(),
          $0.PutIndexPolicyResponse.fromBuffer);
  static final _$putIntegration =
      $grpc.ClientMethod<$0.PutIntegrationRequest, $0.PutIntegrationResponse>(
          '/logs.CloudWatchLogsService/PutIntegration',
          ($0.PutIntegrationRequest value) => value.writeToBuffer(),
          $0.PutIntegrationResponse.fromBuffer);
  static final _$putLogEvents =
      $grpc.ClientMethod<$0.PutLogEventsRequest, $0.PutLogEventsResponse>(
          '/logs.CloudWatchLogsService/PutLogEvents',
          ($0.PutLogEventsRequest value) => value.writeToBuffer(),
          $0.PutLogEventsResponse.fromBuffer);
  static final _$putLogGroupDeletionProtection =
      $grpc.ClientMethod<$0.PutLogGroupDeletionProtectionRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutLogGroupDeletionProtection',
          ($0.PutLogGroupDeletionProtectionRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putMetricFilter =
      $grpc.ClientMethod<$0.PutMetricFilterRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutMetricFilter',
          ($0.PutMetricFilterRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putQueryDefinition = $grpc.ClientMethod<
          $0.PutQueryDefinitionRequest, $0.PutQueryDefinitionResponse>(
      '/logs.CloudWatchLogsService/PutQueryDefinition',
      ($0.PutQueryDefinitionRequest value) => value.writeToBuffer(),
      $0.PutQueryDefinitionResponse.fromBuffer);
  static final _$putResourcePolicy = $grpc.ClientMethod<
          $0.PutResourcePolicyRequest, $0.PutResourcePolicyResponse>(
      '/logs.CloudWatchLogsService/PutResourcePolicy',
      ($0.PutResourcePolicyRequest value) => value.writeToBuffer(),
      $0.PutResourcePolicyResponse.fromBuffer);
  static final _$putRetentionPolicy =
      $grpc.ClientMethod<$0.PutRetentionPolicyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutRetentionPolicy',
          ($0.PutRetentionPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putSubscriptionFilter =
      $grpc.ClientMethod<$0.PutSubscriptionFilterRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutSubscriptionFilter',
          ($0.PutSubscriptionFilterRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putTransformer =
      $grpc.ClientMethod<$0.PutTransformerRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/PutTransformer',
          ($0.PutTransformerRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$startLiveTail =
      $grpc.ClientMethod<$0.StartLiveTailRequest, $0.StartLiveTailResponse>(
          '/logs.CloudWatchLogsService/StartLiveTail',
          ($0.StartLiveTailRequest value) => value.writeToBuffer(),
          $0.StartLiveTailResponse.fromBuffer);
  static final _$startQuery =
      $grpc.ClientMethod<$0.StartQueryRequest, $0.StartQueryResponse>(
          '/logs.CloudWatchLogsService/StartQuery',
          ($0.StartQueryRequest value) => value.writeToBuffer(),
          $0.StartQueryResponse.fromBuffer);
  static final _$stopQuery =
      $grpc.ClientMethod<$0.StopQueryRequest, $0.StopQueryResponse>(
          '/logs.CloudWatchLogsService/StopQuery',
          ($0.StopQueryRequest value) => value.writeToBuffer(),
          $0.StopQueryResponse.fromBuffer);
  static final _$tagLogGroup =
      $grpc.ClientMethod<$0.TagLogGroupRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/TagLogGroup',
          ($0.TagLogGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$testMetricFilter = $grpc.ClientMethod<
          $0.TestMetricFilterRequest, $0.TestMetricFilterResponse>(
      '/logs.CloudWatchLogsService/TestMetricFilter',
      ($0.TestMetricFilterRequest value) => value.writeToBuffer(),
      $0.TestMetricFilterResponse.fromBuffer);
  static final _$testTransformer =
      $grpc.ClientMethod<$0.TestTransformerRequest, $0.TestTransformerResponse>(
          '/logs.CloudWatchLogsService/TestTransformer',
          ($0.TestTransformerRequest value) => value.writeToBuffer(),
          $0.TestTransformerResponse.fromBuffer);
  static final _$untagLogGroup =
      $grpc.ClientMethod<$0.UntagLogGroupRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/UntagLogGroup',
          ($0.UntagLogGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAnomaly =
      $grpc.ClientMethod<$0.UpdateAnomalyRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/UpdateAnomaly',
          ($0.UpdateAnomalyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateDeliveryConfiguration = $grpc.ClientMethod<
          $0.UpdateDeliveryConfigurationRequest,
          $0.UpdateDeliveryConfigurationResponse>(
      '/logs.CloudWatchLogsService/UpdateDeliveryConfiguration',
      ($0.UpdateDeliveryConfigurationRequest value) => value.writeToBuffer(),
      $0.UpdateDeliveryConfigurationResponse.fromBuffer);
  static final _$updateLogAnomalyDetector =
      $grpc.ClientMethod<$0.UpdateLogAnomalyDetectorRequest, $1.Empty>(
          '/logs.CloudWatchLogsService/UpdateLogAnomalyDetector',
          ($0.UpdateLogAnomalyDetectorRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateScheduledQuery = $grpc.ClientMethod<
          $0.UpdateScheduledQueryRequest, $0.UpdateScheduledQueryResponse>(
      '/logs.CloudWatchLogsService/UpdateScheduledQuery',
      ($0.UpdateScheduledQueryRequest value) => value.writeToBuffer(),
      $0.UpdateScheduledQueryResponse.fromBuffer);
}

@$pb.GrpcServiceName('logs.CloudWatchLogsService')
abstract class CloudWatchLogsServiceBase extends $grpc.Service {
  $core.String get $name => 'logs.CloudWatchLogsService';

  CloudWatchLogsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AssociateKmsKeyRequest, $1.Empty>(
        'AssociateKmsKey',
        associateKmsKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateKmsKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.AssociateSourceToS3TableIntegrationRequest,
            $0.AssociateSourceToS3TableIntegrationResponse>(
        'AssociateSourceToS3TableIntegration',
        associateSourceToS3TableIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateSourceToS3TableIntegrationRequest.fromBuffer(value),
        ($0.AssociateSourceToS3TableIntegrationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelExportTaskRequest, $1.Empty>(
        'CancelExportTask',
        cancelExportTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelExportTaskRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelImportTaskRequest,
            $0.CancelImportTaskResponse>(
        'CancelImportTask',
        cancelImportTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelImportTaskRequest.fromBuffer(value),
        ($0.CancelImportTaskResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDeliveryRequest,
            $0.CreateDeliveryResponse>(
        'CreateDelivery',
        createDelivery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDeliveryRequest.fromBuffer(value),
        ($0.CreateDeliveryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateExportTaskRequest,
            $0.CreateExportTaskResponse>(
        'CreateExportTask',
        createExportTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateExportTaskRequest.fromBuffer(value),
        ($0.CreateExportTaskResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateImportTaskRequest,
            $0.CreateImportTaskResponse>(
        'CreateImportTask',
        createImportTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateImportTaskRequest.fromBuffer(value),
        ($0.CreateImportTaskResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateLogAnomalyDetectorRequest,
            $0.CreateLogAnomalyDetectorResponse>(
        'CreateLogAnomalyDetector',
        createLogAnomalyDetector_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateLogAnomalyDetectorRequest.fromBuffer(value),
        ($0.CreateLogAnomalyDetectorResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateLogGroupRequest, $1.Empty>(
        'CreateLogGroup',
        createLogGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateLogGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateLogStreamRequest, $1.Empty>(
        'CreateLogStream',
        createLogStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateLogStreamRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateScheduledQueryRequest,
            $0.CreateScheduledQueryResponse>(
        'CreateScheduledQuery',
        createScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateScheduledQueryRequest.fromBuffer(value),
        ($0.CreateScheduledQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAccountPolicyRequest, $1.Empty>(
        'DeleteAccountPolicy',
        deleteAccountPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAccountPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteDataProtectionPolicyRequest, $1.Empty>(
            'DeleteDataProtectionPolicy',
            deleteDataProtectionPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteDataProtectionPolicyRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDeliveryRequest, $1.Empty>(
        'DeleteDelivery',
        deleteDelivery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDeliveryRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteDeliveryDestinationRequest, $1.Empty>(
            'DeleteDeliveryDestination',
            deleteDeliveryDestination_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteDeliveryDestinationRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDeliveryDestinationPolicyRequest,
            $1.Empty>(
        'DeleteDeliveryDestinationPolicy',
        deleteDeliveryDestinationPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDeliveryDestinationPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDeliverySourceRequest, $1.Empty>(
        'DeleteDeliverySource',
        deleteDeliverySource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDeliverySourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDestinationRequest, $1.Empty>(
        'DeleteDestination',
        deleteDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDestinationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIndexPolicyRequest,
            $0.DeleteIndexPolicyResponse>(
        'DeleteIndexPolicy',
        deleteIndexPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIndexPolicyRequest.fromBuffer(value),
        ($0.DeleteIndexPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIntegrationRequest,
            $0.DeleteIntegrationResponse>(
        'DeleteIntegration',
        deleteIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIntegrationRequest.fromBuffer(value),
        ($0.DeleteIntegrationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteLogAnomalyDetectorRequest, $1.Empty>(
            'DeleteLogAnomalyDetector',
            deleteLogAnomalyDetector_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteLogAnomalyDetectorRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteLogGroupRequest, $1.Empty>(
        'DeleteLogGroup',
        deleteLogGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteLogGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteLogStreamRequest, $1.Empty>(
        'DeleteLogStream',
        deleteLogStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteLogStreamRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMetricFilterRequest, $1.Empty>(
        'DeleteMetricFilter',
        deleteMetricFilter_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMetricFilterRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteQueryDefinitionRequest,
            $0.DeleteQueryDefinitionResponse>(
        'DeleteQueryDefinition',
        deleteQueryDefinition_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteQueryDefinitionRequest.fromBuffer(value),
        ($0.DeleteQueryDefinitionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyRequest, $1.Empty>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRetentionPolicyRequest, $1.Empty>(
        'DeleteRetentionPolicy',
        deleteRetentionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRetentionPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteScheduledQueryRequest,
            $0.DeleteScheduledQueryResponse>(
        'DeleteScheduledQuery',
        deleteScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteScheduledQueryRequest.fromBuffer(value),
        ($0.DeleteScheduledQueryResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteSubscriptionFilterRequest, $1.Empty>(
            'DeleteSubscriptionFilter',
            deleteSubscriptionFilter_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteSubscriptionFilterRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTransformerRequest, $1.Empty>(
        'DeleteTransformer',
        deleteTransformer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTransformerRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAccountPoliciesRequest,
            $0.DescribeAccountPoliciesResponse>(
        'DescribeAccountPolicies',
        describeAccountPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAccountPoliciesRequest.fromBuffer(value),
        ($0.DescribeAccountPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeConfigurationTemplatesRequest,
            $0.DescribeConfigurationTemplatesResponse>(
        'DescribeConfigurationTemplates',
        describeConfigurationTemplates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeConfigurationTemplatesRequest.fromBuffer(value),
        ($0.DescribeConfigurationTemplatesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDeliveriesRequest,
            $0.DescribeDeliveriesResponse>(
        'DescribeDeliveries',
        describeDeliveries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDeliveriesRequest.fromBuffer(value),
        ($0.DescribeDeliveriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDeliveryDestinationsRequest,
            $0.DescribeDeliveryDestinationsResponse>(
        'DescribeDeliveryDestinations',
        describeDeliveryDestinations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDeliveryDestinationsRequest.fromBuffer(value),
        ($0.DescribeDeliveryDestinationsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDeliverySourcesRequest,
            $0.DescribeDeliverySourcesResponse>(
        'DescribeDeliverySources',
        describeDeliverySources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDeliverySourcesRequest.fromBuffer(value),
        ($0.DescribeDeliverySourcesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDestinationsRequest,
            $0.DescribeDestinationsResponse>(
        'DescribeDestinations',
        describeDestinations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDestinationsRequest.fromBuffer(value),
        ($0.DescribeDestinationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeExportTasksRequest,
            $0.DescribeExportTasksResponse>(
        'DescribeExportTasks',
        describeExportTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeExportTasksRequest.fromBuffer(value),
        ($0.DescribeExportTasksResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeFieldIndexesRequest,
            $0.DescribeFieldIndexesResponse>(
        'DescribeFieldIndexes',
        describeFieldIndexes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeFieldIndexesRequest.fromBuffer(value),
        ($0.DescribeFieldIndexesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeImportTaskBatchesRequest,
            $0.DescribeImportTaskBatchesResponse>(
        'DescribeImportTaskBatches',
        describeImportTaskBatches_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeImportTaskBatchesRequest.fromBuffer(value),
        ($0.DescribeImportTaskBatchesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeImportTasksRequest,
            $0.DescribeImportTasksResponse>(
        'DescribeImportTasks',
        describeImportTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeImportTasksRequest.fromBuffer(value),
        ($0.DescribeImportTasksResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeIndexPoliciesRequest,
            $0.DescribeIndexPoliciesResponse>(
        'DescribeIndexPolicies',
        describeIndexPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeIndexPoliciesRequest.fromBuffer(value),
        ($0.DescribeIndexPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeLogGroupsRequest,
            $0.DescribeLogGroupsResponse>(
        'DescribeLogGroups',
        describeLogGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeLogGroupsRequest.fromBuffer(value),
        ($0.DescribeLogGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeLogStreamsRequest,
            $0.DescribeLogStreamsResponse>(
        'DescribeLogStreams',
        describeLogStreams_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeLogStreamsRequest.fromBuffer(value),
        ($0.DescribeLogStreamsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeMetricFiltersRequest,
            $0.DescribeMetricFiltersResponse>(
        'DescribeMetricFilters',
        describeMetricFilters_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMetricFiltersRequest.fromBuffer(value),
        ($0.DescribeMetricFiltersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeQueriesRequest,
            $0.DescribeQueriesResponse>(
        'DescribeQueries',
        describeQueries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeQueriesRequest.fromBuffer(value),
        ($0.DescribeQueriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeQueryDefinitionsRequest,
            $0.DescribeQueryDefinitionsResponse>(
        'DescribeQueryDefinitions',
        describeQueryDefinitions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeQueryDefinitionsRequest.fromBuffer(value),
        ($0.DescribeQueryDefinitionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeResourcePoliciesRequest,
            $0.DescribeResourcePoliciesResponse>(
        'DescribeResourcePolicies',
        describeResourcePolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeResourcePoliciesRequest.fromBuffer(value),
        ($0.DescribeResourcePoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeSubscriptionFiltersRequest,
            $0.DescribeSubscriptionFiltersResponse>(
        'DescribeSubscriptionFilters',
        describeSubscriptionFilters_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeSubscriptionFiltersRequest.fromBuffer(value),
        ($0.DescribeSubscriptionFiltersResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisassociateKmsKeyRequest, $1.Empty>(
        'DisassociateKmsKey',
        disassociateKmsKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateKmsKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DisassociateSourceFromS3TableIntegrationRequest,
            $0.DisassociateSourceFromS3TableIntegrationResponse>(
        'DisassociateSourceFromS3TableIntegration',
        disassociateSourceFromS3TableIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateSourceFromS3TableIntegrationRequest.fromBuffer(
                value),
        ($0.DisassociateSourceFromS3TableIntegrationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.FilterLogEventsRequest,
            $0.FilterLogEventsResponse>(
        'FilterLogEvents',
        filterLogEvents_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.FilterLogEventsRequest.fromBuffer(value),
        ($0.FilterLogEventsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDataProtectionPolicyRequest,
            $0.GetDataProtectionPolicyResponse>(
        'GetDataProtectionPolicy',
        getDataProtectionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDataProtectionPolicyRequest.fromBuffer(value),
        ($0.GetDataProtectionPolicyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetDeliveryRequest, $0.GetDeliveryResponse>(
            'GetDelivery',
            getDelivery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetDeliveryRequest.fromBuffer(value),
            ($0.GetDeliveryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeliveryDestinationRequest,
            $0.GetDeliveryDestinationResponse>(
        'GetDeliveryDestination',
        getDeliveryDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeliveryDestinationRequest.fromBuffer(value),
        ($0.GetDeliveryDestinationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeliveryDestinationPolicyRequest,
            $0.GetDeliveryDestinationPolicyResponse>(
        'GetDeliveryDestinationPolicy',
        getDeliveryDestinationPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeliveryDestinationPolicyRequest.fromBuffer(value),
        ($0.GetDeliveryDestinationPolicyResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeliverySourceRequest,
            $0.GetDeliverySourceResponse>(
        'GetDeliverySource',
        getDeliverySource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeliverySourceRequest.fromBuffer(value),
        ($0.GetDeliverySourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIntegrationRequest,
            $0.GetIntegrationResponse>(
        'GetIntegration',
        getIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetIntegrationRequest.fromBuffer(value),
        ($0.GetIntegrationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetLogAnomalyDetectorRequest,
            $0.GetLogAnomalyDetectorResponse>(
        'GetLogAnomalyDetector',
        getLogAnomalyDetector_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetLogAnomalyDetectorRequest.fromBuffer(value),
        ($0.GetLogAnomalyDetectorResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetLogEventsRequest, $0.GetLogEventsResponse>(
            'GetLogEvents',
            getLogEvents_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetLogEventsRequest.fromBuffer(value),
            ($0.GetLogEventsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetLogFieldsRequest, $0.GetLogFieldsResponse>(
            'GetLogFields',
            getLogFields_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetLogFieldsRequest.fromBuffer(value),
            ($0.GetLogFieldsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetLogGroupFieldsRequest,
            $0.GetLogGroupFieldsResponse>(
        'GetLogGroupFields',
        getLogGroupFields_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetLogGroupFieldsRequest.fromBuffer(value),
        ($0.GetLogGroupFieldsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetLogObjectRequest, $0.GetLogObjectResponse>(
            'GetLogObject',
            getLogObject_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetLogObjectRequest.fromBuffer(value),
            ($0.GetLogObjectResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetLogRecordRequest, $0.GetLogRecordResponse>(
            'GetLogRecord',
            getLogRecord_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetLogRecordRequest.fromBuffer(value),
            ($0.GetLogRecordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueryResultsRequest,
            $0.GetQueryResultsResponse>(
        'GetQueryResults',
        getQueryResults_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueryResultsRequest.fromBuffer(value),
        ($0.GetQueryResultsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetScheduledQueryRequest,
            $0.GetScheduledQueryResponse>(
        'GetScheduledQuery',
        getScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetScheduledQueryRequest.fromBuffer(value),
        ($0.GetScheduledQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetScheduledQueryHistoryRequest,
            $0.GetScheduledQueryHistoryResponse>(
        'GetScheduledQueryHistory',
        getScheduledQueryHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetScheduledQueryHistoryRequest.fromBuffer(value),
        ($0.GetScheduledQueryHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTransformerRequest,
            $0.GetTransformerResponse>(
        'GetTransformer',
        getTransformer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTransformerRequest.fromBuffer(value),
        ($0.GetTransformerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAggregateLogGroupSummariesRequest,
            $0.ListAggregateLogGroupSummariesResponse>(
        'ListAggregateLogGroupSummaries',
        listAggregateLogGroupSummaries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAggregateLogGroupSummariesRequest.fromBuffer(value),
        ($0.ListAggregateLogGroupSummariesResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListAnomaliesRequest, $0.ListAnomaliesResponse>(
            'ListAnomalies',
            listAnomalies_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListAnomaliesRequest.fromBuffer(value),
            ($0.ListAnomaliesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListIntegrationsRequest,
            $0.ListIntegrationsResponse>(
        'ListIntegrations',
        listIntegrations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListIntegrationsRequest.fromBuffer(value),
        ($0.ListIntegrationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListLogAnomalyDetectorsRequest,
            $0.ListLogAnomalyDetectorsResponse>(
        'ListLogAnomalyDetectors',
        listLogAnomalyDetectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListLogAnomalyDetectorsRequest.fromBuffer(value),
        ($0.ListLogAnomalyDetectorsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListLogGroupsRequest, $0.ListLogGroupsResponse>(
            'ListLogGroups',
            listLogGroups_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListLogGroupsRequest.fromBuffer(value),
            ($0.ListLogGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListLogGroupsForQueryRequest,
            $0.ListLogGroupsForQueryResponse>(
        'ListLogGroupsForQuery',
        listLogGroupsForQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListLogGroupsForQueryRequest.fromBuffer(value),
        ($0.ListLogGroupsForQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListScheduledQueriesRequest,
            $0.ListScheduledQueriesResponse>(
        'ListScheduledQueries',
        listScheduledQueries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListScheduledQueriesRequest.fromBuffer(value),
        ($0.ListScheduledQueriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSourcesForS3TableIntegrationRequest,
            $0.ListSourcesForS3TableIntegrationResponse>(
        'ListSourcesForS3TableIntegration',
        listSourcesForS3TableIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSourcesForS3TableIntegrationRequest.fromBuffer(value),
        ($0.ListSourcesForS3TableIntegrationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsLogGroupRequest,
            $0.ListTagsLogGroupResponse>(
        'ListTagsLogGroup',
        listTagsLogGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsLogGroupRequest.fromBuffer(value),
        ($0.ListTagsLogGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountPolicyRequest,
            $0.PutAccountPolicyResponse>(
        'PutAccountPolicy',
        putAccountPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountPolicyRequest.fromBuffer(value),
        ($0.PutAccountPolicyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutBearerTokenAuthenticationRequest, $1.Empty>(
            'PutBearerTokenAuthentication',
            putBearerTokenAuthentication_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutBearerTokenAuthenticationRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDataProtectionPolicyRequest,
            $0.PutDataProtectionPolicyResponse>(
        'PutDataProtectionPolicy',
        putDataProtectionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDataProtectionPolicyRequest.fromBuffer(value),
        ($0.PutDataProtectionPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDeliveryDestinationRequest,
            $0.PutDeliveryDestinationResponse>(
        'PutDeliveryDestination',
        putDeliveryDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDeliveryDestinationRequest.fromBuffer(value),
        ($0.PutDeliveryDestinationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDeliveryDestinationPolicyRequest,
            $0.PutDeliveryDestinationPolicyResponse>(
        'PutDeliveryDestinationPolicy',
        putDeliveryDestinationPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDeliveryDestinationPolicyRequest.fromBuffer(value),
        ($0.PutDeliveryDestinationPolicyResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDeliverySourceRequest,
            $0.PutDeliverySourceResponse>(
        'PutDeliverySource',
        putDeliverySource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDeliverySourceRequest.fromBuffer(value),
        ($0.PutDeliverySourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDestinationRequest,
            $0.PutDestinationResponse>(
        'PutDestination',
        putDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDestinationRequest.fromBuffer(value),
        ($0.PutDestinationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDestinationPolicyRequest, $1.Empty>(
        'PutDestinationPolicy',
        putDestinationPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDestinationPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutIndexPolicyRequest,
            $0.PutIndexPolicyResponse>(
        'PutIndexPolicy',
        putIndexPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutIndexPolicyRequest.fromBuffer(value),
        ($0.PutIndexPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutIntegrationRequest,
            $0.PutIntegrationResponse>(
        'PutIntegration',
        putIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutIntegrationRequest.fromBuffer(value),
        ($0.PutIntegrationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutLogEventsRequest, $0.PutLogEventsResponse>(
            'PutLogEvents',
            putLogEvents_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutLogEventsRequest.fromBuffer(value),
            ($0.PutLogEventsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutLogGroupDeletionProtectionRequest, $1.Empty>(
            'PutLogGroupDeletionProtection',
            putLogGroupDeletionProtection_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutLogGroupDeletionProtectionRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutMetricFilterRequest, $1.Empty>(
        'PutMetricFilter',
        putMetricFilter_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutMetricFilterRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutQueryDefinitionRequest,
            $0.PutQueryDefinitionResponse>(
        'PutQueryDefinition',
        putQueryDefinition_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutQueryDefinitionRequest.fromBuffer(value),
        ($0.PutQueryDefinitionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyRequest,
            $0.PutResourcePolicyResponse>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyRequest.fromBuffer(value),
        ($0.PutResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRetentionPolicyRequest, $1.Empty>(
        'PutRetentionPolicy',
        putRetentionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutRetentionPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutSubscriptionFilterRequest, $1.Empty>(
        'PutSubscriptionFilter',
        putSubscriptionFilter_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutSubscriptionFilterRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutTransformerRequest, $1.Empty>(
        'PutTransformer',
        putTransformer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutTransformerRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartLiveTailRequest, $0.StartLiveTailResponse>(
            'StartLiveTail',
            startLiveTail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartLiveTailRequest.fromBuffer(value),
            ($0.StartLiveTailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartQueryRequest, $0.StartQueryResponse>(
        'StartQuery',
        startQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.StartQueryRequest.fromBuffer(value),
        ($0.StartQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopQueryRequest, $0.StopQueryResponse>(
        'StopQuery',
        stopQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.StopQueryRequest.fromBuffer(value),
        ($0.StopQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagLogGroupRequest, $1.Empty>(
        'TagLogGroup',
        tagLogGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagLogGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceRequest, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestMetricFilterRequest,
            $0.TestMetricFilterResponse>(
        'TestMetricFilter',
        testMetricFilter_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestMetricFilterRequest.fromBuffer(value),
        ($0.TestMetricFilterResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestTransformerRequest,
            $0.TestTransformerResponse>(
        'TestTransformer',
        testTransformer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestTransformerRequest.fromBuffer(value),
        ($0.TestTransformerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagLogGroupRequest, $1.Empty>(
        'UntagLogGroup',
        untagLogGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagLogGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceRequest, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAnomalyRequest, $1.Empty>(
        'UpdateAnomaly',
        updateAnomaly_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAnomalyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDeliveryConfigurationRequest,
            $0.UpdateDeliveryConfigurationResponse>(
        'UpdateDeliveryConfiguration',
        updateDeliveryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDeliveryConfigurationRequest.fromBuffer(value),
        ($0.UpdateDeliveryConfigurationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateLogAnomalyDetectorRequest, $1.Empty>(
            'UpdateLogAnomalyDetector',
            updateLogAnomalyDetector_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateLogAnomalyDetectorRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateScheduledQueryRequest,
            $0.UpdateScheduledQueryResponse>(
        'UpdateScheduledQuery',
        updateScheduledQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateScheduledQueryRequest.fromBuffer(value),
        ($0.UpdateScheduledQueryResponse value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> associateKmsKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AssociateKmsKeyRequest> $request) async {
    return associateKmsKey($call, await $request);
  }

  $async.Future<$1.Empty> associateKmsKey(
      $grpc.ServiceCall call, $0.AssociateKmsKeyRequest request);

  $async.Future<$0.AssociateSourceToS3TableIntegrationResponse>
      associateSourceToS3TableIntegration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.AssociateSourceToS3TableIntegrationRequest>
              $request) async {
    return associateSourceToS3TableIntegration($call, await $request);
  }

  $async.Future<$0.AssociateSourceToS3TableIntegrationResponse>
      associateSourceToS3TableIntegration($grpc.ServiceCall call,
          $0.AssociateSourceToS3TableIntegrationRequest request);

  $async.Future<$1.Empty> cancelExportTask_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CancelExportTaskRequest> $request) async {
    return cancelExportTask($call, await $request);
  }

  $async.Future<$1.Empty> cancelExportTask(
      $grpc.ServiceCall call, $0.CancelExportTaskRequest request);

  $async.Future<$0.CancelImportTaskResponse> cancelImportTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelImportTaskRequest> $request) async {
    return cancelImportTask($call, await $request);
  }

  $async.Future<$0.CancelImportTaskResponse> cancelImportTask(
      $grpc.ServiceCall call, $0.CancelImportTaskRequest request);

  $async.Future<$0.CreateDeliveryResponse> createDelivery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDeliveryRequest> $request) async {
    return createDelivery($call, await $request);
  }

  $async.Future<$0.CreateDeliveryResponse> createDelivery(
      $grpc.ServiceCall call, $0.CreateDeliveryRequest request);

  $async.Future<$0.CreateExportTaskResponse> createExportTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateExportTaskRequest> $request) async {
    return createExportTask($call, await $request);
  }

  $async.Future<$0.CreateExportTaskResponse> createExportTask(
      $grpc.ServiceCall call, $0.CreateExportTaskRequest request);

  $async.Future<$0.CreateImportTaskResponse> createImportTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateImportTaskRequest> $request) async {
    return createImportTask($call, await $request);
  }

  $async.Future<$0.CreateImportTaskResponse> createImportTask(
      $grpc.ServiceCall call, $0.CreateImportTaskRequest request);

  $async.Future<$0.CreateLogAnomalyDetectorResponse>
      createLogAnomalyDetector_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateLogAnomalyDetectorRequest> $request) async {
    return createLogAnomalyDetector($call, await $request);
  }

  $async.Future<$0.CreateLogAnomalyDetectorResponse> createLogAnomalyDetector(
      $grpc.ServiceCall call, $0.CreateLogAnomalyDetectorRequest request);

  $async.Future<$1.Empty> createLogGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateLogGroupRequest> $request) async {
    return createLogGroup($call, await $request);
  }

  $async.Future<$1.Empty> createLogGroup(
      $grpc.ServiceCall call, $0.CreateLogGroupRequest request);

  $async.Future<$1.Empty> createLogStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateLogStreamRequest> $request) async {
    return createLogStream($call, await $request);
  }

  $async.Future<$1.Empty> createLogStream(
      $grpc.ServiceCall call, $0.CreateLogStreamRequest request);

  $async.Future<$0.CreateScheduledQueryResponse> createScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateScheduledQueryRequest> $request) async {
    return createScheduledQuery($call, await $request);
  }

  $async.Future<$0.CreateScheduledQueryResponse> createScheduledQuery(
      $grpc.ServiceCall call, $0.CreateScheduledQueryRequest request);

  $async.Future<$1.Empty> deleteAccountPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAccountPolicyRequest> $request) async {
    return deleteAccountPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteAccountPolicy(
      $grpc.ServiceCall call, $0.DeleteAccountPolicyRequest request);

  $async.Future<$1.Empty> deleteDataProtectionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDataProtectionPolicyRequest> $request) async {
    return deleteDataProtectionPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteDataProtectionPolicy(
      $grpc.ServiceCall call, $0.DeleteDataProtectionPolicyRequest request);

  $async.Future<$1.Empty> deleteDelivery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDeliveryRequest> $request) async {
    return deleteDelivery($call, await $request);
  }

  $async.Future<$1.Empty> deleteDelivery(
      $grpc.ServiceCall call, $0.DeleteDeliveryRequest request);

  $async.Future<$1.Empty> deleteDeliveryDestination_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDeliveryDestinationRequest> $request) async {
    return deleteDeliveryDestination($call, await $request);
  }

  $async.Future<$1.Empty> deleteDeliveryDestination(
      $grpc.ServiceCall call, $0.DeleteDeliveryDestinationRequest request);

  $async.Future<$1.Empty> deleteDeliveryDestinationPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDeliveryDestinationPolicyRequest> $request) async {
    return deleteDeliveryDestinationPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteDeliveryDestinationPolicy(
      $grpc.ServiceCall call,
      $0.DeleteDeliveryDestinationPolicyRequest request);

  $async.Future<$1.Empty> deleteDeliverySource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDeliverySourceRequest> $request) async {
    return deleteDeliverySource($call, await $request);
  }

  $async.Future<$1.Empty> deleteDeliverySource(
      $grpc.ServiceCall call, $0.DeleteDeliverySourceRequest request);

  $async.Future<$1.Empty> deleteDestination_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDestinationRequest> $request) async {
    return deleteDestination($call, await $request);
  }

  $async.Future<$1.Empty> deleteDestination(
      $grpc.ServiceCall call, $0.DeleteDestinationRequest request);

  $async.Future<$0.DeleteIndexPolicyResponse> deleteIndexPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteIndexPolicyRequest> $request) async {
    return deleteIndexPolicy($call, await $request);
  }

  $async.Future<$0.DeleteIndexPolicyResponse> deleteIndexPolicy(
      $grpc.ServiceCall call, $0.DeleteIndexPolicyRequest request);

  $async.Future<$0.DeleteIntegrationResponse> deleteIntegration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteIntegrationRequest> $request) async {
    return deleteIntegration($call, await $request);
  }

  $async.Future<$0.DeleteIntegrationResponse> deleteIntegration(
      $grpc.ServiceCall call, $0.DeleteIntegrationRequest request);

  $async.Future<$1.Empty> deleteLogAnomalyDetector_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteLogAnomalyDetectorRequest> $request) async {
    return deleteLogAnomalyDetector($call, await $request);
  }

  $async.Future<$1.Empty> deleteLogAnomalyDetector(
      $grpc.ServiceCall call, $0.DeleteLogAnomalyDetectorRequest request);

  $async.Future<$1.Empty> deleteLogGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteLogGroupRequest> $request) async {
    return deleteLogGroup($call, await $request);
  }

  $async.Future<$1.Empty> deleteLogGroup(
      $grpc.ServiceCall call, $0.DeleteLogGroupRequest request);

  $async.Future<$1.Empty> deleteLogStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteLogStreamRequest> $request) async {
    return deleteLogStream($call, await $request);
  }

  $async.Future<$1.Empty> deleteLogStream(
      $grpc.ServiceCall call, $0.DeleteLogStreamRequest request);

  $async.Future<$1.Empty> deleteMetricFilter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteMetricFilterRequest> $request) async {
    return deleteMetricFilter($call, await $request);
  }

  $async.Future<$1.Empty> deleteMetricFilter(
      $grpc.ServiceCall call, $0.DeleteMetricFilterRequest request);

  $async.Future<$0.DeleteQueryDefinitionResponse> deleteQueryDefinition_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteQueryDefinitionRequest> $request) async {
    return deleteQueryDefinition($call, await $request);
  }

  $async.Future<$0.DeleteQueryDefinitionResponse> deleteQueryDefinition(
      $grpc.ServiceCall call, $0.DeleteQueryDefinitionRequest request);

  $async.Future<$1.Empty> deleteResourcePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyRequest> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyRequest request);

  $async.Future<$1.Empty> deleteRetentionPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRetentionPolicyRequest> $request) async {
    return deleteRetentionPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteRetentionPolicy(
      $grpc.ServiceCall call, $0.DeleteRetentionPolicyRequest request);

  $async.Future<$0.DeleteScheduledQueryResponse> deleteScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteScheduledQueryRequest> $request) async {
    return deleteScheduledQuery($call, await $request);
  }

  $async.Future<$0.DeleteScheduledQueryResponse> deleteScheduledQuery(
      $grpc.ServiceCall call, $0.DeleteScheduledQueryRequest request);

  $async.Future<$1.Empty> deleteSubscriptionFilter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteSubscriptionFilterRequest> $request) async {
    return deleteSubscriptionFilter($call, await $request);
  }

  $async.Future<$1.Empty> deleteSubscriptionFilter(
      $grpc.ServiceCall call, $0.DeleteSubscriptionFilterRequest request);

  $async.Future<$1.Empty> deleteTransformer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTransformerRequest> $request) async {
    return deleteTransformer($call, await $request);
  }

  $async.Future<$1.Empty> deleteTransformer(
      $grpc.ServiceCall call, $0.DeleteTransformerRequest request);

  $async.Future<$0.DescribeAccountPoliciesResponse> describeAccountPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAccountPoliciesRequest> $request) async {
    return describeAccountPolicies($call, await $request);
  }

  $async.Future<$0.DescribeAccountPoliciesResponse> describeAccountPolicies(
      $grpc.ServiceCall call, $0.DescribeAccountPoliciesRequest request);

  $async.Future<$0.DescribeConfigurationTemplatesResponse>
      describeConfigurationTemplates_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeConfigurationTemplatesRequest>
              $request) async {
    return describeConfigurationTemplates($call, await $request);
  }

  $async.Future<$0.DescribeConfigurationTemplatesResponse>
      describeConfigurationTemplates($grpc.ServiceCall call,
          $0.DescribeConfigurationTemplatesRequest request);

  $async.Future<$0.DescribeDeliveriesResponse> describeDeliveries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeDeliveriesRequest> $request) async {
    return describeDeliveries($call, await $request);
  }

  $async.Future<$0.DescribeDeliveriesResponse> describeDeliveries(
      $grpc.ServiceCall call, $0.DescribeDeliveriesRequest request);

  $async.Future<$0.DescribeDeliveryDestinationsResponse>
      describeDeliveryDestinations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeDeliveryDestinationsRequest>
              $request) async {
    return describeDeliveryDestinations($call, await $request);
  }

  $async.Future<$0.DescribeDeliveryDestinationsResponse>
      describeDeliveryDestinations($grpc.ServiceCall call,
          $0.DescribeDeliveryDestinationsRequest request);

  $async.Future<$0.DescribeDeliverySourcesResponse> describeDeliverySources_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeDeliverySourcesRequest> $request) async {
    return describeDeliverySources($call, await $request);
  }

  $async.Future<$0.DescribeDeliverySourcesResponse> describeDeliverySources(
      $grpc.ServiceCall call, $0.DescribeDeliverySourcesRequest request);

  $async.Future<$0.DescribeDestinationsResponse> describeDestinations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeDestinationsRequest> $request) async {
    return describeDestinations($call, await $request);
  }

  $async.Future<$0.DescribeDestinationsResponse> describeDestinations(
      $grpc.ServiceCall call, $0.DescribeDestinationsRequest request);

  $async.Future<$0.DescribeExportTasksResponse> describeExportTasks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeExportTasksRequest> $request) async {
    return describeExportTasks($call, await $request);
  }

  $async.Future<$0.DescribeExportTasksResponse> describeExportTasks(
      $grpc.ServiceCall call, $0.DescribeExportTasksRequest request);

  $async.Future<$0.DescribeFieldIndexesResponse> describeFieldIndexes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeFieldIndexesRequest> $request) async {
    return describeFieldIndexes($call, await $request);
  }

  $async.Future<$0.DescribeFieldIndexesResponse> describeFieldIndexes(
      $grpc.ServiceCall call, $0.DescribeFieldIndexesRequest request);

  $async.Future<$0.DescribeImportTaskBatchesResponse>
      describeImportTaskBatches_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeImportTaskBatchesRequest> $request) async {
    return describeImportTaskBatches($call, await $request);
  }

  $async.Future<$0.DescribeImportTaskBatchesResponse> describeImportTaskBatches(
      $grpc.ServiceCall call, $0.DescribeImportTaskBatchesRequest request);

  $async.Future<$0.DescribeImportTasksResponse> describeImportTasks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeImportTasksRequest> $request) async {
    return describeImportTasks($call, await $request);
  }

  $async.Future<$0.DescribeImportTasksResponse> describeImportTasks(
      $grpc.ServiceCall call, $0.DescribeImportTasksRequest request);

  $async.Future<$0.DescribeIndexPoliciesResponse> describeIndexPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeIndexPoliciesRequest> $request) async {
    return describeIndexPolicies($call, await $request);
  }

  $async.Future<$0.DescribeIndexPoliciesResponse> describeIndexPolicies(
      $grpc.ServiceCall call, $0.DescribeIndexPoliciesRequest request);

  $async.Future<$0.DescribeLogGroupsResponse> describeLogGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeLogGroupsRequest> $request) async {
    return describeLogGroups($call, await $request);
  }

  $async.Future<$0.DescribeLogGroupsResponse> describeLogGroups(
      $grpc.ServiceCall call, $0.DescribeLogGroupsRequest request);

  $async.Future<$0.DescribeLogStreamsResponse> describeLogStreams_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeLogStreamsRequest> $request) async {
    return describeLogStreams($call, await $request);
  }

  $async.Future<$0.DescribeLogStreamsResponse> describeLogStreams(
      $grpc.ServiceCall call, $0.DescribeLogStreamsRequest request);

  $async.Future<$0.DescribeMetricFiltersResponse> describeMetricFilters_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeMetricFiltersRequest> $request) async {
    return describeMetricFilters($call, await $request);
  }

  $async.Future<$0.DescribeMetricFiltersResponse> describeMetricFilters(
      $grpc.ServiceCall call, $0.DescribeMetricFiltersRequest request);

  $async.Future<$0.DescribeQueriesResponse> describeQueries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeQueriesRequest> $request) async {
    return describeQueries($call, await $request);
  }

  $async.Future<$0.DescribeQueriesResponse> describeQueries(
      $grpc.ServiceCall call, $0.DescribeQueriesRequest request);

  $async.Future<$0.DescribeQueryDefinitionsResponse>
      describeQueryDefinitions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeQueryDefinitionsRequest> $request) async {
    return describeQueryDefinitions($call, await $request);
  }

  $async.Future<$0.DescribeQueryDefinitionsResponse> describeQueryDefinitions(
      $grpc.ServiceCall call, $0.DescribeQueryDefinitionsRequest request);

  $async.Future<$0.DescribeResourcePoliciesResponse>
      describeResourcePolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeResourcePoliciesRequest> $request) async {
    return describeResourcePolicies($call, await $request);
  }

  $async.Future<$0.DescribeResourcePoliciesResponse> describeResourcePolicies(
      $grpc.ServiceCall call, $0.DescribeResourcePoliciesRequest request);

  $async.Future<$0.DescribeSubscriptionFiltersResponse>
      describeSubscriptionFilters_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeSubscriptionFiltersRequest> $request) async {
    return describeSubscriptionFilters($call, await $request);
  }

  $async.Future<$0.DescribeSubscriptionFiltersResponse>
      describeSubscriptionFilters($grpc.ServiceCall call,
          $0.DescribeSubscriptionFiltersRequest request);

  $async.Future<$1.Empty> disassociateKmsKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisassociateKmsKeyRequest> $request) async {
    return disassociateKmsKey($call, await $request);
  }

  $async.Future<$1.Empty> disassociateKmsKey(
      $grpc.ServiceCall call, $0.DisassociateKmsKeyRequest request);

  $async.Future<$0.DisassociateSourceFromS3TableIntegrationResponse>
      disassociateSourceFromS3TableIntegration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisassociateSourceFromS3TableIntegrationRequest>
              $request) async {
    return disassociateSourceFromS3TableIntegration($call, await $request);
  }

  $async.Future<$0.DisassociateSourceFromS3TableIntegrationResponse>
      disassociateSourceFromS3TableIntegration($grpc.ServiceCall call,
          $0.DisassociateSourceFromS3TableIntegrationRequest request);

  $async.Future<$0.FilterLogEventsResponse> filterLogEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.FilterLogEventsRequest> $request) async {
    return filterLogEvents($call, await $request);
  }

  $async.Future<$0.FilterLogEventsResponse> filterLogEvents(
      $grpc.ServiceCall call, $0.FilterLogEventsRequest request);

  $async.Future<$0.GetDataProtectionPolicyResponse> getDataProtectionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDataProtectionPolicyRequest> $request) async {
    return getDataProtectionPolicy($call, await $request);
  }

  $async.Future<$0.GetDataProtectionPolicyResponse> getDataProtectionPolicy(
      $grpc.ServiceCall call, $0.GetDataProtectionPolicyRequest request);

  $async.Future<$0.GetDeliveryResponse> getDelivery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDeliveryRequest> $request) async {
    return getDelivery($call, await $request);
  }

  $async.Future<$0.GetDeliveryResponse> getDelivery(
      $grpc.ServiceCall call, $0.GetDeliveryRequest request);

  $async.Future<$0.GetDeliveryDestinationResponse> getDeliveryDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDeliveryDestinationRequest> $request) async {
    return getDeliveryDestination($call, await $request);
  }

  $async.Future<$0.GetDeliveryDestinationResponse> getDeliveryDestination(
      $grpc.ServiceCall call, $0.GetDeliveryDestinationRequest request);

  $async.Future<$0.GetDeliveryDestinationPolicyResponse>
      getDeliveryDestinationPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDeliveryDestinationPolicyRequest>
              $request) async {
    return getDeliveryDestinationPolicy($call, await $request);
  }

  $async.Future<$0.GetDeliveryDestinationPolicyResponse>
      getDeliveryDestinationPolicy($grpc.ServiceCall call,
          $0.GetDeliveryDestinationPolicyRequest request);

  $async.Future<$0.GetDeliverySourceResponse> getDeliverySource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDeliverySourceRequest> $request) async {
    return getDeliverySource($call, await $request);
  }

  $async.Future<$0.GetDeliverySourceResponse> getDeliverySource(
      $grpc.ServiceCall call, $0.GetDeliverySourceRequest request);

  $async.Future<$0.GetIntegrationResponse> getIntegration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetIntegrationRequest> $request) async {
    return getIntegration($call, await $request);
  }

  $async.Future<$0.GetIntegrationResponse> getIntegration(
      $grpc.ServiceCall call, $0.GetIntegrationRequest request);

  $async.Future<$0.GetLogAnomalyDetectorResponse> getLogAnomalyDetector_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogAnomalyDetectorRequest> $request) async {
    return getLogAnomalyDetector($call, await $request);
  }

  $async.Future<$0.GetLogAnomalyDetectorResponse> getLogAnomalyDetector(
      $grpc.ServiceCall call, $0.GetLogAnomalyDetectorRequest request);

  $async.Future<$0.GetLogEventsResponse> getLogEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogEventsRequest> $request) async {
    return getLogEvents($call, await $request);
  }

  $async.Future<$0.GetLogEventsResponse> getLogEvents(
      $grpc.ServiceCall call, $0.GetLogEventsRequest request);

  $async.Future<$0.GetLogFieldsResponse> getLogFields_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogFieldsRequest> $request) async {
    return getLogFields($call, await $request);
  }

  $async.Future<$0.GetLogFieldsResponse> getLogFields(
      $grpc.ServiceCall call, $0.GetLogFieldsRequest request);

  $async.Future<$0.GetLogGroupFieldsResponse> getLogGroupFields_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogGroupFieldsRequest> $request) async {
    return getLogGroupFields($call, await $request);
  }

  $async.Future<$0.GetLogGroupFieldsResponse> getLogGroupFields(
      $grpc.ServiceCall call, $0.GetLogGroupFieldsRequest request);

  $async.Future<$0.GetLogObjectResponse> getLogObject_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogObjectRequest> $request) async {
    return getLogObject($call, await $request);
  }

  $async.Future<$0.GetLogObjectResponse> getLogObject(
      $grpc.ServiceCall call, $0.GetLogObjectRequest request);

  $async.Future<$0.GetLogRecordResponse> getLogRecord_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLogRecordRequest> $request) async {
    return getLogRecord($call, await $request);
  }

  $async.Future<$0.GetLogRecordResponse> getLogRecord(
      $grpc.ServiceCall call, $0.GetLogRecordRequest request);

  $async.Future<$0.GetQueryResultsResponse> getQueryResults_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueryResultsRequest> $request) async {
    return getQueryResults($call, await $request);
  }

  $async.Future<$0.GetQueryResultsResponse> getQueryResults(
      $grpc.ServiceCall call, $0.GetQueryResultsRequest request);

  $async.Future<$0.GetScheduledQueryResponse> getScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetScheduledQueryRequest> $request) async {
    return getScheduledQuery($call, await $request);
  }

  $async.Future<$0.GetScheduledQueryResponse> getScheduledQuery(
      $grpc.ServiceCall call, $0.GetScheduledQueryRequest request);

  $async.Future<$0.GetScheduledQueryHistoryResponse>
      getScheduledQueryHistory_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetScheduledQueryHistoryRequest> $request) async {
    return getScheduledQueryHistory($call, await $request);
  }

  $async.Future<$0.GetScheduledQueryHistoryResponse> getScheduledQueryHistory(
      $grpc.ServiceCall call, $0.GetScheduledQueryHistoryRequest request);

  $async.Future<$0.GetTransformerResponse> getTransformer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTransformerRequest> $request) async {
    return getTransformer($call, await $request);
  }

  $async.Future<$0.GetTransformerResponse> getTransformer(
      $grpc.ServiceCall call, $0.GetTransformerRequest request);

  $async.Future<$0.ListAggregateLogGroupSummariesResponse>
      listAggregateLogGroupSummaries_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListAggregateLogGroupSummariesRequest>
              $request) async {
    return listAggregateLogGroupSummaries($call, await $request);
  }

  $async.Future<$0.ListAggregateLogGroupSummariesResponse>
      listAggregateLogGroupSummaries($grpc.ServiceCall call,
          $0.ListAggregateLogGroupSummariesRequest request);

  $async.Future<$0.ListAnomaliesResponse> listAnomalies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAnomaliesRequest> $request) async {
    return listAnomalies($call, await $request);
  }

  $async.Future<$0.ListAnomaliesResponse> listAnomalies(
      $grpc.ServiceCall call, $0.ListAnomaliesRequest request);

  $async.Future<$0.ListIntegrationsResponse> listIntegrations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListIntegrationsRequest> $request) async {
    return listIntegrations($call, await $request);
  }

  $async.Future<$0.ListIntegrationsResponse> listIntegrations(
      $grpc.ServiceCall call, $0.ListIntegrationsRequest request);

  $async.Future<$0.ListLogAnomalyDetectorsResponse> listLogAnomalyDetectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListLogAnomalyDetectorsRequest> $request) async {
    return listLogAnomalyDetectors($call, await $request);
  }

  $async.Future<$0.ListLogAnomalyDetectorsResponse> listLogAnomalyDetectors(
      $grpc.ServiceCall call, $0.ListLogAnomalyDetectorsRequest request);

  $async.Future<$0.ListLogGroupsResponse> listLogGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListLogGroupsRequest> $request) async {
    return listLogGroups($call, await $request);
  }

  $async.Future<$0.ListLogGroupsResponse> listLogGroups(
      $grpc.ServiceCall call, $0.ListLogGroupsRequest request);

  $async.Future<$0.ListLogGroupsForQueryResponse> listLogGroupsForQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListLogGroupsForQueryRequest> $request) async {
    return listLogGroupsForQuery($call, await $request);
  }

  $async.Future<$0.ListLogGroupsForQueryResponse> listLogGroupsForQuery(
      $grpc.ServiceCall call, $0.ListLogGroupsForQueryRequest request);

  $async.Future<$0.ListScheduledQueriesResponse> listScheduledQueries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListScheduledQueriesRequest> $request) async {
    return listScheduledQueries($call, await $request);
  }

  $async.Future<$0.ListScheduledQueriesResponse> listScheduledQueries(
      $grpc.ServiceCall call, $0.ListScheduledQueriesRequest request);

  $async.Future<$0.ListSourcesForS3TableIntegrationResponse>
      listSourcesForS3TableIntegration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListSourcesForS3TableIntegrationRequest>
              $request) async {
    return listSourcesForS3TableIntegration($call, await $request);
  }

  $async.Future<$0.ListSourcesForS3TableIntegrationResponse>
      listSourcesForS3TableIntegration($grpc.ServiceCall call,
          $0.ListSourcesForS3TableIntegrationRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTagsLogGroupResponse> listTagsLogGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsLogGroupRequest> $request) async {
    return listTagsLogGroup($call, await $request);
  }

  $async.Future<$0.ListTagsLogGroupResponse> listTagsLogGroup(
      $grpc.ServiceCall call, $0.ListTagsLogGroupRequest request);

  $async.Future<$0.PutAccountPolicyResponse> putAccountPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutAccountPolicyRequest> $request) async {
    return putAccountPolicy($call, await $request);
  }

  $async.Future<$0.PutAccountPolicyResponse> putAccountPolicy(
      $grpc.ServiceCall call, $0.PutAccountPolicyRequest request);

  $async.Future<$1.Empty> putBearerTokenAuthentication_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBearerTokenAuthenticationRequest> $request) async {
    return putBearerTokenAuthentication($call, await $request);
  }

  $async.Future<$1.Empty> putBearerTokenAuthentication(
      $grpc.ServiceCall call, $0.PutBearerTokenAuthenticationRequest request);

  $async.Future<$0.PutDataProtectionPolicyResponse> putDataProtectionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutDataProtectionPolicyRequest> $request) async {
    return putDataProtectionPolicy($call, await $request);
  }

  $async.Future<$0.PutDataProtectionPolicyResponse> putDataProtectionPolicy(
      $grpc.ServiceCall call, $0.PutDataProtectionPolicyRequest request);

  $async.Future<$0.PutDeliveryDestinationResponse> putDeliveryDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutDeliveryDestinationRequest> $request) async {
    return putDeliveryDestination($call, await $request);
  }

  $async.Future<$0.PutDeliveryDestinationResponse> putDeliveryDestination(
      $grpc.ServiceCall call, $0.PutDeliveryDestinationRequest request);

  $async.Future<$0.PutDeliveryDestinationPolicyResponse>
      putDeliveryDestinationPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutDeliveryDestinationPolicyRequest>
              $request) async {
    return putDeliveryDestinationPolicy($call, await $request);
  }

  $async.Future<$0.PutDeliveryDestinationPolicyResponse>
      putDeliveryDestinationPolicy($grpc.ServiceCall call,
          $0.PutDeliveryDestinationPolicyRequest request);

  $async.Future<$0.PutDeliverySourceResponse> putDeliverySource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutDeliverySourceRequest> $request) async {
    return putDeliverySource($call, await $request);
  }

  $async.Future<$0.PutDeliverySourceResponse> putDeliverySource(
      $grpc.ServiceCall call, $0.PutDeliverySourceRequest request);

  $async.Future<$0.PutDestinationResponse> putDestination_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutDestinationRequest> $request) async {
    return putDestination($call, await $request);
  }

  $async.Future<$0.PutDestinationResponse> putDestination(
      $grpc.ServiceCall call, $0.PutDestinationRequest request);

  $async.Future<$1.Empty> putDestinationPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutDestinationPolicyRequest> $request) async {
    return putDestinationPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putDestinationPolicy(
      $grpc.ServiceCall call, $0.PutDestinationPolicyRequest request);

  $async.Future<$0.PutIndexPolicyResponse> putIndexPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutIndexPolicyRequest> $request) async {
    return putIndexPolicy($call, await $request);
  }

  $async.Future<$0.PutIndexPolicyResponse> putIndexPolicy(
      $grpc.ServiceCall call, $0.PutIndexPolicyRequest request);

  $async.Future<$0.PutIntegrationResponse> putIntegration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutIntegrationRequest> $request) async {
    return putIntegration($call, await $request);
  }

  $async.Future<$0.PutIntegrationResponse> putIntegration(
      $grpc.ServiceCall call, $0.PutIntegrationRequest request);

  $async.Future<$0.PutLogEventsResponse> putLogEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutLogEventsRequest> $request) async {
    return putLogEvents($call, await $request);
  }

  $async.Future<$0.PutLogEventsResponse> putLogEvents(
      $grpc.ServiceCall call, $0.PutLogEventsRequest request);

  $async.Future<$1.Empty> putLogGroupDeletionProtection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutLogGroupDeletionProtectionRequest> $request) async {
    return putLogGroupDeletionProtection($call, await $request);
  }

  $async.Future<$1.Empty> putLogGroupDeletionProtection(
      $grpc.ServiceCall call, $0.PutLogGroupDeletionProtectionRequest request);

  $async.Future<$1.Empty> putMetricFilter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutMetricFilterRequest> $request) async {
    return putMetricFilter($call, await $request);
  }

  $async.Future<$1.Empty> putMetricFilter(
      $grpc.ServiceCall call, $0.PutMetricFilterRequest request);

  $async.Future<$0.PutQueryDefinitionResponse> putQueryDefinition_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutQueryDefinitionRequest> $request) async {
    return putQueryDefinition($call, await $request);
  }

  $async.Future<$0.PutQueryDefinitionResponse> putQueryDefinition(
      $grpc.ServiceCall call, $0.PutQueryDefinitionRequest request);

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyRequest> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyRequest request);

  $async.Future<$1.Empty> putRetentionPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRetentionPolicyRequest> $request) async {
    return putRetentionPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putRetentionPolicy(
      $grpc.ServiceCall call, $0.PutRetentionPolicyRequest request);

  $async.Future<$1.Empty> putSubscriptionFilter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutSubscriptionFilterRequest> $request) async {
    return putSubscriptionFilter($call, await $request);
  }

  $async.Future<$1.Empty> putSubscriptionFilter(
      $grpc.ServiceCall call, $0.PutSubscriptionFilterRequest request);

  $async.Future<$1.Empty> putTransformer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutTransformerRequest> $request) async {
    return putTransformer($call, await $request);
  }

  $async.Future<$1.Empty> putTransformer(
      $grpc.ServiceCall call, $0.PutTransformerRequest request);

  $async.Future<$0.StartLiveTailResponse> startLiveTail_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartLiveTailRequest> $request) async {
    return startLiveTail($call, await $request);
  }

  $async.Future<$0.StartLiveTailResponse> startLiveTail(
      $grpc.ServiceCall call, $0.StartLiveTailRequest request);

  $async.Future<$0.StartQueryResponse> startQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartQueryRequest> $request) async {
    return startQuery($call, await $request);
  }

  $async.Future<$0.StartQueryResponse> startQuery(
      $grpc.ServiceCall call, $0.StartQueryRequest request);

  $async.Future<$0.StopQueryResponse> stopQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StopQueryRequest> $request) async {
    return stopQuery($call, await $request);
  }

  $async.Future<$0.StopQueryResponse> stopQuery(
      $grpc.ServiceCall call, $0.StopQueryRequest request);

  $async.Future<$1.Empty> tagLogGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagLogGroupRequest> $request) async {
    return tagLogGroup($call, await $request);
  }

  $async.Future<$1.Empty> tagLogGroup(
      $grpc.ServiceCall call, $0.TagLogGroupRequest request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.TestMetricFilterResponse> testMetricFilter_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestMetricFilterRequest> $request) async {
    return testMetricFilter($call, await $request);
  }

  $async.Future<$0.TestMetricFilterResponse> testMetricFilter(
      $grpc.ServiceCall call, $0.TestMetricFilterRequest request);

  $async.Future<$0.TestTransformerResponse> testTransformer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestTransformerRequest> $request) async {
    return testTransformer($call, await $request);
  }

  $async.Future<$0.TestTransformerResponse> testTransformer(
      $grpc.ServiceCall call, $0.TestTransformerRequest request);

  $async.Future<$1.Empty> untagLogGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagLogGroupRequest> $request) async {
    return untagLogGroup($call, await $request);
  }

  $async.Future<$1.Empty> untagLogGroup(
      $grpc.ServiceCall call, $0.UntagLogGroupRequest request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$1.Empty> updateAnomaly_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAnomalyRequest> $request) async {
    return updateAnomaly($call, await $request);
  }

  $async.Future<$1.Empty> updateAnomaly(
      $grpc.ServiceCall call, $0.UpdateAnomalyRequest request);

  $async.Future<$0.UpdateDeliveryConfigurationResponse>
      updateDeliveryConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateDeliveryConfigurationRequest> $request) async {
    return updateDeliveryConfiguration($call, await $request);
  }

  $async.Future<$0.UpdateDeliveryConfigurationResponse>
      updateDeliveryConfiguration($grpc.ServiceCall call,
          $0.UpdateDeliveryConfigurationRequest request);

  $async.Future<$1.Empty> updateLogAnomalyDetector_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateLogAnomalyDetectorRequest> $request) async {
    return updateLogAnomalyDetector($call, await $request);
  }

  $async.Future<$1.Empty> updateLogAnomalyDetector(
      $grpc.ServiceCall call, $0.UpdateLogAnomalyDetectorRequest request);

  $async.Future<$0.UpdateScheduledQueryResponse> updateScheduledQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateScheduledQueryRequest> $request) async {
    return updateScheduledQuery($call, await $request);
  }

  $async.Future<$0.UpdateScheduledQueryResponse> updateScheduledQuery(
      $grpc.ServiceCall call, $0.UpdateScheduledQueryRequest request);
}
